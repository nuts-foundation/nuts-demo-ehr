package transfer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	sqlUtil "github.com/nuts-foundation/nuts-demo-ehr/sql"
	"github.com/nuts-foundation/nuts-node/vcr/credential"
	"github.com/sirupsen/logrus"
)

// SenderServiceName contains the name of the eOverdracht sender compound-service
const SenderServiceName = "eOverdracht-sender"

func (s service) Create(ctx context.Context, customerID int, dossierID string, description string, transferDate time.Time) (*domain.Transfer, error) {
	composition := fhir.BuildAdvanceNotice()
	err := s.localFHIRClientFactory(fhir.WithTenant(customerID)).CreateOrUpdate(ctx, composition)
	if err != nil {
		return nil, err
	}
	transfer, err := s.transferRepo.Create(ctx, customerID, dossierID, description, transferDate, composition["id"].(string))
	if err != nil {
		return nil, err
	}
	return transfer, nil
}

func (s service) CreateNegotiation(ctx context.Context, customerID int, transferID, organizationDID string, transferDate time.Time) (*domain.TransferNegotiation, error) {
	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil || customer.Did == nil {
		return nil, err
	}

	var (
		negotiation          *domain.TransferNegotiation
	)

	_, err = s.transferRepo.Update(ctx, customerID, transferID, func(transfer *domain.Transfer) (*domain.Transfer, error) {
		// Validate transfer
		if transfer.Status == domain.TransferStatusCancelled ||
			transfer.Status == domain.TransferStatusCompleted ||
			transfer.Status == domain.TransferStatusAssigned {
			return nil, errors.New("can't start new transfer negotiation when status is 'cancelled', 'assigned' or 'completed'")
		}

		// Create negotiation and share it to the other party
		// TODO: Share transaction to this repository call as well
		var err error
		// Pre-emptively resolve the receiver organization's notification endpoint to reduce clutter, avoiding to make FHIR tasks when the receiving party eOverdracht registration is faulty.
		_, err = s.registry.GetCompoundServiceEndpoint(ctx, organizationDID, ReceiverServiceName, "notification")
		if err != nil {
			return nil, err
		}

		compositionPath := fmt.Sprintf("/Composition/%s", transfer.FhirAdvanceNoticeComposition)
		transferTask := fhir.BuildNewTask(fhir.TaskProperties{
			RequesterID: *customer.Did,
			OwnerID:     organizationDID,
			Status:      REQUESTED_STATE,
			Input: []resources.TaskInputOutput{
			{
					Type:           &fhir.LoincAdvanceNoticeType,
					ValueReference: &datatypes.Reference{Reference: fhir.ToStringPtr(compositionPath)},
				},
			},
		})

		err = s.localFHIRClientFactory(fhir.WithTenant(customerID)).CreateOrUpdate(ctx, transferTask)
		if err != nil {
			return nil, err
		}

		if err := s.vcr.CreateAuthorizationCredential(ctx, SenderServiceName, *customer.Did, organizationDID, []credential.Resource{
			{
				Path:       fmt.Sprintf("/Task/%s", fhir.FromIDPtr(transferTask.ID)),
				Operations: []string{"read", "update"},
				UserContext: true,
			},
			{
				Path:       compositionPath,
				Operations: []string{"read", "document"},
				UserContext: true,
			},
		}); err != nil {
			return nil, err
		}

		negotiation, err = s.transferRepo.CreateNegotiation(ctx, customerID, transferID, organizationDID, transfer.TransferDate.Time, fhir.FromIDPtr(transferTask.ID))
		if err != nil {
			return nil, err
		}

		// Update transfer.Status = requested
		//transfer.Status = domain.TransferStatusRequested
		return transfer, nil
	})
	if err == nil {
		// Commit here, otherwise notifications to this server will deadlock on the uncommitted tx.
		tm, _ := sqlUtil.GetTransactionManager(ctx)
		if commitErr := tm.Commit(); commitErr != nil {
			return negotiation, commitErr
		}

		if err = s.sendNotification(ctx, customer, organizationDID); err != nil {
			// TODO: What to do here? Should we maybe rollback?
			logrus.Errorf("Unable to notify receiving care organization of updated FHIR task (did=%s): %s", organizationDID, err)
		}
	}

	return negotiation, err
}

