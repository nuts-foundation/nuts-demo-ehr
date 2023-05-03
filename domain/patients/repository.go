package patients

import (
	"context"
	"github.com/google/uuid"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
)

type Repository interface {
	FindByID(ctx context.Context, customerID int, id string) (*types.Patient, error)
	Update(ctx context.Context, customerID int, id string, updateFn func(c types.Patient) (*types.Patient, error)) (*types.Patient, error)
	NewPatient(ctx context.Context, customerID int, patient types.PatientProperties) (*types.Patient, error)
	All(ctx context.Context, customerID int, name *string) ([]types.Patient, error)
}

type Factory struct{}

// NewPatient creates a new patient from a list of properties. It generates a new UUID for the patientID.
func (f Factory) NewPatient(properties types.PatientProperties) (*types.Patient, error) {
	if properties.Gender == "" {
		properties.Gender = types.Unknown
	}
	return &types.Patient{
		BaseProps: types.BaseProps{
			ObjectID: uuid.NewString(),
		},
		PatientProperties: properties,
	}, nil
}
