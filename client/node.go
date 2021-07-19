package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/nuts-foundation/nuts-demo-ehr/client/auth"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

type HTTPClient struct {
	NutsNodeAddress string
}

func (client HTTPClient) CreateIrmaSession(customer domain.Customer) ([]byte, error) {
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
	contractRespBody, err := client.client().DrawUpContract(ctx, contractBody)
	if err != nil {
		return nil, err
	}
	contractResp, err := testAndParseResponse(http.StatusOK, contractRespBody)
	if err != nil {
		return nil, err
	}
	contract := auth.ContractResponse{}
	json.Unmarshal(contractResp, &contract)

	body := auth.CreateSignSessionJSONRequestBody{
		Means:   "irma",
		Payload: contract.Message,
	}

	resp, err := client.client().CreateSignSession(ctx, body)
	if err != nil {
		return nil, err
	}
	return testAndParseResponse(http.StatusCreated, resp)
}

func (client HTTPClient) GetIrmaSessionResult(sessionToken string) ([]byte, error) {
	// todo set user session

	ctx, _ := context.WithTimeout(context.Background(), 5 * time.Second)
	resp, err := client.client().GetSignSessionStatus(ctx, sessionToken)
	if err != nil {
		return nil, err
	}
	return testAndParseResponse(http.StatusOK, resp)
}

func (client HTTPClient) client() auth.ClientInterface {
	url := client.NutsNodeAddress
	if !strings.Contains(url, "http") {
		url = fmt.Sprintf("http://%v", client.NutsNodeAddress)
	}

	response, err := auth.NewClientWithResponses(url)
	if err != nil {
		panic(err)
	}
	return response
}

func testAndParseResponse(status int, response *http.Response) ([]byte, error) {
	if response.StatusCode == http.StatusNotFound {
		return nil, errors.New("not found")
	}
	if err := testResponseCode(status, response); err != nil {
		return nil, err
	}
	return io.ReadAll(response.Body)
}

func testResponseCode(expectedStatusCode int, response *http.Response) error {
	if response.StatusCode != expectedStatusCode {
		responseData, _ := io.ReadAll(response.Body)
		return fmt.Errorf("server returned HTTP %d (expected: %d), response: %s",
			response.StatusCode, expectedStatusCode, string(responseData))
	}
	return nil
}
