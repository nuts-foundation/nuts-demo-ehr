package transfer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nuts-foundation/nuts-demo-ehr/http/auth"
	registry2 "github.com/nuts-foundation/nuts-demo-ehr/nuts/registry"
	"github.com/nuts-foundation/nuts-node/vcr/credential"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/sirupsen/logrus"

	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir/eoverdracht"
	sqlUtil "github.com/nuts-foundation/nuts-demo-ehr/sql"
)

// ReceiverServiceName contains the name of the eOverdracht receiver compound-service
const ReceiverServiceName = "eOverdracht-receiver"

// SenderServiceName contains the name of the eOverdracht sender compound-service
const SenderServiceName = "eOverdracht-sender"

type Service interface {
	// Create creates a new transfer
	Create(ctx context.Context, customerID string, dossierID string, description string, transferDate time.Time) (*domain.Transfer, error)

	CreateNegotiation(ctx context.Context, customerID, transferID, organizationDID string, transferDate time.Time) (*domain.TransferNegotiation, error)

	// ProposeAlternateDate updates the date on the domain.TransferNegotiation indicated by the negotiationID.
	// It updates the status to ON_HOLD_STATE
	ProposeAlternateDate(ctx context.Context, customerID, negotiationID string) (*domain.TransferNegotiation, error)

	// ConfirmNegotiation confirms the negotiation indicated by the negotiationID.
	// The updates the status to ACCEPTED_STATE.
	// It automatically cancels other negotiations of the domain.Transfer indicated by the transferID
	// by setting their status to CANCELLED_STATE.
	ConfirmNegotiation(ctx context.Context, customerID, negotiationID string) (*domain.TransferNegotiation, error)

	CancelNegotiation(ctx context.Context, customerID, negotiationID string) (*domain.TransferNegotiation, error)

	// GetTransferRequest tries to retrieve a transfer request from requesting care organization's FHIR server.
	GetTransferRequest(ctx context.Context, customerID, requestorDID string, fhirTaskID string) (*domain.TransferRequest, error)

	// AcceptTransferRequest accepts the given transfer request, allowing the sending organization to assign the patient transfer to the local organization
	AcceptTransferRequest(ctx context.Context, customerID, requestorDID, fhirTaskID string) error
}

type service struct {
	transferRepo           Repository
	auth                   auth.Service
	localFHIRClientFactory fhir.Factory // client for interacting with the local FHIR server
	customerRepo           customers.Repository
	registry               registry2.OrganizationRegistry
	vcr                    registry2.VerifiableCredentialRegistry
	notifier               Notifier
}

func NewTransferService(authService auth.Service, localFHIRClientFactory fhir.Factory, transferRepository Repository, customerRepository customers.Repository, organizationRegistry registry2.OrganizationRegistry, vcr registry2.VerifiableCredentialRegistry) *service {
	return &service{
		auth:                   authService,
		localFHIRClientFactory: localFHIRClientFactory,
		transferRepo:           transferRepository,
		customerRepo:           customerRepository,
		registry:               organizationRegistry,
		vcr:                    vcr,
		notifier:               fireAndForgetNotifier{},
	}
}

func (s service) CreateNegotiation(ctx context.Context, customerID, transferID, organizationDID string, transferDate time.Time) (*domain.TransferNegotiation, error) {
	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil || customer.Did == nil {
		return nil, err
	}

	var (
		negotiation          *domain.TransferNegotiation
		notificationEndpoint string
	)

	_, err = s.transferRepo.Update(ctx, customerID, transferID, func(transfer domain.Transfer) (*domain.Transfer, error) {
		// Validate transfer
		if transfer.Status == domain.TransferStatusCancelled ||
			transfer.Status == domain.TransferStatusCompleted ||
			transfer.Status == domain.TransferStatusAssigned {
			return nil, errors.New("can't start new transfer negotiation when status is 'cancelled', 'assigned' or 'completed'")
		}

		// Create negotiation and share it to the other party
		// TODO: Share transaction to this repository call as well
		var err error
		// Pre-emptively resolve the receiver organization's notification endpoint to reduce clutter, avoiding to make FHIR tasks when the receiving party eOverdracht registration is faulty.
		notificationEndpoint, err = s.registry.GetCompoundServiceEndpoint(ctx, organizationDID, ReceiverServiceName, "notification")
		if err != nil {
			return nil, err
		}

		transferTask := fhir.BuildNewTask(fhir.TaskProperties{
			RequesterID: *customer.Did,
			OwnerID:     organizationDID,
			Status:      REQUESTED_STATE,
			Input: []resources.TaskInputOutput{
				{
					Type:           &fhir.LoincAdvanceNoticeType,
					ValueReference: &datatypes.Reference{Reference: fhir.ToStringPtr(transfer.FhirAdvanceNoticeComposition)},
				},
			},
		})

		err = s.localFHIRClientFactory(fhir.WithTenant(customerID)).CreateOrUpdate(ctx, transferTask)
		if err != nil {
			return nil, err
		}

		if err := s.vcr.CreateAuthorizationCredential(ctx, "eOverdracht-receiver", *customer.Did, organizationDID, []credential.Resource{
			{
				Path:       fmt.Sprintf("/Task/%s", fhir.FromIDPtr(transferTask.ID)),
				Operations: []string{"update"},
			},
		}); err != nil {
			return nil, err
		}

		negotiation, err = s.transferRepo.CreateNegotiation(ctx, customerID, transferID, organizationDID, transfer.TransferDate.Time, fhir.FromIDPtr(transferTask.ID))
		if err != nil {
			return nil, err
		}

		// Update transfer.Status = requested
		//transfer.Status = domain.TransferStatusRequested
		return &transfer, nil
	})
	if err == nil {
		// Commit here, otherwise notifications to this server will deadlock on the uncommitted tx.
		tm, _ := sqlUtil.GetTransactionManager(ctx)
		if commitErr := tm.Commit(); commitErr != nil {
			return negotiation, commitErr
		}

		tokenResponse, err := s.auth.RequestAccessToken(ctx, *customer.Did, organizationDID, ReceiverServiceName, nil)
		if err != nil {
			return nil, err
		}

		if err = s.notifier.Notify(tokenResponse.AccessToken, notificationEndpoint); err != nil {
			// TODO: What to do here? Should we maybe rollback?
			logrus.Errorf("Unable to notify receiving care organization of updated FHIR task (did=%s): %s", organizationDID, err)
		}
	}

	return negotiation, err
}

