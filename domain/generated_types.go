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

// Defines values for TransferStatus.
const (
	TransferStatusAccepted TransferStatus = "accepted"

	TransferStatusCancelled TransferStatus = "cancelled"

	TransferStatusCompleted TransferStatus = "completed"

	TransferStatusCreated TransferStatus = "created"

	TransferStatusRequested TransferStatus = "requested"
)

// API request to accept the transfer of a patient, to a care organization that accepted transfer. negotiationID contains the ID of the negotiation that will complete the transfer.
type AcceptTransferRequest struct {

	// An internal object UUID which can be used as unique identifier for entities.
	NegotiationID ObjectID `json:"negotiationID"`
}

// CreateTransferRequest defines model for CreateTransferRequest.
type CreateTransferRequest struct {
	Description  string             `json:"description"`
	TransferDate openapi_types.Date `json:"transferDate"`
}

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

// An internal object UUID which can be used as unique identifier for entities.
type ObjectID string

// A care organization available through the Nuts Network to exchange information.
type Organization struct {

	// City where the care organization is located.
	City string `json:"city"`

	// Decentralized Identifier which uniquely identifies the care organization on the Nuts Network.
	Did string `json:"did"`

	// Name of the care organization.
	Name string `json:"name"`
}

// PasswordAuthenticateRequest defines model for PasswordAuthenticateRequest.
type PasswordAuthenticateRequest struct {

	// Internal ID of the customer for which is being logged in
	CustomerID string `json:"customerID"`
	Password   string `json:"password"`
}

// Patient defines model for Patient.
type Patient struct {
	// Embedded struct due to allOf(#/components/schemas/ObjectID)
	ObjectID `yaml:",inline"`
	// Embedded struct due to allOf(#/components/schemas/PatientProperties)
	PatientProperties `yaml:",inline"`
}

// A patient in the EHR system. Containing the basic information about the like name, adress, dob etc.
type PatientProperties struct {

	// Date of birth.
	Dob *openapi_types.Date `json:"dob,omitempty"`

	// Primary email address.
	Email *openapi_types.Email `json:"email,omitempty"`

	// Given name
	FirstName string `json:"firstName"`

	// Gender of the person according to https://www.hl7.org/fhir/valueset-administrative-gender.html.
	Gender PatientPropertiesGender `json:"gender"`

	// Social security number
	Ssn *string `json:"ssn,omitempty"`

	// Family name. Must include prefixes like "van der".
	Surname string `json:"surname"`

	// The zipcode formatted in dutch form. Can be used to find local care providers.
	Zipcode string `json:"zipcode"`
}

// Gender of the person according to https://www.hl7.org/fhir/valueset-administrative-gender.html.
type PatientPropertiesGender string

// Result of a signing session.
type SessionToken struct {

	// the result from a signing session. It's an updated JWT.
	Token string `json:"token"`
}

// A dossier for transferring a patient to another care organization. It is composed of negotiations with specific care organizations. The patient can be transferred to one of the care organizations that accepted the transfer. TODO: proposing/handling alternative transfer dates is not supported yet.
type Transfer struct {

	// Accompanying text sent to care organizations to assess the patient transfer. It is populated/updated by the last negotiation that was started.
	Description string `json:"description"`

	// An internal object UUID which can be used as unique identifier for entities.
	Id ObjectID `json:"id"`

	// Status of the transfer. If the state is "completed" or "cancelled" the transfer dossier becomes read-only. In that case no additional negotiations can be sent (for this transfer) or accepted. Possible values: - Created: the new transfer dossier is created, but no requests were sent (to receiving care organizations) yet. - Requested: one or more requests were sent to care organizations - Accepted: one of the requests, accepted by the receiving care organizations is accepted by the sending care organization. - Completed: the patient transfer is completed and marked as such by the receiving care organization. - Cancelled: the transfer is cancelled by the sending care organization.
	Status TransferStatus `json:"status"`

	// Transfer date as proposed by the sending XIS. It is populated/updated by the last negotiation that was started.
	TransferDate openapi_types.Date `json:"transferDate"`
}

