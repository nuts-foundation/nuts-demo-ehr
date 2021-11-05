package dossier

import (
	"context"

	"github.com/google/uuid"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
)

type Repository interface {
	FindByID(ctx context.Context, customerID int, id string) (*types.Dossier, error)
	Create(ctx context.Context, customerID int, name, patientID string) (*types.Dossier, error)
	AllByPatient(ctx context.Context, customerID int, patientID string) ([]types.Dossier, error)
}

type Factory struct{}

func (Factory) NewDossier(patientID, name string) *types.Dossier {
	return &types.Dossier{
		Id:        types.ObjectID(uuid.NewString()),
		Name:      name,
		PatientID: types.ObjectID(patientID),
	}
}
