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

// Defines values for TransferNegotiationStatusStatus.
const (
	TransferNegotiationStatusStatusAccepted TransferNegotiationStatusStatus = "accepted"

	TransferNegotiationStatusStatusCancelled TransferNegotiationStatusStatus = "cancelled"

	TransferNegotiationStatusStatusCompleted TransferNegotiationStatusStatus = "completed"

	TransferNegotiationStatusStatusInProgress TransferNegotiationStatusStatus = "in-progress"

	TransferNegotiationStatusStatusOnHold TransferNegotiationStatusStatus = "on-hold"

	TransferNegotiationStatusStatusRequested TransferNegotiationStatusStatus = "requested"
)

// Request used to assign a transfer to a specific negotiation indicated by the negotiationID.
type AssignTransferRequest struct {

	// An internal object UUID which can be used as unique identifier for entities.
	NegotiationID ObjectID `json:"negotiationID"`
}

// API request to create a dossier for a patient.
type CreateDossierRequest struct {
	Name string `json:"name"`

	// An internal object UUID which can be used as unique identifier for entities.
	PatientID ObjectID `json:"patientID"`
}

// An request object to create a new transfer negotiation.
type CreateTransferNegotiationRequest struct {

	// Decentralized Identifier of the organization to which transfer of a patient is requested.
	OrganizationDID string `json:"organizationDID"`

	// Transfer date subject of the negotiation. Can be altered by both sending and receiving care organization.
	TransferDate openapi_types.Date `json:"transferDate"`
}

