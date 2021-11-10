package zorginzage

import (
	"context"
	"fmt"
	"time"

	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
)

type Service interface {
	CreateEpisode(ctx context.Context, patientID string, request types.CreateEpisodeRequest) (*fhir.EpisodeOfCare, error)
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

func (service *service) CreateEpisode(ctx context.Context, patientID string, request types.CreateEpisodeRequest) (*fhir.EpisodeOfCare, error) {
	periodStart := datatypes.DateTime(request.Period.Start.Format(time.RFC3339))
	episode := &fhir.EpisodeOfCare{
		Base: resources.Base{
			ID:           fhir.ToIDPtr(string(request.DossierID)),
			ResourceType: "EpisodeOfCare",
		},
		Type: []datatypes.CodeableConcept{{
			Text: fhir.ToStringPtr(request.Diagnosis),
		}},
		Patient: datatypes.Reference{
			Reference: fhir.ToStringPtr(fmt.Sprintf("Patient/%s", patientID)),
		},
		Period: &datatypes.Period{
			Start: &periodStart,
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
