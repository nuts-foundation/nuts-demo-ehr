package fhir

import "github.com/monarko/fhirgo/STU3/datatypes"

/* Coding systems */

var (
	SnomedCodingSystem datatypes.URI = "http://snomed.info/sct"
	LoincCodingSystem  datatypes.URI = "http://loinc.org"
	NutsCodingSystem   datatypes.URI = "http://nuts.nl"
	UZICodingSystem    datatypes.URI = "http://fhir.nl/fhir/NamingSystem/uzi-nr-pers"
)

/* Codes */
type Code string

var (
	SnomedTransferCode        datatypes.Code = "308292007"
	LoincAdvanceNoticeCode    datatypes.Code = "57830-2"
	SnomedAlternaticeDateCode datatypes.Code = "146851000146105"
	SnomedNursingHandoffCode  datatypes.Code = "371535009"
)

var NursingHandoffDisplay datatypes.String = "verslag van zorg"
var TransferDisplay datatypes.String = "Overdracht van zorg"

/* Short-hand types */
var LoincAdvanceNoticeType = datatypes.CodeableConcept{
	Coding: []datatypes.Coding{{
		System: &LoincCodingSystem,
		Code:   &LoincAdvanceNoticeCode,
	}},
	Text: ToStringPtr("Aanmeldbericht"),
}

var SnomedNursingHandoffType = datatypes.CodeableConcept{
	Coding: []datatypes.Coding{{
		System:  &SnomedCodingSystem,
		Code:    &SnomedNursingHandoffCode,
		Display: &NursingHandoffDisplay,
	}},
}
