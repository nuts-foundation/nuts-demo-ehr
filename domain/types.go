package domain

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir/eoverdracht"
	"github.com/tidwall/gjson"
)

const DobFormat = "2006-01-02"
const BsnSystem = "http://fhir.nl/fhir/NamingSystem/bsn"

type IncomingTransfer struct {
	Id         ObjectID                  `json:"id"`
	FhirTaskID string                    `json:"fhirTaskID"`
	Sender     Organization              `json:"sender"`
	Status     TransferNegotiationStatus `json:"status"`
	CreatedAt  time.Time                 `json:"createdAt"`
}

func FHIRConditionToDomainProblem(condition resources.Condition) Problem {
	var notes []string
	for _, note := range condition.Note {
		notes = append(notes, fhir.FromStringPtr(note.Text))
	}
	return Problem{
		// TODO: get the name from the value instead of the notes
		Name:   strings.Join(notes, ","),
		Status: "active",
	}
}

func FHIRProcedureToDomainIntervention(procedure eoverdracht.Procedure) Intervention {
	return Intervention{Comment: fhir.FromStringPtr(procedure.Note[0].Text)}
}

func FHIRPatientToDomainPatient(fhirPatient resources.Patient) Patient {
	asJSON, _ := json.Marshal(fhirPatient)
	p := gjson.ParseBytes(asJSON)

	dob, _ := time.Parse(DobFormat, p.Get("birthDate").String())
	gender := PatientPropertiesGenderUnknown
	fhirGender := p.Get("gender").String()
	switch fhirGender {
	case string(PatientPropertiesGenderMale):
		gender = PatientPropertiesGenderMale
	case string(PatientPropertiesGenderFemale):
		gender = PatientPropertiesGenderFemale
	}
	ssn := p.Get(fmt.Sprintf(`identifier.#(system==%s).value`, BsnSystem)).String()
	avatar := p.Get(`photo.0.url`).String()
	return Patient{
		ObjectID: ObjectID(p.Get("id").String()),
		PatientProperties: PatientProperties{
			Dob:       &openapi_types.Date{Time: dob},
			Email:     nil,
			FirstName: p.Get(`name.0.given.0`).String(),
			Gender:    gender,
			Ssn:       &ssn,
			Surname:   p.Get(`name.0.family`).String(),
			Zipcode:   p.Get(`address.0.postalCode`).String(),
		},
		AvatarUrl: &avatar,
	}
}

func FHIRAdvanceNoticeToDomainTransfer(notice eoverdracht.AdvanceNotice) (TransferProperties, error) {
	patient := FHIRPatientToDomainPatient(notice.Patient)
	adminData, err := eoverdracht.FilterCompositionSectionByType(notice.Composition.Section, eoverdracht.AdministrativeDocCode)
	if err != nil {
		return TransferProperties{}, fmt.Errorf("administrativeData section missing in advance notice: %w", err)
	}
	transferDate, _ := time.Parse(time.RFC3339, *(*string)(adminData.Extension[0].ValueDateTime))

	domainTransfer := TransferProperties{
		Patient:      patient,
		CarePlan:     CarePlan{PatientProblems: []PatientProblem{}},
		TransferDate: openapi_types.Date{Time: transferDate},
	}

	for _, condition := range notice.Problems {
		var interventions []Intervention
		for _, procedure := range notice.Interventions {
			if fhir.FromStringPtr(procedure.ReasonReference[0].Reference) == "Condition/"+fhir.FromIDPtr(condition.ID) {
				interventions = append(interventions, FHIRProcedureToDomainIntervention(procedure))
			}
		}
		domainTransfer.CarePlan.PatientProblems = append(domainTransfer.CarePlan.PatientProblems,
			PatientProblem{
				Interventions: interventions,
				Problem:       FHIRConditionToDomainProblem(condition),
			})

	}
	return domainTransfer, nil
}

func FHIRNursingHandoffToDomainTransfer(notice eoverdracht.NursingHandoff) (TransferProperties, error) {
	patient := FHIRPatientToDomainPatient(notice.Patient)
	adminData, err := eoverdracht.FilterCompositionSectionByType(notice.Composition.Section, eoverdracht.AdministrativeDocCode)
	if err != nil {
		return TransferProperties{}, fmt.Errorf("administrativeData section missing in advance notice: %w", err)
	}
	transferDate, _ := time.Parse(time.RFC3339, *(*string)(adminData.Extension[0].ValueDateTime))

	domainTransfer := TransferProperties{
		Patient:      patient,
		CarePlan:     CarePlan{PatientProblems: []PatientProblem{}},
		TransferDate: openapi_types.Date{Time: transferDate},
	}

	for _, condition := range notice.Problems {
		var interventions []Intervention
		for _, procedure := range notice.Interventions {
			if fhir.FromStringPtr(procedure.ReasonReference[0].Reference) == "Condition/"+fhir.FromIDPtr(condition.ID) {
				interventions = append(interventions, FHIRProcedureToDomainIntervention(procedure))
			}
		}
		domainTransfer.CarePlan.PatientProblems = append(domainTransfer.CarePlan.PatientProblems,
			PatientProblem{
				Interventions: interventions,
				Problem:       FHIRConditionToDomainProblem(condition),
			})

	}
	return domainTransfer, nil
}
