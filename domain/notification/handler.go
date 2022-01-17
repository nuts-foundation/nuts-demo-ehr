package notification

import (
	"context"
	"fmt"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer/receiver"
	"github.com/nuts-foundation/nuts-demo-ehr/http/auth"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/registry"
)

type Notification struct {
	TaskID      string
	SenderDID   string
	CustomerDID string
	CustomerID  int
}

type Handler interface {
	Handle(ctx context.Context, notification Notification) error
}

type handler struct {
	auth                   auth.Service
	localFHIRClientFactory fhir.Factory
	transferService        receiver.TransferService
	registry               registry.OrganizationRegistry
}

func NewHandler(auth auth.Service, localFHIRClientFactory fhir.Factory, tranferReceiverService receiver.TransferService, registry registry.OrganizationRegistry) Handler {
	return &handler{
		auth:                   auth,
		localFHIRClientFactory: localFHIRClientFactory,
		transferService:        tranferReceiverService,
		registry:               registry,
	}
}

// Handle handles an incoming notification about an updated Task for one of its customers.
func (service *handler) Handle(ctx context.Context, notification Notification) error {
	fhirServer, err := service.registry.GetCompoundServiceEndpoint(ctx, notification.SenderDID, "eOverdracht-sender", "fhir")
	if err != nil {
		return fmt.Errorf("error while looking up custodian's FHIR server (did=%s): %w", notification.SenderDID, err)
	}

	accessToken, err := service.auth.RequestAccessToken(ctx, notification.CustomerDID, notification.SenderDID, "eOverdracht-sender", nil, nil)
	if err != nil {
		return err
	}

	task := &resources.Task{}
	client := fhir.NewFactory(fhir.WithURL(fhirServer), fhir.WithAuthToken(accessToken.AccessToken))

	// FIXME: add query params to filter on the owner so to only process the customer addressed in the notification
	err = client().ReadOne(ctx, fmt.Sprintf("/Task/%s", notification.TaskID), &task)
	if err != nil {
		return err
	}

	// check if Task owner is set
	owner := task.Owner
	if owner == nil {
		return nil
	}

	// check if Task requester is set
	requester := task.Requester
	if requester == nil {
		return nil
	}

	requesterDID := fhir.FromStringPtr(requester.Agent.Identifier.Value)
	ownerDID := fhir.FromStringPtr(owner.Identifier.Value)

	// Check if the requester is the same as the sender of the notification
	if requesterDID != notification.SenderDID {
		return nil
	}

	// Don't update tasks for other customers since we do not have this customers ID.
	if ownerDID != notification.CustomerDID {
		return nil
	}

	err = service.transferService.CreateOrUpdate(ctx, fhir.FromCodePtr(task.Status), notification.CustomerID, requesterDID, fhir.FromIDPtr(task.ID))
	if err != nil {
		return err
	}

	return nil
}
