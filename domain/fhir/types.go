package fhir

import (
	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
)

/* Coding systems */
var (
	SnomedCodingSystem datatypes.URI = "http://snomed.info/sct"
	LoincCodingSystem  datatypes.URI = "http://loinc.org"
	NutsCodingSystem   datatypes.URI = "http://nuts.nl"
	UZICodingSystem    datatypes.URI = "http://fhir.nl/fhir/NamingSystem/uzi-nr-pers"
)

type TaskProperties struct {
	ID        *string
	Status    string
	PatientID string
	// nuts DID of the placer
	RequesterID string
	// nuts DID of the filler
	OwnerID string
	Input   []resources.TaskInputOutput
	Output  []resources.TaskInputOutput
}

// Procedure defines a basic FHIR STU3 Procedure resource which is currently not included in the FHIR library.
type Procedure struct {
	resources.Domain
	Identifier      []datatypes.Identifier `json:"identifier,omitempty"`
	Code            datatypes.Code         `json:"code,omitempty"`
	Subject         datatypes.Reference    `json:"subject,omitempty"`
	ReasonReference []datatypes.Reference  `json:"reasonReference,omitempty"`
	Note            []datatypes.Annotation `json:"note,omitempty"`
}

type CompositionSection struct {
	datatypes.BackboneElement
	Code    datatypes.CodeableConcept `json:"code"`
	Title   *datatypes.String         `json:"title,omitempty"`
	Section []CompositionSection      `json:"section,omitempty"`
	Entry   []datatypes.Reference     `json:"entry,omitempty"`
}

// Composition defines a basic FHIR STU3 Composition resource which is currently not included in the FHIR library.
type Composition struct {
	resources.Base
	Identifier []datatypes.Identifier    `json:"identifier,omitempty"`
	Type       datatypes.CodeableConcept `json:"type"`
	Status     datatypes.Code            `json:"status,omitempty"`
	Subject    datatypes.Reference       `json:"subject"`
	Date       datatypes.DateTime        `json:"date,omitempty"`
	Author     []datatypes.Reference     `json:"author,omitempty"`
	Title      datatypes.String          `json:"title,omitempty"`
	Section    []CompositionSection      `json:"section,omitempty"`
}
