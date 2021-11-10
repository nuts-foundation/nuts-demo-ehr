package episode

import (
	"context"
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir/zorginzage"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
)

type Service interface {
	Create(ctx context.Context, customerID int, patientID string, request types.CreateEpisodeRequest) (*types.Episode, error)
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
