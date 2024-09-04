package sender

import (
	"context"
	"time"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
)

type TransferRepository interface {
	FindByID(ctx context.Context, customerID, transferID string) (*types.Transfer, error)
	FindByPatientID(ctx context.Context, customerID, patientID string) ([]types.Transfer, error)
	Create(ctx context.Context, customerID, dossierID string, date time.Time, fhirAdvanceNoticeCompositionID string) (*types.Transfer, error)

	FindNegotiationByID(ctx context.Context, customerID, negotiationID string) (*types.TransferNegotiation, error)
	FindNegotiationByTaskID(ctx context.Context, customerID, taskID string) (*types.TransferNegotiation, error)

	// Update finds the correct Transfer applies the updateFn and then stores the Transfer
	// it uses the Transaction from the context but does not commit it.
	Update(ctx context.Context, customerID, transferID string, updateFn func(c *types.Transfer) (*types.Transfer, error)) (*types.Transfer, error)

	// Cancel cancels the indicated domain.Transfer and all its domain.TransferNegotiation`s
	Cancel(ctx context.Context, customerID, transferID string) (*types.Transfer, error)

	// CreateNegotiation creates a new domain.TransferNegotiation for the indicated domain.Transfer and
	// the care organisation indicated by the organisationDID.
	// The status will be set to REQUESTED_STATE.
	// It fails when their exists another domain.TransferNegotiation for this transfer with the same organisationDID and
	// a status other than CANCELLED_STATE.
	CreateNegotiation(ctx context.Context, customerID, transferID, organizationDID string, transferDate time.Time, taskID string) (*types.TransferNegotiation, error)

	// ProposeAlternateDate updates the date on the domain.TransferNegotiation indicated by the negotiationID.
	// It updates the status to ON_HOLD_STATE
	ProposeAlternateDate(ctx context.Context, customerID, negotiationID string) (*types.TransferNegotiation, error)

	// ConfirmNegotiation confirms the negotiation indicated by the negotiationID.
	// The updates the status to ACCEPTED_STATE.
	// It automatically cancels other negotiations of the domain.Transfer indicated by the transferID
	// by setting their status to CANCELLED_STATE.
	ConfirmNegotiation(ctx context.Context, customerID, negotiationID string) (*types.TransferNegotiation, error)

	CancelNegotiation(ctx context.Context, customerID, negotiationID string) (*types.TransferNegotiation, error)

	// UpdateNegotiationState updates the negotiation with the new state.
	UpdateNegotiationState(ctx context.Context, customerID, negotiationID string, newState types.FHIRTaskStatus) (*types.TransferNegotiation, error)

	// ListNegotiations returns a list of negotiations for the indicated transfer
	ListNegotiations(ctx context.Context, customerID, transferID string) ([]types.TransferNegotiation, error)
}
