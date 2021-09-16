package sender

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir/eoverdracht"
	"github.com/sirupsen/logrus"

	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/dossier"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/patients"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer"
	"github.com/nuts-foundation/nuts-demo-ehr/http/auth"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/registry"
	sqlUtil "github.com/nuts-foundation/nuts-demo-ehr/sql"
	"github.com/nuts-foundation/nuts-node/vcr/credential"
)

type TransferService interface {
	// Create creates a new transfer
	Create(ctx context.Context, customerID int, request domain.CreateTransferRequest) (*domain.Transfer, error)

	CreateNegotiation(ctx context.Context, customerID int, transferID, organizationDID string, transferDate time.Time) (*domain.TransferNegotiation, error)

	GetTransferByID(ctx context.Context, customerID int, transferID string) (domain.Transfer, error)

	// ConfirmNegotiation confirms the negotiation indicated by the negotiationID.
	// The updates the status to ACCEPTED_STATE.
	// It automatically cancels other negotiations of the domain.Transfer indicated by the transferID
	// by setting their status to CANCELLED_STATE.
	ConfirmNegotiation(ctx context.Context, customerID int, transferID, negotiationID string) (*domain.TransferNegotiation, error)

	// CancelNegotiation withdraws the negotiation/organization from the transfer. This is done by the sending party
	// It updates the status to CANCELLED_STATE, updates the FHIR Task and sends out a notification
	CancelNegotiation(ctx context.Context, customerID int, negotiationID string) (*domain.TransferNegotiation, error)

	// GetTransferRequest tries to retrieve a transfer request from requesting care organization's FHIR server.
	GetTransferRequest(ctx context.Context, customerID int, requestorDID string, fhirTaskID string) (*domain.TransferRequest, error)

	// UpdateTaskState updates the Task resource at the sender side. It updates the local DB, checks the statemachine, updates the FHIR record and sends a notification.
	UpdateTaskState(ctx context.Context, customer domain.Customer, taskID string, newState string) error
}

type service struct {
	transferRepo           TransferRepository
	auth                   auth.Service
	localFHIRClientFactory fhir.Factory // client for interacting with the local FHIR server
	customerRepo           customers.Repository
	dossierRepo            dossier.Repository
	patientRepo            patients.Repository
	registry               registry.OrganizationRegistry
	vcr                    registry.VerifiableCredentialRegistry
	notifier               transfer.Notifier
}

func NewTransferService(authService auth.Service, localFHIRClientFactory fhir.Factory, transferRepository TransferRepository, customerRepository customers.Repository, dossierRepo dossier.Repository, patientRepo patients.Repository, organizationRegistry registry.OrganizationRegistry, vcr registry.VerifiableCredentialRegistry) TransferService {
	return &service{
		auth:                   authService,
		localFHIRClientFactory: localFHIRClientFactory,
		transferRepo:           transferRepository,
		customerRepo:           customerRepository,
		dossierRepo:            dossierRepo,
		patientRepo:            patientRepo,
		registry:               organizationRegistry,
		vcr:                    vcr,
		notifier:               transfer.FireAndForgetNotifier{},
	}
}

func (s service) Create(ctx context.Context, customerID int, request domain.CreateTransferRequest) (*domain.Transfer, error) {
	advanceNotice := domain.BuildAdvanceNotice2(request)
	//logrus.WithFields(logrus.Fields{"advanceNotice": advanceNotice}).Info("advance notice build")

	for _, problem := range advanceNotice.Problems {
		err := s.localFHIRClientFactory(fhir.WithTenant(customerID)).CreateOrUpdate(ctx, problem)
		if err != nil {
			return nil, err
		}
	}
	for _, intervention := range advanceNotice.Interventions {
		err := s.localFHIRClientFactory(fhir.WithTenant(customerID)).CreateOrUpdate(ctx, intervention)
		if err != nil {
			return nil, err
		}
	}
	err := s.localFHIRClientFactory(fhir.WithTenant(customerID)).CreateOrUpdate(ctx, advanceNotice.Composition)
	if err != nil {
		return nil, err
	}

	transfer, err := s.transferRepo.Create(ctx, customerID, string(request.DossierID), request.TransferDate.Time, fhir.FromIDPtr(advanceNotice.Composition.ID))
	if err != nil {
		return nil, err
	}
	return transfer, nil
}

