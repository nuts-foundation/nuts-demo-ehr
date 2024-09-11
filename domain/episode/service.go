package episode

import (
	"context"
	"errors"
	"fmt"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/acl"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir/zorginzage"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/client"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/registry"
	"github.com/sirupsen/logrus"
	"net/url"
)

type Service interface {
	Create(ctx context.Context, customerID, patientID string, request types.CreateEpisodeRequest) (*types.Episode, error)
	Get(ctx context.Context, customerID, dossierID string) (*types.Episode, error)
	GetReports(ctx context.Context, customerID, patientSSN string) ([]types.Report, error)
	CreateCollaboration(ctx context.Context, customerID, dossierID, patientSSN, senderDID string, client fhir.Client) error
	GetCollaborations(ctx context.Context, customerID, dossierID, patientSSN string, client fhir.Client) ([]types.Collaboration, error)
}

func ssnURN(ssn string) string {
	return fmt.Sprintf("urn:oid:2.16.840.1.113883.2.4.6.3:%s", ssn)
}

type service struct {
	factory       fhir.Factory
	nutsClient    *client.HTTPClient
	aclRepository *acl.Repository
	registry      registry.OrganizationRegistry
}

func NewService(factory fhir.Factory, nutsClient *client.HTTPClient, registry registry.OrganizationRegistry, aclRepository *acl.Repository) Service {
	return &service{factory: factory, nutsClient: nutsClient, registry: registry, aclRepository: aclRepository}
}

func (service *service) Create(ctx context.Context, customerID, patientID string, request types.CreateEpisodeRequest) (*types.Episode, error) {
	svc := zorginzage.NewService(service.factory(fhir.WithTenant(customerID)))

	episode, err := svc.CreateEpisode(ctx, patientID, request)
	if err != nil {
		return nil, err
	}

	return zorginzage.ToEpisode(episode), nil
}

func (service *service) Get(ctx context.Context, customerID, dossierID string) (*types.Episode, error) {
	svc := zorginzage.NewService(service.factory(fhir.WithTenant(customerID)))

	episode, err := svc.GetEpisode(ctx, dossierID)
	if err != nil {
		return nil, err
	}

	return zorginzage.ToEpisode(episode), nil
}

func (service *service) CreateCollaboration(ctx context.Context, customerDID, dossierID, patientSSN, senderDID string, client fhir.Client) error {
	// TODO: Need to formalize this in a use case specification
	type authorizedResource struct {
		Path       string
		Parameters map[string]string
		Operations []string
	}
	authorizedResources := []authorizedResource{
		{
			Path:       "Patient",
			Parameters: map[string]string{"identifier": patientSSN},
			Operations: []string{"read"},
		},
		// TODO: Do we need this particular one?
		{
			Path:       fmt.Sprintf("/EpisodeOfCare/%s", dossierID),
			Operations: []string{"read"},
		},
		{
			Path:       "/EpisodeOfCare",
			Parameters: map[string]string{"patient.identifier": patientSSN},
			Operations: []string{"read"},
		},
		{
			Path:       "/Observation",
			Parameters: map[string]string{"patient.identifier": patientSSN},
			Operations: []string{"read"},
		},
	}
	for _, resource := range authorizedResources {
		resourceURL := buildResourcePath(client, resource.Path, resource.Parameters)
		for _, op := range resource.Operations {
			if err := service.aclRepository.GrantAccess(ctx, customerDID, senderDID, op, resourceURL); err != nil {
				return err
			}
		}
	}
	return nil
}

func buildResourcePath(client fhir.Client, resourcePath string, query map[string]string) string {
	resourceURL := client.BuildRequestURI(resourcePath)
	resourceURL.Host = ""
	resourceURL.Scheme = ""
	if len(query) > 0 {
		values := url.Values{}
		for key, value := range query {
			values.Add(key, value)
		}
		resourceURL.RawQuery = values.Encode()
	}
	return resourceURL.String()
}

