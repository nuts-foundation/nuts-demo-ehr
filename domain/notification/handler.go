package notification

import (
	"context"
	"fmt"
	nutsClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client"

	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer/receiver"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/registry"
)

type Notification struct {
	TaskID     string
	SenderDID  string
	CustomerID string
}

type Handler interface {
	Handle(ctx context.Context, notification Notification) error
}

type handler struct {
	nutsClient             *nutsClient.HTTPClient
	localFHIRClientFactory fhir.Factory
	transferService        receiver.TransferService
	registry               registry.OrganizationRegistry
}

func NewHandler(
	nutsClient *nutsClient.HTTPClient,
	localFHIRClientFactory fhir.Factory,
	transferReceiverService receiver.TransferService,
	registry registry.OrganizationRegistry,
) Handler {
	return &handler{
		nutsClient:             nutsClient,
		localFHIRClientFactory: localFHIRClientFactory,
		transferService:        transferReceiverService,
		registry:               registry,
	}
}

// Handle handles an incoming notification about an updated Task for one of its customers.
func (service *handler) Handle(ctx context.Context, notification Notification) error {
	fhirServer, err := service.registry.GetCompoundServiceEndpoint(ctx, notification.SenderDID, transfer.SenderServiceName, "fhir")
	if err != nil {
		return fmt.Errorf("error while looking up custodian's FHIR server (did=%s): %w", notification.SenderDID, err)
	}

	taskPath := fmt.Sprintf("/Task/%s", notification.TaskID)

	accessToken, err := service.nutsClient.RequestServiceAccessToken(ctx, notification.CustomerID, notification.SenderDID, "eOverdracht-sender")
	if err != nil {
		return err
	}

	task := &resources.Task{}
	client := fhir.NewFactory(fhir.WithURL(fhirServer), fhir.WithAuthToken(accessToken))

	// FIXME: add query params to filter on the owner so to only process the customer addressed in the notification
	err = client().ReadOne(ctx, taskPath, &task)
	if err != nil {
		return err
	}

	return service.transferService.CreateOrUpdate(ctx, fhir.FromCodePtr(task.Status), notification.CustomerID, notification.SenderDID, fhir.FromIDPtr(task.ID))
}
