package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir/eoverdracht"
)

func BuildNewTask(props fhir.TaskProperties) resources.Task {
	return resources.Task{
		Domain: resources.Domain{
			Base: resources.Base{
				ResourceType: "Task",
				ID:           fhir.ToIDPtr(generateResourceID()),
			},
		},
		Status: fhir.ToCodePtr(props.Status),
		Code: &datatypes.CodeableConcept{Coding: []datatypes.Coding{{
			System:  &fhir.SnomedCodingSystem,
			Code:    &fhir.SnomedTransferCode,
			Display: &fhir.TransferDisplay,
		}}},
		Requester: &resources.TaskRequester{
			Agent: &datatypes.Reference{
				Identifier: &datatypes.Identifier{
					System: &fhir.NutsCodingSystem,
					Value:  fhir.ToStringPtr(props.RequesterID),
				},
			},
		},
		Owner: &datatypes.Reference{
			Identifier: &datatypes.Identifier{
				System: &fhir.NutsCodingSystem,
				Value:  fhir.ToStringPtr(props.OwnerID),
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

func BuildAdvanceNotice(createRequest CreateTransferRequest) eoverdracht.AdvanceNotice {
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

func buildAdministrativeData(request CreateTransferRequest) eoverdracht.CompositionSection {
	transferDate := request.TransferDate.Format(time.RFC3339)
	return eoverdracht.CompositionSection{
		BackboneElement: datatypes.BackboneElement{
			Element: datatypes.Element{
				Extension: []datatypes.Extension{{
					URL:           (*datatypes.URI)(fhir.ToStringPtr("http://nictiz.nl/fhir/StructureDefinition/eOverdracht-TransferDate")),
					ValueDateTime: (*datatypes.DateTime)(fhir.ToStringPtr(transferDate)),
				}},
			},
		},
		Title: fhir.ToStringPtr("Administrative data"),
		Code: datatypes.CodeableConcept{
			Coding: []datatypes.Coding{{
				System:  &fhir.SnomedCodingSystem,
				Code:    fhir.ToCodePtr(eoverdracht.AdministrativeDocCode),
				Display: fhir.ToStringPtr("Administrative documentation (record artifact)"),
			}}},
	}

}

func buildNursingHandoffComposition(administrativeData, careplan eoverdracht.CompositionSection, patient resources.Patient) eoverdracht.Composition {
	return eoverdracht.Composition{
		Base: resources.Base{
			ResourceType: "Composition",
			ID:           fhir.ToIDPtr(generateResourceID()),
		},
		Type: datatypes.CodeableConcept{
			Coding: []datatypes.Coding{{System: &fhir.SnomedCodingSystem, Code: fhir.ToCodePtr("371535009"), Display: fhir.ToStringPtr("verslag van overdracht")}},
		},
		Subject: datatypes.Reference{Reference: fhir.ToStringPtr("Patient/" + fhir.FromIDPtr(patient.ID))},
		Title:   "Nursing handoff",
		Section: []eoverdracht.CompositionSection{administrativeData, careplan},
	}
}

func buildAdvanceNoticeComposition(administrativeData, careplan eoverdracht.CompositionSection) eoverdracht.Composition {

	return eoverdracht.Composition{
		Base: resources.Base{
			ResourceType: "Composition",
			ID:           fhir.ToIDPtr(generateResourceID()),
		},
		Type: datatypes.CodeableConcept{
			Coding: []datatypes.Coding{{System: &fhir.LoincCodingSystem, Code: fhir.ToCodePtr("57830-2")}},
		},
		Title:   "Advance notice",
		Section: []eoverdracht.CompositionSection{administrativeData, careplan},
	}
}

func buildCarePlan(carePlan CarePlan) (problems []resources.Condition, interventions []eoverdracht.Procedure, section eoverdracht.CompositionSection) {
	for _, cpPatientProblems := range carePlan.PatientProblems {
		newProblem := buildConditionFromProblem(cpPatientProblems.Problem)
		problems = append(problems, newProblem)
		for _, i := range cpPatientProblems.Interventions {
			interventions = append(interventions, buildProcedureFromIntervention(i, fhir.FromIDPtr(newProblem.ID)))
		}
	}

	// new patientProblems
	patientProblems := eoverdracht.CompositionSection{
		Title: fhir.ToStringPtr("Current patient problems"),
		Code: datatypes.CodeableConcept{
			Coding: []datatypes.Coding{{
				System:  &fhir.SnomedCodingSystem,
				Code:    fhir.ToCodePtr("86644006"),
				Display: fhir.ToStringPtr("Nursing diagnosis"),
			}},
		},
	}

	// Add the problems as a section
	for _, p := range problems {
		patientProblems.Entry = append(patientProblems.Entry, datatypes.Reference{Reference: fhir.ToStringPtr("Condition/" + fhir.FromIDPtr(p.ID))})
	}
	for _, i := range interventions {
		patientProblems.Entry = append(patientProblems.Entry, datatypes.Reference{Reference: fhir.ToStringPtr("Procedure/" + fhir.FromIDPtr(i.ID))})
	}

	// Start with empty care plan
	careplan := eoverdracht.CompositionSection{
		Code: datatypes.CodeableConcept{
			Coding: []datatypes.Coding{{
				System:  &fhir.SnomedCodingSystem,
				Code:    fhir.ToCodePtr(eoverdracht.CarePlanCode),
				Display: fhir.ToStringPtr("Nursing care plan (record artifact)"),
			}}},
		Section: []eoverdracht.CompositionSection{
			patientProblems,
		},
	}
	return problems, interventions, careplan
}

func buildProcedureFromIntervention(intervention Intervention, problemID string) eoverdracht.Procedure {
	return eoverdracht.Procedure{
		Domain: resources.Domain{
			Base: resources.Base{
				ResourceType: "Procedure",
				ID:           fhir.ToIDPtr(generateResourceID()),
			},
		},
		ReasonReference: []datatypes.Reference{{Reference: fhir.ToStringPtr("Condition/" + problemID)}},
		Note:            []datatypes.Annotation{{Text: fhir.ToStringPtr(intervention.Comment)}},
	}
}

func buildConditionFromProblem(problem Problem) resources.Condition {
	return resources.Condition{
		Domain: resources.Domain{
			Base: resources.Base{
				ResourceType: "Condition",
				ID:           fhir.ToIDPtr(generateResourceID()),
			},
		},
		Note: []datatypes.Annotation{{Text: fhir.ToStringPtr(problem.Name)}},
	}
}

func BuildNursingHandoffComposition(patient *Patient, advanceNotice eoverdracht.AdvanceNotice) (eoverdracht.Composition, error) {

	careplan, err := eoverdracht.FilterCompositionSectionByType(advanceNotice.Composition.Section, eoverdracht.CarePlanCode)
	if err != nil {
		return eoverdracht.Composition{}, err
	}

	administrativeData, err := eoverdracht.FilterCompositionSectionByType(advanceNotice.Composition.Section, eoverdracht.AdministrativeDocCode)
	if err != nil {
		return eoverdracht.Composition{}, err
	}

	fhirPatient := resources.Patient{Domain: resources.Domain{Base: resources.Base{ID: fhir.ToIDPtr(string(patient.ObjectID))}}}

	return buildNursingHandoffComposition(administrativeData, careplan, fhirPatient), nil
}

func generateResourceID() string {
	return uuid.New().String()
}
