package transfer

import (
	"context"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
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
		rows, err := repo.db.Queryx(query, "c1")
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

func TestSQLiteTransferRepository_FindByPatientID(t *testing.T) {
	t.Run("simple select", func(t *testing.T) {
		db := sqlx.MustConnect("sqlite3", ":memory:")
		repo := NewSQLiteTransferRepository(Factory{}, db)
		repo.db.MustExec("INSERT INTO transfer (`id`, `customer_id`, `status`, `description`, `dossier_id`) VALUES('123', 'c1', 'created', 'the description', 'd1')")
		transfers, err := repo.FindByPatientID(context.Background(), "c1", "p1")

		if !assert.NoError(t, err) {
			return
		}
		assert.Len(t, transfers, 1)
		assert.Equal(t, "the description", transfers[0].Description)
		assert.Equal(t, domain.TransferStatusCreated, transfers[0].Status)
	})
}

func TestSQLiteTransferRepository_CreateNegotiation(t *testing.T) {
	t.Run("create new negotiation", func(t *testing.T) {
		db := sqlx.MustConnect("sqlite3", ":memory:")
		repo := NewSQLiteTransferRepository(Factory{}, db)
		now := time.Now().UTC().Round(time.Minute)
		transfer, _ := repo.Create(context.Background(), "c1", "14", "foo", now)
		negotiation, err := repo.CreateNegotiation(context.Background(), "c1", string(transfer.Id), "foo", now)
		assert.NoError(t, err)
		assert.NotNil(t, negotiation)
	})
}
