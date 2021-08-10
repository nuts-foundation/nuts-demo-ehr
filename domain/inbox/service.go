package inbox

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/registry"
)

type Service struct {
	customerRepository customers.Repository
	repository         Repository
	orgRegistry        registry.OrganizationRegistry
}

func NewService(customerRepository customers.Repository, repository Repository, organizationRegistry registry.OrganizationRegistry) *Service {
	return &Service{customerRepository: customerRepository, repository: repository, orgRegistry: organizationRegistry}
}

func (s Service) RegisterNotification(ctx context.Context, customerID, senderDID string) error {
	return s.repository.registerNotification(ctx, customerID, senderDID)
}

func (s Service) List(ctx context.Context, customer *domain.Customer) ([]domain.InboxEntry, error) {
	notifications, err := s.repository.getAll(ctx, customer.Id)
	if err != nil {
		return nil, err
	}
	remoteFHIRServers := make(map[string]string, 0)
	for _, not := range notifications {
		if remoteFHIRServers[not.SenderDID] != "" {
			continue
		}
		fhirServer, err := s.orgRegistry.GetCompoundServiceEndpoint(ctx, not.SenderDID, "eOverdracht-sender", "fhir")
		if err != nil {
			logrus.Errorf("Unable to retrieve FHIR tasks from remote FHIR server (server=%s,did=%s): %v", fhirServer, not.SenderDID, err)
			continue
		}
		remoteFHIRServers[not.SenderDID] = fhirServer
	}

	var results []domain.InboxEntry
	for senderDID, fhirServer := range remoteFHIRServers {
		sendingOrg, err := s.orgRegistry.Get(ctx, senderDID)
		if err != nil {
			return nil, fmt.Errorf("error while looking up sender for inbox entry (did=%s): %w", senderDID, err)
		}
		entries, err := getInboxEntries(fhir.NewClient(fhirServer), *sendingOrg, *customer.Did)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve tasks from XIS (did=%s,url=%s): %w", senderDID, fhirServer, err)
		}
		results = append(results, entries...)
	}
	return results, nil
}

func getInboxEntries(client fhir.Client, sender domain.Organization, receiverDID string) ([]domain.InboxEntry, error) {
	// TODO: add _lastUpdated query paramater as required by Nictiz spec (https://informatiestandaarden.nictiz.nl/wiki/vpk:V4.0_FHIR_eOverdracht#Task_invocations)
	// But we might need some persistence for that, which we don't have right now.
	tasks, err := client.GetResources("/Task", map[string]string{
		"code": fmt.Sprintf("%s|%s", fhir.SnomedCodingSystem, fhir.SnomedTransferCode),
	})
	if err != nil {
		return nil, err
	}
	// Filter on current customer's DID (which is the receiver).
	// Should be done by the remote XIS, but it might return more than just tasks for this particular DID, so we do it client-side as well.
	tasks = fhir.Filter(tasks, func(resource gjson.Result) bool {
		return resource.Get("owner.identifier.value").String() == receiverDID
	})

	var results []domain.InboxEntry
	for _, resource := range tasks {
		results = append(results, domain.InboxEntry{
			Date:       resource.Get("meta.lastUpdated").String(),
			Sender:     sender,
			Title:      resource.Get("code.coding.0.display").String(),
			Type:       "transferRequest",
			ResourceID: resource.Get("id").String(),
		})
	}
	return results, nil
}
