package sender

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	transfer2 "github.com/nuts-foundation/nuts-demo-ehr/domain/transfer"
	"time"

	sqlUtil "github.com/nuts-foundation/nuts-demo-ehr/sql"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/jmoiron/sqlx"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

type sqlTransfer struct {
	ID                            string         `db:"id"`
	CustomerID                    int            `db:"customer_id"`
	Date                          sql.NullTime   `db:"date"`
	Status                        string         `db:"status"`
	DossierID                     string         `db:"dossier_id"`
	Description                   string         `db:"description"`
	FHIRAdvanceNoticeComposition  string         `db:"fhir_advancenotice_composition"`
	FHIRNursingHandoffComposition sql.NullString `db:"fhir_nursinghandoff_composition"`
}

type sqlNegotiation struct {
	ID              string    `db:"id"`
	TransferID      string    `db:"transfer_id"`
	OrganizationDID string    `db:"organization_did"`
	CustomerID      int       `db:"customer_id"`
	Date            time.Time `db:"date"`
	Status          string    `db:"status"`
	TaskID          string    `db:"task_id"`
}

func (dbNegotiation sqlNegotiation) MarshalToDomainNegotiation() (*domain.TransferNegotiation, error) {
	return &domain.TransferNegotiation{
		Id:                        domain.ObjectID(dbNegotiation.ID),
		OrganizationDID:           dbNegotiation.OrganizationDID,
		TransferNegotiationStatus: domain.TransferNegotiationStatus{Status: domain.TransferNegotiationStatusStatus(dbNegotiation.Status)},
		TransferDate:              openapi_types.Date{Time: dbNegotiation.Date},
		TransferID:                domain.ObjectID(dbNegotiation.TransferID),
		TaskID:                    dbNegotiation.TaskID,
	}, nil
}

func (dbNegotiation *sqlNegotiation) UnmarshalFromDomainNegotiation(customerID int, negotiation domain.TransferNegotiation) error {
	*dbNegotiation = sqlNegotiation{
		ID:              string(negotiation.Id),
		TransferID:      string(negotiation.TransferID),
		OrganizationDID: negotiation.OrganizationDID,
		CustomerID:      customerID,
		Date:            negotiation.TransferDate.Time,
		Status:          string(negotiation.Status),
		TaskID:          negotiation.TaskID,
	}
	return nil
}

