package registry

import (
	"context"
	"errors"
	"fmt"
	"github.com/nuts-foundation/nuts-node/vcr/credential"
	"sync"
	"time"

	"github.com/nuts-foundation/nuts-demo-ehr/nuts/client"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/client/didman"
)

const cacheMaxAge = 10 * time.Second

type OrganizationRegistry interface {
	Search(ctx context.Context, query string, didServiceType *string) ([]NutsOrganization, error)
	Get(ctx context.Context, organizationDID string) (*NutsOrganization, error)
	GetCompoundServiceEndpoint(ctx context.Context, organizationDID, serviceType string, field string) (string, error)
}

func NewOrganizationRegistry(client *client.HTTPClient) OrganizationRegistry {
	return &remoteOrganizationRegistry{
		client:   client,
		cache:    map[string]entry{},
		cacheMux: &sync.Mutex{},
	}
}

// NutsOrganization models the credentialSubject of a NutsOrganizationCredential.
type NutsOrganization struct {
	ID      string              `json:"id"`
	Details OrganizationDetails `json:"organization"`
}

type OrganizationDetails struct {
	Name string `json:"name"`
	City string `json:"city"`
}

type remoteOrganizationRegistry struct {
	client   *client.HTTPClient
	cache    map[string]entry
	cacheMux *sync.Mutex
}

type entry struct {
	organization NutsOrganization
	writeTime    time.Time
}

func (r remoteOrganizationRegistry) Search(ctx context.Context, query string, didServiceType *string) ([]NutsOrganization, error) {
	organizations, err := r.client.SearchOrganizations(ctx, query, didServiceType)
	if err != nil {
		return nil, err
	}
	results := make([]NutsOrganization, len(organizations))
	for i, curr := range organizations {
		results[i] = organizationSearchResultToDomain(curr)
	}
	r.toCache(results...)
	return results, nil
}

func (r remoteOrganizationRegistry) Get(ctx context.Context, organizationDID string) (*NutsOrganization, error) {
	cached := r.fromCache(organizationDID)
	if cached != nil {
		return cached, nil
	}

	query := client.GetNutsCredentialTemplate(*credential.NutsOrganizationCredentialTypeURI)
	query.CredentialSubject = []interface{}{}
	credentials, err := r.client.FindCredentials(ctx, query, false)
	if err != nil {
		return nil, err
	}
	if len(credentials) == 0 {
		return nil, errors.New("organization not found")
	}
	if len(credentials) > 1 {
		// TODO: Get latest issued VC, or maybe all of them?
		return nil, errors.New("multiple organizations found (not supported yet)")
	}
	var result NutsOrganization
	err = credentials[0].UnmarshalCredentialSubject(&result)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal NutsOrganizationCredential subject: %w", err)
	}
	r.toCache(result)
	return &result, nil
}

func (r remoteOrganizationRegistry) GetCompoundServiceEndpoint(ctx context.Context, organizationDID, serviceType string, field string) (string, error) {
	return r.client.GetCompoundServiceEndpoint(ctx, organizationDID, serviceType, field)
}

func organizationSearchResultToDomain(result didman.OrganizationSearchResult) NutsOrganization {
	org := result.Organization
	return NutsOrganization{
		ID: result.DIDDocument.ID.String(),
		Details: OrganizationDetails{
			Name: org["name"].(string),
			City: org["city"].(string),
		},
	}
}

func (r remoteOrganizationRegistry) toCache(organizations ...NutsOrganization) {
	r.cacheMux.Lock()
	defer r.cacheMux.Unlock()
	for _, organization := range organizations {
		r.cache[organization.ID] = entry{
			organization: organization,
			writeTime:    time.Now(),
		}
	}
}

func (r remoteOrganizationRegistry) fromCache(organizationDID string) *NutsOrganization {
	r.cacheMux.Lock()
	defer r.cacheMux.Unlock()
	if cached, ok := r.cache[organizationDID]; ok && time.Now().Before(cached.writeTime.Add(cacheMaxAge)) {
		return &cached.organization
	}
	return nil
}
