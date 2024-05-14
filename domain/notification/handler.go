package notification

import (
	"context"
	"errors"
	"fmt"
	nutsClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client"

	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer/receiver"
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
	nutsClient             *nutsClient.HTTPClient
	localFHIRClientFactory fhir.Factory
	transferService        receiver.TransferService
	registry               registry.OrganizationRegistry
	vcr                    registry.VerifiableCredentialRegistry
}

func NewHandler(
	nutsClient *nutsClient.HTTPClient,
	localFHIRClientFactory fhir.Factory,
	transferReceiverService receiver.TransferService,
	registry registry.OrganizationRegistry,
	vcr registry.VerifiableCredentialRegistry,
) Handler {
	return &handler{
		nutsClient:             nutsClient,
		localFHIRClientFactory: localFHIRClientFactory,
		transferService:        transferReceiverService,
		registry:               registry,
		vcr:                    vcr,
	}
}

// Handle handles an incoming notification about an updated Task for one of its customers.
func (service *handler) Handle(ctx context.Context, notification Notification) error {
	fhirServer, err := service.registry.GetCompoundServiceEndpoint(ctx, notification.SenderDID, transfer.SenderServiceName, "fhir")
	if err != nil {
		return fmt.Errorf("error while looking up custodian's FHIR server (did=%s): %w", notification.SenderDID, err)
	}

	taskPath := fmt.Sprintf("/Task/%s", notification.TaskID)

	credentials, err := service.vcr.FindAuthorizationCredentials(
		ctx,
		&registry.VCRSearchParams{
			PurposeOfUse: transfer.SenderServiceName,
			Issuer:       notification.SenderDID,
			SubjectID:    notification.CustomerDID,
			ResourcePath: taskPath,
		},
	)
	if err != nil {
		return err
	}

	if len(credentials) == 0 {
		return errors.New("no NutsAuthorizationCredential found to retrieve the Task resource")
	}

	accessToken, err := service.nutsClient.RequestServiceAccessToken(ctx, notification.CustomerDID, notification.SenderDID, "eOverdracht-sender")
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

	return service.transferService.CreateOrUpdate(ctx, fhir.FromCodePtr(task.Status), notification.CustomerID, requesterDID, fhir.FromIDPtr(task.ID))
}
