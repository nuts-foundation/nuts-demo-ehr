package client

import (
	"context"
	"encoding/json"
	nutsIamClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client/iam"
	"net/http"
	"time"
)

type Iam interface {
	CreateAuthenticationRequest(customerDID string) (*nutsIamClient.RedirectResponseWithID, error)
	GetAuthenticationResult(token string) (*nutsIamClient.TokenResponse, error)
}

func (c HTTPClient) CreateAuthenticationRequest(customerDID string) (*nutsIamClient.RedirectResponseWithID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.iam().RequestUserAccessToken(ctx, customerDID, nutsIamClient.RequestUserAccessTokenJSONRequestBody{
		RedirectUri: "http://localhost:1304/#/close",
		Scope:       "test",
		UserId:      "test",
		Verifier:    customerDID,
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
