package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nuts-foundation/go-did"
	"github.com/nuts-foundation/go-did/vc"
	"github.com/nuts-foundation/nuts-node/vcr/credential"
	"github.com/nuts-foundation/nuts-node/vcr/holder"
	"net/http"
	"time"

	"github.com/nuts-foundation/nuts-demo-ehr/nuts/client/vcr"
)

type VCRClient interface {
	CreateVC(ctx context.Context, typeName, issuer string, credentialSubject map[string]interface{}, expirationDate *time.Time) error
	FindCredentials(ctx context.Context, credential vc.VerifiableCredential, untrusted bool) ([]vc.VerifiableCredential, error)
	FindCredentialIDs(ctx context.Context, credential vc.VerifiableCredential, untrusted bool) ([]string, error)
	RevokeCredential(ctx context.Context, credentialID string) error
	ResolveCredential(ctx context.Context, credentialType ssi.URI, credentialID string, untrusted bool) (*vc.VerifiableCredential, error)
}

func (c HTTPClient) CreateVC(ctx context.Context, typeName, issuer string, credentialSubject map[string]interface{}, expirationDate *time.Time) error {
	var exp *string

	if expirationDate != nil {
		formatted := expirationDate.Format(time.RFC3339)
		exp = &formatted
	}

	response, err := c.vcr().IssueVC(ctx, vcr.IssueVCJSONRequestBody{
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

func (c HTTPClient) FindCredentials(ctx context.Context, credential vc.VerifiableCredential, untrusted bool) ([]vc.VerifiableCredential, error) {
	return c.search(ctx, credential, untrusted)
}

func (c HTTPClient) FindCredentialIDs(ctx context.Context, credential vc.VerifiableCredential, untrusted bool) ([]string, error) {
	credentials, err := c.search(ctx, credential, untrusted)
	if err != nil {
		return nil, err
	}
	var credentialIDs = make([]string, len(credentials))
	j := 0
	for _, curr := range credentials {
		if curr.ID != nil {
			credentialIDs[j] = curr.ID.String()
			j++
		}
	}
	return credentialIDs[:j], nil
}

func (c HTTPClient) RevokeCredential(ctx context.Context, credentialID string) error {
	response, err := c.vcr().RevokeVC(ctx, credentialID)
	if err != nil {
		return err
	}
	_, err = testAndReadResponse(http.StatusOK, response)
	return err
}

func (c HTTPClient) ResolveCredential(ctx context.Context, credentialType ssi.URI, credentialID string, untrusted bool) (*vc.VerifiableCredential, error) {
	query := GetNutsCredentialTemplate(credentialType)
	query.ID, _ = ssi.ParseURI(credentialID)
	result, err := c.search(ctx, query, untrusted)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("credential not found with ID: %s", credentialID)
	}
	if len(result) > 1 {
		return nil, fmt.Errorf("multiple credentials found with ID: %s", credentialID)
	}
	return &result[0], nil
}

func GetNutsCredentialTemplate(credentialType ssi.URI) vc.VerifiableCredential {
	return vc.VerifiableCredential{
		Context:           []ssi.URI{holder.VerifiableCredentialLDContextV1, *credential.NutsContextURI},
		Type:              []ssi.URI{vc.VerifiableCredentialTypeV1URI(), credentialType},
	}
}

func (c HTTPClient) search(ctx context.Context, credential vc.VerifiableCredential, untrusted bool) ([]vc.VerifiableCredential, error) {
	response, err := c.vcr().SearchVCs(ctx, vcr.SearchVCsJSONRequestBody{Query: credential, SearchOptions: &vcr.SearchOptions{
		AllowUntrustedIssuer: &untrusted,
	}})
	if err != nil {
		return nil, err
	}
	data, err := testAndReadResponse(http.StatusOK, response)
	if err != nil {
		return nil, err
	}
	var result []vc.VerifiableCredential
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
