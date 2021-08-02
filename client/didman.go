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

func (c HTTPClient) didman() didman.ClientInterface {
	response, err := didman.NewClientWithResponses(c.getNodeURL())
	if err != nil {
		panic(err)
	}
	return response
}
