package patients

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestSQLitePatientRepository_FindByID(t *testing.T) {
	t.Run("no results", func(t *testing.T) {
		db := sqlx.MustConnect("sqlite3", ":memory:")
		repo := NewSQLitePatientRepository(db)
		result, err := repo.FindByID(context.Background(), "c1", "p1")
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("1 result", func(t *testing.T) {
		db := sqlx.MustConnect("sqlite3", ":memory:")
		repo := NewSQLitePatientRepository(db)
		db.MustExec("INSERT INTO `patient` (`customer_id`, `id`, `first_name`, `internal_id`) VALUES('c1', 'p1', 'Henk', 'c1-patient-1')")
		db.MustExec("INSERT INTO `patient` (`customer_id`, `id`, `first_name`, `internal_id`) VALUES('c2', 'p1', 'Peter', 'c2-patient-1')")
		result, err := repo.FindByID(context.Background(), "c1", "p1")
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t,"Henk", *result.FirstName)
	})
}

func TestSQLitePatientRepository_All(t *testing.T) {
	t.Run("all patient", func(t *testing.T) {
		db := sqlx.MustConnect("sqlite3", ":memory:")
		repo := NewSQLitePatientRepository(db)
		db.MustExec("INSERT INTO `patient` (`customer_id`, `id`, `first_name`, `internal_id`) VALUES('c1', 'p1', 'Henk', 'c1-patient-1')")
		db.MustExec("INSERT INTO `patient` (`customer_id`, `id`, `first_name`, `internal_id`) VALUES('c1', 'p2', 'Peter', 'c1-patient-2')")
		result, err := repo.All(context.Background(), "c1")
		if !assert.NoError(t, err) {
			return
		}
		assert.Len(t, result, 2)
		assert.Equal(t,"Henk", *result[0].FirstName)
		assert.Equal(t,"Peter", *result[1].FirstName)

	})
}