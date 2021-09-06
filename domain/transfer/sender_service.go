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
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir/eoverdracht"
	sqlUtil "github.com/nuts-foundation/nuts-demo-ehr/sql"
	"github.com/nuts-foundation/nuts-node/vcr/credential"
	"github.com/sirupsen/logrus"
)

// SenderServiceName contains the name of the eOverdracht sender compound-service
const SenderServiceName = "eOverdracht-sender"

func (s service) Create(ctx context.Context, customerID int, dossierID string, description string, transferDate time.Time) (*domain.Transfer, error) {
	elements := map[string]interface{}{
		"title": "Aanmeldbericht",
		"type":  fhir.LoincAdvanceNoticeType,
		// TODO: patient seems mandatory in the spec, but can only be sent when placer already
		// has patient in care to protect the identity of the patient during the negotiation phase.
		//"subject":  fhir.Reference{Reference: "Patient/Anonymous"},
		"author": eoverdracht.Practitioner{
			// TODO: Derive from authenticated user?
			Identifier: datatypes.Identifier{
				System: &fhir.UZICodingSystem,
				Value:  fhir.ToStringPtr("12345"),
			},
			Name: &datatypes.HumanName{
				Family: fhir.ToStringPtr("Demo EHR"),
				Given:  []datatypes.String{"Nuts"},
			},
		},
		// TODO: sections
	}
	composition := fhir.BuildNewComposition(elements)
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
		notificationEndpoint string
	)

	_, err = s.transferRepo.Update(ctx, customerID, transferID, func(transfer domain.Transfer) (*domain.Transfer, error) {
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
		notificationEndpoint, err = s.registry.GetCompoundServiceEndpoint(ctx, organizationDID, ReceiverServiceName, "notification")
		if err != nil {
			return nil, err
		}

		transferTask := fhir.BuildNewTask(fhir.TaskProperties{
			RequesterID: *customer.Did,
			OwnerID:     organizationDID,
			Status:      REQUESTED_STATE,
			Input: []resources.TaskInputOutput{
				{
					Type:           &fhir.LoincAdvanceNoticeType,
					ValueReference: &datatypes.Reference{Reference: fhir.ToStringPtr(transfer.FhirAdvanceNoticeComposition)},
				},
			},
		})

		err = s.localFHIRClientFactory(fhir.WithTenant(customerID)).CreateOrUpdate(ctx, transferTask)
		if err != nil {
			return nil, err
		}

		if err := s.vcr.CreateAuthorizationCredential(ctx, "eOverdracht-receiver", *customer.Did, organizationDID, []credential.Resource{
			{
				Path:       fmt.Sprintf("/Task/%s", fhir.FromIDPtr(transferTask.ID)),
				Operations: []string{"update"},
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
		return &transfer, nil
	})
	if err == nil {
		// Commit here, otherwise notifications to this server will deadlock on the uncommitted tx.
		tm, _ := sqlUtil.GetTransactionManager(ctx)
		if commitErr := tm.Commit(); commitErr != nil {
			return negotiation, commitErr
		}

		tokenResponse, err := s.auth.RequestAccessToken(ctx, *customer.Did, organizationDID, ReceiverServiceName, nil)
		if err != nil {
			return nil, err
		}

		if err = s.notifier.Notify(tokenResponse.AccessToken, notificationEndpoint); err != nil {
			// TODO: What to do here? Should we maybe rollback?
			logrus.Errorf("Unable to notify receiving care organization of updated FHIR task (did=%s): %s", organizationDID, err)
		}
	}

	return negotiation, err
}

func (s service) ConfirmNegotiation(ctx context.Context, customerID int, negotiationID string) (*domain.TransferNegotiation, error) {
	panic("implement me")
}
