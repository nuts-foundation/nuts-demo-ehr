package dossier

import (
	"context"
	"testing"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	"github.com/nuts-foundation/nuts-demo-ehr/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestNewSQLiteDossierRepository(t *testing.T) {
	t.Run("create database", func(t *testing.T) {
		db := sqlx.MustConnect("sqlite3", ":memory:")
		_ = NewSQLiteDossierRepository(Factory{}, db)
		assert.NoError(t, db.Ping())
	})
}

func TestSQLiteDossierRepository_Create(t *testing.T) {
	db := sqlx.MustConnect("sqlite3", ":memory:")
	repo := NewSQLiteDossierRepository(Factory{}, db)

	var newDossier *types.Dossier
	var err error
	err = sql.ExecuteTransactional(db, func(ctx context.Context) error {
		newDossier, err = repo.Create(ctx, 1, "Broken leg", "p1")
		return err
	})

	if !assert.NoError(t, err) || !assert.NotNil(t, newDossier) {
		return
	}
	assert.NotEmpty(t, newDossier.Id)

	query := "SELECT * FROM `dossier` WHERE customer_id = ? ORDER BY id ASC"
	rows, err := db.Queryx(query, 1)
	if !assert.NoError(t, err) {
		return
	}
	dbDossier := &sqlDossier{}
	rows.Next()
	if !assert.NoError(t, rows.StructScan(dbDossier)) {
		return
	}
	assert.Equal(t, string(newDossier.Id), dbDossier.ID)
	assert.Equal(t, "Broken leg", dbDossier.Name)
	assert.Equal(t, "p1", dbDossier.PatientID)
	assert.Equal(t, 1, dbDossier.CustomerID)
}
