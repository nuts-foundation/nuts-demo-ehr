package patients

import (
	"context"
	"testing"

	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/stretchr/testify/assert"
)

func TestMemoryPatientRepository_NewPatient(t *testing.T) {
	t.Run("New patient", func(t *testing.T) {
		repo := NewMemoryPatientRepository()
		fistName := "henk"
		patient := domain.Patient{
			PatientID: "p1",
			PatientProperties: domain.PatientProperties{
				FirstName: &fistName,
			},
		}
		_, err := repo.NewPatient(context.Background(), "c1", patient)
		assert.NoError(t, err)
		assert.Len(t, repo.patients, 1)
		assert.Equal(t, patient, repo.patients["c1"]["p1"])
	})
}

func TestMemoryPatientRepository_Update(t *testing.T) {
	repo := NewMemoryPatientRepository()
	fistName := "henk"
	patient := domain.Patient{
		PatientID: "p1",
		PatientProperties: domain.PatientProperties{
			FirstName: &fistName,
		},
	}
	repo.patients["c1"] = map[domain.PatientID]domain.Patient{ "p1": patient}

	_, err := repo.Update(context.Background(), "c1", "p1", func(c domain.Patient) (*domain.Patient, error) {
		newName := "Peter"
		c.FirstName = &newName
		return &c, nil
	})
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "Peter", *repo.patients["c1"]["p1"].FirstName)
}
