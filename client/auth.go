package client

import (
	"context"
	"encoding/json"
	"github.com/nuts-foundation/nuts-demo-ehr/client/auth"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"net/http"
	"time"
)

func (c HTTPClient) CreateIrmaSession(customer domain.Customer) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	le := auth.LegalEntity(*customer.Did)
	t := time.Now().Format(time.RFC3339)
	contractBody := auth.DrawUpContractJSONRequestBody{
		Language:      "NL",
		LegalEntity:   le,
		Type:          "BehandelaarLogin",
		ValidFrom:     &t,
		Version:       "v3",
	}
	contractRespBody, err := c.auth().DrawUpContract(ctx, contractBody)
	if err != nil {
		return nil, err
	}
	contractResp, err := testAndReadResponse(http.StatusOK, contractRespBody)
	if err != nil {
		return nil, err
	}
	contract := auth.ContractResponse{}
	json.Unmarshal(contractResp, &contract)

	body := auth.CreateSignSessionJSONRequestBody{
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

	ctx, _ := context.WithTimeout(context.Background(), 5 * time.Second)
	resp, err := c.auth().GetSignSessionStatus(ctx, sessionToken)
	if err != nil {
		return nil, err
	}
	return testAndReadResponse(http.StatusOK, resp)
}

func (c HTTPClient) auth() auth.ClientInterface {
	response, err := auth.NewClientWithResponses(c.getNodeURL())
	if err != nil {
		panic(err)
	}
	return response
}