package dossier

import (
	"context"

	"github.com/google/uuid"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
)

type Repository interface {
	FindByID(ctx context.Context, customerID string, id string) (*types.Dossier, error)
	Create(ctx context.Context, customerID string, name, patientID string) (*types.Dossier, error)
	AllByPatient(ctx context.Context, customerID string, patientID string) ([]types.Dossier, error)
}

type Factory struct{}

func (Factory) NewDossier(patientID, name string) *types.Dossier {
	return &types.Dossier{
		Id:        types.ObjectID(uuid.NewString()),
		Name:      name,
		PatientID: types.ObjectID(patientID),
	}
}
