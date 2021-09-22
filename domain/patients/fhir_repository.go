package patients

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/tidwall/gjson"
)

const dobFormat = "2006-01-02"
const bsnSystem = "http://fhir.nl/fhir/NamingSystem/bsn"

func ToDomainPatient(fhirPatient resources.Patient) domain.Patient {
	asJSON, _ := json.Marshal(fhirPatient)
	p := gjson.ParseBytes(asJSON)

	dob, _ := time.Parse(dobFormat, p.Get("birthDate").String())
	gender := domain.PatientPropertiesGenderUnknown
	fhirGender := p.Get("gender").String()
	switch fhirGender {
	case string(domain.PatientPropertiesGenderMale):
		gender = domain.PatientPropertiesGenderMale
	case string(domain.PatientPropertiesGenderFemale):
		gender = domain.PatientPropertiesGenderFemale
	}
	ssn := p.Get(fmt.Sprintf(`identifier.#(system==%s).value`, bsnSystem)).String()
	avatar := p.Get(`photo.0.url`).String()
	return domain.Patient{
		ObjectID: domain.ObjectID(p.Get("id").String()),
		PatientProperties: domain.PatientProperties{
			Dob:       &openapi_types.Date{Time: dob},
			Email:     nil,
			FirstName: p.Get(`name.#(use=="official").given.0`).String(),
			Gender:    gender,
			Ssn:       &ssn,
			Surname:   p.Get(`name.#(use=="official").family`).String(),
			Zipcode:   p.Get(`address.0.postalCode`).String(),
		},
		AvatarUrl: &avatar,
	}
}

func ToFHIRPatient(domainPatient domain.Patient) resources.Patient {
	dob := datatypes.Date(domainPatient.Dob.Format(dobFormat))

	fhirPatient := resources.Patient{
		Domain: resources.Domain{
			Base:      resources.Base{
				ResourceType: "Patient",
				ID: fhir.ToIDPtr(string(domainPatient.ObjectID)),
			},
		},
		Name: []datatypes.HumanName{{
			Use:     fhir.ToCodePtr("official"),
			Family:  fhir.ToStringPtr(domainPatient.Surname),
			Given:   []datatypes.String{datatypes.String(domainPatient.FirstName)},
		}},
		BirthDate: &dob,
		Gender: fhir.ToCodePtr(string(domainPatient.Gender)),
	}

	if domainPatient.Ssn != nil {
		fhirPatient.Identifier = append(fhirPatient.Identifier, datatypes.Identifier{System: fhir.ToUriPtr(bsnSystem), Value: fhir.ToStringPtr(*domainPatient.Ssn)})
	}
	if domainPatient.AvatarUrl != nil {
		fhirPatient.Photo = append(fhirPatient.Photo, datatypes.Attachment{URL: fhir.ToUriPtr(*domainPatient.AvatarUrl) })
	}
	if domainPatient.Zipcode != "" {
		fhirPatient.Address = append(fhirPatient.Address, datatypes.Address{
			PostalCode: fhir.ToStringPtr(domainPatient.Zipcode),
		})
	}

	return fhirPatient
}

type FHIRPatientRepository struct {
	fhirClientFactory fhir.Factory
	factory           Factory
}

func NewFHIRPatientRepository(factory Factory, fhirClientFactory fhir.Factory) *FHIRPatientRepository {
	return &FHIRPatientRepository{
		fhirClientFactory: fhirClientFactory,
		factory:           factory,
	}
}

func (r FHIRPatientRepository) FindByID(ctx context.Context, customerID int, id string) (*domain.Patient, error) {
	patient := resources.Patient{}
	err := r.fhirClientFactory(fhir.WithTenant(customerID)).ReadOne(ctx, "Patient/"+id, &patient)
	if err != nil {
		return nil, err
	}
	result := ToDomainPatient(patient)
	return &result, nil
}

func (r FHIRPatientRepository) Update(ctx context.Context, customerID int, id string, updateFn func(c domain.Patient) (*domain.Patient, error)) (*domain.Patient, error) {
	domainPatient, err := r.FindByID(ctx, customerID, id)
	if err != nil {
		return nil, fmt.Errorf("could not update patient: could not read current patient from FHIR store: %w", err)
	}
	updatedDomainPatient, err := updateFn(*domainPatient)
	if err != nil {
		return nil, err
	}
	updatedFHIRPatient := ToFHIRPatient(*updatedDomainPatient)
	return updatedDomainPatient, fhirClient.CreateOrUpdate(ctx, updatedFHIRPatient)
}

func (r FHIRPatientRepository) NewPatient(ctx context.Context, customerID int, patientProperties domain.PatientProperties) (*domain.Patient, error) {
	patient, err := r.factory.NewPatientWithAvatar(patientProperties)
	if err != nil {
		return nil, err
	}
	err = r.fhirClientFactory(fhir.WithTenant(customerID)).CreateOrUpdate(ctx, ToFHIRPatient(*patient))
	if err != nil {
		return nil, err
	}
	return patient, nil
}

func (r FHIRPatientRepository) All(ctx context.Context, customerID int, name *string) ([]domain.Patient, error) {
	var params map[string]string
	if name != nil {
		params = map[string]string{"name": *name}
	}
	fhirPatients := []resources.Patient{}
	err := r.fhirClientFactory(fhir.WithTenant(customerID)).ReadMultiple(ctx, "Patient", params, &fhirPatients)
	if err != nil {
		return nil, err
	}

	patients := make([]domain.Patient, 0)
	for _, patient := range fhirPatients {
		patients = append(patients, ToDomainPatient(patient))
	}

	return patients, nil
}
