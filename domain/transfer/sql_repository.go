package transfer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sql2 "github.com/nuts-foundation/nuts-demo-ehr/domain/sql"
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

type sqlNegotiation struct {
	TransferID      string    `db:"transfer_id"`
	OrganizationDID string    `db:"organization_did"`
	CustomerID      string    `db:"customer_id"`
	Date            time.Time `db:"date"`
	Status          string    `db:"status"`
}

func (n sqlNegotiation) MarshalToDomainTransfer() (*domain.TransferNegotiation, error) {
	return &domain.TransferNegotiation{
		OrganizationDID: n.OrganizationDID,
		Status:          domain.TransferNegotiationStatus(n.Status),
		TransferDate:    openapi_types.Date{Time: n.Date},
	}, nil
}

func (dbTransfer *sqlTransfer) UnmarshalFromDomainTransfer(customerID string, transfer domain.Transfer) error {
	*dbTransfer = sqlTransfer{
		ID:          string(transfer.Id),
		CustomerID:  customerID,
		DossierID:   string(transfer.DossierID),
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
		DossierID:    domain.ObjectID(dbTransfer.DossierID),
		Description:  dbTransfer.Description,
		Status:       status,
		TransferDate: openapi_types.Date{},
	}, nil
}

const transferSchema = `
	CREATE TABLE transfer (
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

const negotiationSchema = `
	CREATE TABLE transfer_negotiation (
		organization_did varchar(200) NOT NULL,
		transfer_id char(36) NOT NULL,
		customer_id varchar(100) NOT NULL,
		date DATETIME DEFAULT NULL,
		status char(10) NOT NULL DEFAULT 'requested',
		PRIMARY KEY (organization_did, transfer_id),
		FOREIGN KEY (transfer_id) REFERENCES transfer(id)
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
	db.MustExec(transferSchema)
	db.MustExec(negotiationSchema)

	return &SQLiteTransferRepository{
		factory: factory,
		db:      db,
	}
}

func (r SQLiteTransferRepository) findByID(ctx context.Context, db sql2.SQLContextGetter, customerID, id string) (*domain.Transfer, error) {
	// TODO: filter on patient by dossier
	const query = `SELECT * FROM transfer WHERE customer_id = ? AND id = ? ORDER BY id ASC`

	dbTransfer := sqlTransfer{}
	err := db.GetContext(ctx, &dbTransfer, query, customerID, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return dbTransfer.MarshalToDomainTransfer()
}

func (r SQLiteTransferRepository) FindByID(ctx context.Context, customerID, id string) (*domain.Transfer, error) {
	return r.findByID(ctx, r.db, customerID, id)
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
	transfer := r.factory.NewTransfer(description, date, domain.ObjectID(dossierID))
	dbTransfer := sqlTransfer{}
	if err := dbTransfer.UnmarshalFromDomainTransfer(customerID, *transfer); err != nil {
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

func (r SQLiteTransferRepository) Update(ctx context.Context, customerID, transferID string, updateFn func(c domain.Transfer) (*domain.Transfer, error)) (entity *domain.Transfer, err error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("%w, unable to start transaction", err)
	}

	defer func() {
		if err == nil {
			err = tx.Commit()
		} else {
			tx.Rollback()
			entity = nil
		}
		if err != nil {
			entity = nil
		}
	}()

	entity, err = r.findByID(ctx, tx, customerID, transferID)
	if err != nil {
		return
	}
	updated, err := updateFn(*entity)
	if err != nil {
		return
	}

	dbEntity := sqlTransfer{}
	err = dbEntity.UnmarshalFromDomainTransfer(customerID, *updated)
	if err != nil {
		return
	}

	const query = `
	UPDATE transfer SET
		date = :date,
		status = :status,
		description = :description
	WHERE customer_id = :customer_id AND id = :id
`
	_, err = tx.NamedExec(query, dbEntity)

	return
}

func (r SQLiteTransferRepository) Cancel(ctx context.Context, customerID, id string) {
	panic("implement me")
}

func (r SQLiteTransferRepository) CreateNegotiation(ctx context.Context, customerID, transferID, organizationDID string, date time.Time) (*domain.TransferNegotiation, error) {
	negotiation := sqlNegotiation{
		TransferID:      transferID,
		OrganizationDID: organizationDID,
		CustomerID:      customerID,
		Date:            date,
	}
	const query = `INSERT INTO transfer_negotiation 
		(transfer_id, organization_did, customer_id, date)
		values(:transfer_id, :organization_did, :customer_id, :date)
`

	if _, err := r.db.NamedExec(query, negotiation); err != nil {
		return nil, err
	}
	return negotiation.MarshalToDomainTransfer()
}

func (r SQLiteTransferRepository) ListNegotiations(ctx context.Context, customerID, transferID string) ([]domain.TransferNegotiation, error) {
	const query = `SELECT * FROM transfer_negotiation WHERE customer_id = ? AND transfer_id = ? ORDER BY organization_did ASC`

	dbNegotiations := []sqlNegotiation{}
	err := r.db.SelectContext(ctx, &dbNegotiations, query, customerID, transferID)
	if errors.Is(err, sql.ErrNoRows) {
		return []domain.TransferNegotiation{}, nil
	} else if err != nil {
		return nil, err
	}

	result := make([]domain.TransferNegotiation, len(dbNegotiations))

	for idx, dbNegotiation := range dbNegotiations {
		item, err := dbNegotiation.MarshalToDomainTransfer()
		if err != nil {
			return nil, err
		}
		result[idx] = *item
	}

	return result, nil
}
