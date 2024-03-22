package client

import (
	"context"
	"encoding/json"
	"github.com/nuts-foundation/go-did/did"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts"
	nutsDiscoveryClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client/discovery"
	"github.com/nuts-foundation/nuts-node/vdr/didweb"
	"io"
	"net/http"
	"time"
)

var _ Discovery = (*HTTPClient)(nil)

type Discovery interface {
	SearchService(ctx context.Context, organizationSearchParam string, discoveryServiceID string, didServiceType *string) ([]nuts.NutsOrganization, error)
}

func (c HTTPClient) SearchService(ctx context.Context, organizationSearchParam string, discoveryServiceID string, didServiceType *string) ([]nuts.NutsOrganization, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	searchParams := map[string]interface{}{
		"credentialSubject.organization.name": organizationSearchParam + "*",
	}
	params := nutsDiscoveryClient.SearchPresentationsParams{Query: &searchParams}

	resp, err := c.discovery().SearchPresentations(ctx, discoveryServiceID, &params)
	if err != nil {
		return nil, err
	}
	respData, err := testAndReadResponse(http.StatusOK, resp)
	if err != nil {
		return nil, err
	}
	response := make([]nutsDiscoveryClient.SearchResult, 0)
	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, err
	}

	// resolve all DIDs from .subjectId and filter on given didServiceType if given
	results := make([]nuts.NutsOrganization, 0)
	for _, searchResult := range response {
		if didServiceType != nil {
			// parse did and convert did:web to url
			doc, err := resolve(ctx, searchResult.SubjectId)
			if err != nil {
				return nil, err
			}
			// check if the didServiceType is in the service array
			serviceFound := false
			for _, service := range doc.Service {
				if service.Type == *didServiceType {
					serviceFound = true
					break
				}
			}
			if !serviceFound {
				continue
			}
		}
		// convert searchResult to NutsOrganization and use the toCache method to store the results in the cache
		results = append(results, organizationSearchResultToDomain(searchResult))
	}

	return results, nil
}

// resolve a did:web document. This is a helper function to resolve a did:web document to a did.Document.
// construct the right URL, call it and parse the result to a did.Document.
func resolve(ctx context.Context, didStr string) (*did.Document, error) {
	id, err := did.ParseDID(didStr)
	if err != nil {
		return nil, err
	}
	didURL, err := didweb.DIDToURL(*id)
	if err != nil {
		return nil, err
	}
	didURL = didURL.JoinPath("did.json")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, didURL.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return did.ParseDocument(string(bytes))
}

func organizationSearchResultToDomain(searchResult nutsDiscoveryClient.SearchResult) nuts.NutsOrganization {
	return nuts.NutsOrganization{
		ID: searchResult.SubjectId,
		Details: nuts.OrganizationDetails{
			Name: searchResult.Fields["organization_name"].(string),
			City: searchResult.Fields["organization_city"].(string),
		},
	}
}

func (c HTTPClient) discovery() nutsDiscoveryClient.ClientInterface {
	var response nutsDiscoveryClient.ClientInterface
	var err error

	if c.Authorizer != nil {
		requestEditorFn := nutsDiscoveryClient.RequestEditorFn(c.Authorizer.RequestEditorFn())
		response, err = nutsDiscoveryClient.NewClientWithResponses(c.getNodeURL(), nutsDiscoveryClient.WithRequestEditorFn(requestEditorFn))
	} else {
		response, err = nutsDiscoveryClient.NewClientWithResponses(c.getNodeURL())
	}

	if err != nil {
		panic(err)
	}
	return response
}
