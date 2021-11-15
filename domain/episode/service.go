package episode

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/go-did/vc"
	reports "github.com/nuts-foundation/nuts-demo-ehr/domain/reports"
	"github.com/nuts-foundation/nuts-demo-ehr/http/auth"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/client/vcr"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir/zorginzage"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/registry"
	"github.com/nuts-foundation/nuts-node/vcr/credential"
)

type Service interface {
	Create(ctx context.Context, customerID int, patientID string, request types.CreateEpisodeRequest) (*types.Episode, error)
	Get(ctx context.Context, customerID int, dossierID string) (*types.Episode, error)
	GetReports(ctx context.Context, customerDID, patientSSN string) ([]types.Report, error)
	CreateCollaboration(ctx context.Context, customerDID, dossierID, patientSSN, senderDID string) error
	GetCollaborations(ctx context.Context, customerDID, dossierID, patientSSN string) ([]types.Collaboration, error)
}

func ssnURN(ssn string) string {
	return fmt.Sprintf("urn:oid:2.16.840.1.113883.2.4.6.3:%s", ssn)
}

func parseAuthCredentialSubject(credentialResponse vcr.VerifiableCredential) (*credential.NutsAuthorizationCredentialSubject, error) {
	bytes, err := json.Marshal(credentialResponse)
	if err != nil {
		return nil, err
	}

	authCredential := vc.VerifiableCredential{}
	if err = json.Unmarshal(bytes, &authCredential); err != nil {
		return nil, err
	}

	subject := make([]credential.NutsAuthorizationCredentialSubject, 0)

	if err := authCredential.UnmarshalCredentialSubject(&subject); err != nil {
		return nil, fmt.Errorf("invalid content for NutsAuthorizationCredential credentialSubject: %w", err)
	}

	return &subject[0], nil
}

type service struct {
	factory  fhir.Factory
	auth     auth.Service
	registry registry.OrganizationRegistry
	vcr      registry.VerifiableCredentialRegistry
}

func NewService(factory fhir.Factory, auth auth.Service, registry registry.OrganizationRegistry, vcr registry.VerifiableCredentialRegistry) Service {
	return &service{factory: factory, auth: auth, registry: registry, vcr: vcr}
}

func parseEpisodeOfCareID(authCredential vcr.VerifiableCredential) (string, error) {
	subject, err := parseAuthCredentialSubject(authCredential)
	if err != nil {
		return "", err
	}

	for _, resource := range subject.Resources {
		if strings.HasPrefix(resource.Path, "/EpisodeOfCare/") {
			return resource.Path[len("/EpisodeOfCare/"):], nil
		}
	}

	return "", errors.New("no episode found in credential")
}

func toEpisode(episode *fhir.EpisodeOfCare) *types.Episode {
	status := types.EpisodeStatus(episode.Status)
	periodStart := time.Time{}
	if episode.Period != nil {
		if episode.Period.Start != nil {
			periodStart, _ = time.Parse(time.RFC3339, string(*episode.Period.Start))
		}
	}

	diagnosis := ""
	if len(episode.Type) > 0 {
		diagnosis = fhir.FromStringPtr(episode.Type[0].Text)
	}

	return &types.Episode{
		Id:        types.ObjectID(fhir.FromIDPtr(episode.ID)),
		Status:    &status,
		Period:    types.Period{Start: &openapi_types.Date{Time: periodStart}},
		Diagnosis: diagnosis,
	}
}

func (service *service) Create(ctx context.Context, customerID int, patientID string, request types.CreateEpisodeRequest) (*types.Episode, error) {
	svc := zorginzage.NewService(service.factory(fhir.WithTenant(customerID)))

	episode, err := svc.CreateEpisode(ctx, patientID, request)
	if err != nil {
		return nil, err
	}

	return toEpisode(episode), nil
}

func (service *service) Get(ctx context.Context, customerID int, dossierID string) (*types.Episode, error) {
	svc := zorginzage.NewService(service.factory(fhir.WithTenant(customerID)))

	episode, err := svc.GetEpisode(ctx, dossierID)
	if err != nil {
		return nil, err
	}

	return toEpisode(episode), nil
}

