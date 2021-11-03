package reports

import (
	"context"
	"fmt"
	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"strings"
)

type Repository interface {
	AllByPatient(ctx context.Context, customerID int, patientID string) ([]domain.Report, error)
}

type fhirRepository struct {
	factory fhir.Factory
}

func NewFHIRRepository(factory fhir.Factory) *fhirRepository {
	return &fhirRepository{
		factory: factory,
	}
}

func renderQuantity(quantity *datatypes.Quantity) string {
	return fmt.Sprintf("%f %s", *quantity.Value, fhir.FromStringPtr(quantity.Unit))
}

func convertToDomain(observation *resources.Observation, id string) domain.Report {
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

	return domain.Report{
		Type:      fhir.FromStringPtr(observation.Code.Coding[0].Display),
		Id:        domain.ObjectID(fhir.FromIDPtr(observation.ID)),
		Source:    source,
		PatientID: domain.ObjectID(id),
		Value:     value,
	}
}

func (repo *fhirRepository) AllByPatient(ctx context.Context, customerID int, patientID string) ([]domain.Report, error) {
	observations := []resources.Observation{}

	if err := repo.factory(fhir.WithTenant(customerID)).ReadMultiple(ctx, "Observation", map[string]string{
		"patient": patientID,
	}, &observations); err != nil {
		return nil, err
	}

	var reports []domain.Report

	for _, observation := range observations {
		ref := fhir.FromStringPtr(observation.Subject.Reference)

		if !strings.HasPrefix(ref, "Patient/") {
			continue
		}

		reports = append(reports, convertToDomain(&observation, ref[len("Patient/"):]))
	}

	return reports, nil
}
