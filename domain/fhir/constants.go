package fhir

/* Coding systems */
type CodingSystem string

const (
	SnomedCodingSystem CodingSystem = "http://snomed.info/sct"
	LoincCodingSystem               = "http://loinc.org"
	NutsCodingSystem                = "http://nuts.nl"
)

/* Codes */
type Code string

const (
	SnomedTransferCode        Code = "308292007"
	TransferDisplay                = "Overdracht van zorg"
	LoincAdvanceNoticeCode         = "57830-2"
	SnomedAlternaticeDateCode      = "146851000146105"
	SnomedNursingHandoffCode       = "371535009"
)

/* Short-hand types */
var LoincAdvanceNoticeType = CodeableConcept{
	Coding: Coding{
		System: LoincCodingSystem,
		Code:   LoincAdvanceNoticeCode,
	},
	Text: "Aanmeldbericht",
}
