package reports

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
)

type Repository interface {
	AllByPatient(ctx context.Context, customerID int, patientID string, episodeID *string) ([]types.Report, error)
	Create(ctx context.Context, customerID int, patientID string, report types.Report) error
}

type fhirRepository struct {
	factory fhir.Factory
}

func NewFHIRRepository(factory fhir.Factory) *fhirRepository {
	return &fhirRepository{
		factory: factory,
	}
}

func (repo *fhirRepository) Create(ctx context.Context, customerID int, patientID string, report types.Report) error {
	if report.Id == "" {
		report.Id = types.ObjectID(uuid.NewString())
	}

	observation, err := convertToFHIR(report)
	if err != nil {
		return fmt.Errorf("unable to convert report to FHIR observation: %w", err)
	}

	err = repo.factory(fhir.WithTenant(customerID)).CreateOrUpdate(ctx, observation)
	if err != nil {
		return fmt.Errorf("unable to write observation to FHIR store: %w", err)
	}

	return nil
}

func renderQuantity(quantity *datatypes.Quantity) string {
	return fmt.Sprintf("%f %s", *quantity.Value, fhir.FromStringPtr(quantity.Unit))
}

func convertToFHIR(report types.Report) (*resources.Observation, error) {
	if report.Type == "heartRate" {
		value, err := strconv.ParseFloat(report.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse report value as float: %w", err)
		}
		valueDecimal := datatypes.Decimal(value)
		observation := &resources.Observation{
			Domain: resources.Domain{
				Base: resources.Base{
					ID:           fhir.ToIDPtr(string(report.Id)),
					ResourceType: "Observation",
				},
			},
			Code: &datatypes.CodeableConcept{
				Coding: []datatypes.Coding{{
					System:  &fhir.LoincCodingSystem,
					Code:    fhir.ToCodePtr("8893-0"),
					Display: fhir.ToStringPtr("Heart rate Peripheral artery by Palpation"),
				}},
			},
			Subject: &datatypes.Reference{Reference: fhir.ToStringPtr("Patient/" + string(report.PatientID))},
			ValueQuantity: &datatypes.Quantity{
				Value: &valueDecimal,
			},
			EffectiveDateTime: fhir.ToDateTimePtr(time.Now().Format(fhir.DateTimeLayout)),
		}
		if report.EpisodeID != nil {
			observation.Context = &datatypes.Reference{Reference: fhir.ToStringPtr("EpisodeOfCare/" + string(*report.EpisodeID))}
		}
		return observation, nil
	}

	return nil, errors.New("unknown report type")
}

func ConvertToDomain(observation *resources.Observation, patientID string) types.Report {
	var value string

	switch {
	case observation.ValueString != nil:
		value = fhir.FromStringPtr(observation.ValueString)
	case observation.ValueQuantity != nil:
		value = renderQuantity(observation.ValueQuantity)
	case observation.Component != nil:
		var values []string
		for _, component := range observation.Component {
			if component.ValueString != nil {
				values = append(values, fhir.FromStringPtr(component.ValueString))
			} else if component.ValueQuantity != nil {
				values = append(values, renderQuantity(component.ValueQuantity))
			}
		}
		value = strings.Join(values, ", ")
	}

	source := "Unknown"

	if len(observation.Performer) > 0 {
		source = fhir.FromStringPtr(observation.Performer[0].Display)
	}

	report := types.Report{
		Type:      fhir.FromStringPtr(observation.Code.Coding[0].Display),
		Id:        types.ObjectID(fhir.FromIDPtr(observation.ID)),
		Source:    source,
		PatientID: types.ObjectID(patientID),
		Value:     value,
	}

	if observation.Context != nil {
		id := types.ObjectID(strings.Split(fhir.FromStringPtr(observation.Context.Reference), "/")[1])
		report.EpisodeID = &id
	}

	return report
}

func (repo *fhirRepository) AllByPatient(ctx context.Context, customerID int, patientID string, episodeID *string) ([]types.Report, error) {
	observations := []resources.Observation{}

	queryMap := map[string]string{
		"subject": fmt.Sprintf("Patient/%s", patientID),
	}
	if episodeID != nil {
		queryMap["context"] = fmt.Sprintf("EpisodeOfCare/%s", *episodeID)
	}

	if err := repo.factory(fhir.WithTenant(customerID)).ReadMultiple(ctx, "Observation", queryMap, &observations); err != nil {
		return nil, err
	}

	reports := []types.Report{}

	for _, observation := range observations {
		ref := fhir.FromStringPtr(observation.Subject.Reference)

		if !strings.HasPrefix(ref, "Patient/") {
			continue
		}

		reports = append(reports, ConvertToDomain(&observation, ref[len("Patient/"):]))
	}

	return reports, nil
}
