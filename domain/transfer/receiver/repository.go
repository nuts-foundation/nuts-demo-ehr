package receiver

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	sqlUtil "github.com/nuts-foundation/nuts-demo-ehr/sql"
	"time"
)

const transferSchema = `
	CREATE TABLE IF NOT EXISTS incoming_transfers (
		id char(36) NOT NULL,
		status VARCHAR(100) CHECK (status IN (
		    'accepted',
			'cancelled',
		    'completed',
		    'in-progress',
		    'on-hold',
		    'requested'
		)) NOT NULL DEFAULT 'requested',
	    task_id VARCHAR(100) NOT NULL,
		customer_id VARCHAR(100) NOT NULL,
		sender_did VARCHAR(100) NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		PRIMARY KEY (id),
		
		CONSTRAINT idx_task_id UNIQUE (task_id)
	);
`

type TransferRepository interface {
	GetCount(ctx context.Context, customerID int) (int, error)
	GetAll(ctx context.Context, customerID int) ([]domain.IncomingTransfer, error)
	CreateOrUpdate(ctx context.Context, status, taskID string, customerID int, senderDID string) (*domain.IncomingTransfer, error)
}

func NewTransferRepository(db *sqlx.DB) TransferRepository {
	tx, _ := db.Beginx()
	tx.MustExec(transferSchema)

	if err := tx.Commit(); err != nil {
		panic(err)
	}

	return repository{db: db}
}

type sqlTransfer struct {
	ID         string    `db:"id"`
	TaskID     string    `db:"task_id"`
	Status     string    `db:"status"`
	CustomerID int       `db:"customer_id"`
	SenderDID  string    `db:"sender_did"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func (transfer sqlTransfer) marshalToDomain() domain.IncomingTransfer {
	return domain.IncomingTransfer{
		Id:         domain.ObjectID(transfer.ID),
		FhirTaskID: transfer.TaskID,
		// @TODO: should we resolve organization details here?
		Sender: domain.Organization{
			Did: transfer.SenderDID,
		},
		CreatedAt: transfer.CreatedAt,
		Status:    domain.TransferNegotiationStatus{Status: domain.TransferNegotiationStatusStatus(transfer.Status)},
	}
}

type repository struct {
	db *sqlx.DB
}

func (f repository) GetAll(ctx context.Context, customerID int) ([]domain.IncomingTransfer, error) {
	const query = `SELECT * FROM incoming_transfers WHERE customer_id = ? ORDER BY updated_at DESC`

	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return nil, err
	}

	var transfers []sqlTransfer

	if err := tx.SelectContext(ctx, &transfers, query, customerID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []domain.IncomingTransfer{}, nil
		}

		return nil, err
	}

	var results []domain.IncomingTransfer

	for _, transfer := range transfers {
		results = append(results, transfer.marshalToDomain())
	}

	return results, nil
}

func (f repository) GetCount(ctx context.Context, customerID int) (int, error) {
	const query = `SELECT COUNT(*) FROM incoming_transfers WHERE customer_id = ?`

	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return 0, err
	}

	var count int

	if err := tx.GetContext(ctx, &count, query, customerID); err != nil {
		return 0, err
	}

	return count, nil
}

func (f repository) CreateOrUpdate(ctx context.Context, status, taskID string, customerID int, senderDID string) (*domain.IncomingTransfer, error) {
	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return nil, err
	}

	const query = `INSERT INTO incoming_transfers (id, created_at, status, task_id, customer_id, sender_did, updated_at)
		VALUES(:id, :created_at, :status, :task_id, :customer_id, :sender_did, :updated_at)
		ON CONFLICT(task_id) DO
		UPDATE SET updated_at = :updated_at`

	transfer := &sqlTransfer{
		ID:         uuid.New().String(),
		Status:     status,
		TaskID:     taskID,
		CustomerID: customerID,
		SenderDID:  senderDID,
		UpdatedAt:  time.Now(),
		CreatedAt:  time.Now(),
	}

	_, err = tx.NamedExecContext(ctx, query, transfer)
	if err != nil {
		return nil, err
	}

	incomingTransfer := transfer.marshalToDomain()

	return &incomingTransfer, nil
}
