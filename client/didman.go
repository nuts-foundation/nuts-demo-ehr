package client

import (
	"context"
	"github.com/nuts-foundation/nuts-demo-ehr/client/didman"
	"net/http"
)

const organizationConcept = "organization"

func (c HTTPClient) SearchOrganizations(ctx context.Context, query string, didServiceType *string) ([]didman.OrganizationSearchResult, error) {
	response, err := c.didman().SearchOrganizations(ctx, &didman.SearchOrganizationsParams{
		Query:          query,
		DidServiceType: didServiceType,
	})
	if err != nil {
		return nil, err
	}
	err = testResponseCode(http.StatusOK, response)
	if err != nil {
		return nil, err
	}
	searchResponse, err := didman.ParseSearchOrganizationsResponse(response)
	if err != nil {
		return nil, err
	}
	return *searchResponse.JSON200, nil
}

func (c HTTPClient) GetCompoundServiceEndpoint(ctx context.Context, organizationDID, serviceType string, field string) (string, error) {
	resolve := true
	response, err := c.didman().GetCompoundServiceEndpoint(ctx, organizationDID, serviceType, field, &didman.GetCompoundServiceEndpointParams{Resolve: &resolve})
	if err != nil {
		return "", err
	}
	err = testResponseCode(http.StatusOK, response)
	if err != nil {
		return "", err
	}
	parsedResponse, err := didman.ParseGetCompoundServiceEndpointResponse(response)
	if err != nil {
		return "", err
	}
	return *parsedResponse.JSON200, nil
}

func (c HTTPClient) didman() didman.ClientInterface {
	response, err := didman.NewClientWithResponses(c.getNodeURL())
	if err != nil {
		panic(err)
	}
	return response
}
