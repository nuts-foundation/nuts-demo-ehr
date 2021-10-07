package eoverdracht

import (
	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
)

const AdministrativeDocCode = "405624007"
const CarePlanCode = "773130005"
const NursingDiagnosisCode = "86644006"

type TransferTask struct {
	ID               string
	Status           string
	AdvanceNoticeID  *string
	NursingHandoffID *string
}

// Practitioner models https://simplifier.net/packages/nictiz.fhir.nl.stu3.zib2017/2.1.1/files/361872
type Practitioner struct {
	datatypes.Element
	Identifier datatypes.Identifier `json:"identifier"`
	Name       *datatypes.HumanName `json:"name,omitempty"`
}

// AdvanceNotice is a container to hold all FHIR resources associated with a Transfer advance notice
type AdvanceNotice struct {
	Composition   fhir.Composition
	Patient       resources.Patient
	Problems      []resources.Condition
	Interventions []fhir.Procedure
}

// NursingHandoff is a container to hold all FHIR resources associated with a Transfers Nursing Handoff.
// Note: Currently it contains exactly the same content as a advance notice, but this can be extended.
type NursingHandoff struct {
	Composition   fhir.Composition
	Patient       resources.Patient
	Problems      []resources.Condition
	Interventions []fhir.Procedure
}
