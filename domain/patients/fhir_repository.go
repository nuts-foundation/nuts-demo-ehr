package patients

import (
	"context"
	"fmt"

	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir/eoverdracht"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
)

func ToDomainPatient(fhirPatient resources.Patient) types.Patient {
	return eoverdracht.ToDomainPatient(fhirPatient)
}

func ToFHIRPatient(domainPatient types.Patient) resources.Patient {
	dob := datatypes.Date(domainPatient.Dob.Format(types.DobFormat))

	fhirPatient := resources.Patient{
		Domain: resources.Domain{
			Base: resources.Base{
				ResourceType: "Patient",
				ID:           fhir.ToIDPtr(string(domainPatient.ObjectID)),
			},
		},
		Name: []datatypes.HumanName{{
			Use:    fhir.ToCodePtr("official"),
			Family: fhir.ToStringPtr(domainPatient.Surname),
			Given:  []datatypes.String{datatypes.String(domainPatient.FirstName)},
		}},
		BirthDate: &dob,
		Gender:    fhir.ToCodePtr(string(domainPatient.Gender)),
	}

	if domainPatient.Ssn != nil {
		fhirPatient.Identifier = append(fhirPatient.Identifier, datatypes.Identifier{System: fhir.ToUriPtr(types.BsnSystem), Value: fhir.ToStringPtr(*domainPatient.Ssn)})
	}
	if domainPatient.AvatarUrl != nil {
		fhirPatient.Photo = append(fhirPatient.Photo, datatypes.Attachment{URL: fhir.ToUriPtr(*domainPatient.AvatarUrl)})
	}
	if domainPatient.Zipcode != "" {
		fhirPatient.Address = append(fhirPatient.Address, datatypes.Address{
			PostalCode: fhir.ToStringPtr(domainPatient.Zipcode),
		})
	}

	return fhirPatient
}

type FHIRPatientRepository struct {
	fhirClientFactory fhir.Factory
	factory           Factory
}

func NewFHIRPatientRepository(factory Factory, fhirClientFactory fhir.Factory) *FHIRPatientRepository {
	return &FHIRPatientRepository{
		fhirClientFactory: fhirClientFactory,
		factory:           factory,
	}
}

func (r FHIRPatientRepository) FindByID(ctx context.Context, customerID int, id string) (*types.Patient, error) {
	patient := resources.Patient{}
	err := r.fhirClientFactory(fhir.WithTenant(customerID)).ReadOne(ctx, "Patient/"+id, &patient)
	if err != nil {
		return nil, err
	}
	result := ToDomainPatient(patient)
	return &result, nil
}

func (r FHIRPatientRepository) Update(ctx context.Context, customerID int, id string, updateFn func(c types.Patient) (*types.Patient, error)) (*types.Patient, error) {
	domainPatient, err := r.FindByID(ctx, customerID, id)
	if err != nil {
		return nil, fmt.Errorf("could not update patient: could not read current patient from FHIR store: %w", err)
	}
	updatedDomainPatient, err := updateFn(*domainPatient)
	if err != nil {
		return nil, err
	}
	fhirClient := r.fhirClientFactory(fhir.WithTenant(customerID))
	updatedFHIRPatient := ToFHIRPatient(*updatedDomainPatient)
	return updatedDomainPatient, fhirClient.CreateOrUpdate(ctx, updatedFHIRPatient)
}

func (r FHIRPatientRepository) NewPatient(ctx context.Context, customerID int, patientProperties types.PatientProperties) (*types.Patient, error) {
	patient, err := r.factory.NewPatientWithAvatar(patientProperties)
	if err != nil {
		return nil, err
	}
	err = r.fhirClientFactory(fhir.WithTenant(customerID)).CreateOrUpdate(ctx, ToFHIRPatient(*patient))
	if err != nil {
		return nil, err
	}
	return patient, nil
}

func (r FHIRPatientRepository) All(ctx context.Context, customerID int, name *string) ([]types.Patient, error) {
	var params map[string]string
	if name != nil {
		params = map[string]string{"name": *name}
	} else {
		// Filter patients by having a name. This filters out the anonymous patients created just for the eOverdracht advance notice.
		params = map[string]string{"name:above": "_"}
	}
	fhirPatients := []resources.Patient{}
	err := r.fhirClientFactory(fhir.WithTenant(customerID)).ReadMultiple(ctx, "Patient", params, &fhirPatients)
	if err != nil {
		return nil, err
	}

	patients := make([]types.Patient, 0)
	for _, patient := range fhirPatients {
		patients = append(patients, ToDomainPatient(patient))
	}

	return patients, nil
}
