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
	problems, interventions, careplan := buildCarePlan(createRequest.CarePlan)
	administrativeData := buildAdministrativeData(createRequest)

	an := eoverdracht.AdvanceNotice{
		Patient:       resources.Patient{},
		Problems:      problems,
		Interventions: interventions,
	}

	composition := buildAdvanceNoticeComposition(administrativeData, careplan)
	an.Composition = composition

	return an
}

func buildAdministrativeData(request domain.CreateTransferRequest) eoverdracht.CompositionSection {
	transferDate := request.TransferDate.Format(time.RFC3339)
	return eoverdracht.CompositionSection{
		BackboneElement: datatypes.BackboneElement{
			Element: datatypes.Element{
				Extension: []datatypes.Extension{{
					URL:           (*datatypes.URI)(ToStringPtr("http://nictiz.nl/fhir/StructureDefinition/eOverdracht-TransferDate")),
					ValueDateTime: (*datatypes.DateTime)(ToStringPtr(transferDate)),
				}},
			},
		},
		Title: ToStringPtr("Administrative data"),
		Code: datatypes.CodeableConcept{
			Coding: []datatypes.Coding{{
				System:  &SnomedCodingSystem,
				Code:    ToCodePtr(eoverdracht.AdministrativeDocCode),
				Display: ToStringPtr("Administrative documentation (record artifact)"),
			}}},
	}

}

func buildAdvanceNoticeComposition(administrativeData, careplan eoverdracht.CompositionSection) eoverdracht.Composition {

	return eoverdracht.Composition{
		Base: resources.Base{
			ResourceType: "Composition",
			ID:           ToIDPtr(generateResourceID()),
		},
		Type: datatypes.CodeableConcept{
			Coding: []datatypes.Coding{{System: &LoincCodingSystem, Code: ToCodePtr("57830-2")}},
		},
		Title:   "Advance notice",
		Section: []eoverdracht.CompositionSection{administrativeData,careplan},
	}
}

func buildCarePlan(carePlan domain.CarePlan) (problems []resources.Condition, interventions []eoverdracht.Procedure, section eoverdracht.CompositionSection) {
	for _, cpPatientProblems := range carePlan.PatientProblems {
		newProblem := buildConditionFromProblem(cpPatientProblems.Problem)
		problems = append(problems, newProblem)
		for _, i := range cpPatientProblems.Interventions {
			interventions = append(interventions, buildProcedureFromIntervention(i, FromIDPtr(newProblem.ID)))
		}
	}

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
	for _, p := range problems {
		patientProblems.Entry = append(patientProblems.Entry, datatypes.Reference{Reference: ToStringPtr("Condition/" + FromIDPtr(p.ID))})
	}
	for _, i := range interventions {
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
	return problems, interventions, careplan
}

func buildProcedureFromIntervention(intervention domain.Intervention, problemID string) eoverdracht.Procedure {
	return eoverdracht.Procedure{
		Domain: resources.Domain{
			Base: resources.Base{
				ResourceType: "Procedure",
				ID:           ToIDPtr(generateResourceID()),
			},
		},
		ReasonReference: []datatypes.Reference{{Reference: ToStringPtr("Condition/" + problemID)}},
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
