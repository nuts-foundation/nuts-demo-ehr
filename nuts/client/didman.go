package client

import (
	"context"
	"net/http"

	nutsDIDManClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client/didman"
)

type DIDManClient interface {
	SearchOrganizations(ctx context.Context, query string, didServiceType *string) ([]nutsDIDManClient.OrganizationSearchResult, error)
	GetCompoundServiceEndpoint(ctx context.Context, organizationDID, serviceType string, field string) (string, error)
}

func (c HTTPClient) SearchOrganizations(ctx context.Context, query string, didServiceType *string) ([]nutsDIDManClient.OrganizationSearchResult, error) {
	response, err := c.didman().SearchOrganizations(ctx, &nutsDIDManClient.SearchOrganizationsParams{
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
	searchResponse, err := nutsDIDManClient.ParseSearchOrganizationsResponse(response)
	if err != nil {
		return nil, err
	}
	return *searchResponse.JSON200, nil
}

func (c HTTPClient) GetCompoundServiceEndpoint(ctx context.Context, organizationDID, serviceType string, field string) (string, error) {
	resolve := true
	response, err := c.didman().GetCompoundServiceEndpoint(ctx, organizationDID, serviceType, field, &nutsDIDManClient.GetCompoundServiceEndpointParams{Resolve: &resolve})
	if err != nil {
		return "", err
	}
	err = testResponseCode(http.StatusOK, response)
	if err != nil {
		return "", err
	}
	parsedResponse, err := nutsDIDManClient.ParseGetCompoundServiceEndpointResponse(response)
	if err != nil {
		return "", err
	}
	return parsedResponse.JSON200.Endpoint, nil
}

func (c HTTPClient) didman() nutsDIDManClient.ClientInterface {
	response, err := nutsDIDManClient.NewClientWithResponses(c.getNodeURL())
	if err != nil {
		panic(err)
	}
	return response
}
