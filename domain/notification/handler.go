package notification

import (
	"context"
	"fmt"
	"time"

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

// Handle handles an incoming notification about an updated Task for one of its customers.
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

	// FIXME: add query params to filter on the owner so to only process the customer addressed in the notification
	err = client().ReadMultiple(ctx, "/Task", map[string]string{
		"code":         fmt.Sprintf("%s|%s", fhir.SnomedCodingSystem, fhir.SnomedTransferCode), // filter on transfer tasks
		"_lastUpdated": fmt.Sprintf("ge%sT00:00:00", time.Now().Format("2006-01-02")),          // filter on date
		"_count":       "80",                                                                   // prevent having to fetch multiple pages
	}, &tasks)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		// check if Task owner is set
		owner := task.Owner
		if owner == nil {
			continue
		}

		// check if Task requester is set
		requester := task.Requester
		if requester == nil {
			continue
		}

		requesterDID := fhir.FromStringPtr(requester.Agent.Identifier.Value)
		ownerDID := fhir.FromStringPtr(owner.Identifier.Value)

		// Check if the requester is the same as the sender of the notification
		if requesterDID != notification.SenderDID {
			continue
		}

		// Don't update tasks for other customers since we do not have this customers ID.
		if ownerDID != notification.CustomerDID {
			continue
		}

		err = service.transferService.CreateOrUpdate(ctx, notification.CustomerID, requesterDID, fhir.FromIDPtr(task.ID))
		if err != nil {
			return err
		}
	}

	return nil
}
