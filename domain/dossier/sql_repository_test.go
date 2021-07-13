package dossier

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestNewSQLiteDossierRepository(t *testing.T) {
	t.Run("create database", func(t *testing.T) {
		db := sqlx.MustConnect("sqlite3", ":memory:")
		repo := NewSQLiteDossierRepository(Factory{}, db)
		assert.NoError(t, repo.db.Ping())
	})
}
