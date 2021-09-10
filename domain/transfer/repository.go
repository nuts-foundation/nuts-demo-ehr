package transfer

import (
	"context"
	"time"

	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

// All possible states as described by the Nictiz eOverdracht v4.0:
// https://informatiestandaarden.nictiz.nl/wiki/vpk:V4.0_FHIR_eOverdracht#Using_Task_to_manage_the_workflow
const REQUESTED_STATE = "requested"
const ACCEPTED_STATE = "accepted"
const REJECTED_STATE = "rejected"
const ON_HOLD_STATE = "on-hold"
const CANCELLED_STATE = "cancelled"
const IN_PROGRESS_STATE = "in-progress"
const COMPLETED_STATE = "completed"

type Repository interface {
	FindByID(ctx context.Context, customerID int, transferID string) (*domain.Transfer, error)
	FindByPatientID(ctx context.Context, customerID int, patientID string) ([]domain.Transfer, error)
	Create(ctx context.Context, customerID int, dossierID, description string, date time.Time, fhirAdvanceNoticeComposition string) (*domain.Transfer, error)

	FindNegotiationByID(ctx context.Context, customerID int, negotiationID string) (*domain.TransferNegotiation, error)
	FindNegotiationByTaskID(ctx context.Context, customerID int, taskID string) (*domain.TransferNegotiation, error)

	// Update finds the correct Transfer applies the updateFn and then stores the Transfer
	// it uses the Transaction from the context but does not commit it.
	Update(ctx context.Context, customerID int, transferID string, updateFn func(c *domain.Transfer) (*domain.Transfer, error)) (*domain.Transfer, error)

	// Cancel cancels the indicated domain.Transfer and all its domain.TransferNegotiation`s
	Cancel(ctx context.Context, customerID int, transferID string) (*domain.Transfer, error)

	// CreateNegotiation creates a new domain.TransferNegotiation for the indicated domain.Transfer and
	// the care organisation indicated by the organisationDID.
	// The status will be set to REQUESTED_STATE.
	// It fails when their exists another domain.TransferNegotiation for this transfer with the same organisationDID and
	// a status other than CANCELLED_STATE.
	CreateNegotiation(ctx context.Context, customerID int, transferID, organizationDID string, transferDate time.Time, taskID string) (*domain.TransferNegotiation, error)

	// ProposeAlternateDate updates the date on the domain.TransferNegotiation indicated by the negotiationID.
	// It updates the status to ON_HOLD_STATE
	ProposeAlternateDate(ctx context.Context, customerID int, negotiationID string) (*domain.TransferNegotiation, error)

	// ConfirmNegotiation confirms the negotiation indicated by the negotiationID.
	// The updates the status to ACCEPTED_STATE.
	// It automatically cancels other negotiations of the domain.Transfer indicated by the transferID
	// by setting their status to CANCELLED_STATE.
	ConfirmNegotiation(ctx context.Context, customerID int, negotiationID string) (*domain.TransferNegotiation, error)

	CancelNegotiation(ctx context.Context, customerID int, negotiationID string) (*domain.TransferNegotiation, error)

	// UpdateNegotiationState updates the negotiation with the new state.
	UpdateNegotiationState(ctx context.Context, customerID int, negotiationID string, newState domain.TransferNegotiationStatusStatus) (*domain.TransferNegotiation, error)

	// ListNegotiations returns a list of negotiations for the indicated transfer
	ListNegotiations(ctx context.Context, customerID int, transferID string) ([]domain.TransferNegotiation, error)
}
