package sender

import (
	"context"
	"errors"
	"fmt"
	nutsClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"

	sqlUtil "github.com/nuts-foundation/nuts-demo-ehr/sql"
	"github.com/sirupsen/logrus"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/dossier"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir/eoverdracht"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/patients"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/registry"
)

type TransferService interface {
	// AssignTransfer assigns a transfer directly to a single organization
	AssignTransfer(ctx context.Context, customerID int, transferID, organizationDID string) (*types.TransferNegotiation, error)

	// CreateTransfer creates a new transfer
	CreateTransfer(ctx context.Context, customerID int, request types.CreateTransferRequest) (*types.Transfer, error)

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

	// UpdateTaskState updates the Task resource. It updates the local DB, checks the statemachine, updates the FHIR record and sends a notification.
	UpdateTaskState(ctx context.Context, customer types.Customer, taskID string, newState string) error
}

type service struct {
	transferRepo           TransferRepository
	nutsClient             *nutsClient.HTTPClient
	localFHIRClientFactory fhir.Factory // client for interacting with the local FHIR server
	customerRepo           customers.Repository
	dossierRepo            dossier.Repository
	patientRepo            patients.Repository
	registry               registry.OrganizationRegistry
	vcr                    registry.VerifiableCredentialRegistry
	notifier               transfer.Notifier
}

func NewTransferService(nutsClient *nutsClient.HTTPClient, localFHIRClientFactory fhir.Factory, transferRepository TransferRepository, customerRepository customers.Repository, dossierRepo dossier.Repository, patientRepo patients.Repository, organizationRegistry registry.OrganizationRegistry, vcr registry.VerifiableCredentialRegistry, notifier transfer.Notifier) TransferService {
	return &service{
		nutsClient:             nutsClient,
		localFHIRClientFactory: localFHIRClientFactory,
		transferRepo:           transferRepository,
		customerRepo:           customerRepository,
		dossierRepo:            dossierRepo,
		patientRepo:            patientRepo,
		registry:               organizationRegistry,
		vcr:                    vcr,
		notifier:               notifier,
	}
}