func (s service) GetTransferRequest(ctx context.Context, customerID string, requestorDID string, fhirTaskID string) (*domain.TransferRequest, error) {
	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil || customer.Did == nil {
		return nil, err
	}

	client, err := s.getRemoteFHIRClient(ctx, requestorDID, *customer.Did)
	if err != nil {
		return nil, err
	}

	task, err := s.getTransferTask(ctx, client, fhirTaskID)
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

func (s service) AcceptTransferRequest(ctx context.Context, customerID, requestorDID, fhirTaskID string) error {
	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil || customer.Did == nil {
		return err
	}

	client, err := s.getRemoteFHIRClient(ctx, requestorDID, *customer.Did)
	if err != nil {
		return err
	}

	task, err := s.getTransferTask(ctx, client, fhirTaskID)
	if err != nil {
		return err
	}
	task.Status = fhir.ToCodePtr(ACCEPTED_STATE)
	// TODO: This doesn't work yet: the access token is missing NutsAuthorizationCredential
	return client().CreateOrUpdate(ctx, task)
}

func (s service) Create(ctx context.Context, customerID string, dossierID string, description string, transferDate time.Time) (*domain.Transfer, error) {
	elements := map[string]interface{}{
		"title": "Aanmeldbericht",
		"type":  fhir.LoincAdvanceNoticeType,
		// TODO: patient seems mandatory in the spec, but can only be sent when placer already
		// has patient in care to protect the identity of the patient during the negotiation phase.
		//"subject":  fhir.Reference{Reference: "Patient/Anonymous"},
		"author": eoverdracht.Practitioner{
			// TODO: Derive from authenticated user?
			Identifier: datatypes.Identifier{
				System: &fhir.UZICodingSystem,
				Value:  fhir.ToStringPtr("12345"),
			},
			Name: &datatypes.HumanName{
				Family: fhir.ToStringPtr("Demo EHR"),
				Given:  []datatypes.String{"Nuts"},
			},
		},
		// TODO: sections
	}
	composition := fhir.BuildNewComposition(elements)
	err := s.localFHIRClientFactory(fhir.WithTenant(customerID)).CreateOrUpdate(ctx, composition)
	if err != nil {
		return nil, err
	}
	transfer, err := s.transferRepo.Create(ctx, customerID, dossierID, description, transferDate, composition["id"].(string))
	if err != nil {
		return nil, err
	}
	return transfer, nil
}

func (s service) ProposeAlternateDate(ctx context.Context, customerID, negotiationID string) (*domain.TransferNegotiation, error) {
	panic("implement me")
}

func (s service) ConfirmNegotiation(ctx context.Context, customerID, negotiationID string) (*domain.TransferNegotiation, error) {
	panic("implement me")
}

func (s service) CancelNegotiation(ctx context.Context, customerID, negotiationID string) (*domain.TransferNegotiation, error) {
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

func (s service) getTransferTask(ctx context.Context, client fhir.Factory, fhirTaskID string) (resources.Task, error) {
	// TODO: Read AdvanceNotification here instead of the transfer task
	task := resources.Task{}
	err := client().ReadOne(ctx, "/Task/"+fhirTaskID, &task)
	if err != nil {
		return resources.Task{}, fmt.Errorf("error while looking up transfer task (task-id=%s): %w", fhirTaskID, err)
	}
	return task, nil
}
