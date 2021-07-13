// Package domain provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.7.1 DO NOT EDIT.
package domain

import (
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
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
	TransferStatusAssigned TransferStatus = "assigned"

	TransferStatusCancelled TransferStatus = "cancelled"

	TransferStatusCompleted TransferStatus = "completed"

	TransferStatusCreated TransferStatus = "created"

	TransferStatusRequested TransferStatus = "requested"
)

// Defines values for TransferNegotiationStatus.
const (
	TransferNegotiationStatusAccepted TransferNegotiationStatus = "accepted"

	TransferNegotiationStatusCompleted TransferNegotiationStatus = "completed"

	TransferNegotiationStatusInProgress TransferNegotiationStatus = "in-progress"

	TransferNegotiationStatusOnHold TransferNegotiationStatus = "on-hold"

	TransferNegotiationStatusRequested TransferNegotiationStatus = "requested"
)

// API request to accept the transfer of a patient, to a care organization that accepted transfer. negotiationID contains the ID of the negotiation that will complete the transfer.
type AcceptTransferRequest struct {

	// An internal object UUID which can be used as unique identifier for entities.
	NegotiationID ObjectID `json:"negotiationID"`
}

// API request to create a dossier for a patient.
type CreateDossierRequest struct {
	Name string `json:"name"`

	// An internal object UUID which can be used as unique identifier for entities.
	PatientID ObjectID `json:"patientID"`
}

// Create a new transfer for a specific dossier with a date and description.
type CreateTransferRequest struct {
	Description string `json:"description"`

	// An internal object UUID which can be used as unique identifier for entities.
	DossierID    ObjectID           `json:"dossierID"`
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

// Dossier defines model for Dossier.
type Dossier struct {

	// An internal object UUID which can be used as unique identifier for entities.
	Id   ObjectID `json:"id"`
	Name string   `json:"name"`

	// An internal object UUID which can be used as unique identifier for entities.
	PatientID ObjectID `json:"patientID"`
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

	// Status of the transfer. If the state is "completed" or "cancelled" the transfer dossier becomes read-only. In that case no additional negotiations can be sent (for this transfer) or accepted. Possible values: - Created: the new transfer dossier is created, but no requests were sent (to receiving care organizations) yet. - Requested: one or more requests were sent to care organizations - Assigned: The transfer is assigned to one the receiving care organizations thet accepted the transfer. - Completed: the patient transfer is completed and marked as such by the receiving care organization. - Cancelled: the transfer is cancelled by the sending care organization.
	Status TransferStatus `json:"status"`

	// Transfer date as proposed by the sending XIS. It is populated/updated by the last negotiation that was started.
	TransferDate openapi_types.Date `json:"transferDate"`
}

// Status of the transfer. If the state is "completed" or "cancelled" the transfer dossier becomes read-only. In that case no additional negotiations can be sent (for this transfer) or accepted. Possible values: - Created: the new transfer dossier is created, but no requests were sent (to receiving care organizations) yet. - Requested: one or more requests were sent to care organizations - Assigned: The transfer is assigned to one the receiving care organizations thet accepted the transfer. - Completed: the patient transfer is completed and marked as such by the receiving care organization. - Cancelled: the transfer is cancelled by the sending care organization.
type TransferStatus string

// A negotiation with a specific care organization to transfer a patient.
type TransferNegotiation struct {

	// Decentralized Identifier of the organization to which transfer of a patient is requested.
	OrganizationDID string `json:"organizationDID"`

	// Status of the negotiation, maps to FHIR eOverdracht task states (https://informatiestandaarden.nictiz.nl/wiki/vpk:V4.0_FHIR_eOverdracht#Using_Task_to_manage_the_workflow).
	Status TransferNegotiationStatus `json:"status"`

	// Transfer date subject of the negotiation. Can be altered by both sending and receiving care organization.
	TransferDate openapi_types.Date `json:"transferDate"`
}

// Status of the negotiation, maps to FHIR eOverdracht task states (https://informatiestandaarden.nictiz.nl/wiki/vpk:V4.0_FHIR_eOverdracht#Using_Task_to_manage_the_workflow).
type TransferNegotiationStatus string

// SetCustomerJSONBody defines parameters for SetCustomer.
type SetCustomerJSONBody Customer

// AuthenticateWithIRMAJSONBody defines parameters for AuthenticateWithIRMA.
type AuthenticateWithIRMAJSONBody IRMAAuthenticationRequest

// AuthenticateWithPasswordJSONBody defines parameters for AuthenticateWithPassword.
type AuthenticateWithPasswordJSONBody PasswordAuthenticateRequest

// GetDossierParams defines parameters for GetDossier.
type GetDossierParams struct {

	// The patient ID
	PatientID string `json:"patientID"`
}

// CreateDossierJSONBody defines parameters for CreateDossier.
type CreateDossierJSONBody CreateDossierRequest

// SearchOrganizationsParams defines parameters for SearchOrganizations.
type SearchOrganizationsParams struct {

	// Keyword for finding care organizations.
	Query string `json:"query"`

	// Filters care organizations on service, only returning care organizations have a service in their DID Document which' type matches the given didServiceType. If not supplied, care organizations aren't filtered on service.
	DidServiceType *string `json:"didServiceType,omitempty"`
}

// UpdatePatientJSONBody defines parameters for UpdatePatient.
type UpdatePatientJSONBody PatientProperties

// NewPatientJSONBody defines parameters for NewPatient.
type NewPatientJSONBody PatientProperties

// GetPatientTransfersParams defines parameters for GetPatientTransfers.
type GetPatientTransfersParams struct {

	// The patient ID
	PatientID string `json:"patientID"`
}

// CreateTransferJSONBody defines parameters for CreateTransfer.
type CreateTransferJSONBody CreateTransferRequest

// AssignTransferNegotiationJSONBody defines parameters for AssignTransferNegotiation.
type AssignTransferNegotiationJSONBody AcceptTransferRequest

// SetCustomerJSONRequestBody defines body for SetCustomer for application/json ContentType.
type SetCustomerJSONRequestBody SetCustomerJSONBody

// AuthenticateWithIRMAJSONRequestBody defines body for AuthenticateWithIRMA for application/json ContentType.
type AuthenticateWithIRMAJSONRequestBody AuthenticateWithIRMAJSONBody

// AuthenticateWithPasswordJSONRequestBody defines body for AuthenticateWithPassword for application/json ContentType.
type AuthenticateWithPasswordJSONRequestBody AuthenticateWithPasswordJSONBody

// CreateDossierJSONRequestBody defines body for CreateDossier for application/json ContentType.
type CreateDossierJSONRequestBody CreateDossierJSONBody

// UpdatePatientJSONRequestBody defines body for UpdatePatient for application/json ContentType.
type UpdatePatientJSONRequestBody UpdatePatientJSONBody

// NewPatientJSONRequestBody defines body for NewPatient for application/json ContentType.
type NewPatientJSONRequestBody NewPatientJSONBody

// CreateTransferJSONRequestBody defines body for CreateTransfer for application/json ContentType.
type CreateTransferJSONRequestBody CreateTransferJSONBody

// AssignTransferNegotiationJSONRequestBody defines body for AssignTransferNegotiation for application/json ContentType.
type AssignTransferNegotiationJSONRequestBody AssignTransferNegotiationJSONBody
