package dossier

import (
	"context"

	"github.com/google/uuid"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

type Repository interface {
	FindByID(ctx context.Context, customerID, id string) (*domain.Dossier, error)
	Create(ctx context.Context, customerID ,name, patientID string) (*domain.Dossier, error)
}

type Factory struct{}

func (Factory) NewDossier(patientID, name string) *domain.Dossier {
	return &domain.Dossier{
		Id:        domain.ObjectID(uuid.NewString()),
		Name:      name,
		PatientID: domain.ObjectID(patientID),
	}
}
