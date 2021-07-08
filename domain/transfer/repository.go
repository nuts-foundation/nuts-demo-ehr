package transfer

import (
	"context"
	"time"

	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

type Repository interface {
	FindByID(ctx context.Context, customerID, id string) (*domain.Transfer, error)
	Create(ctx context.Context, customerID, dossierID, date, description string)
	Update(ctx context.Context, customerID, description string, date time.Time)
	Cancel(ctx context.Context, customerID, id string)
}