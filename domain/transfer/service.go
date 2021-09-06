package transfer

import (
	"context"
	"time"

	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/http/auth"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/registry"
)

type Service interface {
	// Create creates a new transfer
	Create(ctx context.Context, customerID int, dossierID string, description string, transferDate time.Time) (*domain.Transfer, error)

	CreateNegotiation(ctx context.Context, customerID int, transferID, organizationDID string, transferDate time.Time) (*domain.TransferNegotiation, error)

	// ProposeAlternateDate updates the date on the domain.TransferNegotiation indicated by the negotiationID.
	// It updates the status to ON_HOLD_STATE
	ProposeAlternateDate(ctx context.Context, customerID int, negotiationID string) (*domain.TransferNegotiation, error)

	// ConfirmNegotiation confirms the negotiation indicated by the negotiationID.
	// The updates the status to ACCEPTED_STATE.
	// It automatically cancels other negotiations of the domain.Transfer indicated by the transferID
	// by setting their status to CANCELLED_STATE.
	ConfirmNegotiation(ctx context.Context, customerID int, transferID, negotiationID string) (*domain.TransferNegotiation, error)

	// CancelNegotiation withdraws the negotiation/organization from the transfer. This is done by the sending party
	// It updates the status to CANCELLED_STATE, updates the FHIR Task and sends out a notification
	CancelNegotiation(ctx context.Context, customerID int, negotiationID string) (*domain.TransferNegotiation, error)

	// RejectNegotiation rejects the proposed transfer. This is done by the receiving party
	RejectNegotiation(ctx context.Context, customerID int, negotiationID string) (*domain.TransferNegotiation, error)

	// GetTransferRequest tries to retrieve a transfer request from requesting care organization's FHIR server.
	GetTransferRequest(ctx context.Context, customerID int, requestorDID string, fhirTaskID string) (*domain.TransferRequest, error)

	// AcceptTransferRequest accepts the given transfer request, allowing the sending organization to assign the patient transfer to the local organization
	AcceptTransferRequest(ctx context.Context, customerID int, requestorDID, fhirTaskID string) error
}

type service struct {
	transferRepo           Repository
	auth                   auth.Service
	localFHIRClientFactory fhir.Factory // client for interacting with the local FHIR server
	customerRepo customers.Repository
	registry     registry.OrganizationRegistry
	vcr          registry.VerifiableCredentialRegistry
	notifier     Notifier
}

func NewTransferService(authService auth.Service, localFHIRClientFactory fhir.Factory, transferRepository Repository, customerRepository customers.Repository, organizationRegistry registry.OrganizationRegistry, vcr registry.VerifiableCredentialRegistry) *service {
	return &service{
		auth:                   authService,
		localFHIRClientFactory: localFHIRClientFactory,
		transferRepo:           transferRepository,
		customerRepo:           customerRepository,
		registry:               organizationRegistry,
		vcr:                    vcr,
		notifier:               fireAndForgetNotifier{},
	}
}