func (s service) ConfirmNegotiation(ctx context.Context, customerID int, transferID, negotiationID string) (*domain.TransferNegotiation, error) {
	// find transfer
	transfer, err := s.transferRepo.FindByID(ctx, customerID, transferID)
	if err != nil {
		return nil, err
	}
	if transfer == nil {
		return nil, fmt.Errorf("transfer with ID: %s, not found", transferID)
	}

	var (
		negotiation   *domain.TransferNegotiation
		dossier       *domain.Dossier
		patient       *domain.Patient
		customer      *domain.Customer
		notifications []*notification
	)

	_, err = s.transferRepo.Update(ctx, customerID, transferID, func(transfer *domain.Transfer) (*domain.Transfer, error) {
		negotiations, err := s.transferRepo.ListNegotiations(ctx, customerID, transferID)
		if err != nil {
			return nil, err
		}

		advanceNoticePath := fmt.Sprintf("/Composition/%s", transfer.FhirAdvanceNoticeComposition)

		// cancel other negotiations + tasks + batch notifications
		for _, n := range negotiations {
			if negotiationID != string(n.Id) {
				// this also handles the FHIR and notification stuff
				_, notification, err := s.cancelNegotiation(ctx, customerID, string(n.Id), advanceNoticePath)
				if err != nil {
					return nil, err
				}
				notifications = append(notifications, notification)
			}
		}

		// alter state to in-progress in DB
		if negotiation, err = s.transferRepo.ConfirmNegotiation(ctx, customerID, negotiationID); err != nil {
			return nil, err
		}

		// retrieve patient
		if dossier, err = s.dossierRepo.FindByID(ctx, customerID, string(transfer.DossierID)); err != nil {
			return nil, err
		}
		if patient, err = s.patientRepo.FindByID(ctx, customerID, string(dossier.PatientID)); err != nil {
			return nil, err
		}
		// customer
		if customer, err = s.customerRepo.FindByID(customerID); err != nil {
			return nil, err
		}

		// create eTransfer composition connect to transfer
		composition := fhir.BuildNursingHandoff(patient)
		compositionID := composition["id"].(string)
		compositionPath := fmt.Sprintf("/Composition/%s", compositionID)
		transfer.FhirNursingHandoffComposition = &compositionID
		transfer.Status = domain.TransferStatusAssigned

		// update task
		task, err := s.getLocalTransferTask(ctx, customerID, negotiation.TaskID)
		if err != nil {
			return nil, err
		}
		task.Status = fhir.ToCodePtr(IN_PROGRESS_STATE)
		task.Input = []resources.TaskInputOutput{
			{
				Type:           &fhir.SnomedNursingHandoffType,
				ValueReference: &datatypes.Reference{Reference: fhir.ToStringPtr(compositionPath)},
			},
		}
		if err = s.localFHIRClientFactory(fhir.WithTenant(customerID)).CreateOrUpdate(ctx, task); err != nil {
			return nil, err
		}
		if err = s.localFHIRClientFactory(fhir.WithTenant(customerID)).CreateOrUpdate(ctx, composition); err != nil {
			return nil, err
		}

		// update authorization credential
		// todo referenced resources from within composition
		if err = s.vcr.RevokeAuthorizationCredential(ctx, SenderServiceName, negotiation.OrganizationDID, advanceNoticePath); err != nil {
			return nil, err
		}
		if err := s.vcr.CreateAuthorizationCredential(ctx, SenderServiceName, *customer.Did, negotiation.OrganizationDID, []credential.Resource{
			{
				Path:       fmt.Sprintf("/Task/%s", fhir.FromIDPtr(task.ID)),
				Operations: []string{"read", "update"},
				UserContext: true,
			},
			{
				Path:       compositionPath,
				Operations: []string{"read", "document"},
				UserContext: true,
			},
		}); err != nil {
			return nil, err
		}
		notifications = append(notifications, &notification{
			customer:        customer,
			organizationDID: negotiation.OrganizationDID,
		})

		return transfer, nil
	})
	if err == nil {
		// Commit here, otherwise notifications to this server will deadlock on the uncommitted tx.
		tm, _ := sqlUtil.GetTransactionManager(ctx)
		if commitErr := tm.Commit(); commitErr != nil {
			return negotiation, commitErr
		}

		for _, n := range notifications {
			s.sendNotification(ctx, n.customer, n.organizationDID)
		}
	}



	return negotiation, err
}

