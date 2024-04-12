package types

import (
	"github.com/nuts-foundation/nuts-demo-ehr/nuts"
	r4 "github.com/samply/golang-fhir-models/fhir-models/fhir"
	"time"
)

const DobFormat = "2006-01-02"
const BsnSystem = "http://fhir.nl/fhir/NamingSystem/bsn"

type IncomingTransfer struct {
	Id         ObjectID                  `json:"id"`
	FhirTaskID string                    `json:"fhirTaskID"`
	Sender     Organization              `json:"sender"`
	Status     TransferNegotiationStatus `json:"status"`
	CreatedAt  time.Time                 `json:"createdAt"`
}

// SharedCarePlan is a HL7 FHIR R4 CarePlan
type SharedCarePlan struct {
	DossierID         string             `json:"dossierID"`
	FHIRCarePlan      r4.CarePlan        `json:"fhirCarePlan"`
	FHIRCarePlanURL   string             `json:"fhirCarePlanURL"`
	FHIRActivityTasks map[string]r4.Task `json:"fhirActivityTasks"`
	Participants      []Organization     `json:"participants"`
}

type SharedCarePlanNotifyRequest struct {
	CarePlanURL string     `json:"carePlanURL"`
	Patient     r4.Patient `json:"patient"`
	Task        r4.Task    `json:"task"`
}

func FromNutsOrganization(src nuts.NutsOrganization) Organization {
	return Organization{
		Did:         src.ID,
		Name:        src.Details.Name,
		City:        src.Details.City,
		Identifiers: map[string]string{},
	}
}

type FHIRCodeableConcept = r4.CodeableConcept

type FHIRIdentifier = r4.Identifier
