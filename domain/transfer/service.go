package transfer

import (
	"context"
	"errors"
	"fmt"
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/sirupsen/logrus"

	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/auth"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir/eoverdracht"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/registry"
	sqlUtil "github.com/nuts-foundation/nuts-demo-ehr/sql"
	"github.com/nuts-foundation/nuts-node/vcr/credential"
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
	GetTransferRequest(ctx context.Context, requestorDID string, fhirTaskID string) (*domain.TransferRequest, error)
}

type service struct {
	transferRepo Repository
	auth         auth.Service
	fhirClient   fhir.Client
	customerRepo customers.Repository
	registry     registry.OrganizationRegistry
	vcr          registry.VerifiableCredentialRegistry
	notifier     Notifier
}

func NewTransferService(authService auth.Service, fhirClient fhir.Client, transferRepository Repository, customerRepository customers.Repository, organizationRegistry registry.OrganizationRegistry, vcr registry.VerifiableCredentialRegistry) *service {
	return &service{
		vcr:          vcr,
		auth:         authService,
		fhirClient:   fhirClient,
		transferRepo: transferRepository,
		customerRepo: customerRepository,
		registry:     organizationRegistry,
		notifier:     fireAndForgetNotifier{},
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
		taskProperties := fhir.TaskProperties{
			RequesterID: *customer.Did,
			OwnerID:     organizationDID,
			Input: []resources.TaskInputOutput{
				{
					Type:           &fhir.LoincAdvanceNoticeType,
					ValueReference: &datatypes.Reference{Reference: fhir.ToStringPtr(transfer.FhirAdvanceNoticeComposition)},
				},
			},
		}

		// Pre-emptively resolve the receiver organization's notification endpoint to reduce clutter, avoiding to make FHIR tasks when the receiving party eOverdracht registration is faulty.
		notificationEndpoint, err = s.registry.GetCompoundServiceEndpoint(ctx, organizationDID, ReceiverServiceName, "notification")
		if err != nil {
			return nil, err
		}

		transferTask := fhir.BuildNewTask(taskProperties)
		err = s.fhirClient.CreateOrUpdate(ctx, transferTask)
		if err != nil {
			return nil, err
		}

		if err := s.vcr.CreateAuthorizationCredential(ctx, "eOverdracht-receiver", *customer.Did, organizationDID, []credential.Resource{
			{
				Path:       fmt.Sprintf("/Task/%s", transferTask.ID),
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

		if err = s.notifier.Notify(tokenResponse.AccessToken, notificationEndpoint, organizationDID); err != nil {
			// TODO: What to do here? Should we maybe rollback?
			logrus.Errorf("Unable to notify receiving care organization of updated FHIR task (did=%s): %w", organizationDID, err)
		}
	}

	return negotiation, err
}

func (s service) GetTransferRequest(ctx context.Context, requestorDID string, fhirTaskID string) (*domain.TransferRequest, error) {
	task, err := s.getTransferTask(ctx, requestorDID, fhirTaskID)
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
	}, nil
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
	err := s.fhirClient.CreateOrUpdate(ctx, composition)
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

func (s service) getTransferTask(ctx context.Context, requestorDID string, fhirTaskID string) (resources.Task, error) {
	fhirServer, err := s.registry.GetCompoundServiceEndpoint(ctx, requestorDID, SenderServiceName, "fhir")
	if err != nil {
		return resources.Task{}, fmt.Errorf("error while looking up sender's FHIR server (did=%s): %w", requestorDID, err)
	}
	// TODO: Read AdvanceNotification here instead of the transfer task
	task := resources.Task{}
	err = fhir.NewClient(fhirServer).ReadOne(ctx, "/Task/"+fhirTaskID, &task)
	if err != nil {
		return resources.Task{}, fmt.Errorf("error while looking up transfer task (fhir-server=%s, task-id=%d): %w", fhirServer, fhirTaskID, err)
	}
	return task, nil
}
