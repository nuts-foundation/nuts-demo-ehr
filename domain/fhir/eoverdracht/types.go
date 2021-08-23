package eoverdracht

import "github.com/monarko/fhirgo/STU3/datatypes"

// Practitioner models https://simplifier.net/packages/nictiz.fhir.nl.stu3.zib2017/2.1.1/files/361872
type Practitioner struct {
	datatypes.Element
	Identifier datatypes.Identifier `json:"identifier"`
	Name       *datatypes.HumanName `json:"name,omitempty"`
}
