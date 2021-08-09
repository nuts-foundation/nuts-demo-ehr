package transfer

import (
	"context"
	"errors"
	"fmt"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"net/http"
	"time"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/registry"
	sqlUtil "github.com/nuts-foundation/nuts-demo-ehr/sql"
	"github.com/sirupsen/logrus"

	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/task"
)

type Service interface {
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
	GetTransferRequest(requestorDID string, fhirTaskID string) (*domain.TransferRequest, error)
}

type service struct {
	transferRepo Repository
	taskRepo     task.Repository
	customerRepo customers.Repository
	registry     registry.OrganizationRegistry
	notifier     Notifier
}

func NewTransferService(taskRespository task.Repository, transferRepository Repository, customerRepository customers.Repository, organizationRegistry registry.OrganizationRegistry) *service {
	return &service{
		taskRepo:     taskRespository,
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
		if transfer.Status == domain.TransferStatusCancelled || transfer.Status == domain.TransferStatusCompleted || transfer.Status == domain.TransferStatusAssigned {
			return nil, errors.New("can't start new transfer negotiation when status is 'cancelled', 'assigned' or 'completed'")
		}
		// Create negotiation and share it to the other party
		// TODO: Share transaction to this repository call as well
		var err error
		taskProperties := domain.TaskProperties{
			RequesterID: *customer.Did,
			OwnerID:     organizationDID,
		}

		// Pre-emptively resolve the receiver organization's notification endpoint to reduce clutter, avoiding to make FHIR tasks when the receiving party eOverdracht registration is faulty.
		notificationEndpoint, err = s.registry.GetCompoundServiceEndpoint(ctx, organizationDID, "eOverdracht-receiver", "notification")
		if err != nil {
			return nil, err
		}

		transferTask, err := s.taskRepo.Create(ctx, taskProperties)
		if err != nil {
			return nil, err
		}

		negotiation, err = s.transferRepo.CreateNegotiation(ctx, customerID, transferID, organizationDID, transfer.TransferDate.Time, transferTask.ID)
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

		if err = s.notifier.Notify(notificationEndpoint, *customer.Did, organizationDID); err != nil {
			// TODO: What to do here? Should we maybe rollback?
			logrus.Errorf("Unable to notify receiving care organization of updated FHIR task (did=%s): %w", organizationDID, err)
		}
	}
	return negotiation, err
}

func (s service) GetTransferRequest(ctx context.Context, requestorDID string, fhirTaskID string) (*domain.TransferRequest, error) {
	fhirServer, err := s.registry.GetCompoundServiceEndpoint(ctx, requestorDID, "eOverdracht-sender", "fhir")
	if err != nil {
		return nil, fmt.Errorf("error while looking up sender's FHIR server (did=%s): %w", requestorDID, err)
	}
	// TODO: Read AdvanceNotification here instead of the transfer task
	task, err := fhir.NewClient(fhirServer).GetResource("/Task/" + fhirTaskID)
	if err != nil {
		return nil, fmt.Errorf("error while looking up transfer task (fhir-server=%s, task-id=%d): %w", fhirServer, fhirTaskID, err)
	}
	
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
