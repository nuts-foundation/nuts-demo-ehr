package receiver

import (
	"context"
	"fmt"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir/eoverdracht"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/client"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/registry"
)

type TransferService interface {
	// CreateOrUpdate creates or updates an incoming transfer record in the local storage
	CreateOrUpdate(ctx context.Context, status string, customerID int, senderDID, fhirTaskID string) error
	UpdateTransferRequestState(ctx context.Context, customerID int, requesterDID, fhirTaskID string, newState string) error
	GetTransferRequest(ctx context.Context, customerID int, requesterDID string, fhirTaskID string, token string) (*types.TransferRequest, error)
}

type service struct {
	transferRepo           TransferRepository
	notifier               transfer.Notifier
	nutsClient             *client.HTTPClient
	localFHIRClientFactory fhir.Factory // client for interacting with the local FHIR server
	customerRepo           customers.Repository
	registry               registry.OrganizationRegistry
}

func NewTransferService(nutsClient *client.HTTPClient, localFHIRClientFactory fhir.Factory, transferRepository TransferRepository, customerRepository customers.Repository, organizationRegistry registry.OrganizationRegistry, notifier transfer.Notifier) TransferService {
	return &service{
		nutsClient:             nutsClient,
		localFHIRClientFactory: localFHIRClientFactory,
		transferRepo:           transferRepository,
		customerRepo:           customerRepository,
		registry:               organizationRegistry,
		notifier:               notifier,
	}
}

func (s service) CreateOrUpdate(ctx context.Context, status string, customerID int, senderDID, fhirTaskID string) error {
	_, err := s.transferRepo.CreateOrUpdate(ctx, status, fhirTaskID, customerID, senderDID)
	return err
}

func (s service) UpdateTransferRequestState(ctx context.Context, customerID int, requesterDID, fhirTaskID string, newState string) error {
	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil || customer.Did == nil {
		return err
	}

	fhirClient, err := s.getServiceFHIRClient(ctx, requesterDID, *customer.Did)
	if err != nil {
		return err
	}

	fhirService := eoverdracht.NewFHIRTransferService(fhirClient)
	task, err := fhirService.GetTask(ctx, fhirTaskID)
	if err != nil {
		return err
	}

	// state machine
	if (task.Status == transfer.InProgressState && newState == transfer.CompletedState) ||
		(task.Status == transfer.RequestedState && newState == transfer.AcceptedState) {
		err = fhirService.UpdateTaskStatus(ctx, fhirTaskID, newState)
		if err != nil {
			return err
		}
		// update was a success. Get the remote task again and update the local transfer_request
		task, err = fhirService.GetTask(ctx, fhirTaskID)
		if err != nil {
			return err
		}
		_, err = s.transferRepo.CreateOrUpdate(ctx, task.Status, fhirTaskID, customerID, requesterDID)
		if err != nil {
			return fmt.Errorf("could update incomming transfers with new state")
		}
		return nil
	}

	return fmt.Errorf("invalid state change from %s to %s", task.Status, newState)
}

func (s service) GetTransferRequest(ctx context.Context, customerID int, requesterDID string, fhirTaskID string, accessToken string) (*types.TransferRequest, error) {
	const getTransferRequestErr = "unable to get transferRequest: %w"

	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil || customer.Did == nil {
		return nil, fmt.Errorf("unable to find customer: %w", err)
	}

	// First get the task, this uses a separate task auth credential
	fhirTaskClient, err := s.getUserFHIRClient(ctx, requesterDID, *customer.Did, accessToken)
	if err != nil {
		return nil, err
	}
	fhirTaskReceiverService := eoverdracht.NewFHIRTransferService(fhirTaskClient)
	task, err := fhirTaskReceiverService.GetTask(ctx, fhirTaskID)
	if err != nil {
		return nil, fmt.Errorf(getTransferRequestErr, err)
	}

	organization, err := s.registry.Get(ctx, requesterDID)
	if err != nil {
		return nil, fmt.Errorf("unable to get organization from registry: %w", err)
	}

	transferRequest := types.TransferRequest{
		Sender: types.FromNutsOrganization(*organization),
		Status: task.Status,
	}

	if task.Status == transfer.CompletedState || task.Status == transfer.CancelledState {
		return &transferRequest, nil
	}

	if task.AdvanceNoticeID != nil {
		fhirCompositionClient, err := s.getServiceFHIRClient(ctx, requesterDID, *customer.Did)
		if err != nil {
			return nil, err
		}
		fhirCompositionService := eoverdracht.NewFHIRTransferService(fhirCompositionClient)

		// if it contains an AdvanceNotice
		advanceNotice, err := fhirCompositionService.GetAdvanceNotice(ctx, *task.AdvanceNoticeID)
		if err != nil {
			return nil, fmt.Errorf("unable to get advance notice: %w", err)
		}
		transferRequest.AdvanceNotice, err = eoverdracht.AdvanceNoticeToDomainTransfer(advanceNotice)
		if err != nil {
			return nil, err
		}
	}

	// If the task input contains the nursing handoff
	if task.NursingHandoffID != nil {
		fhirCompositionClient, err := s.getServiceFHIRClient(ctx, requesterDID, *customer.Did)
		if err != nil {
			return nil, err
		}
		fhirCompositionService := eoverdracht.NewFHIRTransferService(fhirCompositionClient)
		nursingHandoff, err := fhirCompositionService.GetNursingHandoff(ctx, *task.NursingHandoffID)
		if err != nil {
			return nil, fmt.Errorf("unable to get nursing handoff: %w", err)
		}
		domainNursingHandoff, err := eoverdracht.NursingHandoffToDomainTransfer(nursingHandoff)
		if err != nil {
			return nil, err
		}
		transferRequest.NursingHandoff = &domainNursingHandoff
	}

	return &transferRequest, nil
}

func (s service) getServiceFHIRClient(ctx context.Context, authorizerDID string, localRequesterDID string) (fhir.Client, error) {
	fhirServer, err := s.registry.GetCompoundServiceEndpoint(ctx, authorizerDID, transfer.SenderServiceName, "fhir")
	if err != nil {
		return nil, fmt.Errorf("error while looking up authorizer's FHIR server (did=%s): %w", authorizerDID, err)
	}

	// TODO: This should be the user access token instead when medical data is involved,
	// but this depends on the scope mapping which then has to change for v6
	accessToken, err := s.nutsClient.RequestServiceAccessToken(ctx, localRequesterDID, authorizerDID, transfer.SenderServiceName)
	if err != nil {
		return nil, err
	}
	return fhir.NewFactory(fhir.WithURL(fhirServer), fhir.WithAuthToken(accessToken))(), nil
}

func (s service) getUserFHIRClient(ctx context.Context, authorizerDID string, localRequesterDID string, accessToken string) (fhir.Client, error) {
	fhirServer, err := s.registry.GetCompoundServiceEndpoint(ctx, authorizerDID, transfer.SenderServiceName, "fhir")
	if err != nil {
		return nil, fmt.Errorf("error while looking up authorizer's FHIR server (did=%s): %w", authorizerDID, err)
	}

	return fhir.NewFactory(fhir.WithURL(fhirServer), fhir.WithAuthToken(accessToken))(), nil
}
