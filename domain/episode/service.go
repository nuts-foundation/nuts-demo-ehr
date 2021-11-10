package episode

import (
	"context"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir/zorginzage"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
)

type Service interface {
	Create(ctx context.Context, customerID int, dossierID, patientID string) (*types.Episode, error)
	Get(ctx context.Context, customerID int, dossierID string) (*types.Episode, error)
}

type service struct {
	factory fhir.Factory
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
