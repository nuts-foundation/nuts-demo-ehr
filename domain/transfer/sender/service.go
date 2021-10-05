package sender

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir/eoverdracht"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	"github.com/sirupsen/logrus"

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
	Create(ctx context.Context, customerID int, request types.CreateTransferRequest) (*types.Transfer, error)

	CreateNegotiation(ctx context.Context, customerID int, transferID, organizationDID string) (*types.TransferNegotiation, error)

	GetTransferByID(ctx context.Context, customerID int, transferID string) (types.Transfer, error)

	// ConfirmNegotiation confirms the negotiation indicated by the negotiationID.
	// The updates the status to in progress
	// It automatically cancels other negotiations of the domain.Transfer indicated by the transferID
	// by setting their status to CANCELLED_STATE.
	ConfirmNegotiation(ctx context.Context, customerID int, transferID, negotiationID string) (*types.TransferNegotiation, error)

	// CancelNegotiation withdraws the negotiation/organization from the transfer. This is done by the sending party
	// It updates the status to CANCELLED_STATE, updates the FHIR Task and sends out a notification
	CancelNegotiation(ctx context.Context, customerID int, transferID, negotiationID string) (*types.TransferNegotiation, error)

	// UpdateTaskState updates the Task resource at the sender side. It updates the local DB, checks the statemachine, updates the FHIR record and sends a notification.
	UpdateTaskState(ctx context.Context, customer types.Customer, taskID string, newState string) error
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

func (s service) Create(ctx context.Context, customerID int, request types.CreateTransferRequest) (*types.Transfer, error) {
	// Fetch the patient
	patient, err := s.findPatientByDossierID(ctx, customerID, string(request.DossierID))
	if err != nil {
		return nil, fmt.Errorf("could not create transfer: %w", err)
	}

	// Build the advance notice resources
	advanceNotice := eoverdracht.NewFHIRBuilder().BuildAdvanceNotice(request, patient)

	// Save the resources to the fhir storage
	err = s.saveAdvanceNoticeFHIRResources(ctx, customerID, advanceNotice)
	if err != nil {
		return nil, fmt.Errorf("could not create transfer: unable to store advance notification fhir resources: %w", err)
	}

	// Create the database transfer
	return s.transferRepo.Create(ctx, customerID, string(request.DossierID), request.TransferDate.Time, fhir.FromIDPtr(advanceNotice.Composition.ID))
}

