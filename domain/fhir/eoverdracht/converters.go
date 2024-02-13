package eoverdracht

import (
	"encoding/json"
	"fmt"
	openapiTypes "github.com/oapi-codegen/runtime/types"
	"strings"
	"time"

	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	"github.com/tidwall/gjson"
)

func ToDomainProblem(condition resources.Condition) types.Problem {
	var notes []string
	for _, note := range condition.Note {
		notes = append(notes, fhir.FromStringPtr(note.Text))
	}
	return types.Problem{
		// TODO: get the name from the value instead of the notes
		Name:   strings.Join(notes, ","),
		Status: "active",
	}
}

func ToDomainIntervention(procedure fhir.Procedure) types.Intervention {
	return types.Intervention{Comment: fhir.FromStringPtr(procedure.Note[0].Text)}
}

func ToDomainPatient(fhirPatient resources.Patient) types.Patient {
	asJSON, _ := json.Marshal(fhirPatient)
	p := gjson.ParseBytes(asJSON)

	dob, _ := time.Parse(types.DobFormat, p.Get("birthDate").String())
	gender := types.Unknown
	fhirGender := p.Get("gender").String()
	switch fhirGender {
	case string(types.Male):
		gender = types.Male
	case string(types.Female):
		gender = types.Female
	}
	ssn := p.Get(fmt.Sprintf(`identifier.#(system==%s).value`, types.BsnSystem)).String()
	avatar := p.Get(`photo.0.url`).String()
	return types.Patient{
		ObjectID:  p.Get("id").String(),
		Dob:       &openapiTypes.Date{Time: dob},
		Email:     nil,
		FirstName: p.Get(`name.0.given.0`).String(),
		Gender:    gender,
		Ssn:       &ssn,
		Surname:   p.Get(`name.0.family`).String(),
		Zipcode:   p.Get(`address.0.postalCode`).String(),
		AvatarUrl: &avatar,
	}
}

func AdvanceNoticeToDomainTransfer(notice AdvanceNotice) (types.TransferProperties, error) {
	patient := ToDomainPatient(notice.Patient)
	adminData, err := FilterCompositionSectionByType(notice.Composition.Section, AdministrativeDocCode)
	if err != nil {
		return types.TransferProperties{}, fmt.Errorf("administrativeData section missing in advance notice: %w", err)
	}
	transferDate, _ := time.Parse(time.RFC3339, *(*string)(adminData.Extension[0].ValueDateTime))

	domainTransfer := types.TransferProperties{
		Patient:      patient,
		CarePlan:     types.CarePlan{PatientProblems: []types.PatientProblem{}},
		TransferDate: openapiTypes.Date{Time: transferDate},
	}

	for _, condition := range notice.Problems {
		var interventions []types.Intervention
		for _, procedure := range notice.Interventions {
			if fhir.FromStringPtr(procedure.ReasonReference[0].Reference) == "Condition/"+fhir.FromIDPtr(condition.ID) {
				interventions = append(interventions, ToDomainIntervention(procedure))
			}
		}
		domainTransfer.CarePlan.PatientProblems = append(domainTransfer.CarePlan.PatientProblems,
			types.PatientProblem{
				Interventions: interventions,
				Problem:       ToDomainProblem(condition),
			})

	}
	return domainTransfer, nil
}

func NursingHandoffToDomainTransfer(nursingHandoff NursingHandoff) (types.TransferProperties, error) {
	patient := ToDomainPatient(nursingHandoff.Patient)
	adminData, err := FilterCompositionSectionByType(nursingHandoff.Composition.Section, AdministrativeDocCode)
	if err != nil {
		return types.TransferProperties{}, fmt.Errorf("administrativeData section missing in advance notice: %w", err)
	}
	transferDate, _ := time.Parse(time.RFC3339, *(*string)(adminData.Extension[0].ValueDateTime))

	domainTransfer := types.TransferProperties{
		Patient:      patient,
		CarePlan:     types.CarePlan{PatientProblems: []types.PatientProblem{}},
		TransferDate: openapiTypes.Date{Time: transferDate},
	}

	for _, condition := range nursingHandoff.Problems {
		var interventions []types.Intervention
		for _, procedure := range nursingHandoff.Interventions {
			if fhir.FromStringPtr(procedure.ReasonReference[0].Reference) == "Condition/"+fhir.FromIDPtr(condition.ID) {
				interventions = append(interventions, ToDomainIntervention(procedure))
			}
		}
		domainTransfer.CarePlan.PatientProblems = append(domainTransfer.CarePlan.PatientProblems,
			types.PatientProblem{
				Interventions: interventions,
				Problem:       ToDomainProblem(condition),
			})

	}
	return domainTransfer, nil
}
