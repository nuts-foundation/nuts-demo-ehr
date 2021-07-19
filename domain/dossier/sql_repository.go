package dossier

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

type sqlDossier struct {
	ID         string `db:"id"`
	Name       string `db:"name"`
	CustomerID string `db:"customer_id"`
	PatientID  string `db:"patient_id"`
}

func (dbDossier *sqlDossier) UnmarshalFromDomainDossier(customerID string, dossier *domain.Dossier) error {
	*dbDossier = sqlDossier{
		ID:         string(dossier.Id),
		Name:       dossier.Name,
		CustomerID: customerID,
		PatientID:  string(dossier.PatientID),
	}
	return nil
}

func (dbDossier sqlDossier) MarshalToDomainDossier() (*domain.Dossier, error) {
	return &domain.Dossier{
		Id:        domain.ObjectID(dbDossier.ID),
		Name:      dbDossier.Name,
		PatientID: domain.ObjectID(dbDossier.PatientID),
	}, nil
}

const schema = `
	CREATE TABLE IF NOT EXISTS dossier (
		id char(36) NOT NULL,
		customer_id char(36) NOT NULL,
	    patient_id char(36) NOT NULL,
		name varchar(20) NOT NULL,
		PRIMARY KEY (id),
		UNIQUE(customer_id, id)
	);
`

type SQLiteDossierRepository struct {
	factory Factory
	db      *sqlx.DB
}

func NewSQLiteDossierRepository(factory Factory, db *sqlx.DB) *SQLiteDossierRepository {
	if db == nil {
		panic("missing db for DossierRepository")
	}
	db.MustExec(schema)

	return &SQLiteDossierRepository{
		factory: factory,
		db:      db,
	}
}

func (r SQLiteDossierRepository) FindByID(ctx context.Context, customerID, id string) (*domain.Dossier, error) {
	panic("implement me")
}

func (r SQLiteDossierRepository) Create(ctx context.Context, customerID, name, patientID string) (*domain.Dossier, error) {
	dossier := r.factory.NewDossier(patientID, name)
	dbDossier := sqlDossier{}
	if err := dbDossier.UnmarshalFromDomainDossier(customerID, dossier); err != nil {
		return nil, err
	}

	const query = `INSERT INTO dossier
		(id, customer_id, patient_id, name)
		VALUES (:id, :customer_id, :patient_id, :name)`

	if _, err := r.db.NamedExec(query, dbDossier); err != nil {
		return nil, err
	}
	return dossier, nil
}

func (r SQLiteDossierRepository) AllByPatient(ctx context.Context, customerID, patientID string) ([]domain.Dossier, error) {
	const query = `SELECT * FROM dossier WHERE patient_id = ? and customer_id = ? ORDER BY id ASC`
	dbDossiers := []sqlDossier{}
	err := r.db.SelectContext(ctx, &dbDossiers, query, patientID, customerID)
	if errors.Is(err, sql.ErrNoRows) {
		return []domain.Dossier{}, nil
	} else if err != nil {
		return nil, err
	}

	result := make([]domain.Dossier, len(dbDossiers))
	for idx, dbDossier := range dbDossiers {
		dossier, _ := dbDossier.MarshalToDomainDossier()
		result[idx] = *dossier
	}
	return result, nil
}

