package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nuts-foundation/go-did/did"
	"github.com/nuts-foundation/nuts-demo-ehr/client/vdr"
	"net/http"
	"net/url"
)

func (c HTTPClient) GetCompoundService(ctx context.Context, holderDID string, didServiceType string) (map[string]string, error) {
	// TODO: Find out how this complex methods can be moved to Nuts Node
	result, err := c.resolveDIDDocument(ctx, holderDID)
	if err != nil {
		return nil, err
	}

	for _, service := range result.Service {
		if service.Type == didServiceType {
			result := make(map[string]string, 0)
			if err := service.UnmarshalServiceEndpoint(&result); err != nil {
				return nil, fmt.Errorf("DID document service is not a compound service (service=%s): %w", service.ID, err)
			}
			return result, nil
		}
	}
	return nil, fmt.Errorf("DID document doesn't contain specified service (did=%s,service-type=%s)", holderDID, didServiceType)
}

func (c HTTPClient) ResolveEndpoint(ctx context.Context, referenceStr string) (string, error) {
	reference, err := did.ParseDIDURL(referenceStr)
	if err != nil {
		return "", err
	}
	query, err := url.ParseQuery(reference.Query)
	if err != nil {
		return "", err
	}
	referencedType := query.Get("type")
	if referencedType == "" {
		return "", errors.New("reference does not contain a type")
	}

	holderDID := *reference
	holderDID.Query = ""
	holderDID.Fragment = ""

	document, err := c.resolveDIDDocument(ctx, holderDID.String())
	if err != nil {
		return "", fmt.Errorf("unable to resolve holder DID document for endpoint reference (ref=%s): %w", referenceStr, err)
	}
	_, result, err := document.ResolveEndpointURL(referencedType)
	return result, err
}

func (c HTTPClient) resolveDIDDocument(ctx context.Context, holderDID string) (did.Document, error) {
	response, err := c.vdr().GetDID(ctx, holderDID)
	if err != nil {
		return did.Document{}, err
	}
	data, err := testAndReadResponse(http.StatusOK, response)
	if err != nil {
		return did.Document{}, err
	}
	result := did.Document{}
	if err := json.Unmarshal(data, &result); err != nil {
		return did.Document{}, err
	}
	return result, nil
}

func (c HTTPClient) vdr() vdr.ClientInterface {
	response, err := vdr.NewClientWithResponses(c.getNodeURL())
	if err != nil {
		panic(err)
	}
	return response
}
