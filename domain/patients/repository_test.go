package patients

import (
	"context"
	"testing"

	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/stretchr/testify/assert"
)

func TestMemoryPatientRepository_NewPatient(t *testing.T) {
	t.Run("New patient", func(t *testing.T) {
		repo := NewMemoryPatientRepository(Factory{})
		pProps := domain.PatientProperties{
			FirstName: "henk",
			Gender:    "unknown",
		}
		newPatient, err := repo.NewPatient(context.Background(), "c1", pProps)
		assert.NoError(t, err)
		assert.Len(t, repo.patients, 1)
		expectedPatient := domain.Patient{
			ObjectID:         newPatient.ObjectID,
			PatientProperties: pProps,
		}
		assert.Equal(t, expectedPatient, repo.patients["c1"][newPatient.ObjectID])
	})
}

func TestMemoryPatientRepository_Update(t *testing.T) {
	repo := NewMemoryPatientRepository(Factory{})
	patient := domain.Patient{
		ObjectID: "p1",
		PatientProperties: domain.PatientProperties{
			FirstName: "henk",
		},
	}
	repo.patients["c1"] = map[domain.ObjectID]domain.Patient{"p1": patient}

	_, err := repo.Update(context.Background(), "c1", "p1", func(c domain.Patient) (*domain.Patient, error) {
		c.FirstName = "Peter"
		return &c, nil
	})
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "Peter", repo.patients["c1"]["p1"].FirstName)
}
