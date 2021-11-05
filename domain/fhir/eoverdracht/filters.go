package eoverdracht

import (
	"fmt"

	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
)

func FilterCompositionSectionByType(sections []fhir.CompositionSection, typeFilter string) (fhir.CompositionSection, error) {
	for _, c := range sections {
		if *c.Code.Coding[0].Code == datatypes.Code(typeFilter) {
			return c, nil
		}
	}
	return fhir.CompositionSection{}, fmt.Errorf("composition not found")
}