func (dbTransfer *sqlTransfer) UnmarshalFromDomainTransfer(customerID int, transfer domain.Transfer) error {
	*dbTransfer = sqlTransfer{
		ID:                            string(transfer.Id),
		CustomerID:                    customerID,
		DossierID:                     string(transfer.DossierID),
		Date:                          sql.NullTime{Time: transfer.TransferDate.Time, Valid: !transfer.TransferDate.IsZero()},
		Status:                        string(transfer.Status),
		Description:                   transfer.Description,
		FHIRAdvanceNoticeComposition:  transfer.FhirAdvanceNoticeComposition,
		FHIRNursingHandoffComposition: toNullString(transfer.FhirNursingHandoffComposition),
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

	var transferTime openapi_types.Date
	if dbTransfer.Date.Valid {
		transferTime = openapi_types.Date{Time: dbTransfer.Date.Time}
	}

	return &domain.Transfer{
		Id:        domain.ObjectID(dbTransfer.ID),
		DossierID: domain.ObjectID(dbTransfer.DossierID),
		Status:    status,
		TransferProperties: domain.TransferProperties{
			Description:  dbTransfer.Description,
			TransferDate: transferTime,
		},
		FhirAdvanceNoticeComposition:  dbTransfer.FHIRAdvanceNoticeComposition,
		FhirNursingHandoffComposition: fromNullString(dbTransfer.FHIRNursingHandoffComposition),
	}, nil
}

func fromNullString(input sql.NullString) *string {
	if input.Valid {
		return &input.String
	}
	return nil
}

func toNullString(input *string) sql.NullString {
	if input == nil {
		return sql.NullString{}
	}
	return sql.NullString{
		String: *input,
		Valid:  true,
	}
}

const transferSchema = `
	CREATE TABLE IF NOT EXISTS transfer (
		id char(36) NOT NULL,
		customer_id integer(11) NOT NULL,
		date DATETIME DEFAULT NULL,
		status char(10) NOT NULL DEFAULT 'created',
		description varchar(200) NOT NULL,
	    dossier_id char(36) NOT NULL,
	    fhir_advancenotice_composition VARCHAR(100) NOT NULL,
	    fhir_nursinghandoff_composition VARCHAR(100) NULL,
		PRIMARY KEY (id),
		UNIQUE(customer_id, id)
	);
`

const negotiationSchema = `
	CREATE TABLE IF NOT EXISTS transfer_negotiation (
	    id char(36) NOT NULL,
		organization_did varchar(200) NOT NULL,
		transfer_id char(36) NOT NULL,
		customer_id integer(11) NOT NULL,
		date DATETIME DEFAULT NULL,
		status char(10) NOT NULL DEFAULT 'requested',
		task_id char(36) NOT NULL,
		PRIMARY KEY (id),
		FOREIGN KEY (transfer_id) REFERENCES transfer(id)
	);
`

type SQLiteTransferRepository struct {
}

func NewTransferRepository(db *sqlx.DB) *SQLiteTransferRepository {
	if db == nil {
		panic("missing db")
	}

	tx, _ := db.Beginx()
	tx.MustExec(transferSchema)
	tx.MustExec(negotiationSchema)
	if err := tx.Commit(); err != nil {
		panic(err)
	}

	return &SQLiteTransferRepository{}
}

func (r SQLiteTransferRepository) findByID(ctx context.Context, tx *sqlx.Tx, customerID int, id string) (*domain.Transfer, error) {
	// TODO: filter on patient by dossier
	const query = `SELECT * FROM transfer WHERE customer_id = ? AND id = ? ORDER BY id ASC`

	dbTransfer := sqlTransfer{}
	err := tx.GetContext(ctx, &dbTransfer, query, customerID, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return dbTransfer.MarshalToDomainTransfer()
}

func (r SQLiteTransferRepository) updateTransfer(ctx context.Context, tx *sqlx.Tx, customerID int, transfer domain.Transfer) error {
	const query = `
	UPDATE transfer SET
		date = :date,
		status = :status,
		description = :description
	WHERE customer_id = :customer_id AND id = :id
`
	dbEntity := sqlTransfer{}
	err := dbEntity.UnmarshalFromDomainTransfer(customerID, transfer)
	if err != nil {
		return err
	}
	_, err = tx.NamedExecContext(ctx, query, dbEntity)
	if err != nil {
		return fmt.Errorf("unable to update the Transfer: %w", err)
	}
	return nil
}

func (r SQLiteTransferRepository) FindNegotiationByID(ctx context.Context, customerID int, negotiationID string) (*domain.TransferNegotiation, error) {
	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return nil, err
	}

	return r.findNegotiationByID(ctx, tx, customerID, negotiationID)
}

func (r SQLiteTransferRepository) findNegotiationByID(ctx context.Context, tx *sqlx.Tx, customerID int, negotiationID string) (*domain.TransferNegotiation, error) {
	const query = `SELECT * FROM transfer_negotiation WHERE customer_id = ? AND id = ?`

	dbNegotiation := sqlNegotiation{}
	err := tx.GetContext(ctx, &dbNegotiation, query, customerID, negotiationID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("unable to find negotiation by id: %w", err)
	}

	return dbNegotiation.MarshalToDomainNegotiation()
}

func (r SQLiteTransferRepository) FindNegotiationByTaskID(ctx context.Context, customerID int, taskID string) (*domain.TransferNegotiation, error) {
	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return nil, err
	}
	const query = `SELECT * FROM transfer_negotiation WHERE customer_id = ? AND task_id = ?`

	dbNegotiation := sqlNegotiation{}
	err = tx.GetContext(ctx, &dbNegotiation, query, customerID, taskID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("unable to find negotiation by task_id: %w", err)
	}

	return dbNegotiation.MarshalToDomainNegotiation()
}

