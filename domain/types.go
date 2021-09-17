package domain

import (
	"fmt"
	"strings"
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir/eoverdracht"
)

type IncomingTransfer struct {
	Id         ObjectID                  `json:"id"`
	FhirTaskID string                    `json:"fhirTaskID"`
	Sender     Organization              `json:"sender"`
	Status     TransferNegotiationStatus `json:"status"`
	CreatedAt  time.Time                 `json:"createdAt"`
}

// Transfer data types as described by the ART Decor app from the Nictiz
// https://decor.nictiz.nl/pub/eoverdracht/e-overdracht-html-20210510T093529/tr-2.16.840.1.113883.2.4.3.11.60.30.4.63-2021-01-27T000000.html

type AdministrativeData struct {
	TransferDate time.Time
	PersonalInformation
}

type PersonalInformation struct {
	AnonymousPatient
}

type AnonymousPatient struct {
	PostalCode string
}

//type CarePlan struct {
//	Problems []Problem
//}

//// Problem a s defined by https://zibs.nl/wiki/Probleem-v4.1(2017NL)
//type Problem struct {
//	Explanation string // NL-CM:5.1.5
//	Active      bool   // NL-CM:5.1.4
//	// Code from the Dutch Snomed CT Kernset patientproblemen
//	// https://www.nictiz.nl/standaardisatie/terminologiecentrum/referentielijsten/nationale-kernset/nationale-kernset-patientproblemen/
//	Code        string // NL-CM:5.1.3
//	Intervention NursingIntervention // Intervention to perform to treat this problem
//}

// NursingIntervention as defined by https://zibs.nl/wiki/NursingIntervention-v3.1(2017EN)
type NursingIntervention struct {
	Comment string // NL-CM:14.2.5 Comment on the nursing intervention
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

func FHIRAdvanceNoticeToDomainTransfer(notice eoverdracht.AdvanceNotice) (TransferProperties, error) {
	adminData, err := eoverdracht.FilterCompositionSectionByType(notice.Composition.Section, eoverdracht.AdministrativeDocCode)
	if err != nil {
		return TransferProperties{}, fmt.Errorf("administrativeData section missing in advance notice: %w", err)
	}
	transferDate, _ := time.Parse(time.RFC3339, *(*string)(adminData.Extension[0].ValueDateTime))

	domainTransfer := TransferProperties{
		CarePlan:     CarePlan{PatientProblems: []PatientProblem{}},
		TransferDate: openapi_types.Date{transferDate},
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
