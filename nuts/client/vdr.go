package client

import (
	"context"
	"errors"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/client/discovery"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/client/vdr_v2"
	"net/http"
)

func (c HTTPClient) ResolveServiceEndpoint(ctx context.Context, verifierID string, serviceType string, endpointType string) (string, error) {
	// search on discovery service where credentialSubject.authServerURL == verifierID
	response, err := c.discovery().SearchPresentations(ctx, serviceType, &discovery.SearchPresentationsParams{Query: &map[string]interface{}{"credentialSubject.authServerURL": verifierID}})
	if err != nil {
		return "", err
	}
	if err = testResponseCode(http.StatusOK, response); err != nil {
		return "", err
	}
	searchReponse, err := discovery.ParseSearchPresentationsResponse(response)
	if err != nil {
		return "", err
	}
	searchResults := searchReponse.JSON200
	// find service
	for _, service := range *searchResults {
		if value, ok := service.Parameters[endpointType]; ok {
			return value.(string), nil
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
