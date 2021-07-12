package transfer

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

type sqlTransfer struct {
	ID          string       `db:"id"`
	CustomerID  string       `db:"customer_id"`
	Date        sql.NullTime `db:"date"`
	Status      string       `db:"status"`
	Description string       `db:"description"`
}

func (dbTransfer *sqlTransfer) UnmarshalFromDomainTransfer(customerID string, transfer domain.Transfer) error {
	*dbTransfer = sqlTransfer{
		ID:          string(transfer.Id),
		Date:        sql.NullTime{Time: transfer.TransferDate.Time, Valid: !transfer.TransferDate.IsZero()},
		CustomerID:  customerID,
		Status:      string(transfer.Status),
		Description: transfer.Description,
	}
	return nil
}

const schema = `
	CREATE TABLE IF NOT EXISTS transfer (
		id char(36) NOT NULL,
		customer_id varchar(100) NOT NULL,
		date DATETIME DEFAULT NULL,
		status char(10) NOT NULL DEFAULT 'created',
		description varchar(200) NOT NULL,
		PRIMARY KEY (id),
		UNIQUE(customer_id, id)
	);
`

type SQLiteTransferRepository struct {
	factory Factory
	db      *sqlx.DB
}

func NewSQLiteTransferRepository(factory Factory, db *sqlx.DB) *SQLiteTransferRepository {
	if db == nil {
		panic("missing db")
	}
	db.MustExec(schema)

	return &SQLiteTransferRepository{
		factory: factory,
		db:      db,
	}
}

func (r SQLiteTransferRepository) FindByID(ctx context.Context, customerID, id string) (*domain.Transfer, error) {
	panic("implement me")
}

func (r SQLiteTransferRepository) FindByPatientID(ctx context.Context, customerID, patientID string) []domain.Transfer {
	panic("implement me")
}

func (r SQLiteTransferRepository) Create(ctx context.Context, customerID, dossierID, description string, date time.Time) (*domain.Transfer, error) {
	transfer := r.factory.NewTransfer(description, date)
	dbTransfer := sqlTransfer{}
	if err := dbTransfer.UnmarshalFromDomainTransfer(customerID, *transfer); err != nil {
		return nil, err
	}
	const query = `INSERT INTO transfer 
		(id, customer_id, date, status, description)
		values(:id, :customer_id, :date, :status, :description)
`

	if _, err := r.db.NamedExec(query, dbTransfer); err != nil {
		return nil, err
	}
	return transfer, nil
}

func (r SQLiteTransferRepository) Update(ctx context.Context, customerID, description string, date time.Time, state domain.TransferStatus) error {
	panic("implement me")
}

func (r SQLiteTransferRepository) Cancel(ctx context.Context, customerID, id string) {
	panic("implement me")
}

func (r SQLiteTransferRepository) CreateNegotiation(ctx context.Context, transferID string, organizationDID string) (*domain.TransferNegotiation, error) {
	panic("implement me")
}

func (r SQLiteTransferRepository) ListNegotiations(ctx context.Context, customerID, transferID string) ([]domain.TransferNegotiation, error) {
	panic("implement me")
}
