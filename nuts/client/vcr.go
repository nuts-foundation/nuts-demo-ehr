package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nuts-foundation/nuts-demo-ehr/nuts/client/vcr"
)

type VCRClient interface {
	GetOrganization(ctx context.Context, organizationDID string) ([]map[string]interface{}, error)
	CreateVC(ctx context.Context, typeName, issuer string, credentialSubject map[string]interface{}, expirationDate *time.Time) error
	FindAuthorizationCredentials(ctx context.Context, params map[string]string) ([]vcr.VerifiableCredential, error)
	FindAuthorizationCredentialIDs(ctx context.Context, params map[string]string) ([]string, error)
	RevokeCredential(ctx context.Context, credentialID string) error
	ResolveVerifiableCredential(ctx context.Context, credentialID string, untrusted bool) (*vcr.VerifiableCredential, error)
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

var trueVal = true

func (c HTTPClient) FindAuthorizationCredentials(ctx context.Context, params map[string]string) ([]vcr.VerifiableCredential, error) {
	var vcParams = make([]vcr.KeyValuePair, 0)
	for k, v := range params {
		vcParams = append(vcParams, vcr.KeyValuePair{
			Key:   k,
			Value: v,
		})
	}

	response, err := c.vcr().Search(ctx, authorizationConcept, &vcr.SearchParams{Untrusted: &trueVal}, vcr.SearchJSONRequestBody{Params: vcParams})
	if err != nil {
		return nil, err
	}
	// returns array of raw credentials
	data, err := testAndReadResponse(http.StatusOK, response)
	if err != nil {
		return nil, err
	}
	var credentials []vcr.VerifiableCredential
	if err = json.Unmarshal(data, &credentials); err != nil {
		return nil, err
	}

	return credentials, nil
}

func (c HTTPClient) FindAuthorizationCredentialIDs(ctx context.Context, params map[string]string) ([]string, error) {
	credentials, err := c.FindAuthorizationCredentials(ctx, params)
	if err != nil {
		return nil, err
	}

	var credentialIDs = make([]string, len(credentials))
	j := 0
	for _, c := range credentials {
		if c.Id != nil {
			credentialIDs[j] = *c.Id
			j++
		}
	}
	return credentialIDs[:j], nil
}

func (c HTTPClient) RevokeCredential(ctx context.Context, credentialID string) error {
	response, err := c.vcr().Revoke(ctx, credentialID)
	if err != nil {
		return err
	}
	_, err = testAndReadResponse(http.StatusOK, response)
	return err
}

func (c HTTPClient) ResolveVerifiableCredential(ctx context.Context, credentialID string, untrusted bool) (*vcr.VerifiableCredential, error) {
	response, err := c.vcr().Resolve(ctx, credentialID, nil)
	if err != nil {
		return nil, err
	}
	data, err := testAndReadResponse(http.StatusOK, response)
	result := vcr.ResolutionResult{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	if !untrusted && result.CurrentStatus != "trusted" {
		return nil, fmt.Errorf("credential with ID %s is not trusted (but %s)", credentialID, result.CurrentStatus)
	}
	return &result.VerifiableCredential, nil
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
