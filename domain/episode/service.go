package episode

import (
	"context"
	"fmt"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/dossier"
	"github.com/nuts-foundation/nuts-node/vcr/credential"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir/zorginzage"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/registry"
)

type Service interface {
	Create(ctx context.Context, customerID int, dossierID, patientID string) (*types.Episode, error)
	Get(ctx context.Context, customerID int, dossierID string) (*types.Episode, error)
	CreateCollaboration(ctx context.Context, customerDID, dossierID, senderDID string) error
}

type service struct {
	factory fhir.Factory
	repo    dossier.Repository
	vcr     registry.VerifiableCredentialRegistry
}

func NewService(factory fhir.Factory) Service {
	return &service{factory: factory}
}

func toEpisode(episode *fhir.EpisodeOfCare) *types.Episode {
	status := types.EpisodeStatus(episode.Status)

	return &types.Episode{
		Id:     types.ObjectID(fhir.FromIDPtr(episode.ID)),
		Status: &status,
	}
}

func (service *service) Create(ctx context.Context, customerID int, dossierID, patientID string) (*types.Episode, error) {
	svc := zorginzage.NewService(service.factory(fhir.WithTenant(customerID)))

	episode, err := svc.CreateEpisode(ctx, dossierID, patientID)
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

func (service *service) CreateCollaboration(ctx context.Context, customerDID, dossierID, senderDID string) error {
	subject := "urn"

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
