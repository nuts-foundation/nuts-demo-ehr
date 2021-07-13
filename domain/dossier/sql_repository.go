package dossier

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

type sqlDossier struct {
	ID        string `db:"id"`
	Subject   string `db:"subject"`
	PatientID string `db:"patient_id"`
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

func (S SQLiteDossierRepository) FindByID(ctx context.Context, customerID, id string) (*domain.Dossier, error) {
	panic("implement me")
}

func (S SQLiteDossierRepository) Create(ctx context.Context, name, patientID string) (*domain.Dossier, error) {
	panic("implement me")
}

