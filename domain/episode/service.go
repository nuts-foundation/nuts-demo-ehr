package episode

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nuts-foundation/go-did/vc"
	"time"

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
	CreateCollaboration(ctx context.Context, customerDID, dossierID, patientSSN, senderDID string) error
	GetCollaborations(ctx context.Context, customerDID, dossierID, patientSSN string) ([]types.Collaboration, error)
}

func ssnURN(ssn string) string {
	return fmt.Sprintf("urn:oid:2.16.840.1.113883.2.4.6.3:%s", ssn)
}

type service struct {
	factory fhir.Factory
	vcr     registry.VerifiableCredentialRegistry
}

func NewService(factory fhir.Factory, vcr registry.VerifiableCredentialRegistry) Service {
	return &service{factory: factory, vcr: vcr}
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

func (service *service) GetCollaborations(ctx context.Context, _customerDID, dossierID, patientSSN string) ([]types.Collaboration, error) {
	credentials, err := service.vcr.FindAuthorizationCredentials(
		ctx,
		&registry.VCRSearchParams{
			PurposeOfUse: zorginzage.ServiceName,
			Subject:      ssnURN(patientSSN),
			ResourcePath: fmt.Sprintf("/EpisodeOfCare/%s", dossierID),
		},
	)
	if err != nil {
		return nil, err
	}

	episodeID := types.ObjectID(dossierID)

	var collaborations []types.Collaboration

	for _, credentialResponse := range credentials {
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

		collaborations = append(collaborations, types.Collaboration{
			EpisodeID:       episodeID,
			OrganizationDID: subject[0].ID,
		})
	}

	return collaborations, nil
}
