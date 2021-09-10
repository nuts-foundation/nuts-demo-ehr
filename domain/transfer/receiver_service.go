package transfer

import (
	"context"
	"fmt"
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
)

// ReceiverServiceName contains the name of the eOverdracht receiver compound-service
const ReceiverServiceName = "eOverdracht-receiver"

func (s service) GetTransferRequest(ctx context.Context, customerID int, requestorDID string, fhirTaskID string) (*domain.TransferRequest, error) {
	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil || customer.Did == nil {
		return nil, err
	}

	client, err := s.getRemoteFHIRClient(ctx, requestorDID, *customer.Did)
	if err != nil {
		return nil, err
	}

	task, err := s.getRemoteTransferTask(ctx, client, fhirTaskID)
	if err != nil {
		return nil, err
	}
	organization, err := s.registry.Get(ctx, requestorDID)
	if err != nil {
		return nil, err
	}
	// TODO: Do we need nil checks?
	transferDate, _ := time.Parse(time.RFC3339, string(*task.Meta.LastUpdated))
	return &domain.TransferRequest{
		Description:  "TODO",
		Sender:       *organization,
		TransferDate: openapi_types.Date{Time: transferDate},
		Status:       fhir.FromCodePtr(task.Status),
	}, nil
}

func (s service) AcceptTransferRequest(ctx context.Context, customerID int, requestorDID, fhirTaskID string) error {
	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil || customer.Did == nil {
		return err
	}

	client, err := s.getRemoteFHIRClient(ctx, requestorDID, *customer.Did)
	if err != nil {
		return err
	}

	task, err := s.getRemoteTransferTask(ctx, client, fhirTaskID)
	if err != nil {
		return err
	}
	task.Status = fhir.ToCodePtr(ACCEPTED_STATE)
	// TODO: This doesn't work yet: the access token is missing NutsAuthorizationCredential
	return client().CreateOrUpdate(ctx, task)
}

func (s service) ProposeAlternateDate(ctx context.Context, customerID int, negotiationID string) (*domain.TransferNegotiation, error) {
	panic("implement me")
}

func (s service) RejectNegotiation(ctx context.Context, customerID int, negotiationID string) (*domain.TransferNegotiation, error) {
	panic("implement me")
}

func (s service) getRemoteFHIRClient(ctx context.Context, custodianDID string, localActorDID string) (fhir.Factory, error) {
	fhirServer, err := s.registry.GetCompoundServiceEndpoint(ctx, custodianDID, SenderServiceName, "fhir")
	if err != nil {
		return nil, fmt.Errorf("error while looking up custodian's FHIR server (did=%s): %w", custodianDID, err)
	}
	accessToken, err := s.auth.RequestAccessToken(ctx, localActorDID, custodianDID, SenderServiceName, nil)
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
