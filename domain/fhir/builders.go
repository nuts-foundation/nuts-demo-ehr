package fhir

import (
	"github.com/google/uuid"
	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
	"time"
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
		// TODO: patient seems mandatory in the spec, but can only be sent when placer already
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

func generateResourceID() string {
	return uuid.New().String()
}
