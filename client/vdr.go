package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/nuts-foundation/go-did/did"
	"github.com/nuts-foundation/nuts-demo-ehr/client/vdr"
)

type CompoundService struct {
	ID              string
	Type            string
	ServiceEndpoint map[string]string
}

func (c HTTPClient) GetCompoundService(ctx context.Context, holderDID string, didServiceType string) (*CompoundService, error) {
	// TODO: Find out how this complex methods can be moved to Nuts Node
	result, err := c.resolveDIDDocument(ctx, holderDID)
	if err != nil {
		return nil, err
	}

	for _, service := range result.Service {
		if service.Type == didServiceType {
			compoundService := &CompoundService{}
			if err := service.UnmarshalServiceEndpoint(&compoundService.ServiceEndpoint); err == nil {
				compoundService.ID = service.ID.String()
				compoundService.Type = service.Type
				return compoundService, nil
			}
			var serviceRef string
			if err := service.UnmarshalServiceEndpoint(&serviceRef); err != nil {
				return nil, fmt.Errorf("DID document service is neither a compound service nor a reference (service=%s)", service.ID)
				//continue
			}
			parsedRef, err := did.ParseDIDURL(serviceRef)
			if err != nil {
				return nil, fmt.Errorf("reference id is not a valid DID URL: %w", err)
			}
			refDID := "did:" + parsedRef.Method + ":" + parsedRef.ID
			query, err := url.ParseQuery(parsedRef.Query)
			if err != nil {
				// ignore ill formatted query params and try the next candidate
				continue
			}
			refType := query.Get("type")
			if refType == "" {
				continue
			}

			if refDID == holderDID {
				continue
			}
			return c.GetCompoundService(ctx, refDID, refType)
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
	responseStruct := struct {
		Document did.Document
	}{}
	if err := json.Unmarshal(data, &responseStruct); err != nil {
		return did.Document{}, err
	}
	return responseStruct.Document, nil
}

func (c HTTPClient) vdr() vdr.ClientInterface {
	response, err := vdr.NewClientWithResponses(c.getNodeURL())
	if err != nil {
		panic(err)
	}
	return response
}
