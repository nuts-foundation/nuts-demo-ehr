package zorginzage

import (
	"context"
	"fmt"
	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
)

type Service interface {
	CreateEpisode(ctx context.Context, dossierID, patientID string) (*fhir.EpisodeOfCare, error)
	GetEpisode(ctx context.Context, dossierID string) (*fhir.EpisodeOfCare, error)
}

type service struct {
	fhirClient fhir.Client
}

func NewService(fhirClient fhir.Client) Service {
	return &service{
		fhirClient: fhirClient,
	}
}

func (service *service) CreateEpisode(ctx context.Context, dossierID, patientID string) (*fhir.EpisodeOfCare, error) {
	episode := &fhir.EpisodeOfCare{
		Base: resources.Base{
			ID:           fhir.ToIDPtr(dossierID),
			ResourceType: "EpisodeOfCare",
		},
		Patient: datatypes.Reference{
			Reference: fhir.ToStringPtr(fmt.Sprintf("Patient/%s", patientID)),
		},
		Status: fhir.EpisodeStatusActive,
	}

	if err := service.fhirClient.CreateOrUpdate(ctx, episode); err != nil {
		return nil, err
	}

	return episode, nil
}

func (service *service) GetEpisode(ctx context.Context, dossierID string) (*fhir.EpisodeOfCare, error) {
	episode := &fhir.EpisodeOfCare{}

	if err := service.fhirClient.ReadOne(ctx, fmt.Sprintf("EpisodeOfCare/%s", dossierID), episode); err != nil {
		return nil, err
	}

	return episode, nil
}
