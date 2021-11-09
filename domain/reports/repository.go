package reports

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
)

type Repository interface {
	AllByPatient(ctx context.Context, customerID int, patientID string) ([]types.Report, error)
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
		}
		return observation, nil
	}

	return nil, errors.New("unknown report type")
}

func convertToDomain(observation *resources.Observation, patientId string) types.Report {
	var value string

	if observation.ValueString != nil {
		value = fhir.FromStringPtr(observation.ValueString)
	} else if observation.ValueQuantity != nil {
		value = renderQuantity(observation.ValueQuantity)
	} else if observation.Component != nil {
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

	return types.Report{
		Type:      fhir.FromStringPtr(observation.Code.Coding[0].Display),
		Id:        types.ObjectID(fhir.FromIDPtr(observation.ID)),
		Source:    source,
		PatientID: types.ObjectID(patientId),
		Value:     value,
	}
}

func (repo *fhirRepository) AllByPatient(ctx context.Context, customerID int, patientID string) ([]types.Report, error) {
	observations := []resources.Observation{}

	if err := repo.factory(fhir.WithTenant(customerID)).ReadMultiple(ctx, "Observation", map[string]string{
		"patient": patientID,
	}, &observations); err != nil {
		return nil, err
	}

	var reports []types.Report

	for _, observation := range observations {
		ref := fhir.FromStringPtr(observation.Subject.Reference)

		if !strings.HasPrefix(ref, "Patient/") {
			continue
		}

		reports = append(reports, convertToDomain(&observation, ref[len("Patient/"):]))
	}

	return reports, nil
}
