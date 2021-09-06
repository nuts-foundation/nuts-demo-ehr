package client

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/nuts-foundation/nuts-demo-ehr/nuts/client/vcr"
)

type VCRClient interface {
	GetOrganization(ctx context.Context, organizationDID string) ([]map[string]interface{}, error)
	CreateVC(ctx context.Context, typeName, issuer string, credentialSubject map[string]interface{}, expirationDate *time.Time) error
}

func (c HTTPClient) GetOrganization(ctx context.Context, organizationDID string) ([]map[string]interface{}, error) {
	return c.searchVCR(ctx, organizationConcept, []vcr.KeyValuePair{
		{Key: "subject", Value: organizationDID},
	})
}

func (c HTTPClient) CreateVC(ctx context.Context, typeName, issuer string, credentialSubject map[string]interface{}, expirationDate *time.Time) error {
	var exp *string

	if expirationDate != nil {
		formatted := expirationDate.Format(time.RFC3339)
		exp = &formatted
	}

	response, err := c.vcr().Create(ctx, vcr.CreateJSONRequestBody{
		Type:              typeName,
		Issuer:            issuer,
		CredentialSubject: credentialSubject,
		ExpirationDate:    exp,
	})
	if err != nil {
		return err
	}

	_, err = testAndReadResponse(http.StatusOK, response)
	if err != nil {
		return err
	}

	return nil
}

func (c HTTPClient) searchVCR(ctx context.Context, concept string, params []vcr.KeyValuePair) ([]map[string]interface{}, error) {
	response, err := c.vcr().Search(ctx, concept, &vcr.SearchParams{}, vcr.SearchJSONRequestBody{Params: params})
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
