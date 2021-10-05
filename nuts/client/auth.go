package client

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	nutsAuthClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client/auth"
)

type Auth interface {
	CreateIrmaSession(customer types.Customer) ([]byte, error)
	GetIrmaSessionResult(sessionToken string) (*nutsAuthClient.SignSessionStatusResponse, error)

	CreateDummySession(customer types.Customer) ([]byte, error)
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
	json.Unmarshal(respData, sessionResponse)
	return sessionResponse, nil
}

func (c HTTPClient) CreateIrmaSession(customer types.Customer) ([]byte, error) {
	return c.createSession(customer, nutsAuthClient.SignSessionRequestMeansIrma)
}

func (c HTTPClient) GetIrmaSessionResult(sessionToken string) (*nutsAuthClient.SignSessionStatusResponse, error) {
	return c.getSessionResult(sessionToken)
}

func (c HTTPClient) CreateDummySession(customer types.Customer) ([]byte, error) {
	return c.createSession(customer, nutsAuthClient.SignSessionRequestMeansDummy)
}

func (c HTTPClient) GetDummySessionResult(sessionToken string) (*nutsAuthClient.SignSessionStatusResponse, error) {
	return c.getSessionResult(sessionToken)
}

func (c HTTPClient) createSession(customer types.Customer, means nutsAuthClient.SignSessionRequestMeans) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	le := nutsAuthClient.LegalEntity(*customer.Did)
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
	response, err := nutsAuthClient.NewClientWithResponses(c.getNodeURL())
	if err != nil {
		panic(err)
	}
	return response
}
