package domain

import "time"

type IncomingTransfer struct {
	Id         ObjectID     `json:"id"`
	FhirTaskID string       `json:"fhirTaskID"`
	Sender     Organization `json:"sender"`
}

// Transfer data types as described by the ART Decor app from the Nictiz
// https://decor.nictiz.nl/pub/eoverdracht/e-overdracht-html-20210510T093529/tr-2.16.840.1.113883.2.4.3.11.60.30.4.63-2021-01-27T000000.html

type AdministrativeData struct {
	TransferDate time.Time
	PersonalInformation
}

type PersonalInformation struct {
	AnonymousPatient
}

type AnonymousPatient struct {
	PostalCode string
}

//type CarePlan struct {
//	Problems []Problem
//}

//// Problem a s defined by https://zibs.nl/wiki/Probleem-v4.1(2017NL)
//type Problem struct {
//	Explanation string // NL-CM:5.1.5
//	Active      bool   // NL-CM:5.1.4
//	// Code from the Dutch Snomed CT Kernset patientproblemen
//	// https://www.nictiz.nl/standaardisatie/terminologiecentrum/referentielijsten/nationale-kernset/nationale-kernset-patientproblemen/
//	Code        string // NL-CM:5.1.3
//	Intervention NursingIntervention // Intervention to perform to treat this problem
//}

// NursingIntervention as defined by https://zibs.nl/wiki/NursingIntervention-v3.1(2017EN)
type NursingIntervention struct {
	Comment string // NL-CM:14.2.5 Comment on the nursing intervention
}
