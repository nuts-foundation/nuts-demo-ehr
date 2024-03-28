package domain

import (
	"context"
	"fmt"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/patients"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	nutsClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client"
	"net/url"
)

type ZorginzageService struct {
	NutsClient  *nutsClient.HTTPClient
	FHIRFactory fhir.Factory
}

func (z ZorginzageService) RemotePatient(ctx context.Context, localDID, remotePartyDID string, patientSSN string) (*types.Patient, error) {
	fhirClient, err := z.fhirClient(ctx, localDID, remotePartyDID, "homemonitoring", "homemonitoring")
	if err != nil {
		return nil, err
	}
	patient := resources.Patient{}
	if err = fhirClient.ReadOne(ctx, "/Patient?identifier="+url.QueryEscape(patientSSN), &patient); err != nil {
		return nil, fmt.Errorf("unable to read remote patient: %w", err)
	}
	result := patients.ToDomainPatient(patient)
	return &result, nil
}

func (z ZorginzageService) fhirClient(ctx context.Context, localDID string, remotePartyDID string, scope string, serviceName string) (fhir.Client, error) {
	endpointsInterf, err := z.NutsClient.ResolveServiceEndpoint(ctx, remotePartyDID, serviceName, "object")
	if err != nil {
		return nil, fmt.Errorf("resolve DID service (DID=%s, service=%s): %w", remotePartyDID, serviceName, err)
	}
	endpoints := endpointsInterf.(map[string]string)
	fhirEndpoint := endpoints["fhir"]
	if fhirEndpoint == "" {
		return nil, fmt.Errorf("remote XIS does not have its FHIR endpoint registered (DID=%s)", remotePartyDID)
	}
	accessToken, err := z.NutsClient.RequestServiceAccessToken(ctx, localDID, remotePartyDID, scope)
	if err != nil {
		return nil, fmt.Errorf("unable to get access token (DID=%s,scope=%s): %w", remotePartyDID, scope, err)
	}
	fhirClient := z.FHIRFactory(fhir.WithURL(fhirEndpoint), fhir.WithAuthToken(accessToken))
	return fhirClient, nil
}
