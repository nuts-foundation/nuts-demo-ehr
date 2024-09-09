package dossier

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	sqlUtil "github.com/nuts-foundation/nuts-demo-ehr/sql"

	"github.com/jmoiron/sqlx"
)

type sqlDossier struct {
	ID         string `db:"id"`
	Name       string `db:"name"`
	CustomerID string `db:"customer_id"`
	PatientID  string `db:"patient_id"`
}

func (dbDossier *sqlDossier) UnmarshalFromDomainDossier(customerID string, dossier *types.Dossier) error {
	*dbDossier = sqlDossier{
		ID:         string(dossier.Id),
		Name:       dossier.Name,
		CustomerID: customerID,
		PatientID:  string(dossier.PatientID),
	}
	return nil
}

func (dbDossier sqlDossier) MarshalToDomainDossier() (*types.Dossier, error) {
	return &types.Dossier{
		Id:        dbDossier.ID,
		Name:      dbDossier.Name,
		PatientID: dbDossier.PatientID,
	}, nil
}

const schema = `
	CREATE TABLE IF NOT EXISTS dossier (
		id char(36) NOT NULL,
		customer_id varchar(255) NOT NULL,
	    patient_id char(36) NOT NULL,
		name varchar(20) NOT NULL,
		PRIMARY KEY (id),
		UNIQUE(customer_id, id)
	);
`

type SQLiteDossierRepository struct {
	factory Factory
}

func NewSQLiteDossierRepository(factory Factory, db *sqlx.DB) *SQLiteDossierRepository {
	if db == nil {
		panic("missing db for DossierRepository")
	}
	tx, _ := db.Beginx()
	tx.MustExec(schema)
	if err := tx.Commit(); err != nil {
		panic(err)
	}

	return &SQLiteDossierRepository{
		factory: factory,
	}
}

func (r SQLiteDossierRepository) FindByID(ctx context.Context, customerID string, id string) (*types.Dossier, error) {
	const query = `SELECT * FROM dossier WHERE customer_id = ? AND id = ? ORDER BY id ASC`

	dbDossier := sqlDossier{}
	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return nil, err
	}
	err = tx.GetContext(ctx, &dbDossier, query, customerID, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return dbDossier.MarshalToDomainDossier()
}

func (r SQLiteDossierRepository) Create(ctx context.Context, customerID string, name, patientID string) (*types.Dossier, error) {
	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return nil, err
	}
	dossier := r.factory.NewDossier(patientID, name)
	dbDossier := sqlDossier{}
	if err := dbDossier.UnmarshalFromDomainDossier(customerID, dossier); err != nil {
		return nil, err
	}

	const query = `INSERT INTO dossier
		(id, customer_id, patient_id, name)
		VALUES (:id, :customer_id, :patient_id, :name)`

	if _, err := tx.NamedExec(query, dbDossier); err != nil {
		return nil, err
	}
	return dossier, nil
}

func (r SQLiteDossierRepository) AllByPatient(ctx context.Context, customerID string, patientID string) ([]types.Dossier, error) {
	const query = `SELECT * FROM dossier WHERE patient_id = ? and customer_id = ? ORDER BY id ASC`
	dbDossiers := []sqlDossier{}
	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return nil, err
	}
	err = tx.SelectContext(ctx, &dbDossiers, query, patientID, customerID)
	if errors.Is(err, sql.ErrNoRows) {
		return []types.Dossier{}, nil
	} else if err != nil {
		return nil, err
	}

	result := make([]types.Dossier, len(dbDossiers))
	for idx, dbDossier := range dbDossiers {
		dossier, _ := dbDossier.MarshalToDomainDossier()
		result[idx] = *dossier
	}
	return result, nil
}
