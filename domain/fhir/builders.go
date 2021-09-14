package fhir

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir/eoverdracht"
)

func BuildNewTask(props TaskProperties) resources.Task {
	return resources.Task{
		Domain: resources.Domain{
			Base: resources.Base{
				ResourceType: "Task",
				ID:           ToIDPtr(generateResourceID()),
			},
		},
		Status: toCodePtr(props.Status),
		Code: &datatypes.CodeableConcept{Coding: []datatypes.Coding{{
			System:  &SnomedCodingSystem,
			Code:    &SnomedTransferCode,
			Display: &TransferDisplay,
		}}},
		Requester: &resources.TaskRequester{
			Agent: &datatypes.Reference{
				Identifier: &datatypes.Identifier{
					System: &NutsCodingSystem,
					Value:  ToStringPtr(props.RequesterID),
				},
			},
		},
		Owner: &datatypes.Reference{
			Identifier: &datatypes.Identifier{
				System: &NutsCodingSystem,
				Value:  ToStringPtr(props.OwnerID),
			}},
		// TODO: patient seems mandatory in the spec, but can only be sent when placed already
		// has patient in care to protect the identity of the patient during the negotiation phase.
		//"for": map[string]string{
		//	"reference": fmt.Sprintf("Patient/%s", domainTask.PatientID),
		//},
		Input:  props.Input,
		Output: props.Output,
	}
}

func BuildNewComposition(elements map[string]interface{}) Composition {
	fhirData := map[string]interface{}{
		"resourceType": "Composition",
		"id":           generateResourceID(),
		"status":       "final",
		"date":         time.Now().Format(time.RFC3339),
	}
	for key, value := range elements {
		fhirData[key] = value
	}
	return fhirData
}

func BuildAdvanceNotice2(createRequest domain.CreateTransferRequest) eoverdracht.AdvanceNotice {
	problems, interventions := buildCarePlan(createRequest.CarePlan)

	an := eoverdracht.AdvanceNotice{
		Patient:       resources.Patient{},
		Problems:      problems,
		Interventions: interventions,
	}

	composition := buildAdvanceNoticeComposition(an)
	an.Composition = composition

	return an
}

func buildAdvanceNoticeComposition(an eoverdracht.AdvanceNotice) eoverdracht.Composition {

	// new patientProblems
	patientProblems := eoverdracht.CompositionSection{
		Title: ToStringPtr("Current patient problems"),
		Code: datatypes.CodeableConcept{
			Coding: []datatypes.Coding{{
				System:  &SnomedCodingSystem,
				Code:    ToCodePtr("86644006"),
				Display: ToStringPtr("Nursing diagnosis"),
			}},
		},
	}

	// Add the problems as a section
	for _, p := range an.Problems {
		patientProblems.Entry = append(patientProblems.Entry, datatypes.Reference{Reference: ToStringPtr("Condition/" + FromIDPtr(p.ID))})
	}
	for _, i := range an.Interventions {
		patientProblems.Entry = append(patientProblems.Entry, datatypes.Reference{Reference: ToStringPtr("Procedure/" + FromIDPtr(i.ID))})
	}

	// Start with empty care plan
	careplan := eoverdracht.CompositionSection{
		Code: datatypes.CodeableConcept{
			Coding: []datatypes.Coding{{
				System:  &SnomedCodingSystem,
				Code:    ToCodePtr("773130005"),
				Display: ToStringPtr("Nursing care plan (record artifact)"),
			}}},
		Section: []eoverdracht.CompositionSection{
			patientProblems,
		},
	}


	return eoverdracht.Composition{
		Base: resources.Base{
			ResourceType: "Composition",
			ID:           ToIDPtr(generateResourceID()),
		},
		Type: datatypes.CodeableConcept{
			Coding: []datatypes.Coding{{System: &LoincCodingSystem, Code: ToCodePtr("57830-2")}},
		},
		Title:   "Advance notice",
		Section: []eoverdracht.CompositionSection{careplan},
	}
}

func buildCarePlan(carePlan domain.CarePlan) (problems []resources.Condition, interventions []eoverdracht.Procedure) {
	for _, cpPatientProblems := range carePlan.PatientProblems {
		newProblem := buildConditionFromProblem(cpPatientProblems.Problem)
		problems = append(problems, newProblem)
		for _, i := range cpPatientProblems.Interventions {
			interventions = append(interventions, buildProcedureFromIntervention(i, FromIDPtr(newProblem.ID)))
		}
	}
	return problems, interventions
}

func buildProcedureFromIntervention(intervention domain.Intervention, problemID string) eoverdracht.Procedure {
	return eoverdracht.Procedure{
		Domain: resources.Domain{
			Base: resources.Base{
				ResourceType: "Procedure",
				ID:           ToIDPtr(generateResourceID()),
			},
		},
		ReasonReference: datatypes.Reference{Reference: ToStringPtr("Condition/" + problemID)},
		Note:            []datatypes.Annotation{{Text: ToStringPtr(intervention.Comment)}},
	}
}

func buildConditionFromProblem(problem domain.Problem) resources.Condition {
	return resources.Condition{
		Domain: resources.Domain{
			Base: resources.Base{
				ResourceType: "Condition",
				ID:           ToIDPtr(generateResourceID()),
			},
		},
		Note: []datatypes.Annotation{{Text: ToStringPtr(problem.Name)}},
	}
}

func BuildAdvanceNotice() Composition {
	elements := map[string]interface{}{
		"title": "Advance notice",
		"type":  LoincAdvanceNoticeType,
		// TODO: patient seems mandatory in the spec, but can only be sent when placer already
		// has patient in care to protect the identity of the patient during the negotiation phase.
		//"subject":  fhir.Reference{Reference: "Patient/Anonymous"},
		"author": eoverdracht.Practitioner{
			// TODO: Derive from authenticated user?
			Identifier: datatypes.Identifier{
				System: &UZICodingSystem,
				Value:  ToStringPtr("12345"),
			},
			Name: &datatypes.HumanName{
				Family: ToStringPtr("Demo EHR"),
				Given:  []datatypes.String{"Nuts"},
			},
		},
		// TODO: sections
	}
	return BuildNewComposition(elements)
}

func BuildNursingHandoff(patient *domain.Patient) Composition {
	patientPath := fmt.Sprintf("Patient/%s", string(patient.ObjectID))
	elements := map[string]interface{}{
		"title":   "Nursing handoff",
		"type":    SnomedNursingHandoffType,
		"subject": datatypes.Reference{Reference: ToStringPtr(patientPath)},
		"author": eoverdracht.Practitioner{
			// TODO: Derive from authenticated user?
			Identifier: datatypes.Identifier{
				System: &UZICodingSystem,
				Value:  ToStringPtr("12345"),
			},
			Name: &datatypes.HumanName{
				Family: ToStringPtr("Demo EHR"),
				Given:  []datatypes.String{"Nuts"},
			},
		},
		// TODO: sections
	}
	return BuildNewComposition(elements)
}

func generateResourceID() string {
	return uuid.New().String()
}
