package client

import (
	"context"
	"errors"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/client/vdr_v2"
	"net/http"
)

func (c HTTPClient) ResolveServiceEndpoint(ctx context.Context, verifierDID string, serviceType string, endpointType string) (interface{}, error) {
	// resolve DID
	response, err := c.vdr().ResolveDID(ctx, verifierDID)
	if err != nil {
		return "", err
	}
	if err = testResponseCode(http.StatusOK, response); err != nil {
		return "", err
	}
	resolveReponse, err := vdr_v2.ParseResolveDIDResponse(response)
	if err != nil {
		return "", err
	}
	didDocument := resolveReponse.JSON200.Document
	// find service
	for _, service := range didDocument.Service {
		if service.Type == serviceType {
			serviceDef := make(map[string]interface{})
			if err := service.UnmarshalServiceEndpoint(&serviceDef); err != nil {
				return "", err
			}
			return serviceDef[endpointType], nil
		}
	}

	return "", errors.New("service not found")
}

func (c HTTPClient) ListSubjectDIDs(ctx context.Context, customerID string) ([]string, error) {
	response, err := c.vdr().SubjectDIDs(ctx, customerID)
	if err != nil {
		return nil, err
	}
	if err = testResponseCode(http.StatusOK, response); err != nil {
		return nil, err
	}
	parsedResponse, err := vdr_v2.ParseSubjectDIDsResponse(response)
	if err != nil {
		return nil, err
	}
	return *parsedResponse.JSON200, nil
}

func (c HTTPClient) vdr() vdr_v2.ClientInterface {
	var response vdr_v2.ClientInterface
	var err error

	if c.Authorizer != nil {
		requestEditorFn := vdr_v2.RequestEditorFn(c.Authorizer.RequestEditorFn())
		response, err = vdr_v2.NewClientWithResponses(c.getNodeURL(), vdr_v2.WithRequestEditorFn(requestEditorFn))
	} else {
		response, err = vdr_v2.NewClientWithResponses(c.getNodeURL())
	}

	if err != nil {
		panic(err)
	}
	return response
}
