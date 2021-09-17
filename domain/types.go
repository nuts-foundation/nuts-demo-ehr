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
