package registry

import (
	"context"
	"errors"
	"fmt"
	"github.com/nuts-foundation/nuts-demo-ehr/client"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"strings"
)

type OrganizationRegistry interface {
	Search(ctx context.Context, query string, didServiceType *string) ([]domain.Organization, error)
	Get(ctx context.Context, organizationDID string) (*domain.Organization, error)
	GetCompoundServiceEndpoint(ctx context.Context, organizationDID, serviceType string, field string) (string, error)
}

func NewOrganizationRegistry(client *client.HTTPClient) OrganizationRegistry {
	return &remoteOrganizationRegistry{
		client: client,
	}
}

type remoteOrganizationRegistry struct {
	client *client.HTTPClient
}

func (r remoteOrganizationRegistry) Search(ctx context.Context, query string, didServiceType *string) ([]domain.Organization, error) {
	organizations, err := r.client.SearchOrganizations(ctx, query, didServiceType)
	if err != nil {
		return nil, err
	}
	results := make([]domain.Organization, len(organizations))
	for i, curr := range organizations {
		results[i] = organizationConceptToDomain(curr.Organization)
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

func (r remoteOrganizationRegistry) GetCompoundServiceEndpoint(ctx context.Context, organizationDID, serviceType string, field string) (string, error) {
	endpoints, err := r.client.GetCompoundService(ctx, organizationDID, serviceType)
	if err != nil {
		return "", err
	}
	endpoint := endpoints[field]
	if endpoint == "" {
		return "", fmt.Errorf("DID compound service does not contain the requested endpoint (did=%s,service=%s,name=%s)", organizationDID, serviceType, field)
	}

	if strings.HasPrefix(endpoint, "did:nuts:") {
		// Endpoint is a reference which needs to be resolved
		resolvedEndpoint, err := r.client.ResolveEndpoint(ctx, endpoint)
		if err != nil {
			return "", fmt.Errorf("unable to resolve endpoint reference in DID service (did=%s,service=%s,ref=%s): %w", organizationDID, serviceType, endpoint)
		}
		return resolvedEndpoint, nil
	}
	return endpoint, nil
}

func organizationConceptToDomain(concept map[string]interface{}) domain.Organization {
	inner := concept["organization"].(map[string]interface{})
	return domain.Organization{
		Did:  concept["subject"].(string),
		City: inner["city"].(string),
		Name: inner["name"].(string),
	}
}
