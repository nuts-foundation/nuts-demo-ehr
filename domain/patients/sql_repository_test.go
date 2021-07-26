package patients

import (
	"context"
	"github.com/nuts-foundation/nuts-demo-ehr/sql"
	"testing"
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/stretchr/testify/assert"
)

func TestSQLitePatientRepository_FindByID(t *testing.T) {
	t.Run("no results", func(t *testing.T) {
		db := sqlx.MustConnect("sqlite3", ":memory:")
		repo := NewSQLitePatientRepository(Factory{}, db)

		sql.ExecuteTransactional(db, func(ctx context.Context) error {
			result, err := repo.FindByID(ctx, "c1", "p1")
			assert.NoError(t, err)
			assert.Nil(t, result)
			return nil
		})
	})

	t.Run("1 result", func(t *testing.T) {
		db := sqlx.MustConnect("sqlite3", ":memory:")
		repo := NewSQLitePatientRepository(Factory{}, db)
		db.MustExec("INSERT INTO `patient` (`ssn`, `customer_id`, `id`, `first_name`, `surname`) VALUES('123','c1', 'p1', 'Henk', 'de Vries')")
		db.MustExec("INSERT INTO `patient` (`ssn`, `customer_id`, `id`, `first_name`, `surname`) VALUES('456','c2', 'p1', 'Floris-Jan', 'van Kleppensteyn')")
		sql.ExecuteTransactional(db, func(ctx context.Context) error {
			result, err := repo.FindByID(ctx, "c1", "p1")
			if !assert.NoError(t, err) {
				return nil
			}
			assert.Equal(t, "Henk", result.FirstName)
			return nil
		})
	})
}

func TestSQLitePatientRepository_All(t *testing.T) {
	t.Run("all patient", func(t *testing.T) {
		db := sqlx.MustConnect("sqlite3", ":memory:")
		repo := NewSQLitePatientRepository(Factory{}, db)
		db.MustExec("INSERT INTO `patient` (`ssn`,`customer_id`, `id`, `first_name`, `surname`) VALUES('123','c1', 'p1', 'Fred', 'Klooydonk')")
		db.MustExec("INSERT INTO `patient` (`ssn`,`customer_id`, `id`, `first_name`, `surname`) VALUES('456','c1', 'p2', 'Arie', 'de Eiker')")
		sql.ExecuteTransactional(db, func(ctx context.Context) error {
			result, err := repo.All(ctx, "c1", nil)
			if !assert.NoError(t, err) {
				return nil
			}
			assert.Len(t, result, 2)
			assert.Equal(t, "Fred", result[0].FirstName)
			assert.Equal(t, "Arie", result[1].FirstName)
			return nil
		})
	})
}

func TestSQLitePatientRepository_NewPatient(t *testing.T) {
	t.Run("new patient", func(t *testing.T) {
		db := sqlx.MustConnect("sqlite3", ":memory:")
		repo := NewSQLitePatientRepository(Factory{}, db)
		email := openapi_types.Email("foo@bar.com")
		ssn := "99999909"
		sql.ExecuteTransactional(db, func(ctx context.Context) error {
			newPatient, err := repo.NewPatient(ctx, "c15", domain.PatientProperties{
				Dob:       &openapi_types.Date{Time: time.Now().UTC().Round(time.Minute)},
				Email:     &email,
				FirstName: "Henk",
				Surname:   "de Vries",
				Gender:    domain.PatientPropertiesGenderMale,
				Ssn:       &ssn,
				Zipcode:   "7551AB",
			})
			if !assert.NoError(t, err) || !assert.NotNil(t, newPatient) {
				return nil
			}
			assert.NotEmpty(t, newPatient.ObjectID)

			foundPatient, err := repo.FindByID(ctx, "c15", string(newPatient.ObjectID))
			if !assert.NoError(t, err) || !assert.NotNil(t, newPatient) {
				return nil
			}
			assert.Equal(t, newPatient, foundPatient)
			return nil
		})
	})
}
