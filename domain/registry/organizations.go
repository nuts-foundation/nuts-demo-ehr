package registry

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/nuts-foundation/nuts-demo-ehr/client"
	"github.com/nuts-foundation/nuts-demo-ehr/client/didman"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

const cacheMaxAge = 10 * time.Second

type OrganizationRegistry interface {
	Search(ctx context.Context, query string, didServiceType *string) ([]domain.Organization, error)
	Get(ctx context.Context, organizationDID string) (*domain.Organization, error)
	GetCompoundServiceEndpoint(ctx context.Context, organizationDID, serviceType string, field string) (string, error)
}

func NewOrganizationRegistry(client *client.HTTPClient) OrganizationRegistry {
	return &remoteOrganizationRegistry{
		client:   client,
		cache:    map[string]entry{},
		cacheMux: &sync.Mutex{},
	}
}

type remoteOrganizationRegistry struct {
	client   *client.HTTPClient
	cache    map[string]entry
	cacheMux *sync.Mutex
}

type entry struct {
	organization domain.Organization
	writeTime    time.Time
}

func (r remoteOrganizationRegistry) Search(ctx context.Context, query string, didServiceType *string) ([]domain.Organization, error) {
	organizations, err := r.client.SearchOrganizations(ctx, query, didServiceType)
	if err != nil {
		return nil, err
	}
	results := make([]domain.Organization, len(organizations))
	for i, curr := range organizations {
		results[i] = organizationSearchResultToDomain(curr)
	}
	r.toCache(results...)
	return results, nil
}

func (r remoteOrganizationRegistry) Get(ctx context.Context, organizationDID string) (*domain.Organization, error) {
	cached := r.fromCache(organizationDID)
	if cached != nil {
		return cached, nil
	}

	raw, err := r.client.GetOrganization(ctx, organizationDID)
	if err != nil {
		return nil, err
	}
	if len(raw) == 0 {
		return nil, errors.New("organization not found")
	}
	if len(raw) > 1 {
		// TODO: Get latest issued VC, or maybe all of them?
		return nil, errors.New("multiple organizations found (not supported yet)")
	}
	result := organizationConceptToDomain(raw[0])
	r.toCache(result)
	return &result, nil
}

func (r remoteOrganizationRegistry) GetCompoundServiceEndpoint(ctx context.Context, organizationDID, serviceType string, field string) (string, error) {
	endpoints, err := r.client.GetCompoundService(ctx, organizationDID, serviceType)
	if err != nil {
		return "", err
	}
	endpoint := endpoints.ServiceEndpoint[field]
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

func organizationSearchResultToDomain(result didman.OrganizationSearchResult) domain.Organization {
	org := result.Organization
	return domain.Organization{
		Did:  result.DIDDocument.ID.String(),
		City: org["city"].(string),
		Name: org["name"].(string),
	}
}

func (r remoteOrganizationRegistry) toCache(organizations ...domain.Organization) {
	r.cacheMux.Lock()
	defer r.cacheMux.Unlock()
	for _, organization := range organizations {
		r.cache[organization.Did] = entry{
			organization: organization,
			writeTime:    time.Now(),
		}
	}
}

func (r remoteOrganizationRegistry) fromCache(organizationDID string) *domain.Organization {
	r.cacheMux.Lock()
	defer r.cacheMux.Unlock()
	if cached, ok := r.cache[organizationDID]; ok && time.Now().Before(cached.writeTime.Add(cacheMaxAge)) {
		return &cached.organization
	}
	return nil
}
