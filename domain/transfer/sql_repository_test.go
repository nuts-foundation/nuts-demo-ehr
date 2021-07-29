package transfer

import (
	"context"
	"github.com/nuts-foundation/nuts-demo-ehr/sql"
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
		_ = NewSQLiteTransferRepository(Factory{}, db)
		assert.NoError(t, db.Ping())
	})
}

func TestSQLiteTransferRepository_Create(t *testing.T) {
	t.Run("create new transfer", func(t *testing.T) {
		db := sqlx.MustConnect("sqlite3", ":memory:")
		repo := NewSQLiteTransferRepository(Factory{}, db)
		now := time.Now().UTC().Round(time.Minute)

		var newTransfer *domain.Transfer
		var err error
		sql.ExecuteTransactional(db, func(ctx context.Context) error {
			newTransfer, err = repo.Create(ctx, "c1", "14", "foo", now)
			return err
		})

		if !assert.NoError(t, err) || !assert.NotNil(t, newTransfer) {
			return
		}
		assert.NotEmpty(t, newTransfer.Id)
		query := "SELECT * FROM `transfer` WHERE customer_id = ? ORDER BY id ASC"
		rows, err := db.Queryx(query, "c1")
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
		db.MustExec("INSERT INTO transfer (`id`, `customer_id`, `status`, `description`, `dossier_id`) VALUES('123', 'c1', 'created', 'the description', 'd1')")
		var transfers []domain.Transfer
		var err error
		sql.ExecuteTransactional(db, func(ctx context.Context) error {
			transfers, err = repo.FindByPatientID(ctx, "c1", "p1")
			return err
		})

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
		sql.ExecuteTransactional(db, func(ctx context.Context) error {
			transfer, _ := repo.Create(ctx, "c1", "14", "foo", now)
			negotiation, err := repo.CreateNegotiation(ctx, "c1", string(transfer.Id), "foo", now, "")
			assert.NoError(t, err)
			assert.NotNil(t, negotiation)
			return nil
		})
	})
}
