package registry

import (
	"context"
	"errors"
	"github.com/nuts-foundation/nuts-demo-ehr/client"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

type OrganizationRegistry interface {
	Search(ctx context.Context, query string) ([]domain.Organization, error)
	Get(ctx context.Context, organizationDID string) (*domain.Organization, error)
}

func NewOrganizationRegistry(client *client.HTTPClient) OrganizationRegistry {
	return &remoteOrganizationRegistry{
		client: client,
	}
}

type remoteOrganizationRegistry struct {
	client *client.HTTPClient
}

func (r remoteOrganizationRegistry) Search(ctx context.Context, query string) ([]domain.Organization, error) {
	organizations, err := r.client.SearchOrganizations(ctx, query)
	if err != nil {
		return nil, err
	}
	results := make([]domain.Organization, len(organizations))
	for i, curr := range organizations {
		results[i] = organizationConceptToDomain(curr)
	}
	return results, nil
}

func (r remoteOrganizationRegistry) Get(ctx context.Context, organizationDID string) (*domain.Organization, error) {
	// TODO: Load from cache (use LRU cache)
	results, err := r.client.GetOrganization(ctx, organizationDID)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("organization not found")
	}
	if len(results) > 1 {
		// TODO: Get latest issued VC, or maybe all of them?
		return nil, errors.New("multiple organizations found (not supported yet)")
	}
	result := organizationConceptToDomain(results[0])
	return &result, nil
}

func organizationConceptToDomain(concept map[string]interface{}) domain.Organization {
	inner := concept["organization"].(map[string]interface{})
	return domain.Organization{
		Did:  concept["subject"].(string),
		City: inner["city"].(string),
		Name: inner["name"].(string),
	}
}