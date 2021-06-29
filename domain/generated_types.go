// Package domain provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.7.1 DO NOT EDIT.
package domain

import (
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
)

// Defines values for PatientPropertiesGender.
const (
	PatientPropertiesGenderFemale PatientPropertiesGender = "female"

	PatientPropertiesGenderMale PatientPropertiesGender = "male"

	PatientPropertiesGenderOther PatientPropertiesGender = "other"

	PatientPropertiesGenderUnknown PatientPropertiesGender = "unknown"
)

// A customer object.
type Customer struct {

	// If a VC has been issued for this customer.
	Active bool `json:"active"`

	// Locality for this customer.
	City *string `json:"city,omitempty"`

	// The customer DID.
	Did *string `json:"did,omitempty"`

	// The email domain of the care providers employees, required for logging in.
	Domain *string `json:"domain,omitempty"`

	// The internal customer ID.
	Id string `json:"id"`

	// Internal name for this customer.
	Name string `json:"name"`
}

// IRMAAuthenticationRequest defines model for IRMAAuthenticationRequest.
type IRMAAuthenticationRequest struct {

	// Internal ID of the customer for which is being logged in
	CustomerID string `json:"customerID"`
}

// PasswordAuthenticateRequest defines model for PasswordAuthenticateRequest.
type PasswordAuthenticateRequest struct {

	// Internal ID of the customer for which is being logged in
	CustomerID string `json:"customerID"`
	Password   string `json:"password"`
}

// Patient defines model for Patient.
type Patient struct {
	// Embedded struct due to allOf(#/components/schemas/PatientID)
	PatientID `yaml:",inline"`
	// Embedded struct due to allOf(#/components/schemas/PatientProperties)
	PatientProperties `yaml:",inline"`
}

// The ID of the patient. E.g. UUID or database increment.
type PatientID string

// A patient in the EHR system. Containing the basic information about the like name, adress, dob etc.
type PatientProperties struct {

	// Date of birth. Can include time if known.
	Dob *openapi_types.Date `json:"dob,omitempty"`

	// Primary email address.
	Email *openapi_types.Email `json:"email,omitempty"`

	// Given name
	FirstName string `json:"firstName"`

	// Gender of the person according to https://www.hl7.org/fhir/valueset-administrative-gender.html.
	Gender PatientPropertiesGender `json:"gender"`

	// The internal ID of the Patient. Can be any internal system. Not to be confused by a database ID or a uuid.
	InternalID string `json:"internalID"`

	// Family name. Must include prefixes like "van der".
	Surname string `json:"surname"`

	// The zipcode formatted in dutch form. Can be used to find local care providers.
	Zipcode string `json:"zipcode"`
}

// Gender of the person according to https://www.hl7.org/fhir/valueset-administrative-gender.html.
type PatientPropertiesGender string

// Result of a signing session.
type SessionToken struct {

	// the result from a signing session. It's a base64 encoded Verifiable Presentation
	Token string `json:"token"`
}

// AuthenticateWithIRMAJSONBody defines parameters for AuthenticateWithIRMA.
type AuthenticateWithIRMAJSONBody IRMAAuthenticationRequest

// AuthenticateWithPasswordJSONBody defines parameters for AuthenticateWithPassword.
type AuthenticateWithPasswordJSONBody PasswordAuthenticateRequest

// NewPatientJSONBody defines parameters for NewPatient.
type NewPatientJSONBody PatientProperties

// AuthenticateWithIRMAJSONRequestBody defines body for AuthenticateWithIRMA for application/json ContentType.
type AuthenticateWithIRMAJSONRequestBody AuthenticateWithIRMAJSONBody

// AuthenticateWithPasswordJSONRequestBody defines body for AuthenticateWithPassword for application/json ContentType.
type AuthenticateWithPasswordJSONRequestBody AuthenticateWithPasswordJSONBody

// NewPatientJSONRequestBody defines body for NewPatient for application/json ContentType.
type NewPatientJSONRequestBody NewPatientJSONBody
