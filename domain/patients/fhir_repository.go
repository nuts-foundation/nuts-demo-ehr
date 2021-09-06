package patients

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
			Zipcode:   "",
		},
		AvatarUrl: &avatar,
	}
}

func ToFHIRPatient(domainPatient domain.Patient) map[string]interface{} {
	// TODO: Update to resources.Patient instead of map
	return map[string]interface{}{
		"resourceType": "Patient",
		"id":           domainPatient.ObjectID,
		"name": []map[string]interface{}{
			{
				"use":    "official",
				"family": domainPatient.Surname,
				"given":  []string{domainPatient.FirstName},
			},
		},
		"birthDate":  domainPatient.Dob.Format(dobFormat),
		"gender":     domainPatient.Gender,
		"photo":      []map[string]interface{}{{"url": domainPatient.AvatarUrl}},
		"identifier": []map[string]interface{}{{"system": bsnSystem, "value": domainPatient.Ssn}},
	}
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
	panic("implement me")
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
