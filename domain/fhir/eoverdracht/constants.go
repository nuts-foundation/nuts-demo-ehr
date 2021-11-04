package eoverdracht

import (
	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
)

var (
	SnomedTransferCode        datatypes.Code   = "308292007"
	LoincAdvanceNoticeCode    datatypes.Code   = "57830-2"
	SnomedAlternaticeDateCode datatypes.Code   = "146851000146105"
	SnomedNursingHandoffCode  datatypes.Code   = "371535009"
	NursingHandoffDisplay     datatypes.String = "verslag van zorg"
	TransferDisplay           datatypes.String = "Overdracht van zorg"
)

const (
	AdministrativeDocCode = "405624007"
	CarePlanCode          = "773130005"
	NursingDiagnosisCode  = "86644006"
)

/* Short-hand types */
var AdministrativeDocConcept = datatypes.CodeableConcept{
	Coding: []datatypes.Coding{{
		System:  &fhir.SnomedCodingSystem,
		Code:    fhir.ToCodePtr(AdministrativeDocCode),
		Display: fhir.ToStringPtr("Administrative documentation (record artifact)"),
	}}}

var LoincAdvanceNoticeType = datatypes.CodeableConcept{
	Coding: []datatypes.Coding{{
		System: &fhir.LoincCodingSystem,
		Code:   &LoincAdvanceNoticeCode,
	}},
	Text: fhir.ToStringPtr("Aanmeldbericht"),
}

var CarePlanConcept = datatypes.CodeableConcept{
	Coding: []datatypes.Coding{{
		System:  &fhir.SnomedCodingSystem,
		Code:    fhir.ToCodePtr(CarePlanCode),
		Display: fhir.ToStringPtr("Nursing care plan (record artifact)"),
	}},
}

var SnomedNursingHandoffType = datatypes.CodeableConcept{
	Coding: []datatypes.Coding{{
		System:  &fhir.SnomedCodingSystem,
		Code:    &SnomedNursingHandoffCode,
		Display: &NursingHandoffDisplay,
	}},
}

var SnomedTransferType = datatypes.CodeableConcept{
	Coding: []datatypes.Coding{{
		System:  &fhir.SnomedCodingSystem,
		Code:    &SnomedTransferCode,
		Display: &TransferDisplay,
	}},
}
