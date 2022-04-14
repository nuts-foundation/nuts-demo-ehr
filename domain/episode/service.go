package episode

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/go-did/vc"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir/zorginzage"
	reports "github.com/nuts-foundation/nuts-demo-ehr/domain/reports"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	"github.com/nuts-foundation/nuts-demo-ehr/http/auth"
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

func parseAuthCredentialSubject(authCredential vc.VerifiableCredential) (*credential.NutsAuthorizationCredentialSubject, error) {
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

func parseEpisodeOfCareID(authCredential vc.VerifiableCredential) (string, error) {
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

func (service *service) Create(ctx context.Context, customerID int, patientID string, request types.CreateEpisodeRequest) (*types.Episode, error) {
	svc := zorginzage.NewService(service.factory(fhir.WithTenant(customerID)))

	episode, err := svc.CreateEpisode(ctx, patientID, request)
	if err != nil {
		return nil, err
	}

	return zorginzage.ToEpisode(episode), nil
}

func (service *service) Get(ctx context.Context, customerID int, dossierID string) (*types.Episode, error) {
	svc := zorginzage.NewService(service.factory(fhir.WithTenant(customerID)))

	episode, err := svc.GetEpisode(ctx, dossierID)
	if err != nil {
		return nil, err
	}

	return zorginzage.ToEpisode(episode), nil
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
		org, err := service.registry.Get(ctx, subject.ID)
		if err != nil {
			return nil, err
		}
		collaborations = append(collaborations, types.Collaboration{
			EpisodeID:        episodeID,
			OrganizationDID:  subject.ID,
			OrganizationName: org.Details.Name,
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

	// TODO: loop over all credentials
	issuer := credentials[0].Issuer.String()

	fhirServer, err := service.registry.GetCompoundServiceEndpoint(ctx, issuer, zorginzage.ServiceName, "fhir")
	if err != nil {
		return nil, fmt.Errorf("error while looking up authorizer's FHIR server (did=%s): %w", issuer, err)
	}

	issuerOrg, err := service.registry.Get(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("error while searching organization :%w", err)
	}

	accessToken, err := service.auth.RequestAccessToken(ctx, customerDID, issuer, zorginzage.ServiceName, []vc.VerifiableCredential{credentials[0]}, nil)
	if err != nil {
		return nil, err
	}

	episodeOfCareID, err := parseEpisodeOfCareID(credentials[0])
	if err != nil {
		return nil, err
	}

	fhirClient := fhir.NewFactory(fhir.WithURL(fhirServer), fhir.WithAuthToken(accessToken.AccessToken))()

	fhirEpisode := &fhir.EpisodeOfCare{}
	err = fhirClient.ReadOne(ctx, "/EpisodeOfCare/"+episodeOfCareID, fhirEpisode)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve episode of care: %w", err)
	}
	episode := zorginzage.ToEpisode(fhirEpisode)

	observations := []resources.Observation{}
	if err := fhirClient.ReadMultiple(ctx, "/Observation", map[string]string{
		"context": fmt.Sprintf("EpisodeOfCare/%s", episodeOfCareID),
		//"subject": fmt.Sprintf("Patient/%s", patientSSN),
	}, &observations); err != nil {
		return nil, err
	}

	results := make([]types.Report, len(observations))

	for _, observation := range observations {
		domainObservation := reports.ConvertToDomain(&observation, fhir.FromStringPtr(observation.Subject.ID))
		domainObservation.Source = issuerOrg.Details.Name
		domainObservation.EpisodeName = &episode.Diagnosis
		results = append(results, domainObservation)
	}

	return results, nil
}