// Status of the transfer. If the state is "completed" or "cancelled" the transfer dossier becomes read-only. In that case no additional negotiations can be sent (for this transfer) or accepted. Possible values: - Created: the new transfer dossier is created, but no requests were sent (to receiving care organizations) yet. - Requested: one or more requests were sent to care organizations - Accepted: one of the requests, accepted by the receiving care organizations is accepted by the sending care organization. - Completed: the patient transfer is completed and marked as such by the receiving care organization. - Cancelled: the transfer is cancelled by the sending care organization.
type TransferStatus string

// A negotiation with a specific care organization to transfer a patient.
type TransferNegotiation struct {

	// An internal object UUID which can be used as unique identifier for entities.
	Id ObjectID `json:"id"`

	// Decentralized Identifier of the organization to which transfer of a patient is requested.
	OrganizationDID string `json:"organizationDID"`

	// Status of the negotiation, maps to FHIR eOverdracht task states (https://informatiestandaarden.nictiz.nl/wiki/vpk:V4.0_FHIR_eOverdracht#Using_Task_to_manage_the_workflow).
	Status string `json:"status"`

	// Transfer date subject of the negotiation. Can be altered by both sending and receiving care organization.
	TransferDate openapi_types.Date `json:"transferDate"`
}

// API request to start a negotiation with a specific care organization for transferring a patient.
type TransferNegotiationRequest struct {

	// Decentralized Identifier of the organization to which transfer of a patient is requested.
	OrganizationDID string `json:"organizationDID"`
}

// SetCustomerJSONBody defines parameters for SetCustomer.
type SetCustomerJSONBody Customer

// AuthenticateWithIRMAJSONBody defines parameters for AuthenticateWithIRMA.
type AuthenticateWithIRMAJSONBody IRMAAuthenticationRequest

// AuthenticateWithPasswordJSONBody defines parameters for AuthenticateWithPassword.
type AuthenticateWithPasswordJSONBody PasswordAuthenticateRequest

// SearchOrganizationsParams defines parameters for SearchOrganizations.
type SearchOrganizationsParams struct {

	// Keyword for finding care organizations.
	Query string `json:"query"`

	// Filters care organizations on service, only returning care organizations have a service in their DID Document which' type matches the given didServiceType. If not supplied, care organizations aren't filtered on service.
	DidServiceType *string `json:"didServiceType,omitempty"`
}

// UpdatePatientJSONBody defines parameters for UpdatePatient.
type UpdatePatientJSONBody PatientProperties

// CreateTransferJSONBody defines parameters for CreateTransfer.
type CreateTransferJSONBody CreateTransferRequest

// StartTransferNegotiationJSONBody defines parameters for StartTransferNegotiation.
type StartTransferNegotiationJSONBody TransferNegotiationRequest

// AcceptTransferNegotiationJSONBody defines parameters for AcceptTransferNegotiation.
type AcceptTransferNegotiationJSONBody AcceptTransferRequest

// NewPatientJSONBody defines parameters for NewPatient.
type NewPatientJSONBody PatientProperties

// SetCustomerJSONRequestBody defines body for SetCustomer for application/json ContentType.
type SetCustomerJSONRequestBody SetCustomerJSONBody

// AuthenticateWithIRMAJSONRequestBody defines body for AuthenticateWithIRMA for application/json ContentType.
type AuthenticateWithIRMAJSONRequestBody AuthenticateWithIRMAJSONBody

// AuthenticateWithPasswordJSONRequestBody defines body for AuthenticateWithPassword for application/json ContentType.
type AuthenticateWithPasswordJSONRequestBody AuthenticateWithPasswordJSONBody

// UpdatePatientJSONRequestBody defines body for UpdatePatient for application/json ContentType.
type UpdatePatientJSONRequestBody UpdatePatientJSONBody

// CreateTransferJSONRequestBody defines body for CreateTransfer for application/json ContentType.
type CreateTransferJSONRequestBody CreateTransferJSONBody

// StartTransferNegotiationJSONRequestBody defines body for StartTransferNegotiation for application/json ContentType.
type StartTransferNegotiationJSONRequestBody StartTransferNegotiationJSONBody

// AcceptTransferNegotiationJSONRequestBody defines body for AcceptTransferNegotiation for application/json ContentType.
type AcceptTransferNegotiationJSONRequestBody AcceptTransferNegotiationJSONBody

// NewPatientJSONRequestBody defines body for NewPatient for application/json ContentType.
type NewPatientJSONRequestBody NewPatientJSONBody
