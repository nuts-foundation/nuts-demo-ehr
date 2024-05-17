package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nuts-foundation/go-did/did"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts"
	nutsDiscoveryClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client/discovery"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/client/vdr_v2"
	"github.com/oapi-codegen/runtime"
	"net/http"
	"net/url"
	"time"
)

var _ Discovery = (*HTTPClient)(nil)

// DiscoverySearchResult models a single result for when searching on Discovery Services.
type DiscoverySearchResult struct {
	nuts.NutsOrganization
	ServiceID string `json:"service_id"`
}

type Discovery interface {
	SearchDiscoveryService(ctx context.Context, query map[string]string, discoveryServiceID *string, didServiceType *string) ([]DiscoverySearchResult, error)
}

func (c HTTPClient) SearchDiscoveryService(ctx context.Context, query map[string]string, discoveryServiceID *string, didServiceType *string) ([]DiscoverySearchResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var serviceIDs []string
	if discoveryServiceID != nil {
		serviceIDs = append(serviceIDs, *discoveryServiceID)
	} else {
		// service ID not specified, search all discovery services
		var err error
		serviceIDs, err = c.getDiscoveryServices(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get discovery services: %w", err)
		}
	}
	searchResults := make([]DiscoverySearchResult, 0)
	for _, serviceID := range serviceIDs {
		currResults, err := c.searchDiscoveryService(ctx, query, serviceID, didServiceType)
		if err != nil {
			return nil, fmt.Errorf("failed to get service participants for %s, %s: %w", serviceID, *didServiceType, err)
		}
		searchResults = append(searchResults, currResults...)
	}
	return searchResults, nil
}

func (c HTTPClient) searchDiscoveryService(ctx context.Context, query map[string]string, discoveryServiceID string, didServiceType *string) ([]DiscoverySearchResult, error) {
	queryAsMap := make(map[string]interface{}, 0)
	for key, value := range query {
		queryAsMap[key] = value
	}

	// replace generated code with own client call to avoid oapi runtime bug
	client := c.discovery().(*nutsDiscoveryClient.ClientWithResponses).ClientInterface.(*nutsDiscoveryClient.Client)
	req, err := newSearchPresentationsRequest(client.Server, discoveryServiceID, query)
	if err != nil {
		return nil, fmt.Errorf("failed to construct SearchPresentationsRequest: %w", err)
	}
	req = req.WithContext(ctx)

	resp, err := client.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to query discovery services: %w", err)
	}
	respData, err := testAndReadResponse(http.StatusOK, resp)
	if err != nil {
		return nil, err
	}
	response := make([]nutsDiscoveryClient.SearchResult, 0)
	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result from discovery services: %w", err)
	}

	// resolve all DIDs from .subjectId and filter on given didServiceType if given
	results := make([]DiscoverySearchResult, 0)
	for _, searchResult := range response {
		if didServiceType != nil {
			// parse did and convert did:web to url
			doc, err := c.resolveDID(ctx, searchResult.SubjectId)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve DID %s: %w", searchResult.SubjectId, err)
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
		results = append(results, DiscoverySearchResult{
			NutsOrganization: organizationSearchResultToDomain(searchResult),
			ServiceID:        discoveryServiceID,
		})
	}
	return results, nil
}

// newSearchPresentationsRequest is a replacement for the generated one which has a bug.
func newSearchPresentationsRequest(server string, serviceID string, params map[string]string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "serviceID", runtime.ParamLocationPath, serviceID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/internal/discovery/v1/%s", pathParam0)
	//if operationPath[0] == '/' {
	//	operationPath = "." + operationPath
	//}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	queryValues := queryURL.Query()
	for k, v := range params {
		queryValues.Add(k, v)
	}

	queryURL.RawQuery = queryValues.Encode()
	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c HTTPClient) resolveDID(ctx context.Context, didStr string) (*did.Document, error) {
	response, err := c.vdr().ResolveDID(ctx, didStr)
	if err != nil {
		return nil, nil
	}
	if err := testResponseCode(http.StatusOK, response); err != nil {
		return nil, err
	}
	didResponse, err := vdr_v2.ParseResolveDIDResponse(response)
	if err != nil {
		return nil, err
	}
	return &didResponse.JSON200.Document, nil
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

func (c HTTPClient) getDiscoveryServices(ctx context.Context) ([]string, error) {
	response, err := c.discovery().GetServices(ctx)
	if err != nil {
		return nil, err
	}
	err = testResponseCode(http.StatusOK, response)
	if err != nil {
		return nil, err
	}
	services, err := nutsDiscoveryClient.ParseGetServicesResponse(response)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	for _, service := range *services.JSON200 {
		result = append(result, service.Id)
	}
	return result, nil
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
