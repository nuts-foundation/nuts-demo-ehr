package nutspxp

import (
	"context"
	"errors"
	"fmt"
	"github.com/nuts-foundation/nuts-demo-ehr/nutspxp/client/pip"
	"io"
	"net/http"
	"time"
)

type Client interface {
	AddPIPData(id string, client string, scope string, verifier string, authInput map[string]interface{}) error
	DeletePIPData(id string) error
}

type HTTPClient struct {
	PIPAddress string
}

var _ Client = HTTPClient{}

func (c HTTPClient) AddPIPData(id string, client string, scope string, verifier string, authInput map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client().CreateData(ctx, id, pip.CreateDataJSONRequestBody{
		AuthInput:  authInput,
		ClientId:   client,
		Scope:      scope,
		VerifierId: verifier,
	})
	if err != nil {
		return err
	}
	err = testResponseCode(http.StatusNoContent, resp)
	if err != nil {
		return err
	}
	return nil
}

func (c HTTPClient) DeletePIPData(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client().DeleteData(ctx, id)
	if err != nil {
		return err
	}
	err = testResponseCode(http.StatusNoContent, resp)
	if err != nil {
		return err
	}
	return nil
}

func (c HTTPClient) client() pip.ClientInterface {
	response, err := pip.NewClientWithResponses(c.PIPAddress)

	if err != nil {
		panic(err)
	}
	return response
}

func testAndReadResponse(status int, response *http.Response) ([]byte, error) {
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