func (r SQLiteTransferRepository) FindByID(ctx context.Context, customerID int, id string) (*domain.Transfer, error) {
	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return nil, err
	}
	return r.findByID(ctx, tx, customerID, id)
}

func (r SQLiteTransferRepository) FindByPatientID(ctx context.Context, customerID int, patientID string) ([]domain.Transfer, error) {
	const query = `SELECT transfer.* FROM transfer, dossier WHERE transfer.customer_id = ? AND dossier.id == transfer.dossier_id AND dossier.patient_id = ? ORDER BY id ASC`
	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return nil, err
	}

	dbTransfers := []sqlTransfer{}
	err = tx.SelectContext(ctx, &dbTransfers, query, customerID, patientID)
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

func (r SQLiteTransferRepository) Create(ctx context.Context, customerID int, dossierID, description string, date time.Time, fhirAdvanceNoticeComposition string) (*domain.Transfer, error) {
	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return nil, err
	}
	transfer := &domain.Transfer{
		Id:                           domain.ObjectID(uuid.NewString()),
		DossierID:                    domain.ObjectID(dossierID),
		Status:                       domain.TransferStatusCreated,
		FhirAdvanceNoticeComposition: fhirAdvanceNoticeComposition,
		TransferProperties: domain.TransferProperties{
			Description:  description,
			TransferDate: openapi_types.Date{Time: date},
		},
	}
	dbTransfer := sqlTransfer{}
	if err := dbTransfer.UnmarshalFromDomainTransfer(customerID, *transfer); err != nil {
		return nil, err
	}
	const query = `INSERT INTO transfer 
		(id, customer_id, dossier_id, date, status, description, fhir_advancenotice_composition)
		values(:id, :customer_id, :dossier_id, :date, :status, :description, :fhir_advancenotice_composition)
`

	if _, err := tx.NamedExecContext(ctx, query, dbTransfer); err != nil {
		return nil, err
	}
	return transfer, nil
}

func (r SQLiteTransferRepository) Update(ctx context.Context, customerID int, transferID string, updateFn func(c *domain.Transfer) (*domain.Transfer, error)) (entity *domain.Transfer, err error) {
	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return nil, err
	}

	entity, err = r.findByID(ctx, tx, customerID, transferID)
	if err != nil {
		return
	}
	entity, err = updateFn(entity)
	if err != nil {
		return
	}

	err = r.updateTransfer(ctx, tx, customerID, *entity)

	return
}

func (r SQLiteTransferRepository) Cancel(ctx context.Context, customerID int, transferID string) (*domain.Transfer, error) {
	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return nil, err
	}

	transfer, err := r.findByID(ctx, tx, customerID, transferID)
	if err != nil {
		return nil, nil
	}
	transfer.Status = transfer2.CancelledState

	if err := r.updateTransfer(ctx, tx, customerID, *transfer); err != nil {
		return nil, err
	}

	negotiations, err := r.ListNegotiations(ctx, customerID, string(transfer.Id))
	if err != nil {
		return nil, err
	}
	for _, negotiation := range negotiations {
		negotiation.Status = transfer2.CancelledState
		if err := r.updateNegotiation(ctx, tx, customerID, negotiation); err != nil {
			return nil, err
		}
	}

	return transfer, nil
}

func (r SQLiteTransferRepository) UpdateNegotiationState(ctx context.Context, customerID int, negotiationID string, newState domain.TransferNegotiationStatusStatus) (*domain.TransferNegotiation, error) {
	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return nil, err
	}
	negotiation, err := r.findNegotiationByID(ctx, tx, customerID, negotiationID)
	if err != nil {
		return nil, err
	}
	wrongStateErrStr := "could not update status: invalid state transition from: %s to: %s"
	switch negotiation.Status {
	case transfer2.CancelledState:
	case transfer2.CompletedState:
		return nil, fmt.Errorf(wrongStateErrStr, negotiation.Status, newState)
	}
	negotiation.Status = newState
	if err = r.updateNegotiation(ctx, tx, customerID, *negotiation); err != nil {
		return nil, err
	}
	return negotiation, nil
}

