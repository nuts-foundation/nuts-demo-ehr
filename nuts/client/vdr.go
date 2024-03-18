package client

import (
	"context"
	"fmt"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/client/vdr_v2"
	"net/http"
)

func (c HTTPClient) ResolveServiceEndpoint(ctx context.Context, did string, serviceType string, endpointType string) (interface{}, error) {
	ep := vdr_v2.FilterServicesParamsEndpointType(endpointType)
	response, err := c.vdr().FilterServices(ctx, did, &vdr_v2.FilterServicesParams{EndpointType: &ep, Type: &serviceType})
	if err != nil {
		return "", err
	}
	if err = testResponseCode(http.StatusOK, response); err != nil {
		return "", err
	}
	services, err := vdr_v2.ParseFilterServicesResponse(response)
	if err != nil {
		return "", err
	}
	if len(*services.JSON200) != 1 {
		return "", fmt.Errorf("expected exactly one service (DID=%s, type=%s), got %d", did, serviceType, len(*services.JSON200))
	}
	return (*services.JSON200)[0].ServiceEndpoint, nil
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
