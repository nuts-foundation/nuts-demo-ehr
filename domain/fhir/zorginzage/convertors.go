package zorginzage

import (
	"time"

	types2 "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
)

func ToEpisode(episode *fhir.EpisodeOfCare) *types.Episode {
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
		Period:    types.Period{Start: &types2.Date{Time: periodStart}},
		Diagnosis: diagnosis,
	}
}