func (s service) GetTransferByID(ctx context.Context, customerID int, transferID string) (domain.Transfer, error) {
	transfer, err := s.transferRepo.FindByID(ctx, customerID, transferID)
	if err != nil {
		return domain.Transfer{}, err
	}

	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil || customer.Did == nil {
		return domain.Transfer{}, err
	}
	client := s.localFHIRClientFactory(fhir.WithTenant(customerID))

	advanceNotice, err := s.getAdvanceNotice(ctx, client, "Composition/"+transfer.FhirAdvanceNoticeComposition)
	if err != nil || customer.Did == nil {
		return domain.Transfer{}, err
	}
	domainTransfer, err := domain.FHIRAdvanceNoticeToDomainTransfer(advanceNotice)
	if err != nil || customer.Did == nil {
		return domain.Transfer{}, err
	}

	return domain.Transfer{
		TransferProperties:            domainTransfer,
		DossierID:                     transfer.DossierID,
		FhirAdvanceNoticeComposition:  transfer.FhirAdvanceNoticeComposition,
		FhirNursingHandoffComposition: transfer.FhirNursingHandoffComposition,
		Id:                            transfer.Id,
		Status:                        transfer.Status,
	}, nil
}

func (s service) taskContainsCode(task resources.Task, code datatypes.Code) bool {
	for _, input := range task.Input {
		if fhir.FromCodePtr(input.Type.Coding[0].Code) == string(code) {
			return true
		}
	}

	return false
}

func (s service) GetTransferRequest(ctx context.Context, customerID int, requesterDID string, fhirTaskID string) (*domain.TransferRequest, error) {
	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil || customer.Did == nil {
		return nil, err
	}

	client, err := s.getRemoteFHIRClient(ctx, requesterDID, *customer.Did)
	if err != nil {
		return nil, err
	}

	task, err := s.getRemoteTransferTask(ctx, client, fhirTaskID)
	if err != nil {
		return nil, err
	}

	if !s.taskContainsCode(task, fhir.LoincAdvanceNoticeCode) {
		return nil, fmt.Errorf("invalid task, expected an advanceNotice composition")
	}

	advanceNotice, err := s.getAdvanceNotice(ctx, client(), fhir.FromStringPtr(task.Input[0].ValueReference.Reference))

	organization, err := s.registry.Get(ctx, requesterDID)
	if err != nil {
		return nil, err
	}

	domainTransfer, err := domain.FHIRAdvanceNoticeToDomainTransfer(advanceNotice)

	// TODO: Do we need nil checks?
	return &domain.TransferRequest{
		Sender:        *organization,
		AdvanceNotice: domainTransfer,
		Status:        fhir.FromCodePtr(task.Status),
	}, nil
}

