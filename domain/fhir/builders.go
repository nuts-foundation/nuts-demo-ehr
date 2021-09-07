package fhir

import (
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
	elements := map[string]interface{}{
		"title": "Nursing handoff",
		"type":  SnomedNursingHandoffType,
		"subject": datatypes.Reference{Reference: ToStringPtr(string(patient.ObjectID))},
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