func (r SQLiteTransferRepository) CancelNegotiation(ctx context.Context, customerID int, negotiationID string) (*domain.TransferNegotiation, error) {
	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return nil, err
	}
	negotiation, err := r.findNegotiationByID(ctx, tx, customerID, negotiationID)
	if err != nil {
		return nil, err
	}
	negotiation.Status = transfer2.CancelledState
	if err := r.updateNegotiation(ctx, tx, customerID, *negotiation); err != nil {
		return nil, err
	}
	return negotiation, nil
}

func (r SQLiteTransferRepository) ProposeAlternateDate(ctx context.Context, customerID int, negotiationID string) (*domain.TransferNegotiation, error) {
	panic("implement me")
}

func (r SQLiteTransferRepository) ConfirmNegotiation(ctx context.Context, customerID int, negotiationID string) (*domain.TransferNegotiation, error) {
	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return nil, err
	}
	negotiation, err := r.findNegotiationByID(ctx, tx, customerID, negotiationID)
	if err != nil {
		return nil, err
	}
	negotiation.Status = transfer2.InProgressState
	if err := r.updateNegotiation(ctx, tx, customerID, *negotiation); err != nil {
		return nil, err
	}
	return negotiation, nil
}

func (r SQLiteTransferRepository) updateNegotiation(ctx context.Context, tx *sqlx.Tx, customerID int, negotiation domain.TransferNegotiation) error {
	const query = `
	UPDATE transfer_negotiation SET
		date = :date,
		status = :status
	WHERE customer_id = :customer_id AND id = :id
`

	dbEntity := sqlNegotiation{}
	err := dbEntity.UnmarshalFromDomainNegotiation(customerID, negotiation)
	if err != nil {
		return err
	}
	_, err = tx.NamedExecContext(ctx, query, dbEntity)
	if err != nil {
		return fmt.Errorf("unable to update the negotiation: %w", err)
	}
	return nil
}

func (r SQLiteTransferRepository) CreateNegotiation(ctx context.Context, customerID int, transferID, organizationDID string, transferDate time.Time, taskID string) (*domain.TransferNegotiation, error) {
	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return nil, err
	}

	negotiation := sqlNegotiation{
		ID:              uuid.NewString(),
		TransferID:      transferID,
		OrganizationDID: organizationDID,
		CustomerID:      customerID,
		Date:            transferDate,
		Status:          transfer2.RequestedState,
		TaskID:          taskID,
	}

	const query = `INSERT INTO transfer_negotiation 
		(id, transfer_id, organization_did, customer_id, date, status, task_id)
		values(:id, :transfer_id, :organization_did, :customer_id, :date, :status, :task_id)`

	if _, err := tx.NamedExecContext(ctx, query, negotiation); err != nil {
		return nil, err
	}

	return negotiation.MarshalToDomainNegotiation()
}

func (r SQLiteTransferRepository) ListNegotiations(ctx context.Context, customerID int, transferID string) ([]domain.TransferNegotiation, error) {
	const query = `SELECT * FROM transfer_negotiation WHERE customer_id = ? AND transfer_id = ? ORDER BY organization_did ASC`
	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return nil, err
	}

	dbNegotiations := []sqlNegotiation{}
	err = tx.SelectContext(ctx, &dbNegotiations, query, customerID, transferID)
	if errors.Is(err, sql.ErrNoRows) {
		return []domain.TransferNegotiation{}, nil
	} else if err != nil {
		return nil, err
	}

	result := make([]domain.TransferNegotiation, len(dbNegotiations))

	for idx, dbNegotiation := range dbNegotiations {
		item, err := dbNegotiation.MarshalToDomainNegotiation()
		if err != nil {
			return nil, err
		}
		result[idx] = *item
	}

	return result, nil
}
