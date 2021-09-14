package eoverdracht

import (
	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
)

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
	Code            string
	Subject         datatypes.Reference
	ReasonReference datatypes.Reference
	Note            []datatypes.Annotation
}

type CompositionSection struct {
	datatypes.BackboneElement
	Code    datatypes.CodeableConcept
	Title   *datatypes.String
	Section []CompositionSection
	Entry   []datatypes.Reference
}

// Composition defines a basic FHIR STU3 Composition resource which is currently not included in the FHIR library.
type Composition struct {
	resources.Base
	Identifier []datatypes.Identifier `json:"identifier,omitempty"`
	Type       datatypes.CodeableConcept
	Status     datatypes.Code
	Subject    datatypes.Reference
	Date       datatypes.DateTime
	Author     datatypes.Reference
	Title      datatypes.String
	Section    []CompositionSection
}

type AdministrativeData struct{}

type GeneralPatientContext struct{}

type MedicalContext struct{}

type CarePlan struct{}

type HealthState struct {
}
