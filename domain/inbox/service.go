package inbox

import (
	"context"
	"fmt"

	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/registry"
)

type Service struct {
	customerRepository customers.Repository
	repository         Repository
	orgRegistry           registry.OrganizationRegistry

}

func NewService(customerRepository customers.Repository, repository Repository, organizationRegistry registry.OrganizationRegistry) *Service {
	return &Service{customerRepository: customerRepository, repository: repository, orgRegistry: organizationRegistry}
}

func (s Service) RegisterNotification(ctx context.Context, customerID, senderDID string) error {
	return s.repository.registerNotification(ctx, customerID, senderDID)
}

func (s Service) List(ctx context.Context, customerID string) ([]domain.InboxEntry, error) {
	notifications, err := s.repository.getAll(ctx, customerID)
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
			return nil, err
		}
		remoteFHIRServers[not.SenderDID] = fhirServer
	}

	var results []domain.InboxEntry
	for senderDID, _ := range remoteFHIRServers {
		entries, err := getInboxEntries(fhir.NewClient(senderDID))
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve tasks from XIS (url=%s): %w", senderDID, err)
		}
		results = append(results, entries...)
	}
	return results, nil
}

func getInboxEntries(client fhir.Client) ([]domain.InboxEntry, error) {
	tasks, err := client.GetResources("/Task")
	if err != nil {
		return nil, err
	}
	transferResources := fhir.FilterResources(tasks, fhir.SnomedCodingSystem, fhir.SnomedTransferCode)
	var results []domain.InboxEntry
	for _, resource := range transferResources {
		results = append(results, domain.InboxEntry{Title: resource.Get("code.coding.0.display").String()})
	}
	return results, nil
}