func (s service) CancelNegotiation(ctx context.Context, customerID int, negotiationID string) (*domain.TransferNegotiation, error) {
	// find transfer
	negotiation, err := s.transferRepo.FindNegotiationByID(ctx, customerID, negotiationID)
	if err != nil {
		return nil, err
	}
	transfer, err := s.transferRepo.FindByID(ctx, customerID, string(negotiation.TransferID))
	if err != nil {
		return nil, err
	}

	// update DB, Task and credential state
	negotiation, notification, err := s.cancelNegotiation(ctx, customerID, negotiationID,transfer.FhirAdvanceNoticeComposition)
	if err != nil {
		return nil, err
	}

	return negotiation, s.sendNotification(ctx, notification.customer, notification.organizationDID)
}

type notification struct {
	customer        *domain.Customer
	organizationDID string
}

// cancelNegotiation is like CancelNegotiation but it doesn't send any notification.
// the notification is returned so they can be send as batch.
func (s service) cancelNegotiation(ctx context.Context, customerID int, negotiationID, advanceNoticePath string) (*domain.TransferNegotiation, *notification, error) {
	// update DB state
	negotiation, err := s.transferRepo.CancelNegotiation(ctx, customerID, negotiationID)
	if err != nil {
		return nil, nil, err
	}

	// update local Task
	task, err := s.getLocalTransferTask(ctx, customerID, negotiation.TaskID)
	if err != nil {
		return nil, nil, err
	}
	task.Status = fhir.ToCodePtr(CANCELLED_STATE)
	if err = s.localFHIRClientFactory(fhir.WithTenant(customerID)).CreateOrUpdate(ctx, task); err != nil {
		return nil, nil, err
	}

	// revoke credential, find by AdvanceNotice
	if err = s.vcr.RevokeAuthorizationCredential(ctx, SenderServiceName, negotiation.OrganizationDID, advanceNoticePath); err != nil {
		return nil, nil, err
	}

	// create notification
	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil {
		return nil, nil, err
	}
	return negotiation, &notification{customer: customer, organizationDID: negotiation.OrganizationDID}, nil
}

func (s service) sendNotification(ctx context.Context, customer *domain.Customer, organizationDID string) error {
	notificationEndpoint, err := s.registry.GetCompoundServiceEndpoint(ctx, organizationDID, ReceiverServiceName, "notification")
	if err != nil {
		return err
	}

	tokenResponse, err := s.auth.RequestAccessToken(ctx, *customer.Did, organizationDID, ReceiverServiceName, nil)
	if err != nil {
		return err
	}

	return s.notifier.Notify(tokenResponse.AccessToken, notificationEndpoint)
}


func (s service) getLocalTransferTask(ctx context.Context, customerID int, fhirTaskID string) (resources.Task, error) {
	task := resources.Task{}
	err := s.localFHIRClientFactory(fhir.WithTenant(customerID)).ReadOne(ctx, "/Task/"+fhirTaskID, &task)
	if err != nil {
		return resources.Task{}, fmt.Errorf("error while looking up transfer task locally (task-id=%s): %w", fhirTaskID, err)
	}
	return task, nil
}
