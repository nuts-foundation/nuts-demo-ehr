package receiver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/go-did/vc"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer"
	"github.com/nuts-foundation/nuts-demo-ehr/http/auth"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/registry"
)

type TransferService interface {
	CreateOrUpdate(ctx context.Context, status string, customerID int, senderDID, fhirTaskID string) error
	UpdateTransferRequestState(ctx context.Context, customerID int, requesterDID, fhirTaskID string, newState string) error
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

func NewTransferService(authService auth.Service, localFHIRClientFactory fhir.Factory, transferRepository TransferRepository, customerRepository customers.Repository, organizationRegistry registry.OrganizationRegistry, vcr registry.VerifiableCredentialRegistry) TransferService {
	return &service{
		auth:                   authService,
		localFHIRClientFactory: localFHIRClientFactory,
		transferRepo:           transferRepository,
		customerRepo:           customerRepository,
		registry:               organizationRegistry,
		vcr:                    vcr,
		notifier:               transfer.FireAndForgetNotifier{},
	}
}

func (s service) CreateOrUpdate(ctx context.Context, status string, customerID int, senderDID, fhirTaskID string) error {
	_, err := s.transferRepo.CreateOrUpdate(ctx, status, fhirTaskID, customerID, senderDID)
	if err != nil {
		return err
	}

	return nil
}

func (s service) UpdateTransferRequestState(ctx context.Context, customerID int, requesterDID, fhirTaskID string, newState string) error {
	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil || customer.Did == nil {
		return err
	}

	taskPath := fmt.Sprintf("Task/%s", fhirTaskID)
	client, err := s.getRemoteFHIRClient(ctx, requesterDID, *customer.Did, taskPath)
	if err != nil {
		return err
	}

	task, err := s.getRemoteTransferTask(ctx, client, fhirTaskID)
	if err != nil {
		return err
	}

	// state machine
	if (*task.Status == transfer.InProgressState && newState == transfer.CompletedState) || (*task.Status == transfer.RequestedState && newState == transfer.AcceptedState) {
		task.Status = fhir.ToCodePtr(newState)

		return client().CreateOrUpdate(ctx, task)
	}

	return fmt.Errorf("invalid state change from %s to %s", *task.Status, newState)

}

func (s service) getRemoteFHIRClient(ctx context.Context, custodianDID string, localActorDID string, resource string) (fhir.Factory, error) {
	fhirServer, err := s.registry.GetCompoundServiceEndpoint(ctx, custodianDID, transfer.SenderServiceName, "fhir")
	if err != nil {
		return nil, fmt.Errorf("error while looking up custodian's FHIR server (did=%s): %w", custodianDID, err)
	}
	credentials, err := s.vcr.FindAuthorizationCredentials(ctx, transfer.SenderServiceName, localActorDID, resource)

	var transformed = make([]vc.VerifiableCredential, len(credentials))
	for i, c := range credentials {
		bytes, err := json.Marshal(c)
		if err != nil {
			return nil, err
		}
		tCred := vc.VerifiableCredential{}
		if err = json.Unmarshal(bytes, &tCred); err != nil {
			return nil, err
		}
		transformed[i] = tCred
	}

	accessToken, err := s.auth.RequestAccessToken(ctx, localActorDID, custodianDID, transfer.SenderServiceName, transformed)
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
