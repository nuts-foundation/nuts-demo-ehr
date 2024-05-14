package client

import (
	"context"
	"encoding/json"
	"fmt"
	nutsIamClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client/iam"
	"net/http"
	"time"
)

type Iam interface {
	CreateAuthenticationRequest(customerDID string) (*nutsIamClient.RedirectResponseWithID, error)
	GetAuthenticationResult(token string) (*nutsIamClient.TokenResponse, error)
	IntrospectAccessToken(ctx context.Context, accessToken string) (*nutsIamClient.TokenIntrospectionResponse, error)
}

var _ Iam = HTTPClient{}

func (c HTTPClient) CreateAuthenticationRequest(customerDID string) (*nutsIamClient.RedirectResponseWithID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.iam().RequestUserAccessToken(ctx, customerDID, nutsIamClient.RequestUserAccessTokenJSONRequestBody{
		RedirectUri: "http://localhost:1304/#/close",
		Scope:       "test",
		PreauthorizedUser: &nutsIamClient.UserDetails{
			Id:   "12345",
			Name: "John Doe",
			Role: "Verpleegkundige niveau 4",
		},
		Verifier: customerDID,
	})

	if err != nil {
		return nil, err
	}
	respData, err := testAndReadResponse(http.StatusOK, resp)

	if err != nil {
		return nil, err
	}
	response := &nutsIamClient.RedirectResponseWithID{}

	if err := json.Unmarshal(respData, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (c HTTPClient) GetAuthenticationResult(token string) (*nutsIamClient.TokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.iam().RetrieveAccessToken(ctx, token)
	if err != nil {
		return nil, err
	}

	respData, err := testAndReadResponse(http.StatusOK, resp)
	if err != nil {
		return nil, err
	}

	response := &nutsIamClient.TokenResponse{}
	if err := json.Unmarshal(respData, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (c HTTPClient) RequestServiceAccessToken(ctx context.Context, relyingPartyDID, authorizationServerDID string, scope string) (string, error) {
	response, err := c.iam().RequestServiceAccessToken(ctx, relyingPartyDID, nutsIamClient.RequestServiceAccessTokenJSONRequestBody{
		Scope:    scope,
		Verifier: authorizationServerDID,
	})
	if err != nil {
		return "", err
	}
	tokenResponse, err := nutsIamClient.ParseRequestServiceAccessTokenResponse(response)
	if err != nil {
		return "", err
	}

	if tokenResponse.JSON200 == nil {
		var detail string
		if tokenResponse.ApplicationproblemJSONDefault != nil {
			detail = tokenResponse.ApplicationproblemJSONDefault.Detail
		}
		return "", fmt.Errorf("unable to get access token: %s", detail)
	}
	return tokenResponse.JSON200.AccessToken, nil
}

func (c HTTPClient) IntrospectAccessToken(ctx context.Context, accessToken string) (*nutsIamClient.TokenIntrospectionResponse, error) {
	response, err := c.iam().IntrospectAccessTokenWithFormdataBody(ctx, nutsIamClient.IntrospectAccessTokenFormdataRequestBody{
		Token: accessToken,
	})
	if err != nil {
		return nil, err
	}
	introspectResponse, err := nutsIamClient.ParseIntrospectAccessTokenResponse(response)
	if err != nil {
		return nil, err
	}
	if introspectResponse.JSON200 == nil {
		return nil, fmt.Errorf("unable to introspect access token")
	}
	return introspectResponse.JSON200, nil
}

func (c HTTPClient) iam() nutsIamClient.ClientInterface {
	var response nutsIamClient.ClientInterface
	var err error

	if c.Authorizer != nil {
		requestEditorFn := nutsIamClient.RequestEditorFn(c.Authorizer.RequestEditorFn())
		response, err = nutsIamClient.NewClientWithResponses(c.getNodeURL(), nutsIamClient.WithRequestEditorFn(requestEditorFn))
	} else {
		response, err = nutsIamClient.NewClientWithResponses(c.getNodeURL())
	}

	if err != nil {
		panic(err)
	}
	return response
}
