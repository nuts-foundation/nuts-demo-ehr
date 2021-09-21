package eoverdracht

import (
	"fmt"

	"github.com/monarko/fhirgo/STU3/datatypes"
)

func FilterCompositionSectionByType(sections []CompositionSection, typeFilter string) (CompositionSection, error)  {
	for _, c := range sections {
		if *c.Code.Coding[0].Code == datatypes.Code(typeFilter) {
			return c, nil
		}
	}
	return CompositionSection{}, fmt.Errorf("composition not found")
}