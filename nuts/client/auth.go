package client

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	nutsAuthClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client/auth"
)

type Auth interface {
	CreateIrmaSession(customerDID string) ([]byte, error)
	GetIrmaSessionResult(sessionToken string) (*nutsAuthClient.SignSessionStatusResponse, error)

	CreateDummySession(customerDID string) ([]byte, error)
	GetDummySessionResult(sessionToken string) (*nutsAuthClient.SignSessionStatusResponse, error)
}

func (c HTTPClient) getSessionResult(sessionToken string) (*nutsAuthClient.SignSessionStatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.auth().GetSignSessionStatus(ctx, sessionToken)

	if err != nil {
		return nil, err
	}
	respData, err := testAndReadResponse(http.StatusOK, resp)

	if err != nil {
		return nil, err
	}
	sessionResponse := &nutsAuthClient.SignSessionStatusResponse{}

	if err := json.Unmarshal(respData, sessionResponse); err != nil {
		return nil, err
	}

	return sessionResponse, nil
}

func (c HTTPClient) CreateIrmaSession(customerDID string) ([]byte, error) {
	return c.createSession(customerDID, nutsAuthClient.SignSessionRequestMeansIrma)
}

func (c HTTPClient) GetIrmaSessionResult(sessionToken string) (*nutsAuthClient.SignSessionStatusResponse, error) {
	return c.getSessionResult(sessionToken)
}

func (c HTTPClient) CreateDummySession(customerDID string) ([]byte, error) {
	return c.createSession(customerDID, nutsAuthClient.SignSessionRequestMeansDummy)
}

func (c HTTPClient) GetDummySessionResult(sessionToken string) (*nutsAuthClient.SignSessionStatusResponse, error) {
	return c.getSessionResult(sessionToken)
}

func (c HTTPClient) createSession(customerDID string, means nutsAuthClient.SignSessionRequestMeans) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	le := nutsAuthClient.LegalEntity(customerDID)
	t := time.Now().Format(time.RFC3339)
	contractBody := nutsAuthClient.DrawUpContractJSONRequestBody{
		Language:    "NL",
		LegalEntity: le,
		Type:        "BehandelaarLogin",
		ValidFrom:   &t,
		Version:     "v3",
	}
	contractRespBody, err := c.auth().DrawUpContract(ctx, contractBody)
	if err != nil {
		return nil, err
	}
	contractResp, err := testAndReadResponse(http.StatusOK, contractRespBody)
	if err != nil {
		return nil, err
	}
	contract := nutsAuthClient.ContractResponse{}
	err = json.Unmarshal(contractResp, &contract)
	if err != nil {
		return nil, err
	}

	body := nutsAuthClient.CreateSignSessionJSONRequestBody{
		Means:   means,
		Payload: contract.Message,
	}

	resp, err := c.auth().CreateSignSession(ctx, body)
	if err != nil {
		return nil, err
	}
	return testAndReadResponse(http.StatusCreated, resp)
}

func (c HTTPClient) auth() nutsAuthClient.ClientInterface {
        var response nutsAuthClient.ClientInterface
        var err error
                
        if c.Authorizer != nil {
                requestEditorFn := nutsAuthClient.RequestEditorFn(c.Authorizer.RequestEditorFn())
                response, err = nutsAuthClient.NewClientWithResponses(c.getNodeURL(), nutsAuthClient.WithRequestEditorFn(requestEditorFn))
        } else {
                response, err = nutsAuthClient.NewClientWithResponses(c.getNodeURL())
        }

	if err != nil {
		panic(err)
	}
	return response
}
