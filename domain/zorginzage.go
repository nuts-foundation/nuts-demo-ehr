package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/patients"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	nutsClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client"
	"net/url"
)

type ZorginzageService struct {
	NutsClient *nutsClient.HTTPClient
}

func (z ZorginzageService) RemotePatient(ctx context.Context, subjectDID, remotePartyDID string, patientSSN string) (*types.RemotePatientFile, error) {
	fhirClient, err := z.fhirClient(ctx, subjectDID, remotePartyDID, "homemonitoring", "homemonitoring")
	if err != nil {
		return nil, err
	}

	var result types.RemotePatientFile
	// Load Patient resource
	searchResult := resources.Bundle{}
	if err = fhirClient.ReadOne(ctx, "Patient?identifier="+url.QueryEscape(patientSSN), &searchResult); err != nil {
		return nil, fmt.Errorf("unable to read remote Patient resource: %w", err)
	}
	patientBundle, err := searchResult.Transform("Patient")
	if err != nil {
		return nil, fmt.Errorf("unable to transform search result FHIR Bundle into Bundle with Patients: %w", err)
	}
	if len(patientBundle.Entry) == 0 {
		return nil, errors.New("patient not found at remote FHIR server")
	}
	result.Patient = patients.ToDomainPatient(patientBundle.Entry[0].Resource.(resources.Patient))
	// Load Observation resources
	if err = fhirClient.ReadOne(ctx, "Observation?patient.identifier="+url.QueryEscape(patientSSN), &searchResult); err != nil {
		return nil, fmt.Errorf("unable to read remote Observation resources for patient: %w", err)
	}
	observationBundle, err := searchResult.Transform("Observation")
	if err != nil {
		return nil, fmt.Errorf("unable to transform search result FHIR Bundle into Bundle with Observations: %w", err)
	}
	for _, entry := range observationBundle.Entry {
		// remarshal into map
		observationJSON, _ := json.Marshal(entry.Resource.(resources.Observation))
		observation := make(map[string]interface{}, 0)
		_ = json.Unmarshal(observationJSON, &observation)
		result.Observations = append(result.Observations, observation)
	}
	return &result, nil
}

func (z ZorginzageService) fhirClient(ctx context.Context, localDID string, remotePartyDID string, scope string, serviceName string) (fhir.Client, error) {
	fhirEndpoint, err := z.NutsClient.ResolveServiceEndpoint(ctx, remotePartyDID, serviceName, "fhir")
	if err != nil {
		return nil, fmt.Errorf("resolve DID service (DID=%s, service=%s): %w", remotePartyDID, serviceName, err)
	}
	accessToken, err := z.NutsClient.RequestServiceAccessToken(ctx, localDID, remotePartyDID, scope)
	if err != nil {
		return nil, fmt.Errorf("unable to get access token (DID=%s,scope=%s): %w", remotePartyDID, scope, err)
	}
	return fhir.NewFactory(fhir.WithURL(fhirEndpoint), fhir.WithAuthToken(accessToken))(), nil
}
