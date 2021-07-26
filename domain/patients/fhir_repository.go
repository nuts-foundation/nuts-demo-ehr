package patients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	url2 "net/url"
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/tidwall/gjson"
)

type fhirPatient struct {
	data gjson.Result
}

func newFHIRPatientFromGJSON(data gjson.Result) *fhirPatient {
	return &fhirPatient{data: data}
}

func newFHIRPatientFromJSON(jsonBytes []byte) *fhirPatient {
	return newFHIRPatientFromGJSON(gjson.ParseBytes(jsonBytes))
}

const dobFormat = "2006-01-02"
const bsnSystem = "http://fhir.nl/fhir/NamingSystem/bsn"

func (p fhirPatient) MarshalToDomainPatient() (*domain.Patient, error) {
	dob, _ := time.Parse(dobFormat, p.data.Get("birthDate").String())
	gender := domain.PatientPropertiesGenderUnknown
	fhirGender := p.data.Get("gender").String()
	switch fhirGender {
	case string(domain.PatientPropertiesGenderMale):
		gender = domain.PatientPropertiesGenderMale
	case string(domain.PatientPropertiesGenderFemale):
		gender = domain.PatientPropertiesGenderFemale
	}
	ssn := p.data.Get(fmt.Sprintf(`identifier.#(system==%s).value`, bsnSystem)).String()
	avatar := p.data.Get(`photo.0.url`).String()
	return &domain.Patient{
		ObjectID: domain.ObjectID(p.data.Get("id").String()),
		PatientProperties: domain.PatientProperties{
			Dob:       &openapi_types.Date{Time: dob},
			Email:     nil,
			FirstName: p.data.Get(`name.#(use=="official").given.0`).String(),
			Gender:    gender,
			Ssn:       &ssn,
			Surname:   p.data.Get(`name.#(use=="official").family`).String(),
			Zipcode:   "",
		},
		AvatarUrl: &avatar,
	}, nil
}

func (p *fhirPatient) UnmarshalFromDomainPatient(domainPatient domain.Patient) error {
	fhirData := map[string]interface{}{
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
	jsonData, err := json.Marshal(fhirData)
	if err != nil {
		return fmt.Errorf("error unmarshalling fhirPatient from domain.Patient: %w", err)
	}
	*p = *newFHIRPatientFromJSON(jsonData)
	return nil
}

type FHIRPatientRepository struct {
	url     string
	factory Factory
}

func NewFHIRPatientRepository(factory Factory, url string) *FHIRPatientRepository {
	return &FHIRPatientRepository{
		url:     url,
		factory: factory,
	}
}

func (r FHIRPatientRepository) FindByID(ctx context.Context, customerID, id string) (*domain.Patient, error) {
	client := http.Client{}
	res, err := client.Get(r.url + "/Patient/" + id)
	if err != nil {
		return nil, err
	}
	rawPatient, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	parsedGJSON := gjson.Parse(string(rawPatient))
	p := newFHIRPatientFromGJSON(parsedGJSON)
	patient, _ := p.MarshalToDomainPatient()
	return patient, nil
}

func (r FHIRPatientRepository) Update(ctx context.Context, customerID, id string, updateFn func(c domain.Patient) (*domain.Patient, error)) (*domain.Patient, error) {
	panic("implement me")
}

func (r FHIRPatientRepository) NewPatient(ctx context.Context, customerID string, patientProperties domain.PatientProperties) (*domain.Patient, error) {
	p := fhirPatient{}
	newPatient, err := r.factory.NewPatientWithAvatar(patientProperties)
	if err != nil {
		return nil, err
	}
	err = p.UnmarshalFromDomainPatient(*newPatient)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Post(r.url+"/Patient", "application/json", bytes.NewBuffer([]byte(p.data.Raw)))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		body, ioErr := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("unable to create new patient. Unable to read error response: ioerr: %s", ioErr)
		}
		return nil, fmt.Errorf("unable to create new patient: %s", body)
	}

	return newPatient, nil
}

func (r FHIRPatientRepository) All(ctx context.Context, customerID string, name *string) ([]domain.Patient, error) {
	client := http.Client{}
	url := r.url + "/Patient"
	if name != nil {
		url += "/_search?name=" + url2.QueryEscape(*name)
	}
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	bundle, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	parsedGJSON := gjson.Parse(string(bundle))
	bundleEntries := parsedGJSON.Get("entry.#.resource").Array()

	patients := make([]domain.Patient, len(bundleEntries))

	for idx, bundleEntry := range bundleEntries {
		p := newFHIRPatientFromGJSON(bundleEntry)
		patient, _ := p.MarshalToDomainPatient()
		patients[idx] = *patient
	}

	return patients, nil
}
