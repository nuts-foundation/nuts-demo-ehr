package types

import (
	"github.com/nuts-foundation/nuts-demo-ehr/nuts"
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

func FromNutsOrganization(src nuts.NutsOrganization) Organization {
	return Organization{
		Did:  src.ID,
		Name: src.Details.Name,
		City: src.Details.City,
	}
}
