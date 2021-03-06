package dossier

import (
	"context"

	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

type Repository interface {
	FindByID(ctx context.Context, customerID, id string) (*domain.Dossier, error)
	Create(ctx context.Context, name, patientID string) (*domain.Dossier, error)
}
