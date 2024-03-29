package registry

import (
	"context"
	"errors"
	"fmt"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts"
	"sync"
	"time"

	"github.com/nuts-foundation/nuts-demo-ehr/nuts/client"
)

const cacheMaxAge = 10 * time.Second

type OrganizationRegistry interface {
	Get(ctx context.Context, organizationDID string) (*nuts.NutsOrganization, error)
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
	organization nuts.NutsOrganization
	writeTime    time.Time
}

func (r remoteOrganizationRegistry) Get(ctx context.Context, organizationDID string) (*nuts.NutsOrganization, error) {
	cached := r.fromCache(organizationDID)
	if cached != nil {
		return cached, nil
	}

	searchResults, err := r.client.SearchDiscoveryService(ctx, map[string]string{
		"credentialSubject.id": organizationDID,
	}, nil, nil)
	if err != nil {
		return nil, err
	}

	// find a search result that yields organization_name and organization_city
	var results []nuts.NutsOrganization
	for _, searchResult := range searchResults {
		if searchResult.Details.Name == "" {
			continue
		}
		results = append(results, searchResult.NutsOrganization)
	}

	if len(results) > 1 {
		// TODO: Get latest issued VC, or maybe all of them?
		return nil, errors.New("multiple organizations found (not supported yet)")
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("organization not found on any Discovery Service: %s", organizationDID)
	}
	r.toCache(results[0])
	return &results[0], nil
}

func (r remoteOrganizationRegistry) GetCompoundServiceEndpoint(ctx context.Context, organizationDID, serviceType string, field string) (string, error) {
	return r.client.GetCompoundServiceEndpoint(ctx, organizationDID, serviceType, field)
}

func (r remoteOrganizationRegistry) toCache(organizations ...nuts.NutsOrganization) {
	r.cacheMux.Lock()
	defer r.cacheMux.Unlock()
	for _, organization := range organizations {
		r.cache[organization.ID] = entry{
			organization: organization,
			writeTime:    time.Now(),
		}
	}
}

func (r remoteOrganizationRegistry) fromCache(organizationDID string) *nuts.NutsOrganization {
	r.cacheMux.Lock()
	defer r.cacheMux.Unlock()
	if cached, ok := r.cache[organizationDID]; ok && time.Now().Before(cached.writeTime.Add(cacheMaxAge)) {
		return &cached.organization
	}
	return nil
}
