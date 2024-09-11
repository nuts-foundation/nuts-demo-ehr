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
	Get(ctx context.Context, organizationID string) (*nuts.NutsOrganization, error)
	GetCompoundServiceEndpoint(ctx context.Context, organizationID, serviceType string, field string) (string, error)
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

func (r remoteOrganizationRegistry) Get(ctx context.Context, organizationID string) (*nuts.NutsOrganization, error) {
	cached := r.fromCache(organizationID)
	if cached != nil {
		return cached, nil
	}

	searchResults, err := r.client.SearchDiscoveryService(ctx, map[string]string{
		"credentialSubject.authServerURL": organizationID,
	}, nil)
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
		return nil, fmt.Errorf("organization not found on any Discovery Service: %s", organizationID)
	}
	r.toCache(results[0])
	return &results[0], nil
}

func (r remoteOrganizationRegistry) GetCompoundServiceEndpoint(ctx context.Context, organizationID, serviceType string, field string) (string, error) {
	endpoint, err := r.client.ResolveServiceEndpoint(ctx, organizationID, serviceType, field)
	if err != nil {
		return "", err
	}
	return endpoint, nil
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

func (r remoteOrganizationRegistry) fromCache(organizationID string) *nuts.NutsOrganization {
	r.cacheMux.Lock()
	defer r.cacheMux.Unlock()
	if cached, ok := r.cache[organizationID]; ok && time.Now().Before(cached.writeTime.Add(cacheMaxAge)) {
		return &cached.organization
	}
	return nil
}
