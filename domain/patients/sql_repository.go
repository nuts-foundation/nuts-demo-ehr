package patients

import (
	"context"
	"database/sql"
	"errors"
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/jmoiron/sqlx"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

type sqlPatient struct {
	ID string `db:"id"`

	CustomerID string `db:"customer_id"`

	// Date of birth. Can include time if known.
	Dob sql.NullTime `db:"date_of_birth"`

	// Primary email address.
	Email sql.NullString `db:"email"`

	// Given name
	FirstName string `db:"first_name"`

	// Family name. Must include prefixes like "van der".
	Surname string `db:"surname"`

	// Gender of the person according to https://www.hl7.org/fhir/valueset-administrative-gender.html.
	Gender string `db:"gender"`

	// The internal ID of the Patient. Can be any internal system. Not to be confused by a database ID or a uuid.
	InternalID string `db:"internal_id"`

	// The zipcode formatted in dutch form. Can be used to find local care providers.
	Zipcode string `db:"zipcode"`
}

// MarshalToDomainPatient converts a sqlPatient into the domain.Patient.
// It make sure date and gender are correctly set.
func (dbPatient sqlPatient) MarshalToDomainPatient() (*domain.Patient, error) {

	// Convert gender
	var gender domain.PatientPropertiesGender
	switch dbPatient.Gender {
	case string(domain.PatientPropertiesGenderMale):
		gender = domain.PatientPropertiesGenderMale
	case string(domain.PatientPropertiesGenderFemale):
		gender = domain.PatientPropertiesGenderFemale
	case string(domain.PatientPropertiesGenderOther):
		gender = domain.PatientPropertiesGenderOther
	default:
		gender = domain.PatientPropertiesGenderUnknown
	}

	// Convert email
	var email *openapi_types.Email = nil
	if dbPatient.Email.Valid {
		otypeEmail := openapi_types.Email(dbPatient.Email.String)
		email = &otypeEmail
	}

	// Convert date of birth
	dob := time.Time{}
	if dbPatient.Dob.Valid {
		dob = dbPatient.Dob.Time
	}

	return &domain.Patient{
		PatientID:         domain.PatientID(dbPatient.ID),
		PatientProperties: domain.PatientProperties{
			FirstName: dbPatient.FirstName,
			Surname: dbPatient.Surname,
			Dob:        &openapi_types.Date{Time: dob},
			Email:      email,
			Gender:     gender,
			InternalID: dbPatient.InternalID,
			Zipcode:    dbPatient.Zipcode,
		},
	}, nil
}

// sqlContextGetter is an interface provided both by transaction and standard db connection
type sqlContextGetter interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type SQLitePatientRepository struct {
	factory Factory
	db *sqlx.DB
}

const schema = `
	CREATE TABLE IF NOT EXISTS patient (
		id TEXT NOT NULL,
		customer_id varchar(100) NOT NULL,
		date_of_birth DATE DEFAULT NULL,
		email  varchar(100),
		first_name varchar(100) NOT NULL DEFAULT '',
		surname varchar(100) NOT NULL DEFAULT '',
		gender varchar(10) NOT NULL DEFAULT 'unkown',
		internal_id varchar(100) NOT NULL,
		zipcode varchar(10) NOT NULL DEFAULT "",
		PRIMARY KEY (customer_id, id)
	);
`

func NewSQLitePatientRepository(factory Factory, db *sqlx.DB) *SQLitePatientRepository {
	if db == nil {
		panic("missing db")
	}

	db.MustExec(schema)

	return &SQLitePatientRepository{factory: factory, db: db}
}

func (r SQLitePatientRepository) FindByID(ctx context.Context, customerID, id string) (*domain.Patient, error) {
	return r.getPatient(ctx, r.db, customerID, id)
}

func (r SQLitePatientRepository) Update(ctx context.Context, customerID, id string, updateFn func(c domain.Patient) (*domain.Patient, error)) (*domain.Patient, error) {
	panic("implement me")
}

func (r SQLitePatientRepository) NewPatient(ctx context.Context, customerID string, patient domain.PatientProperties) (*domain.Patient, error) {
	panic("implement me")
}

func (r SQLitePatientRepository) All(ctx context.Context, customerID string) ([]domain.Patient, error) {
	query := "SELECT * FROM `patient` WHERE customer_id = ? ORDER BY internal_id ASC"
	dbPatients := []sqlPatient{}
	err := r.db.SelectContext(ctx, &dbPatients, query, customerID)
	if errors.Is(err, sql.ErrNoRows) {
		return []domain.Patient{}, nil
	} else if err != nil {
		return nil, err
	}

	result := make([]domain.Patient, len(dbPatients))

	for idx, dbPatient := range dbPatients {
		patient, err := dbPatient.MarshalToDomainPatient()
		if err != nil {
			return nil, err
		}
		result[idx] = *patient
	}
	return result, nil
}

func (r SQLitePatientRepository) getPatient(ctx context.Context, db sqlContextGetter, customerID, patientID string) (*domain.Patient, error) {
	query := "SELECT * FROM `patient` WHERE customer_id = ? AND id = ? LIMIT 1"
	dbPatient := &sqlPatient{}
	err := db.GetContext(ctx, dbPatient, query, customerID, patientID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &domain.Patient{
		PatientID: domain.PatientID(dbPatient.ID),
		PatientProperties: domain.PatientProperties{
			FirstName: dbPatient.FirstName,
		},
	}, nil
}
