package sharedcareplan

import (
	"context"
	"github.com/jmoiron/sqlx"
	sqlUtil "github.com/nuts-foundation/nuts-demo-ehr/sql"
)

const schema = `
	CREATE TABLE IF NOT EXISTS shared_careplan (
	    dossier_id char(36) NOT NULL,
		customer_id varchar(255) NOT NULL,
		reference varchar(200) NOT NULL,
		PRIMARY KEY (dossier_id)
	);
`

func NewRepository(db *sqlx.DB) (Repository, error) {
	tx, _ := db.Beginx()
	tx.MustExec(schema)
	if err := tx.Commit(); err != nil {
		panic(err)
	}
	return Repository{}, nil
}

type Repository struct {
}

type sqlSharedCarePlan struct {
	DossierID  string `db:"dossier_id"`
	CustomerID string `db:"customer_id"`
	// Reference is the FHIR Reference, as absolute URL, to the CarePlan on the shared Care Plan Service.
	// It can be used by FHIR clients to resolve the CarePlan.
	Reference string `db:"reference"`
}

func (r Repository) Create(ctx context.Context, customerID, dossierID string, reference string) error {
	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return err
	}
	dbCarePlan := sqlSharedCarePlan{}
	dbCarePlan.CustomerID = customerID
	dbCarePlan.DossierID = dossierID
	dbCarePlan.Reference = reference

	const query = `INSERT INTO shared_careplan
		(customer_id, dossier_id, reference)
		VALUES (:customer_id, :dossier_id, :reference)`
	_, err = tx.NamedExec(query, dbCarePlan)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) FindByDossierID(ctx context.Context, customerID, dossierID string) (*SharedCarePlan, error) {
	tx, err := sqlUtil.GetTransaction(ctx)
	if err != nil {
		return nil, err
	}
	dbCarePlan := sqlSharedCarePlan{}
	const query = `SELECT * FROM shared_careplan WHERE customer_id = :customer_id AND dossier_id = :dossier_id`
	err = tx.Get(&dbCarePlan, query, customerID, dossierID)
	if err != nil {
		return nil, err
	}
	return r.sqlToDomain(dbCarePlan), nil
}

func (r Repository) sqlToDomain(dbCarePlan sqlSharedCarePlan) *SharedCarePlan {
	return &SharedCarePlan{
		DossierID:  dbCarePlan.DossierID,
		CustomerID: dbCarePlan.CustomerID,
		Reference:  dbCarePlan.Reference,
	}
}