func (s service) GetTransferByID(ctx context.Context, customerID int, transferID string) (types.Transfer, error) {
	dbTransfer, err := s.transferRepo.FindByID(ctx, customerID, transferID)
	if err != nil {
		return types.Transfer{}, err
	}

	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil || customer.Did == nil {
		return types.Transfer{}, err
	}
	client := s.localFHIRClientFactory(fhir.WithTenant(customerID))

	advanceNotice, err := s.getAdvanceNotice(ctx, client, "Composition/"+dbTransfer.FhirAdvanceNoticeComposition)
	if err != nil || customer.Did == nil {
		return types.Transfer{}, err
	}
	domainTransfer, err := eoverdracht.FHIRAdvanceNoticeToDomainTransfer(advanceNotice)
	if err != nil || customer.Did == nil {
		return types.Transfer{}, err
	}

	return types.Transfer{
		TransferProperties:            domainTransfer,
		DossierID:                     dbTransfer.DossierID,
		FhirAdvanceNoticeComposition:  dbTransfer.FhirAdvanceNoticeComposition,
		FhirNursingHandoffComposition: dbTransfer.FhirNursingHandoffComposition,
		Id:                            dbTransfer.Id,
		Status:                        dbTransfer.Status,
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

// CreateNegotiation creates a new negotiation(FHIR Task) for a specific transfer and sends the other party a notification.
func (s service) CreateNegotiation(ctx context.Context, customerID int, transferID, organizationDID string) (*types.TransferNegotiation, error) {
	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil {
		return nil, err
	}
	if customer.Did == nil {
		return nil, fmt.Errorf("unable to create negotiation: customer does not have did")
	}

	var negotiation *types.TransferNegotiation

	// Pre-emptively resolve the receiver organization's notification endpoint to check registry configuration.
	// This reduces cleanup code of FHIR task.
	_, err = s.registry.GetCompoundServiceEndpoint(ctx, organizationDID, transfer.ReceiverServiceName, "notification")
	if err != nil {
		return nil, fmt.Errorf("unable to create transfer negotiation: no notification endpoint found: %w", err)
	}

	// Update the transfer
	_, err = s.transferRepo.Update(ctx, customerID, transferID, func(dbTransfer *types.Transfer) (*types.Transfer, error) {
		// Validate if transfer is in correct state to allow new negotiations
		if dbTransfer.Status == types.TransferStatusCancelled ||
			dbTransfer.Status == types.TransferStatusCompleted ||
			dbTransfer.Status == types.TransferStatusAssigned {
			return nil, errors.New("can't start new transfer negotiation when status is 'cancelled', 'assigned' or 'completed'")
		}

		fhirClient := s.localFHIRClientFactory(fhir.WithTenant(customerID))

		compositionPath := fmt.Sprintf("/Composition/%s", dbTransfer.FhirAdvanceNoticeComposition)
		composition := eoverdracht.Composition{}
		err = fhirClient.ReadOne(ctx, compositionPath, &composition)
		if err != nil {
			return nil, fmt.Errorf("could not create transfer negotiation: could not read fhir compositition: %w", err)
		}

		transferTask := eoverdracht.NewFHIRBuilder().BuildNewTask(fhir.TaskProperties{
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

		// Build the list of resources for the authorization credential:
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

		// A list to store all the paths to FHIR resources associated with this advance notice
		// These paths must be included in the authorization credential
		resourcePaths := resourcePathsFromSection(composition.Section, []string{})
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

		negotiation, err = s.transferRepo.CreateNegotiation(ctx, customerID, transferID, organizationDID, dbTransfer.TransferDate.Time, fhir.FromIDPtr(transferTask.ID))
		if err != nil {
			return nil, err
		}

		// Update transfer.Status = requested
		//transfer.Status = domain.TransferStatusRequested
		return dbTransfer, nil
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

// ConfirmNegotiation is executed by the sending organization. It confirms a transfer negotiation and cancels the others.
func (s service) ConfirmNegotiation(ctx context.Context, customerID int, transferID, negotiationID string) (*types.TransferNegotiation, error) {
	var (
		negotiation   *types.TransferNegotiation
		patient       *types.Patient
		customer      *types.Customer
		notifications []*notification
	)

	// Update database transfer
	_, err := s.transferRepo.Update(ctx, customerID, transferID, func(dbTransfer *types.Transfer) (*types.Transfer, error) {
		allNegotiations, err := s.transferRepo.ListNegotiations(ctx, customerID, transferID)
		if err != nil {
			return nil, err
		}

		advanceNoticePath := fmt.Sprintf("/Composition/%s", dbTransfer.FhirAdvanceNoticeComposition)

		// cancel other negotiations + tasks + batch notifications
		for _, n := range allNegotiations {
			if negotiationID != string(n.Id) {
				// this also handles the FHIR and notification stuff
				_, notification, err := s.cancelNegotiation(ctx, customerID, string(n.Id), advanceNoticePath)
				if err != nil {
					return nil, err
				}
				notifications = append(notifications, notification)
			}
		}

		// alter state of this database negotiation to in-progress
		if negotiation, err = s.transferRepo.ConfirmNegotiation(ctx, customerID, negotiationID); err != nil {
			return nil, err
		}

		patient, err = s.findPatientByDossierID(ctx, customerID, string(dbTransfer.DossierID))
		if err != nil {
			return nil, fmt.Errorf("unable to confirm egotiation: could not fetch patient by dossierID: %w", err)
		}

		// retrieve customer
		if customer, err = s.customerRepo.FindByID(customerID); err != nil {
			return nil, err
		}

		// The advance notice contains a lot of the same resources which should also be used in the Nursing Handoff
		// Fetch the advanceNotice FHIR resources
		fhirClient := s.localFHIRClientFactory(fhir.WithTenant(customerID))
		advanceNotice, err := s.getAdvanceNotice(ctx, fhirClient, advanceNoticePath)

		// Create eTransfer composition based on the advanceNotice and patient
		nursingHandoffComposition, err := eoverdracht.NewFHIRBuilder().BuildNursingHandoffComposition(patient, advanceNotice)
		if err != nil {
			return nil, err
		}
		compositionID := nursingHandoffComposition.ID
		compositionPath := fmt.Sprintf("/Composition/%s", fhir.FromIDPtr(compositionID))
		dbTransfer.FhirNursingHandoffComposition = (*string)(compositionID)
		dbTransfer.Status = types.TransferStatusAssigned

		// update Task with new status and composition reference
		task, err := s.getLocalTransferTask(ctx, customerID, negotiation.TaskID)
		if err != nil {
			return nil, err
		}
		task.Status = fhir.ToCodePtr(transfer.InProgressState)
		task.Input = append(task.Input, resources.TaskInputOutput{
			Type:           &fhir.SnomedNursingHandoffType,
			ValueReference: &datatypes.Reference{Reference: fhir.ToStringPtr(compositionPath)},
		})

		// Create nursing handoff composition in the FHIR store
		if err = s.localFHIRClientFactory(fhir.WithTenant(customerID)).CreateOrUpdate(ctx, nursingHandoffComposition); err != nil {
			return nil, err
		}
		// Update the task in the FHIR store
		if err = s.localFHIRClientFactory(fhir.WithTenant(customerID)).CreateOrUpdate(ctx, task); err != nil {
			return nil, err
		}

		// Revoke the old AuthorizationCredential for the Task and AdvanceNotice
		if err = s.vcr.RevokeAuthorizationCredential(ctx, transfer.SenderServiceName, negotiation.OrganizationDID, advanceNoticePath); err != nil {
			return nil, fmt.Errorf("unable to confirm negotiation: could not revoke advance notice authorization credential: %w", err)
		}

		authorizedResources := []credential.Resource{
			{
				Path:        fmt.Sprintf("/Task/%s", fhir.FromIDPtr(task.ID)),
				Operations:  []string{"read", "update"},
				UserContext: true,
			},
		}
		// Add paths of resources of both the advance notice and the nursing handoff
		resourcePaths := resourcePathsFromSection(nursingHandoffComposition.Section, []string{advanceNoticePath, compositionPath})
		resourcePaths = resourcePathsFromSection(advanceNotice.Composition.Section, resourcePaths)
		// Add path to the complete, non-anonymous, patient
		resourcePaths = append(resourcePaths, "/"+fhir.FromStringPtr(nursingHandoffComposition.Subject.Reference))

		// The resourcePaths may contain duplicates, hold a list of processedPaths
		processedPaths := map[string]struct{}{}
		for _, path := range resourcePaths {
			if _, exists := processedPaths[path]; exists {
				continue
			}
			authorizedResources = append(authorizedResources, credential.Resource{
				Path:        path,
				Operations:  []string{"read", "document"},
				UserContext: true,
			})
			processedPaths[path] = struct{}{}
		}

		// Create a new AuthorizationCredential for the Task, AdvanceNotice and NursingHandoff
		if err = s.vcr.CreateAuthorizationCredential(ctx, transfer.SenderServiceName, *customer.Did, negotiation.OrganizationDID, authorizedResources); err != nil {
			return nil, fmt.Errorf("unable to confirm negotiation: could not create authorization credential: %w", err)
		}
		notifications = append(notifications, &notification{
			customer:        customer,
			organizationDID: negotiation.OrganizationDID,
		})

		return dbTransfer, nil
	})
	if err == nil {
		// Commit db transaction here, otherwise notifications to this server will deadlock.
		tm, _ := sqlUtil.GetTransactionManager(ctx)
		if commitErr := tm.Commit(); commitErr != nil {
			return negotiation, commitErr
		}

		for _, n := range notifications {
			_ = s.sendNotification(ctx, n.customer, n.organizationDID)
		}
	}

	return negotiation, err
}

func (s service) CancelNegotiation(ctx context.Context, customerID int, transferID, negotiationID string) (*types.TransferNegotiation, error) {
	// find transfer
	negotiation, err := s.transferRepo.FindNegotiationByID(ctx, customerID, negotiationID)
	if err != nil {
		return nil, err
	}
	dbTransfer, err := s.transferRepo.FindByID(ctx, customerID, string(negotiation.TransferID))
	if err != nil {
		return nil, err
	}

	// update DB, Task and credential state
	negotiation, notification, err := s.cancelNegotiation(ctx, customerID, negotiationID, dbTransfer.FhirAdvanceNoticeComposition)
	if err != nil {
		return nil, err
	}

	return negotiation, s.sendNotification(ctx, notification.customer, notification.organizationDID)
}

func (s service) UpdateTaskState(ctx context.Context, customer types.Customer, taskID string, newState string) error {
	// find negotiation
	negotiation, err := s.transferRepo.FindNegotiationByTaskID(ctx, customer.Id, taskID)
	if err != nil {
		return err
	}

	// check state transition
	if !(negotiation.Status == transfer.RequestedState && newState == transfer.AcceptedState ||
		negotiation.Status == transfer.InProgressState && newState == transfer.CompletedState) {
		// invalid state change
		return fmt.Errorf("invalid task state change: from %s to %s", negotiation.Status, newState)
	}

	if newState == transfer.AcceptedState {
		return s.acceptTask(ctx, customer, negotiation)
	}
	if newState == transfer.CompletedState {
		return s.completeTask(ctx, customer, negotiation)
	}
	return nil
}

// acceptTask sets the negotiation and corresponding task on accepted.
func (s service) acceptTask(ctx context.Context, customer types.Customer, negotiation *types.TransferNegotiation) error {

	// alter state to completed in DB for Task
	if _, err := s.transferRepo.UpdateNegotiationState(ctx, customer.Id, string(negotiation.Id), transfer.AcceptedState); err != nil {
		return err
	}
	// update FHIR task
	task, err := s.getLocalTransferTask(ctx, customer.Id, negotiation.TaskID)
	if err != nil {
		return err
	}
	task.Status = fhir.ToCodePtr(transfer.AcceptedState)
	if err = s.localFHIRClientFactory(fhir.WithTenant(customer.Id)).CreateOrUpdate(ctx, task); err != nil {
		return err
	}

	// create notification
	not := notification{
		customer:        &customer,
		organizationDID: negotiation.OrganizationDID,
	}

	// Commit here, otherwise notifications to this server will deadlock on the uncommitted tx.
	tm, _ := sqlUtil.GetTransactionManager(ctx)
	if commitErr := tm.Commit(); commitErr != nil {
		return commitErr
	}

	_ = s.sendNotification(ctx, not.customer, not.organizationDID)

	return nil
}

// completeTask will also complete the transfer, revoke credential and send a notification
func (s service) completeTask(ctx context.Context, customer types.Customer, negotiation *types.TransferNegotiation) error {
	transferID := string(negotiation.TransferID)

	var not notification

	_, err := s.transferRepo.Update(ctx, customer.Id, transferID, func(transferRecord *types.Transfer) (*types.Transfer, error) {
		var err error
		// alter state to completed in DB for Task
		if negotiation, err = s.transferRepo.UpdateNegotiationState(ctx, customer.Id, string(negotiation.Id), transfer.CompletedState); err != nil {
			return nil, err
		}
		// alter state for transfer to completed as well
		transferRecord.Status = types.TransferStatusCompleted

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

		_ = s.sendNotification(ctx, not.customer, not.organizationDID)
	}

	return err
}

type notification struct {
	customer        *types.Customer
	organizationDID string
}

// cancelNegotiation is like CancelNegotiation but it doesn't send any notification.
// the notification is returned so they can be send as batch.
func (s service) cancelNegotiation(ctx context.Context, customerID int, negotiationID, advanceNoticePath string) (*types.TransferNegotiation, *notification, error) {
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

func (s service) sendNotification(ctx context.Context, customer *types.Customer, organizationDID string) error {
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
	fhirClient := s.localFHIRClientFactory(fhir.WithTenant(customerID))
	fhirRepo := fhir.NewFHIRRepository(fhirClient)
	task, err := fhirRepo.GetTask(ctx, fhirTaskID)
	if err != nil {
		return resources.Task{}, fmt.Errorf("error while looking up transfer task locally (task-id=%s): %w", fhirTaskID, err)
	}
	return task, nil
}

// getAdvanceNotice fetches a complete advance notice from a FHIR server
func (s service) getAdvanceNotice(ctx context.Context, fhirClient fhir.Client, fhirCompositionPath string) (eoverdracht.AdvanceNotice, error) {
	advanceNotice := eoverdracht.AdvanceNotice{}

	err := fhirClient.ReadOne(ctx, "/"+fhirCompositionPath, &advanceNotice.Composition)
	if err != nil {
		return eoverdracht.AdvanceNotice{}, fmt.Errorf("error while fetching the advance notice composition(composition-id=%s): %w", fhirCompositionPath, err)
	}

	if advanceNotice.Composition.Subject.Reference != nil {
		err = fhirClient.ReadOne(ctx, "/"+fhir.FromStringPtr(advanceNotice.Composition.Subject.Reference), &advanceNotice.Patient)
		if err != nil {
			return eoverdracht.AdvanceNotice{}, fmt.Errorf("error while fetching the transfer subject (patient): %w", err)
		}
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
			err = fhirClient.ReadOne(ctx, "/"+conditionID, &condition)
			if err != nil {
				return eoverdracht.AdvanceNotice{}, fmt.Errorf("error while fetching a advance notice condition (condition-id=%s): %w", conditionID, err)
			}
			advanceNotice.Problems = append(advanceNotice.Problems, condition)
		}
		if strings.HasPrefix(fhir.FromStringPtr(entry.Reference), "Procedure") {
			procedureID := fhir.FromStringPtr(entry.Reference)
			procedure := eoverdracht.Procedure{}
			err = fhirClient.ReadOne(ctx, "/"+procedureID, &procedure)
			if err != nil {
				return eoverdracht.AdvanceNotice{}, fmt.Errorf("error while fetching a advance notice procedure (procedure-id=%s): %w", procedureID, err)
			}
			advanceNotice.Interventions = append(advanceNotice.Interventions, procedure)
		}
	}

	return advanceNotice, nil
}

func (s service) findPatientByDossierID(ctx context.Context, customerID int, dossierID string) (*types.Patient, error) {
	transferDossier, err := s.dossierRepo.FindByID(ctx, customerID, dossierID)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch dossier: %w", err)
	}
	if transferDossier == nil {
		return nil, fmt.Errorf("dossier with id %s not found", dossierID)
	}
	patient, err := s.patientRepo.FindByID(ctx, customerID, string(transferDossier.PatientID))
	if err != nil {
		return nil, fmt.Errorf("error while fetching patient: %w", err)
	}
	if patient == nil {
		return nil, fmt.Errorf("patient with id %s not found", string(transferDossier.PatientID))
	}
	return patient, nil
}

func (s service) saveAdvanceNoticeFHIRResources(ctx context.Context, customerID int, advanceNotice eoverdracht.AdvanceNotice) error {
	fhirClient := s.localFHIRClientFactory(fhir.WithTenant(customerID))

	// Save the Patient
	err := fhirClient.CreateOrUpdate(ctx, advanceNotice.Patient)
	if err != nil {
		return err
	}
	// Save the all the problems
	for _, problem := range advanceNotice.Problems {
		err = fhirClient.CreateOrUpdate(ctx, problem)
		if err != nil {
			return err
		}
	}
	// Save all the interventions
	for _, intervention := range advanceNotice.Interventions {
		err = fhirClient.CreateOrUpdate(ctx, intervention)
		if err != nil {
			return err
		}
	}
	// At least save the composition
	err = fhirClient.CreateOrUpdate(ctx, advanceNotice.Composition)
	if err != nil {
		return err
	}
	return nil
}
