package transfer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/jmoiron/sqlx"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

type sqlTransfer struct {
	ID          string       `db:"id"`
	CustomerID  string       `db:"customer_id"`
	Date        sql.NullTime `db:"date"`
	Status      string       `db:"status"`
	DossierID   string       `db:"dossier_id"`
	Description string       `db:"description"`
}

func (dbTransfer *sqlTransfer) UnmarshalFromDomainTransfer(customerID, dossierID string, transfer domain.Transfer) error {
	*dbTransfer = sqlTransfer{
		ID:          string(transfer.Id),
		CustomerID:  customerID,
		DossierID:   dossierID,
		Date:        sql.NullTime{Time: transfer.TransferDate.Time, Valid: !transfer.TransferDate.IsZero()},
		Status:      string(transfer.Status),
		Description: transfer.Description,
	}
	return nil
}

func (dbTransfer sqlTransfer) MarshalToDomainTransfer() (*domain.Transfer, error) {
	var status domain.TransferStatus
	switch dbTransfer.Status {
	case string(domain.TransferStatusCreated):
		status = domain.TransferStatusCreated
	case string(domain.TransferStatusAssigned):
		status = domain.TransferStatusAssigned
	case string(domain.TransferStatusRequested):
		status = domain.TransferStatusRequested
	case string(domain.TransferStatusCompleted):
		status = domain.TransferStatusCompleted
	case string(domain.TransferStatusCancelled):
		status = domain.TransferStatusCancelled
	default:
		return nil, fmt.Errorf("unknown tranfser status: '%s'", dbTransfer.Status)
	}

	return &domain.Transfer{
		Id:           domain.ObjectID(dbTransfer.ID),
		Description:  dbTransfer.Description,
		Status:       status,
		TransferDate: openapi_types.Date{},
	}, nil
}

const schema = `
	CREATE TABLE IF NOT EXISTS transfer (
		id char(36) NOT NULL,
		customer_id varchar(100) NOT NULL,
		date DATETIME DEFAULT NULL,
		status char(10) NOT NULL DEFAULT 'created',
		description varchar(200) NOT NULL,
	    dossier_id char(36) NOT NULL,
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
	// TODO: filter on patient by dossier
	const query = `SELECT * FROM transfer WHERE customer_id = ? AND id = ? ORDER BY id ASC`

	dbTransfer := sqlTransfer{}
	err := r.db.GetContext(ctx, &dbTransfer, query, customerID, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return dbTransfer.MarshalToDomainTransfer()
}

func (r SQLiteTransferRepository) FindByPatientID(ctx context.Context, customerID, patientID string) ([]domain.Transfer, error) {
	// TODO: filter on patient by dossier
	const query = `SELECT * FROM transfer WHERE customer_id = ? ORDER BY id ASC`

	dbTransfers := []sqlTransfer{}
	err := r.db.SelectContext(ctx, &dbTransfers, query, customerID)
	if errors.Is(err, sql.ErrNoRows) {
		return []domain.Transfer{}, nil
	} else if err != nil {
		return nil, err
	}

	result := make([]domain.Transfer, len(dbTransfers))

	for idx, dbTransfer := range dbTransfers {
		patient, err := dbTransfer.MarshalToDomainTransfer()
		if err != nil {
			return nil, err
		}
		result[idx] = *patient
	}
	return result, nil
}

func (r SQLiteTransferRepository) Create(ctx context.Context, customerID, dossierID, description string, date time.Time) (*domain.Transfer, error) {
	transfer := r.factory.NewTransfer(description, date)
	dbTransfer := sqlTransfer{}
	if err := dbTransfer.UnmarshalFromDomainTransfer(customerID, dossierID, *transfer); err != nil {
		return nil, err
	}
	const query = `INSERT INTO transfer 
		(id, customer_id, dossier_id, date, status, description)
		values(:id, :customer_id, :dossier_id, :date, :status, :description)
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
