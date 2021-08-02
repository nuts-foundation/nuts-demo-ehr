package client

import (
	"context"
	"encoding/json"
	"github.com/nuts-foundation/nuts-demo-ehr/client/vcr"
	"net/http"
)

func (c HTTPClient) GetOrganization(ctx context.Context, organizationDID string) ([]map[string]interface{}, error) {
	return c.searchVCR(ctx, []vcr.KeyValuePair{
		{Key: "subject", Value: organizationDID},
	})
}

func (c HTTPClient) searchVCR(ctx context.Context, params []vcr.KeyValuePair) ([]map[string]interface{}, error) {
	response, err := c.vcr().Search(ctx, organizationConcept, vcr.SearchJSONRequestBody{Params: params})
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

func (c HTTPClient) vcr() vcr.ClientInterface {
	response, err := vcr.NewClientWithResponses(c.getNodeURL())
	if err != nil {
		panic(err)
	}
	return response
}
