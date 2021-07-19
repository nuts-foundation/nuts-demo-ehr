package client

import (
	"context"
	"encoding/json"
	"github.com/nuts-foundation/nuts-demo-ehr/client/vcr"
	"net/http"
)

func (client HTTPClient) SearchOrganizations(ctx context.Context, query string) ([]map[string]interface{}, error) {
	response, err := client.vcr().Search(ctx, "organization", vcr.SearchJSONRequestBody{
		Params: []vcr.KeyValuePair{
			{Key: "organization.name", Value: query},
			{Key: "organization.city", Value: ""},
		},
	})
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
