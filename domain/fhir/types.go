package fhir

import (
	"github.com/monarko/fhirgo/STU3/resources"
)

type TaskProperties struct {
	Status    string
	PatientID string
	// nuts DID of the placer
	RequesterID string
	// nuts DID of the filler
	OwnerID string
	Input   []resources.TaskInputOutput
	Output  []resources.TaskInputOutput
}

type Task struct {
	ID string
	TaskProperties
}

type Composition map[string]interface{}
