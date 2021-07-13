package patients

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sql2 "github.com/nuts-foundation/nuts-demo-ehr/domain/sql"
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

type sqlPatient struct {
	ID string `db:"id"`

	SSN sql.NullString `db:"ssn"`

	CustomerID string `db:"customer_id"`

	// Date of birth. Can include time if known.
	Dob sql.NullTime `db:"date_of_birth"`

	AvatarURL sql.NullString `db:"avatar_url"`

	// Primary email address.
	Email sql.NullString `db:"email"`

	// Given name
	FirstName string `db:"first_name"`

	// Family name. Must include prefixes like "van der".
	Surname string `db:"surname"`

	// Gender of the person according to https://www.hl7.org/fhir/valueset-administrative-gender.html.
	Gender string `db:"gender"`

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
	var email *openapi_types.Email
	if dbPatient.Email.Valid {
		otypeEmail := openapi_types.Email(dbPatient.Email.String)
		email = &otypeEmail
	}

	// Convert date of birth
	var dob *openapi_types.Date
	if dbPatient.Dob.Valid {
		dob = &openapi_types.Date{Time: dbPatient.Dob.Time}
	}

	var ssn *string
	if dbPatient.SSN.Valid {
		tmp := dbPatient.SSN.String
		ssn = &tmp
	}

	var avatarURL *string
	if dbPatient.AvatarURL.Valid {
		tmp := dbPatient.AvatarURL.String
		avatarURL = &tmp
	}

	return &domain.Patient{
		ObjectID: domain.ObjectID(dbPatient.ID),
		PatientProperties: domain.PatientProperties{
			Ssn:       ssn,
			FirstName: dbPatient.FirstName,
			Surname:   dbPatient.Surname,
			Dob:       dob,
			Email:     email,
			Gender:    gender,
			Zipcode:   dbPatient.Zipcode,
		},
		AvatarUrl: avatarURL,
	}, nil
}

func (dbPatient *sqlPatient) UnmarshalFromDomainPatient(customerID string, patient domain.Patient) error {
	var (
		email     string
		ssn       string
		dob       time.Time
		avatarURL string
	)
	if patient.Email != nil {
		tmp := *patient.Email
		email = string(tmp)
	}
	if patient.Ssn != nil {
		ssn = *patient.Ssn
	}
	if patient.Dob != nil {
		dob = patient.Dob.Time
	}
	if patient.AvatarUrl != nil {
		avatarURL = *patient.AvatarUrl
	}
	*dbPatient = sqlPatient{
		ID:         string(patient.ObjectID),
		SSN:        sql.NullString{String: ssn, Valid: ssn != ""},
		CustomerID: customerID,
		Dob:        sql.NullTime{Time: dob.UTC(), Valid: !dob.IsZero()},
		Email:      sql.NullString{String: email, Valid: email != ""},
		AvatarURL:  sql.NullString{String: avatarURL, Valid: avatarURL != ""},
		FirstName:  patient.FirstName,
		Surname:    patient.Surname,
		Gender:     string(patient.Gender),
		Zipcode:    patient.Zipcode,
	}
	return nil
}

type SQLitePatientRepository struct {
	factory Factory
	db      *sqlx.DB
}

const schema = `
	CREATE TABLE IF NOT EXISTS patient (
		id char(36) NOT NULL,
		ssn varchar(20) NOT NULL,
		customer_id varchar(100) NOT NULL,
		date_of_birth DATETIME DEFAULT NULL,
		email  varchar(100),
		first_name varchar(100) NOT NULL,
		surname varchar(100) NOT NULL,
		gender varchar(10) NOT NULL DEFAULT 'unknown',
		zipcode varchar(10) NOT NULL DEFAULT '',
	    avatar_url varchar(100),
		PRIMARY KEY (customer_id, id),
		UNIQUE(customer_id, ssn)
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

func (r SQLitePatientRepository) Update(ctx context.Context, customerID, id string, updateFn func(c domain.Patient) (*domain.Patient, error)) (patient *domain.Patient, err error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("%w, unable to start transaction", err)
	}

	defer func() {
		if err == nil {
			err = tx.Commit()
		} else {
			log.Debug(err)
			tx.Rollback()
			patient = nil
		}
		if err != nil {
			patient = nil
		}
	}()

	patient, err = r.getPatient(ctx, tx, customerID, id)
	if err != nil {
		return
	}
	updatedPatient, err := updateFn(*patient)
	if err != nil {
		return
	}

	dbPatient := sqlPatient{}
	err = dbPatient.UnmarshalFromDomainPatient(customerID, *updatedPatient)
	if err != nil {
		return
	}

	const query = `
	UPDATE patient SET
		date_of_birth = :date_of_birth,
		email = :email,
		first_name = :first_name,
		surname = :surname,
		gender = :gender, 
		zipcode = :zipcode
	WHERE customer_id = :customer_id AND id = :id
`
	_, err = tx.NamedExec(query, dbPatient)

	return
}

func (r SQLitePatientRepository) NewPatient(ctx context.Context, customerID string, patientProperties domain.PatientProperties) (patient *domain.Patient, err error) {
	patient, err = r.factory.NewPatientWithAvatar(patientProperties)
	if err != nil {
		return nil, err
	}
	dbPatient := sqlPatient{}
	err = dbPatient.UnmarshalFromDomainPatient(customerID, *patient)
	if err != nil {
		return nil, err
	}
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("%w, unable to start transaction", err)
	}

	defer func() {
		if err == nil {
			err = tx.Commit()
		} else {
			tx.Rollback()
			patient = nil
		}
		if err != nil {
			patient = nil
		}
	}()
	const query = `INSERT INTO patient 
		(id, ssn, customer_id, date_of_birth, email, first_name, surname, gender, zipcode, avatar_url)
		values(:id, :ssn, :customer_id, :date_of_birth, :email, :first_name, :surname, :gender, :zipcode, :avatar_url)
`

	_, err = tx.NamedExec(query, dbPatient)
	return
}

func (r SQLitePatientRepository) All(ctx context.Context, customerID string) ([]domain.Patient, error) {
	query := "SELECT * FROM `patient` WHERE customer_id = ? ORDER BY id ASC"
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

func (r SQLitePatientRepository) getPatient(ctx context.Context, db sql2.SQLContextGetter, customerID, patientID string) (*domain.Patient, error) {
	query := "SELECT * FROM `patient` WHERE customer_id = ? AND id = ? LIMIT 1"
	dbPatient := &sqlPatient{}
	err := db.GetContext(ctx, dbPatient, query, customerID, patientID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	patient, err := dbPatient.MarshalToDomainPatient()
	if err != nil {
		return nil, err
	}
	return patient, nil
}
