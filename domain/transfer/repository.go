package transfer

import (
	"context"
	"time"

	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

type Repository interface {
	FindByID(ctx context.Context, customerID, id string) (*domain.Transfer, error)
	FindByPatientID(ctx context.Context, customerID, patientID string) []domain.Transfer
	Create(ctx context.Context, customerID, dossierID, description string, date time.Time) domain.Transfer
	Update(ctx context.Context, customerID, description string, date time.Time, state domain.TransferStatus) error
	Cancel(ctx context.Context, customerID, id string)
	CreateNegotiation(ctx context.Context, transferID string, organizationDID string) (*domain.TransferNegotiation, error)
	ListNegotiations(ctx context.Context, customerID, transferID string) ([]domain.TransferNegotiation, error)
}

type Factory struct{}
