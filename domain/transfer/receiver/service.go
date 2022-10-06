package receiver

import (
	"context"
	"fmt"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir/eoverdracht"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	"github.com/nuts-foundation/nuts-demo-ehr/http/auth"
	auth2 "github.com/nuts-foundation/nuts-demo-ehr/nuts/client/auth"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/registry"
)

type TransferService interface {
	// CreateOrUpdate creates or updates an incoming transfer record in the local storage
	CreateOrUpdate(ctx context.Context, status string, customerID int, senderDID, fhirTaskID string) error
	UpdateTransferRequestState(ctx context.Context, customerID int, requesterDID, fhirTaskID string, newState string) error
	GetTransferRequest(ctx context.Context, customerID int, requesterDID string, identity auth2.VerifiablePresentation, fhirTaskID string) (*types.TransferRequest, error)
}

type service struct {
	transferRepo           TransferRepository
	notifier               transfer.Notifier
	auth                   auth.Service
	localFHIRClientFactory fhir.Factory // client for interacting with the local FHIR server
	customerRepo           customers.Repository
	registry               registry.OrganizationRegistry
	vcr                    registry.VerifiableCredentialRegistry
}

func NewTransferService(authService auth.Service, localFHIRClientFactory fhir.Factory, transferRepository TransferRepository, customerRepository customers.Repository, organizationRegistry registry.OrganizationRegistry, vcr registry.VerifiableCredentialRegistry, notifier transfer.Notifier) TransferService {
	return &service{
		auth:                   authService,
		localFHIRClientFactory: localFHIRClientFactory,
		transferRepo:           transferRepository,
		customerRepo:           customerRepository,
		registry:               organizationRegistry,
		vcr:                    vcr,
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

	taskPath := fmt.Sprintf("/Task/%s", fhirTaskID)
	fhirClient, err := s.getRemoteFHIRClient(ctx, requesterDID, *customer.Did, taskPath, nil)
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

func (s service) GetTransferRequest(ctx context.Context, customerID int, requesterDID string, identity auth2.VerifiablePresentation, fhirTaskID string) (*types.TransferRequest, error) {
	const getTransferRequestErr = "unable to get transferRequest: %w"

	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil || customer.Did == nil {
		return nil, fmt.Errorf("unable to find customer: %w", err)
	}

	taskPath := fmt.Sprintf("/Task/%s", fhirTaskID)
	fhirClient, err := s.getRemoteFHIRClient(ctx, requesterDID, *customer.Did, taskPath, &identity)
	if err != nil {
		return nil, err
	}

	fhirReceiverService := eoverdracht.NewFHIRTransferService(fhirClient)

	task, err := fhirReceiverService.GetTask(ctx, fhirTaskID)
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

	// if it contains an AdvanceNotice
	if task.AdvanceNoticeID != nil {
		advanceNotice, err := fhirReceiverService.GetAdvanceNotice(ctx, *task.AdvanceNoticeID)
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
		nursingHandoff, err := fhirReceiverService.GetNursingHandoff(ctx, *task.NursingHandoffID)
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

func (s service) getRemoteFHIRClient(ctx context.Context, authorizerDID string, localRequesterDID string, resource string, identity *auth2.VerifiablePresentation) (fhir.Client, error) {
	fhirServer, err := s.registry.GetCompoundServiceEndpoint(ctx, authorizerDID, transfer.SenderServiceName, "fhir")
	if err != nil {
		return nil, fmt.Errorf("error while looking up authorizer's FHIR server (did=%s): %w", authorizerDID, err)
	}

	searchParams := registry.VCRSearchParams{
		PurposeOfUse: transfer.SenderServiceName,
		SubjectID:    localRequesterDID,
		Issuer:       authorizerDID,
		ResourcePath: resource,
	}

	credentials, err := s.vcr.FindAuthorizationCredentials(ctx, &searchParams)
	if err != nil {
		return nil, err
	}
	if len(credentials) == 0 {
		return nil, fmt.Errorf("no credentials found")
	}

	accessToken, err := s.auth.RequestAccessToken(ctx, localRequesterDID, authorizerDID, transfer.SenderServiceName, credentials, identity)

	if err != nil {
		return nil, err
	}

	return fhir.NewFactory(fhir.WithURL(fhirServer), fhir.WithAuthToken(accessToken.AccessToken))(), nil
}
