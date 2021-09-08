package domain

type IncomingTransfer struct {
	Id         ObjectID     `json:"id"`
	FhirTaskID string       `json:"fhirTaskID"`
	Sender     Organization `json:"sender"`
}
