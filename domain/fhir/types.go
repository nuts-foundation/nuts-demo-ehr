package fhir

// Reference models http://hl7.org/fhir/STU3/references.html#Reference
type Reference struct {
	Reference interface{} `json:"reference"`
}

// CodeableConcept models http://hl7.org/fhir/STU3/datatypes.html#CodeableConcept
type CodeableConcept struct {
	Coding Coding `json:"coding"`
	Text   string `json:"text"`
}

// Coding models http://hl7.org/fhir/STU3/datatypes.html#Coding
type Coding struct {
	System  CodingSystem `json:"system"`
	Code    Code         `json:"code"`
	Display string       `json:"display"`
}

// Device models http://hl7.org/fhir/STU3/device.html
type Device struct {
	Manufacturer string `json:"manufacturer"`
}

type Requester struct {
	Agent interface{} `json:"agent"`
}

// Organization models http://hl7.org/fhir/STU3/organization.html
type Organization struct {
	Identifier Identifier `json:"identifier"`
}

// Identifier models http://hl7.org/fhir/STU3/datatypes.html#Identifier
type Identifier struct {
	System CodingSystem `json:"system"`
	Value  string       `json:"value"`
}

// TaskInputOutput is used as container type for input/output elements of tasks
type TaskInputOutput struct {
	Type           CodeableConcept `json:"type"`
	ValueReference Reference       `json:"valueReference"`
}

type TaskProperties struct {
	Status    string
	PatientID string
	// nuts DID of the placer
	RequesterID string
	// nuts DID of the filler
	OwnerID string
	Input   []TaskInputOutput
	Output  []TaskInputOutput
}

type Task struct {
	ID string
	TaskProperties
}

type Composition struct {
	ID        string
	Reference string
}