func (service *service) GetCollaborations(ctx context.Context, customerDID, dossierID, patientSSN string, client fhir.Client) ([]types.Collaboration, error) {
	// Find collaborators given the parties that have access to relevant FHIR resources
	searchResources := []string{
		buildResourcePath(client, "EpisodeOfCare/"+dossierID, nil),
	}
	authorizedDIDs := make(map[string]struct{})
	for _, resource := range searchResources {
		dids, err := service.aclRepository.AuthorizedParties(ctx, customerDID, resource, "read")
		if err != nil {
			return nil, err
		}
		for _, did := range dids {
			authorizedDIDs[did] = struct{}{}
		}
	}

	episodeID := dossierID

	var collaborations []types.Collaboration

	for authorizedDID := range authorizedDIDs {
		org, err := service.registry.Get(ctx, authorizedDID)
		if err != nil {
			logrus.WithError(err).Warn("Error looking up episode collaborator organization")
			collaborations = append(collaborations, types.Collaboration{
				EpisodeID:        episodeID,
				OrganizationID:  authorizedDID,
				OrganizationName: "!ERROR! " + err.Error(),
			})
		} else {
			collaborations = append(collaborations, types.Collaboration{
				EpisodeID:        episodeID,
				OrganizationID:  authorizedDID,
				OrganizationName: org.Details.Name,
			})
		}
	}

	return collaborations, nil
}

func (service *service) GetReports(ctx context.Context, customerDID, patientSSN string) ([]types.Report, error) {
	return nil, errors.New("not implemented")
	//credentials, err := service.vcr.FindAuthorizationCredentials(
	//	ctx,
	//	&registry.VCRSearchParams{
	//		PurposeOfUse: zorginzage.ServiceName,
	//		SubjectID:    customerDID,
	//		Subject:      ssnURN(patientSSN),
	//	},
	//)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if len(credentials) == 0 {
	//	return []types.Report{}, nil
	//}
	//
	//// TODO: loop over all credentials
	//issuer := credentials[0].Issuer.String()
	//
	//fhirServer, err := service.registry.GetCompoundServiceEndpoint(ctx, issuer, zorginzage.ServiceName, "fhir")
	//if err != nil {
	//	return nil, fmt.Errorf("error while looking up authorizer's FHIR server (did=%s): %w", issuer, err)
	//}
	//
	//issuerOrg, err := service.registry.Get(ctx, issuer)
	//if err != nil {
	//	return nil, fmt.Errorf("error while searching organization :%w", err)
	//}
	//
	//// TODO: Should be user access token?
	//accessToken, err := service.nutsClient.RequestServiceAccessToken(ctx, customerDID, issuer, zorginzage.ServiceName)
	//if err != nil {
	//	return nil, err
	//}
	//
	//episodeOfCareID, err := parseEpisodeOfCareID(credentials[0])
	//if err != nil {
	//	return nil, err
	//}
	//
	//fhirClient := fhir.NewFactory(fhir.WithURL(fhirServer), fhir.WithAuthToken(accessToken))()
	//
	//fhirEpisode := &fhir.EpisodeOfCare{}
	//err = fhirClient.ReadOne(ctx, "/EpisodeOfCare/"+episodeOfCareID, fhirEpisode)
	//if err != nil {
	//	return nil, fmt.Errorf("could not retrieve episode of care: %w", err)
	//}
	//episode := zorginzage.ToEpisode(fhirEpisode)
	//
	//observations := []resources.Observation{}
	//if err := fhirClient.ReadMultiple(ctx, "/Observation", map[string]string{
	//	"context": fmt.Sprintf("EpisodeOfCare/%s", episodeOfCareID),
	//	//"subject": fmt.Sprintf("Patient/%s", patientSSN),
	//}, &observations); err != nil {
	//	return nil, err
	//}
	//
	//results := make([]types.Report, len(observations))
	//
	//for _, observation := range observations {
	//	domainObservation := reports.ConvertToDomain(&observation, fhir.FromStringPtr(observation.Subject.ID))
	//	domainObservation.Source = issuerOrg.Details.Name
	//	domainObservation.EpisodeName = &episode.Diagnosis
	//	results = append(results, domainObservation)
	//}
	//
	//return results, nil
}
