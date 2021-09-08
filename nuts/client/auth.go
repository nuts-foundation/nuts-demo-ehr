package client

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	nutsAuthClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client/auth"
)

type Auth interface {
	CreateIrmaSession(customer domain.Customer) ([]byte, error)
	GetIrmaSessionResult(sessionToken string) ([]byte, error)

	CreateDummySession(customer domain.Customer) ([]byte, error)
	GetDummySessionResult(sessionToken string) ([]byte, error)
}

func (c HTTPClient) CreateIrmaSession(customer domain.Customer) ([]byte, error) {
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
		Means:   "irma",
		Payload: contract.Message,
	}

	resp, err := c.auth().CreateSignSession(ctx, body)
	if err != nil {
		return nil, err
	}
	return testAndReadResponse(http.StatusCreated, resp)
}

func (c HTTPClient) GetIrmaSessionResult(sessionToken string) ([]byte, error) {
	// todo set user session

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := c.auth().GetSignSessionStatus(ctx, sessionToken)
	if err != nil {
		return nil, err
	}
	return testAndReadResponse(http.StatusOK, resp)
}

func (c HTTPClient) CreateDummySession(customer domain.Customer) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := nutsAuthClient.CreateSignSessionJSONRequestBody{
		Means:   "dummy",
		Payload: "Ik verklaar te handelen namens " + customer.Name,
	}

	resp, err := c.auth().CreateSignSession(ctx, body)
	if err != nil {
		return nil, err
	}
	return testAndReadResponse(http.StatusCreated, resp)
}

func (c HTTPClient) GetDummySessionResult(sessionToken string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := c.auth().GetSignSessionStatus(ctx, sessionToken)
	if err != nil {
		return nil, err
	}
	return testAndReadResponse(http.StatusOK, resp)
}

func (c HTTPClient) auth() nutsAuthClient.ClientInterface {
	response, err := nutsAuthClient.NewClientWithResponses(c.getNodeURL())
	if err != nil {
		panic(err)
	}
	return response
}
