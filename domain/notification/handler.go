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

func (service *handler) Handle(ctx context.Context, notification Notification) error {
	fhirServer, err := service.registry.GetCompoundServiceEndpoint(ctx, notification.SenderDID, "eOverdracht-sender", "fhir")
	if err != nil {
		return fmt.Errorf("error while looking up custodian's FHIR server (did=%s): %w", notification.SenderDID, err)
	}

	accessToken, err := service.auth.RequestAccessToken(ctx, notification.CustomerDID, notification.SenderDID, "eOverdracht-sender", nil)
	if err != nil {
		return err
	}

	var tasks []resources.Task

	client := fhir.NewFactory(fhir.WithURL(fhirServer), fhir.WithAuthToken(accessToken.AccessToken))

	err = client().ReadMultiple(ctx, "/Task", map[string]string{
		"code": fmt.Sprintf("%s|%s", fhir.SnomedCodingSystem, fhir.SnomedTransferCode),
	}, &tasks)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		// @TODO: Enable this
		//if fhir.FromStringPtr(task.Owner.Identifier.Value) != notification.SenderDID {
		//	continue
		//}

		err = service.transferService.CreateOrUpdate(ctx, notification.CustomerID, notification.SenderDID, fhir.FromIDPtr(task.ID))
		if err != nil {
			return err
		}
	}

	return nil
}