func (s service) CreateNegotiation(ctx context.Context, customerID int, transferID, organizationDID string, transferDate time.Time) (*domain.TransferNegotiation, error) {
	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil || customer.Did == nil {
		return nil, err
	}

	var (
		negotiation *domain.TransferNegotiation
	)

	_, err = s.transferRepo.Update(ctx, customerID, transferID, func(transferRecord *domain.Transfer) (*domain.Transfer, error) {
		// Validate transfer
		if transferRecord.Status == domain.TransferStatusCancelled ||
			transferRecord.Status == domain.TransferStatusCompleted ||
			transferRecord.Status == domain.TransferStatusAssigned {
			return nil, errors.New("can't start new transfer negotiation when status is 'cancelled', 'assigned' or 'completed'")
		}

		// Create negotiation and share it to the other party
		// TODO: Share transaction to this repository call as well
		var err error
		// Pre-emptively resolve the receiver organization's notification endpoint to reduce clutter, avoiding to make FHIR tasks when the receiving party eOverdracht registration is faulty.
		_, err = s.registry.GetCompoundServiceEndpoint(ctx, organizationDID, transfer.ReceiverServiceName, "notification")
		if err != nil {
			return nil, err
		}

		fhirClient := s.localFHIRClientFactory(fhir.WithTenant(customerID))

		compositionPath := fmt.Sprintf("/Composition/%s", transferRecord.FhirAdvanceNoticeComposition)
		composition := eoverdracht.Composition{}
		fhirClient.ReadOne(ctx, compositionPath, &composition)
		if err != nil {
			return nil, fmt.Errorf("could not create FHIR Task: %w", err)
		}

		// A list to store all the paths to FHIR resources associated with this advance notice
		// These paths must be included in the authorization credential
		resourcePaths := resourcePathsFromSection(composition.Section, []string{})

		transferTask := domain.BuildNewTask(fhir.TaskProperties{
			RequesterID: *customer.Did,
			OwnerID:     organizationDID,
			Status:      transfer.RequestedState,
			Input: []resources.TaskInputOutput{
				{
					Type:           &fhir.LoincAdvanceNoticeType,
					ValueReference: &datatypes.Reference{Reference: fhir.ToStringPtr(compositionPath)},
				},
			},
		})

		err = fhirClient.CreateOrUpdate(ctx, transferTask)
		if err != nil {
			return nil, fmt.Errorf("could not create FHIR Task: %w", err)
		}
		authorizedResources := []credential.Resource{
			{
				Path:        fmt.Sprintf("/Task/%s", fhir.FromIDPtr(transferTask.ID)),
				Operations:  []string{"read", "update"},
				UserContext: true,
			},
			{
				Path:        compositionPath,
				Operations:  []string{"read", "document"},
				UserContext: true,
			},
		}

		for _, path := range resourcePaths {
			authorizedResources = append(authorizedResources, credential.Resource{
				Path:        path,
				Operations:  []string{"read", "document"},
				UserContext: true,
			})
		}

		if err := s.vcr.CreateAuthorizationCredential(ctx, transfer.SenderServiceName, *customer.Did, organizationDID, authorizedResources); err != nil {
			return nil, err
		}

		negotiation, err = s.transferRepo.CreateNegotiation(ctx, customerID, transferID, organizationDID, transferRecord.TransferDate.Time, fhir.FromIDPtr(transferTask.ID))
		if err != nil {
			return nil, err
		}

		// Update transfer.Status = requested
		//transfer.Status = domain.TransferStatusRequested
		return transferRecord, nil
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

func resourcePathsFromSection(sections []eoverdracht.CompositionSection, paths []string) []string {
	for _, s := range sections {
		paths = append(paths, resourcePathsFromSection(s.Section, paths)...)
		for _, e := range s.Entry {
			path := fhir.FromStringPtr(e.Reference)
			if path != "" {
				// paths in authorization credential need a / prefix
				paths = append(paths, "/"+path)
			}
		}
	}
	return paths
}

func (s service) ConfirmNegotiation(ctx context.Context, customerID int, transferID, negotiationID string) (*domain.TransferNegotiation, error) {
	// find transfer
	transferRecord, err := s.transferRepo.FindByID(ctx, customerID, transferID)
	if err != nil {
		return nil, err
	}
	if transferRecord == nil {
		return nil, fmt.Errorf("transfer with ID: %s, not found", transferID)
	}

	var (
		negotiation   *domain.TransferNegotiation
		dossier       *domain.Dossier
		patient       *domain.Patient
		customer      *domain.Customer
		notifications []*notification
	)

	_, err = s.transferRepo.Update(ctx, customerID, transferID, func(transferRecord *domain.Transfer) (*domain.Transfer, error) {
		negotiations, err := s.transferRepo.ListNegotiations(ctx, customerID, transferID)
		if err != nil {
			return nil, err
		}

		advanceNoticePath := fmt.Sprintf("/Composition/%s", transferRecord.FhirAdvanceNoticeComposition)

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
		if dossier, err = s.dossierRepo.FindByID(ctx, customerID, string(transferRecord.DossierID)); err != nil {
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
		composition := domain.BuildNursingHandoff(patient)
		compositionID := composition["id"].(string)
		compositionPath := fmt.Sprintf("/Composition/%s", compositionID)
		transferRecord.FhirNursingHandoffComposition = &compositionID
		transferRecord.Status = domain.TransferStatusAssigned

		// update task
		task, err := s.getLocalTransferTask(ctx, customerID, negotiation.TaskID)
		if err != nil {
			return nil, err
		}
		task.Status = fhir.ToCodePtr(transfer.InProgressState)
		task.Input = append(task.Input, resources.TaskInputOutput{
			Type:           &fhir.SnomedNursingHandoffType,
			ValueReference: &datatypes.Reference{Reference: fhir.ToStringPtr(compositionPath)},
		})

		if err = s.localFHIRClientFactory(fhir.WithTenant(customerID)).CreateOrUpdate(ctx, composition); err != nil {
			return nil, err
		}
		if err = s.localFHIRClientFactory(fhir.WithTenant(customerID)).CreateOrUpdate(ctx, task); err != nil {
			return nil, err
		}

		// update authorization credential
		// todo referenced resources from within composition
		if err = s.vcr.RevokeAuthorizationCredential(ctx, transfer.SenderServiceName, negotiation.OrganizationDID, advanceNoticePath); err != nil {
			return nil, err
		}
		if err := s.vcr.CreateAuthorizationCredential(ctx, transfer.SenderServiceName, *customer.Did, negotiation.OrganizationDID, []credential.Resource{
			{
				Path:        fmt.Sprintf("/Task/%s", fhir.FromIDPtr(task.ID)),
				Operations:  []string{"read", "update"},
				UserContext: true,
			},
			{
				Path:        compositionPath,
				Operations:  []string{"read", "document"},
				UserContext: true,
			},
		}); err != nil {
			return nil, err
		}
		notifications = append(notifications, &notification{
			customer:        customer,
			organizationDID: negotiation.OrganizationDID,
		})

		return transferRecord, nil
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
	negotiation, notification, err := s.cancelNegotiation(ctx, customerID, negotiationID, transfer.FhirAdvanceNoticeComposition)
	if err != nil {
		return nil, err
	}

	return negotiation, s.sendNotification(ctx, notification.customer, notification.organizationDID)
}

func (s service) UpdateTaskState(ctx context.Context, customer domain.Customer, taskID string, newState string) error {
	// find negotiation
	negotiation, err := s.transferRepo.FindNegotiationByTaskID(ctx, customer.Id, taskID)
	if err != nil {
		return err
	}

	// check state transition
	// todo this only allows for direct assigned transfers
	if !(negotiation.Status == transfer.InProgressState && newState == transfer.CompletedState) {
		// invalid state change
		return fmt.Errorf("invalid task state change: from %s to %s", negotiation.Status, newState)
	}

	return s.completeTask(ctx, customer, negotiation)
}

// completeTask will also complete the transfer, revoke credential and send a notification
func (s service) completeTask(ctx context.Context, customer domain.Customer, negotiation *domain.TransferNegotiation) error {
	transferID := string(negotiation.TransferID)

	// find transfer
	transferRecord, err := s.transferRepo.FindByID(ctx, customer.Id, transferID)
	if err != nil {
		return err
	}
	if transferRecord == nil {
		return fmt.Errorf("transfer with ID: %s, not found", transferID)
	}

	var not notification

	_, err = s.transferRepo.Update(ctx, customer.Id, transferID, func(transferRecord *domain.Transfer) (*domain.Transfer, error) {
		// alter state to completed in DB for Task
		if negotiation, err = s.transferRepo.UpdateNegotiationState(ctx, customer.Id, string(negotiation.Id), transfer.CompletedState); err != nil {
			return nil, err
		}
		// alter state for transfer to completed as well
		transferRecord.Status = domain.TransferStatusCompleted

		// update FHIR task
		task, err := s.getLocalTransferTask(ctx, customer.Id, negotiation.TaskID)
		if err != nil {
			return nil, err
		}
		task.Status = fhir.ToCodePtr(transfer.CompletedState)
		if err = s.localFHIRClientFactory(fhir.WithTenant(customer.Id)).CreateOrUpdate(ctx, task); err != nil {
			return nil, err
		}

		// reconstruct composition path
		compositionPath := fmt.Sprintf("/Composition/%s", *transferRecord.FhirNursingHandoffComposition)
		// revoke authorization credential
		if err = s.vcr.RevokeAuthorizationCredential(ctx, transfer.SenderServiceName, negotiation.OrganizationDID, compositionPath); err != nil {
			return nil, err
		}

		// create notification
		not = notification{
			customer:        &customer,
			organizationDID: negotiation.OrganizationDID,
		}

		return transferRecord, nil
	})
	if err == nil {
		// Commit here, otherwise notifications to this server will deadlock on the uncommitted tx.
		tm, _ := sqlUtil.GetTransactionManager(ctx)
		if commitErr := tm.Commit(); commitErr != nil {
			return commitErr
		}

		s.sendNotification(ctx, not.customer, not.organizationDID)
	}

	return err
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
	task.Status = fhir.ToCodePtr(transfer.CancelledState)
	if err = s.localFHIRClientFactory(fhir.WithTenant(customerID)).CreateOrUpdate(ctx, task); err != nil {
		return nil, nil, err
	}

	// revoke credential, find by AdvanceNotice
	if err = s.vcr.RevokeAuthorizationCredential(ctx, transfer.SenderServiceName, negotiation.OrganizationDID, advanceNoticePath); err != nil {
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
	notificationEndpoint, err := s.registry.GetCompoundServiceEndpoint(ctx, organizationDID, transfer.ReceiverServiceName, "notification")
	if err != nil {
		return err
	}

	tokenResponse, err := s.auth.RequestAccessToken(ctx, *customer.Did, organizationDID, transfer.ReceiverServiceName, nil)
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

func (s service) getRemoteFHIRClient(ctx context.Context, custodianDID string, localActorDID string) (fhir.Factory, error) {
	fhirServer, err := s.registry.GetCompoundServiceEndpoint(ctx, custodianDID, transfer.SenderServiceName, "fhir")
	if err != nil {
		return nil, fmt.Errorf("error while looking up custodian's FHIR server (did=%s): %w", custodianDID, err)
	}
	accessToken, err := s.auth.RequestAccessToken(ctx, localActorDID, custodianDID, transfer.SenderServiceName, nil)
	if err != nil {
		return nil, err
	}

	return fhir.NewFactory(fhir.WithURL(fhirServer), fhir.WithAuthToken(accessToken.AccessToken)), nil
}

func (s service) getRemoteTransferTask(ctx context.Context, client fhir.Factory, fhirTaskID string) (resources.Task, error) {
	// TODO: Read AdvanceNotification here instead of the transfer task
	task := resources.Task{}
	err := client().ReadOne(ctx, "/Task/"+fhirTaskID, &task)
	if err != nil {
		return resources.Task{}, fmt.Errorf("error while looking up transfer task remotely(task-id=%s): %w", fhirTaskID, err)
	}
	return task, nil
}

// getAdvanceNotice fetches a complete advance notice from the a FHIR server
func (s service) getAdvanceNotice(ctx context.Context, client fhir.Client, fhirCompositionPath string) (eoverdracht.AdvanceNotice, error) {
	advanceNotice := eoverdracht.AdvanceNotice{}

	err := client.ReadOne(ctx, "/"+fhirCompositionPath, &advanceNotice.Composition)
	if err != nil {
		return eoverdracht.AdvanceNotice{}, fmt.Errorf("error while fetching the advance notice composition(composition-id=%s): %w", fhirCompositionPath, err)
	}

	careplan, err := eoverdracht.FilterCompositionSectionByType(advanceNotice.Composition.Section, eoverdracht.CarePlanCode)
	if err != nil {
		return eoverdracht.AdvanceNotice{}, err
	}

	nursingDiagnosis, err := eoverdracht.FilterCompositionSectionByType(careplan.Section, eoverdracht.NursingDiagnosisCode)
	if err != nil {
		return eoverdracht.AdvanceNotice{}, err
	}

	// the nursing diagnosis contains both conditions and procedures
	for _, entry := range nursingDiagnosis.Entry {
		if strings.HasPrefix(fhir.FromStringPtr(entry.Reference), "Condition") {
			conditionID := fhir.FromStringPtr(entry.Reference)
			condition := resources.Condition{}
			err := client.ReadOne(ctx, "/"+conditionID, &condition)
			if err != nil {
				return eoverdracht.AdvanceNotice{}, fmt.Errorf("error while fetching a advance notice condition (condition-id=%s): %w", conditionID, err)
			}
			advanceNotice.Problems = append(advanceNotice.Problems, condition)
		}
		if strings.HasPrefix(fhir.FromStringPtr(entry.Reference), "Procedure") {
			procedureID := fhir.FromStringPtr(entry.Reference)
			procedure := eoverdracht.Procedure{}
			err := client.ReadOne(ctx, "/"+procedureID, &procedure)
			if err != nil {
				return eoverdracht.AdvanceNotice{}, fmt.Errorf("error while fetching a advance notice procedure (procedure-id=%s): %w", procedureID, err)
			}
			advanceNotice.Interventions = append(advanceNotice.Interventions, procedure)
		}
	}

	return advanceNotice, nil
}