func (s service) CreateTransfer(ctx context.Context, customerID int, request types.CreateTransferRequest) (*types.Transfer, error) {
	const createTransferErr = "could not create new transfer: %w"
	// Fetch the patient
	patient, err := s.findPatientByDossierID(ctx, customerID, string(request.DossierID))
	if err != nil {
		return nil, fmt.Errorf(createTransferErr, err)
	}

	// Build the advance notice resources
	advanceNotice := eoverdracht.NewFHIRBuilder().BuildAdvanceNotice(request, patient)
	fhirClient := s.localFHIRClientFactory(fhir.WithTenant(customerID))
	fhirService := eoverdracht.NewFHIRTransferService(fhirClient)

	// Save the resources to the fhir storage
	err = fhirService.CreateAdvanceNotice(ctx, advanceNotice)
	if err != nil {
		return nil, fmt.Errorf(createTransferErr, fmt.Errorf("unable to store advance notification fhir resources: %w", err))
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
	fhirClient := s.localFHIRClientFactory(fhir.WithTenant(customerID))
	fhirService := eoverdracht.NewFHIRTransferService(fhirClient)

	advanceNotice, err := fhirService.GetAdvanceNotice(ctx, dbTransfer.FhirAdvanceNoticeComposition)
	if err != nil {
		return types.Transfer{}, err
	}

	domainTransfer, err := eoverdracht.AdvanceNoticeToDomainTransfer(advanceNotice)
	if err != nil || customer.Did == nil {
		return types.Transfer{}, err
	}

	return types.Transfer{
		CarePlan:                      domainTransfer.CarePlan,
		TransferDate:                  domainTransfer.TransferDate,
		Patient:                       domainTransfer.Patient,
		DossierID:                     dbTransfer.DossierID,
		FhirAdvanceNoticeComposition:  dbTransfer.FhirAdvanceNoticeComposition,
		FhirNursingHandoffComposition: dbTransfer.FhirNursingHandoffComposition,
		Id:                            dbTransfer.Id,
		Status:                        dbTransfer.Status,
	}, nil
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

	fhirClient := s.localFHIRClientFactory(fhir.WithTenant(customerID))
	fhirTransferService := eoverdracht.NewFHIRTransferService(fhirClient)

	// Update the transfer
	_, err = s.transferRepo.Update(ctx, customerID, transferID, func(dbTransfer *types.Transfer) (*types.Transfer, error) {
		// Validate if transfer is in correct state to allow new negotiations
		if dbTransfer.Status == types.Cancelled ||
			dbTransfer.Status == types.Completed ||
			dbTransfer.Status == types.Assigned {
			return nil, errors.New("can't start new transfer negotiation when status is 'cancelled', 'assigned' or 'completed'")
		}

		compositionPath := fmt.Sprintf("/Composition/%s", dbTransfer.FhirAdvanceNoticeComposition)
		composition := fhir.Composition{}
		err = fhirClient.ReadOne(ctx, compositionPath, &composition)
		if err != nil {
			return nil, fmt.Errorf("could not create transfer negotiation: could not read fhir compositition: %w", err)
		}

		transferTask := eoverdracht.TransferTask{
			Status:          transfer.RequestedState,
			ReceiverDID:     organizationDID,
			SenderDID:       *customer.Did,
			AdvanceNoticeID: &dbTransfer.FhirAdvanceNoticeComposition,
		}

		transferTask, err = fhirTransferService.CreateTask(ctx, transferTask)
		if err != nil {
			return nil, fmt.Errorf("could not create FHIR task: %w", err)
		}

		// Build the list of resources for the authorization credential:
		authorizedResources := []registry.Resource{
			{
				Path:       fmt.Sprintf("/Task/%s", transferTask.ID),
				Operations: []string{"read", "update"},
			},
			{
				Path:           compositionPath,
				Operations:     []string{"read", "document"},
				UserContext:    true,
				AssuranceLevel: assuranceLevelLow(),
			},
		}

		// A list to store all the paths to FHIR resources associated with this advance notice
		// These paths must be included in the authorization credential
		resourcePaths := resourcePathsFromSection(composition.Section, []string{})
		// Include subject reference (patient)
		resourcePaths = append(resourcePaths, fhir.FromStringPtr(composition.Subject.Reference))
		for _, path := range resourcePaths {
			authorizedResources = append(authorizedResources, registry.Resource{
				Path:           path,
				Operations:     []string{"read", "document"},
				UserContext:    true,
				AssuranceLevel: assuranceLevelLow(),
			})
		}

		if err := s.vcr.CreateAuthorizationCredential(ctx, *customer.Did, &registry.NutsAuthorizationCredentialSubject{
			ID:           organizationDID,
			PurposeOfUse: transfer.SenderServiceName,
			Resources:    authorizedResources,
		}); err != nil {
			return nil, err
		}

		negotiation, err = s.transferRepo.CreateNegotiation(ctx, customerID, transferID, organizationDID, dbTransfer.TransferDate.Time, transferTask.ID)
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

		if err = s.sendNotification(ctx, customer, organizationDID, negotiation.TaskID); err != nil {
			// TODO: What to do here? Should we maybe rollback?
			logrus.Errorf("Unable to notify receiving care organization of updated FHIR task (did=%s): %s", organizationDID, err)
		}
	}

	return negotiation, err
}

func resourcePathsFromSection(sections []fhir.CompositionSection, paths []string) []string {
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

		fhirClient := s.localFHIRClientFactory(fhir.WithTenant(customerID))
		fhirService := eoverdracht.NewFHIRTransferService(fhirClient)

		// The advance notice contains a lot of the same resources which should also be used in the Nursing Handoff
		// Fetch the advanceNotice FHIR resources
		advanceNotice, err := fhirService.GetAdvanceNotice(ctx, dbTransfer.FhirAdvanceNoticeComposition)
		if err != nil {
			return nil, err
		}

		// Create nursing handoff composition based on the advanceNotice and patient
		nursingHandoffComposition, err := eoverdracht.NewFHIRBuilder().BuildNursingHandoffComposition(patient, advanceNotice)
		if err != nil {
			return nil, err
		}

		// Save nursing handoff composition in the FHIR store
		if err = s.localFHIRClientFactory(fhir.WithTenant(customerID)).CreateOrUpdate(ctx, nursingHandoffComposition, nil); err != nil {
			return nil, err
		}

		compositionID := nursingHandoffComposition.ID
		dbTransfer.FhirNursingHandoffComposition = (*string)(compositionID)
		dbTransfer.Status = types.Assigned

		// Update the task with the new state and nursing handoff composition ID
		if err := fhirService.UpdateTask(ctx, negotiation.TaskID, func(domainTask eoverdracht.TransferTask) eoverdracht.TransferTask {
			domainTask.Status = transfer.InProgressState
			domainTask.NursingHandoffID = dbTransfer.FhirNursingHandoffComposition
			return domainTask
		}); err != nil {
			return nil, fmt.Errorf("could not confirm negotiation: %w", fmt.Errorf("could not update task with in-progress state: %w", err))
		}

		// Revoke the old AuthorizationCredential for the Task and AdvanceNotice
		if err = s.vcr.RevokeAuthorizationCredential(ctx, transfer.SenderServiceName, negotiation.OrganizationDID, advanceNoticePath); err != nil {
			return nil, fmt.Errorf("unable to confirm negotiation: could not revoke advance notice authorization credential: %w", err)
		}

		authorizedResources := []registry.Resource{
			{
				Path:       fmt.Sprintf("/Task/%s", negotiation.TaskID),
				Operations: []string{"read", "update"},
			},
		}
		compositionPath := fmt.Sprintf("/Composition/%s", fhir.FromIDPtr(compositionID))
		// Add paths of resources of both the advance notice and the nursing handoff
		resourcePaths := resourcePathsFromSection(nursingHandoffComposition.Section, []string{advanceNoticePath, compositionPath})
		resourcePaths = resourcePathsFromSection(advanceNotice.Composition.Section, resourcePaths)
		// Add path to the complete, non-anonymous, patient
		resourcePaths = append(resourcePaths, "/"+fhir.FromStringPtr(nursingHandoffComposition.Subject.Reference))
		resourcePaths = append(resourcePaths, "/"+fhir.FromStringPtr(advanceNotice.Composition.Subject.Reference))

		// The resourcePaths may contain duplicates, hold a list of processedPaths
		processedPaths := map[string]struct{}{}
		for _, path := range resourcePaths {
			if _, exists := processedPaths[path]; exists {
				continue
			}
			authorizedResources = append(authorizedResources, registry.Resource{
				Path:           path,
				Operations:     []string{"read", "document"},
				UserContext:    true,
				AssuranceLevel: assuranceLevelLow(),
			})
			processedPaths[path] = struct{}{}
		}

		// Create a new AuthorizationCredential for the Task, AdvanceNotice and NursingHandoff
		if err = s.vcr.CreateAuthorizationCredential(ctx, *customer.Did, &registry.NutsAuthorizationCredentialSubject{
			ID:           negotiation.OrganizationDID,
			PurposeOfUse: transfer.SenderServiceName,
			Resources:    authorizedResources,
		}); err != nil {
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

		var errs []string
		for _, n := range notifications {
			err = s.sendNotification(ctx, n.customer, n.organizationDID, negotiation.TaskID)
			if err != nil {
				errs = append(errs, fmt.Errorf("sending to %s: %w", n.organizationDID, err).Error())
			}
		}
		if len(errs) > 0 {
			return nil, fmt.Errorf("one or more eOverdracht notifications failed: %s", strings.Join(errs, ", "))
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

	return negotiation, s.sendNotification(ctx, notification.customer, notification.organizationDID, negotiation.TaskID)
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

	fhirClient := s.localFHIRClientFactory(fhir.WithTenant(customer.Id))
	fhirService := eoverdracht.NewFHIRTransferService(fhirClient)
	if err := fhirService.UpdateTaskStatus(ctx, negotiation.TaskID, transfer.AcceptedState); err != nil {
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

	_ = s.sendNotification(ctx, not.customer, not.organizationDID, negotiation.TaskID)

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
		transferRecord.Status = types.Completed

		// update FHIR task
		fhirClient := s.localFHIRClientFactory(fhir.WithTenant(customer.Id))
		fhirService := eoverdracht.NewFHIRTransferService(fhirClient)
		if err := fhirService.UpdateTaskStatus(ctx, negotiation.TaskID, transfer.CompletedState); err != nil {
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

		_ = s.sendNotification(ctx, not.customer, not.organizationDID, negotiation.TaskID)
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
	fhirClient := s.localFHIRClientFactory(fhir.WithTenant(customerID))
	fhirService := eoverdracht.NewFHIRTransferService(fhirClient)
	if err := fhirService.UpdateTaskStatus(ctx, negotiation.TaskID, transfer.CancelledState); err != nil {
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

func (s service) sendNotification(ctx context.Context, customer *types.Customer, organizationDID string, fhirTaskID string) error {
	notificationEndpoint, err := s.registry.GetCompoundServiceEndpoint(ctx, organizationDID, transfer.ReceiverServiceName, "notification")
	if err != nil {
		return err
	}

	tokenResponse, err := s.nutsClient.RequestServiceAccessToken(ctx, *customer.Did, organizationDID, transfer.ReceiverServiceName)
	if err != nil {
		return err
	}

	endpoint := notificationEndpoint

	if !strings.HasSuffix(endpoint, "/") {
		endpoint += "/"
	}

	endpoint += fhirTaskID

	return retry.Do(
		func() error {
			return s.notifier.Notify(tokenResponse, endpoint)
		},
		retry.LastErrorOnly(true),
		retry.Attempts(60),
		retry.Delay(1*time.Second),
		retry.DelayType(retry.FixedDelay),
	)
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

func (s service) AssignTransfer(ctx context.Context, customerID int, transferID, organizationDID string) (*types.TransferNegotiation, error) {
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

	fhirClient := s.localFHIRClientFactory(fhir.WithTenant(customerID))
	fhirTransferService := eoverdracht.NewFHIRTransferService(fhirClient)

	// Update the transfer
	_, err = s.transferRepo.Update(ctx, customerID, transferID, func(dbTransfer *types.Transfer) (*types.Transfer, error) {
		// Validate if transfer is in correct state to allow new negotiations
		if dbTransfer.Status == types.Cancelled ||
			dbTransfer.Status == types.Completed ||
			dbTransfer.Status == types.Assigned {
			return nil, errors.New("can't start new transfer negotiation when status is 'cancelled', 'assigned' or 'completed'")
		}
		dbTransfer.Status = types.Assigned

		// the advance notice was created with the dbTransfer
		// it has to be updated to a NursingHandoff
		compositionPath := fmt.Sprintf("/Composition/%s", dbTransfer.FhirAdvanceNoticeComposition)
		composition := fhir.Composition{}
		err = fhirClient.ReadOne(ctx, compositionPath, &composition)
		if err != nil {
			return nil, fmt.Errorf("could not assign transfer negotiation: could not read fhir compositition: %w", err)
		}
		nursingHandoffComposition, err := s.advanceNoticeToNursingHandoff(ctx, customerID, dbTransfer)
		if err != nil {
			return nil, fmt.Errorf("could not assign transfer negotiation: failed to upgrade AdvanceNotice to NursingHandoff: %w", err)
		}
		// Save nursing handoff composition in the FHIR store
		if err = s.localFHIRClientFactory(fhir.WithTenant(customerID)).CreateOrUpdate(ctx, nursingHandoffComposition, nil); err != nil {
			return nil, err
		}
		dbTransfer.FhirNursingHandoffComposition = (*string)(nursingHandoffComposition.ID)

		// create the FHIR task with the Nurse Handoff
		transferTask := eoverdracht.TransferTask{
			Status:           transfer.InProgressState,
			ReceiverDID:      organizationDID,
			SenderDID:        *customer.Did,
			NursingHandoffID: dbTransfer.FhirNursingHandoffComposition,
		}

		transferTask, err = fhirTransferService.CreateTask(ctx, transferTask)
		if err != nil {
			return nil, fmt.Errorf("could not create FHIR task: %w", err)
		}

		if err := s.createAuthCredentials(ctx, &transferTask, nursingHandoffComposition, *customer.Did, organizationDID); err != nil {
			return nil, err
		}

		negotiation, err = s.transferRepo.CreateNegotiation(ctx, customerID, transferID, organizationDID, dbTransfer.TransferDate.Time, transferTask.ID)
		if err != nil {
			return nil, err
		}
		_, err = s.transferRepo.UpdateNegotiationState(ctx, customerID, negotiation.Id, transfer.InProgressState)
		if err != nil {
			return nil, err
		}

		return dbTransfer, nil
	})
	if err == nil {
		// Commit here, otherwise notifications to this server will deadlock on the uncommitted tx.
		tm, _ := sqlUtil.GetTransactionManager(ctx)
		if commitErr := tm.Commit(); commitErr != nil {
			return negotiation, commitErr
		}

		if err = s.sendNotification(ctx, customer, organizationDID, negotiation.TaskID); err != nil {
			// TODO: What to do here? Should we maybe rollback?
			logrus.Errorf("Unable to notify receiving care organization of updated FHIR task (did=%s): %s", organizationDID, err)
		}
	}

	return negotiation, err
}

// createAuthCredentials creates 2 authorization credentials, one for the Task, and one for the nursingHandoffComposition.
func (s service) createAuthCredentials(ctx context.Context, transferTask *eoverdracht.TransferTask, nursingHandoffComposition *fhir.Composition, customerDID, organizationDID string) error {
	// Create an Auth Credential for the Task
	authorizedTask := s.taskForNursingHandoff(transferTask.ID)
	if err := s.vcr.CreateAuthorizationCredential(ctx, customerDID, &registry.NutsAuthorizationCredentialSubject{
		ID:           organizationDID,
		PurposeOfUse: transfer.SenderServiceName,
		Resources:    authorizedTask,
	}); err != nil {
		return err
	}

	// Build the list of resources for the authorization credential:
	authorizedResources := s.resourcesForNursingHandoff(nursingHandoffComposition)

	if err := s.vcr.CreateAuthorizationCredential(ctx, customerDID, &registry.NutsAuthorizationCredentialSubject{
		ID:           organizationDID,
		PurposeOfUse: transfer.SenderServiceName,
		Resources:    authorizedResources,
	}); err != nil {
		return err
	}
	return nil
}

func (s service) advanceNoticeToNursingHandoff(ctx context.Context, customerID int, dbTransfer *types.Transfer) (*fhir.Composition, error) {
	patient, err := s.findPatientByDossierID(ctx, customerID, string(dbTransfer.DossierID))
	if err != nil {
		return nil, fmt.Errorf("could not fetch patient by dossierID: %w", err)
	}

	fhirClient := s.localFHIRClientFactory(fhir.WithTenant(customerID))
	fhirService := eoverdracht.NewFHIRTransferService(fhirClient)

	advanceNotice, err := fhirService.GetAdvanceNotice(ctx, dbTransfer.FhirAdvanceNoticeComposition)
	if err != nil {
		return nil, err
	}

	// Create nursing handoff composition based on the advanceNotice and patient
	nursingHandoffComposition, err := eoverdracht.NewFHIRBuilder().BuildNursingHandoffComposition(patient, advanceNotice)
	return &nursingHandoffComposition, err
}

func (s service) taskForNursingHandoff(taskID string) []registry.Resource {
	return []registry.Resource{
		{
			Path:       fmt.Sprintf("/Task/%s", taskID),
			Operations: []string{"read", "update"},
		},
	}
}

func (s service) resourcesForNursingHandoff(nursingHandoffComposition *fhir.Composition) []registry.Resource {
	authorizedResources := []registry.Resource{
		{
			Path:           fmt.Sprintf("/Composition/%s", fhir.FromIDPtr(nursingHandoffComposition.ID)),
			Operations:     []string{"read", "document"},
			UserContext:    true,
			AssuranceLevel: assuranceLevelLow(),
		},
		{
			Path:           fmt.Sprintf("/%s", fhir.FromStringPtr(nursingHandoffComposition.Subject.Reference)),
			Operations:     []string{"read"},
			UserContext:    true,
			AssuranceLevel: assuranceLevelLow(),
		},
	}
	// Add paths of resources of both the advance notice and the nursing handoff
	resourcePaths := resourcePathsFromSection(nursingHandoffComposition.Section, []string{})
	for _, path := range resourcePaths {
		authorizedResources = append(authorizedResources, registry.Resource{
			Path:           path,
			Operations:     []string{"read"},
			UserContext:    true,
			AssuranceLevel: assuranceLevelLow(),
		})
	}
	return authorizedResources
}

func assuranceLevelLow() *string {
	assuranceLevel := "low"
	return &assuranceLevel
}
