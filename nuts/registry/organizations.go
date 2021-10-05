package registry

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/client"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/client/didman"
)

const cacheMaxAge = 10 * time.Second

type OrganizationRegistry interface {
	Search(ctx context.Context, query string, didServiceType *string) ([]types.Organization, error)
	Get(ctx context.Context, organizationDID string) (*types.Organization, error)
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
	organization types.Organization
	writeTime    time.Time
}

func (r remoteOrganizationRegistry) Search(ctx context.Context, query string, didServiceType *string) ([]types.Organization, error) {
	organizations, err := r.client.SearchOrganizations(ctx, query, didServiceType)
	if err != nil {
		return nil, err
	}
	results := make([]types.Organization, len(organizations))
	for i, curr := range organizations {
		results[i] = organizationSearchResultToDomain(curr)
	}
	r.toCache(results...)
	return results, nil
}

func (r remoteOrganizationRegistry) Get(ctx context.Context, organizationDID string) (*types.Organization, error) {
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
	return r.client.GetCompoundServiceEndpoint(ctx, organizationDID, serviceType, field)
}

func organizationConceptToDomain(concept map[string]interface{}) types.Organization {
	inner := concept["organization"].(map[string]interface{})
	return types.Organization{
		Did:  concept["subject"].(string),
		City: inner["city"].(string),
		Name: inner["name"].(string),
	}
}

func organizationSearchResultToDomain(result didman.OrganizationSearchResult) types.Organization {
	org := result.Organization
	return types.Organization{
		Did:  result.DIDDocument.ID.String(),
		City: org["city"].(string),
		Name: org["name"].(string),
	}
}

func (r remoteOrganizationRegistry) toCache(organizations ...types.Organization) {
	r.cacheMux.Lock()
	defer r.cacheMux.Unlock()
	for _, organization := range organizations {
		r.cache[organization.Did] = entry{
			organization: organization,
			writeTime:    time.Now(),
		}
	}
}

func (r remoteOrganizationRegistry) fromCache(organizationDID string) *types.Organization {
	r.cacheMux.Lock()
	defer r.cacheMux.Unlock()
	if cached, ok := r.cache[organizationDID]; ok && time.Now().Before(cached.writeTime.Add(cacheMaxAge)) {
		return &cached.organization
	}
	return nil
}