// CreateTransferRequest defines model for CreateTransferRequest.
type CreateTransferRequest struct {
	// Embedded struct due to allOf(#/components/schemas/TransferProperties)
	TransferProperties `yaml:",inline"`
	// Embedded fields due to inline allOf schema

	// An internal object UUID which can be used as unique identifier for entities.
	DossierID ObjectID `json:"dossierID"`
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

// InboxEntry defines model for InboxEntry.
type InboxEntry struct {

	// Descriptive title.
	Title string `json:"title"`
}

// InboxInfo defines model for InboxInfo.
type InboxInfo struct {

	// Number of new messages in the inbox.
	MessageCount int `json:"messageCount"`
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
	// Embedded fields due to inline allOf schema
	AvatarUrl *string `json:"avatar_url,omitempty"`
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

// Transfer defines model for Transfer.
type Transfer struct {
	// Embedded struct due to allOf(#/components/schemas/TransferProperties)
	TransferProperties `yaml:",inline"`
	// Embedded fields due to inline allOf schema

	// An internal object UUID which can be used as unique identifier for entities.
	DossierID ObjectID `json:"dossierID"`

	// An internal object UUID which can be used as unique identifier for entities.
	Id ObjectID `json:"id"`

	// Status of the transfer. If the state is "completed" or "cancelled" the transfer dossier becomes read-only. In that case no additional negotiations can be sent (for this transfer) or accepted. Possible values: - Created: the new transfer dossier is created, but no requests were sent (to receiving care organizations) yet. - Requested: one or more requests were sent to care organizations - Assigned: The transfer is assigned to one the receiving care organizations thet accepted the transfer. - Completed: the patient transfer is completed and marked as such by the receiving care organization. - Cancelled: the transfer is cancelled by the sending care organization.
	Status TransferStatus `json:"status"`
}

// Status of the transfer. If the state is "completed" or "cancelled" the transfer dossier becomes read-only. In that case no additional negotiations can be sent (for this transfer) or accepted. Possible values: - Created: the new transfer dossier is created, but no requests were sent (to receiving care organizations) yet. - Requested: one or more requests were sent to care organizations - Assigned: The transfer is assigned to one the receiving care organizations thet accepted the transfer. - Completed: the patient transfer is completed and marked as such by the receiving care organization. - Cancelled: the transfer is cancelled by the sending care organization.
type TransferStatus string

// TransferNegotiation defines model for TransferNegotiation.
type TransferNegotiation struct {
	// Embedded struct due to allOf(#/components/schemas/TransferNegotiationStatus)
	TransferNegotiationStatus `yaml:",inline"`
	// Embedded fields due to inline allOf schema

	// An internal object UUID which can be used as unique identifier for entities.
	Id ObjectID `json:"id"`

	// A care organization available through the Nuts Network to exchange information.
	Organization Organization `json:"organization"`

	// Decentralized Identifier of the organization to which transfer of a patient is requested.
	OrganizationDID string `json:"organizationDID"`

	// The id of the FHIR Task resource which tracks this negotiation.
	TaskID string `json:"taskID"`

	// Transfer date subject of the negotiation. Can be altered by both sending and receiving care organization.
	TransferDate openapi_types.Date `json:"transferDate"`

	// An internal object UUID which can be used as unique identifier for entities.
	TransferID ObjectID `json:"transferID"`
}

// A valid transfer negotiation state.
type TransferNegotiationStatus struct {

	// Status of the negotiation, maps to FHIR eOverdracht task states (https://informatiestandaarden.nictiz.nl/wiki/vpk:V4.0_FHIR_eOverdracht#Using_Task_to_manage_the_workflow).
	Status TransferNegotiationStatusStatus `json:"status"`
}

// Status of the negotiation, maps to FHIR eOverdracht task states (https://informatiestandaarden.nictiz.nl/wiki/vpk:V4.0_FHIR_eOverdracht#Using_Task_to_manage_the_workflow).
type TransferNegotiationStatusStatus string

// Properties of a transfer. These values can be updated over time.
type TransferProperties struct {

	// Accompanying text sent to care organizations to assess the patient transfer. It is populated/updated by the last negotiation that was started.
	Description string `json:"description"`

	// Transfer date as proposed by the sending XIS. It is populated/updated by the last negotiation that was started.
	TransferDate openapi_types.Date `json:"transferDate"`
}

// SetCustomerJSONBody defines parameters for SetCustomer.
type SetCustomerJSONBody Customer

// AuthenticateWithIRMAJSONBody defines parameters for AuthenticateWithIRMA.
type AuthenticateWithIRMAJSONBody IRMAAuthenticationRequest

// AuthenticateWithPasswordJSONBody defines parameters for AuthenticateWithPassword.
type AuthenticateWithPasswordJSONBody PasswordAuthenticateRequest

// NotifyTransferUpdateParams defines parameters for NotifyTransferUpdate.
type NotifyTransferUpdateParams struct {

	// DID of the receiving care organization.
	ReceiverDID string `json:"receiverDID"`
}

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

// GetPatientsParams defines parameters for GetPatients.
type GetPatientsParams struct {

	// Search patients by name
	Name *string `json:"name,omitempty"`
}

// NewPatientJSONBody defines parameters for NewPatient.
type NewPatientJSONBody PatientProperties

// GetPatientTransfersParams defines parameters for GetPatientTransfers.
type GetPatientTransfersParams struct {

	// The patient ID
	PatientID string `json:"patientID"`
}

// CreateTransferJSONBody defines parameters for CreateTransfer.
type CreateTransferJSONBody CreateTransferRequest

// UpdateTransferJSONBody defines parameters for UpdateTransfer.
type UpdateTransferJSONBody TransferProperties

// StartTransferNegotiationJSONBody defines parameters for StartTransferNegotiation.
type StartTransferNegotiationJSONBody CreateTransferNegotiationRequest

// UpdateTransferNegotiationStatusJSONBody defines parameters for UpdateTransferNegotiationStatus.
type UpdateTransferNegotiationStatusJSONBody TransferNegotiationStatus

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

// UpdateTransferJSONRequestBody defines body for UpdateTransfer for application/json ContentType.
type UpdateTransferJSONRequestBody UpdateTransferJSONBody

// StartTransferNegotiationJSONRequestBody defines body for StartTransferNegotiation for application/json ContentType.
type StartTransferNegotiationJSONRequestBody StartTransferNegotiationJSONBody

// UpdateTransferNegotiationStatusJSONRequestBody defines body for UpdateTransferNegotiationStatus for application/json ContentType.
type UpdateTransferNegotiationStatusJSONRequestBody UpdateTransferNegotiationStatusJSONBody
