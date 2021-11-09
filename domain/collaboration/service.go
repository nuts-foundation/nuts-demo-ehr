package collaboration

import (
	"context"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir/zorginzage"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
)

type Service interface {
	CreateCollaboration(ctx context.Context, customerID int, dossierID, patientID string) (*types.Collaboration, error)
}

type service struct {
	factory fhir.Factory
}

func NewService(factory fhir.Factory) Service {
	return &service{factory: factory}
}

func (service *service) CreateCollaboration(ctx context.Context, customerID int, dossierID, patientID string) (*types.Collaboration, error) {
	svc := zorginzage.NewService(service.factory(fhir.WithTenant(customerID)))

	episode, err := svc.CreateEpisode(ctx, dossierID, patientID)
	if err != nil {
		return nil, err
	}

	return &types.Collaboration{
		Id: types.ObjectID(fhir.FromIDPtr(episode.ID)),
	}, nil
}
