package transfer

import (
	"context"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestNewSQLiteTransferRepository(t *testing.T) {
	t.Run("create database", func(t *testing.T) {
		db := sqlx.MustConnect("sqlite3", ":memory:")
		repo := NewSQLiteTransferRepository(Factory{}, db)
		assert.NoError(t, repo.db.Ping())
	})
}

func TestSQLiteTransferRepository_Create(t *testing.T) {
	t.Run("create new transfer", func(t *testing.T) {
		db := sqlx.MustConnect("sqlite3", ":memory:")
		repo := NewSQLiteTransferRepository(Factory{}, db)
		now := time.Now().UTC().Round(time.Minute)
		newTransfer, err := repo.Create(context.Background(), "c1", "14", "foo", now)

		if !assert.NoError(t, err) || !assert.NotNil(t, newTransfer) {
			return
		}
		assert.NotEmpty(t, newTransfer.Id)
		query := "SELECT * FROM `transfer` WHERE customer_id = ? ORDER BY id ASC"
		rows, err := repo.db.Queryx(query, "c1", "14")
		if !assert.NoError(t, err) {
			return
		}
		dbTransfer := &sqlTransfer{}
		rows.Next()
		if !assert.NoError(t, rows.StructScan(dbTransfer)) {
			return
		}
		assert.Equal(t, string(newTransfer.Id), dbTransfer.ID)
		assert.Equal(t, newTransfer.TransferDate.Time, dbTransfer.Date.Time)
		assert.Equal(t, newTransfer.Description, dbTransfer.Description)
	})
}
