package eoverdracht

import (
	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
)

const AdministrativeDocCode = "405624007"
const CarePlanCode = "773130005"
const NursingDiagnosisCode = "86644006"

// Practitioner models https://simplifier.net/packages/nictiz.fhir.nl.stu3.zib2017/2.1.1/files/361872
type Practitioner struct {
	datatypes.Element
	Identifier datatypes.Identifier `json:"identifier"`
	Name       *datatypes.HumanName `json:"name,omitempty"`
}

// AdvanceNotice is a container to hold all FHIR resources associated with a Transfer advance notice
type AdvanceNotice struct {
	Composition   Composition
	Patient       resources.Patient
	Problems      []resources.Condition
	Interventions []Procedure
}

// Procedure defines a basic FHIR STU3 Procedure resource which is currently not included in the FHIR library.
type Procedure struct {
	resources.Domain
	Identifier      []datatypes.Identifier `json:"identifier,omitempty"`
	Code            datatypes.Code         `json:"code,omitempty"`
	Subject         datatypes.Reference    `json:"subject,omitempty"`
	ReasonReference []datatypes.Reference    `json:"reasonReference,omitempty"`
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


type AdministrativeData struct{}

type GeneralPatientContext struct{}

type MedicalContext struct{}

type CarePlan struct{}

type HealthState struct {
}
