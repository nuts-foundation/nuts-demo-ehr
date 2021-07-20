package client

import (
	"context"
	"encoding/json"
	"github.com/nuts-foundation/nuts-demo-ehr/client/vcr"
	"net/http"
)

const organizationConcept = "organization"

func (client HTTPClient) SearchOrganizations(ctx context.Context, query string) ([]map[string]interface{}, error) {
	return client.searchVCR(ctx, []vcr.KeyValuePair{
		{Key: "organization.name", Value: query},
		{Key: "organization.city", Value: ""},
	})
}

func (client HTTPClient) GetOrganization(ctx context.Context, organizationDID string) ([]map[string]interface{}, error) {
	return client.searchVCR(ctx, []vcr.KeyValuePair{
		{Key: "subject", Value: organizationDID},
	})
}

func (client HTTPClient) searchVCR(ctx context.Context, params []vcr.KeyValuePair) ([]map[string]interface{}, error) {
	response, err := client.vcr().Search(ctx, organizationConcept, vcr.SearchJSONRequestBody{Params: params})
	if err != nil {
		return nil, err
	}
	data, err := testAndReadResponse(http.StatusOK, response)
	if err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (client HTTPClient) vcr() vcr.ClientInterface {
	response, err := vcr.NewClientWithResponses(client.getNodeURL())
	if err != nil {
		panic(err)
	}
	return response
}