func (service *service) CreateCollaboration(ctx context.Context, customerDID, dossierID, patientSSN, senderDID string) error {
	subject := ssnURN(patientSSN)

	return service.vcr.CreateAuthorizationCredential(ctx, customerDID, &credential.NutsAuthorizationCredentialSubject{
		ID:      senderDID,
		Subject: &subject,
		LegalBase: credential.LegalBase{
			ConsentType: "implied",
		},
		PurposeOfUse: zorginzage.ServiceName,
		Resources: []credential.Resource{
			{
				Path:       fmt.Sprintf("/EpisodeOfCare/%s", dossierID),
				Operations: []string{"read"},
			},
		},
	})
}

func (service *service) GetCollaborations(ctx context.Context, customerDID, dossierID, patientSSN string) ([]types.Collaboration, error) {
	params := &registry.VCRSearchParams{
		PurposeOfUse: zorginzage.ServiceName,
		Subject:      ssnURN(patientSSN),
		ResourcePath: fmt.Sprintf("/EpisodeOfCare/%s", dossierID),
	}

	if customerDID != "" {
		params.SubjectID = customerDID
	}

	credentials, err := service.vcr.FindAuthorizationCredentials(
		ctx,
		params,
	)
	if err != nil {
		return nil, err
	}

	var subjects []credential.NutsAuthorizationCredentialSubject

	for _, authCredential := range credentials {
		subject, err := parseAuthCredentialSubject(authCredential)
		if err != nil {
			return nil, err
		}

		subjects = append(subjects, *subject)
	}

	episodeID := types.ObjectID(dossierID)

	var collaborations []types.Collaboration

	for _, subject := range subjects {
		collaborations = append(collaborations, types.Collaboration{
			EpisodeID:       episodeID,
			OrganizationDID: subject.ID,
		})
	}

	return collaborations, nil
}

func (service *service) GetReports(ctx context.Context, customerDID, patientSSN string) ([]types.Report, error) {
	credentials, err := service.vcr.FindAuthorizationCredentials(
		ctx,
		&registry.VCRSearchParams{
			PurposeOfUse: zorginzage.ServiceName,
			SubjectID:    customerDID,
			Subject:      ssnURN(patientSSN),
		},
	)
	if err != nil {
		return nil, err
	}

	if len(credentials) == 0 {
		return []types.Report{}, nil
	}

	issuer := string(credentials[0].Issuer)

	fhirServer, err := service.registry.GetCompoundServiceEndpoint(ctx, issuer, zorginzage.ServiceName, "fhir")
	if err != nil {
		return nil, fmt.Errorf("error while looking up authorizer's FHIR server (did=%s): %w", issuer, err)
	}

	org, err := service.registry.Get(ctx, customerDID)
	if err != nil {
		return nil, fmt.Errorf("error while searching organization :%w", err)
	}

	bytes, err := json.Marshal(credentials[0])
	if err != nil {
		return nil, err
	}

	authCredential := &vc.VerifiableCredential{}

	if err := json.Unmarshal(bytes, authCredential); err != nil {
		return nil, err
	}

	accessToken, err := service.auth.RequestAccessToken(ctx, customerDID, issuer, zorginzage.ServiceName, []vc.VerifiableCredential{*authCredential}, nil)
	if err != nil {
		return nil, err
	}

	episodeOfCareID, err := parseEpisodeOfCareID(credentials[0])
	if err != nil {
		return nil, err
	}

	observations := []resources.Observation{}

	fhirClient := fhir.NewFactory(fhir.WithURL(fhirServer), fhir.WithAuthToken(accessToken.AccessToken))()

	if err := fhirClient.ReadMultiple(ctx, "/Observation", map[string]string{
		"context": fmt.Sprintf("EpisodeOfCare/%s", episodeOfCareID),
		//"subject": fmt.Sprintf("Patient/%s", patientSSN),
	}, &observations); err != nil {
		return nil, err
	}

	results := make([]types.Report, len(observations))

	for _, observation := range observations {
		domainObservation := reports.ConvertToDomain(&observation, fhir.FromStringPtr(observation.Subject.ID))
		domainObservation.Source = org.Name
		results = append(results, domainObservation)

	}

	return results, nil
}
