package sharedcareplan

import (
	"context"
	"fmt"
	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/dossier"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/patients"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	"github.com/sirupsen/logrus"
)

type Service struct {
	DossierRepository dossier.Repository
	PatientRepository patients.Repository
	Repository        Repository
	FHIRClient        fhir.Client
}

// Create creates a new shared CarePlan on the Care Plan Service for the given dossierID.
func (s Service) Create(ctx context.Context, customerID int, dossierID string, title string) (*types.SharedCarePlan, error) {
	targetDossier, err := s.DossierRepository.FindByID(ctx, customerID, dossierID)
	if err != nil {
		return nil, err
	}
	patient, err := s.PatientRepository.FindByID(ctx, customerID, targetDossier.PatientID)
	if err != nil {
		return nil, err
	}

	// Create CarePlan at Shared Care Plan Service
	status := datatypes.Code("active")
	intent := datatypes.Code("proposal")
	carePlan := resources.CarePlan{
		Domain: resources.Domain{
			Base: resources.Base{
				ResourceType: "CarePlan",
			},
		},
		Status: &status,
		Intent: &intent,
		Title:  fhir.ToStringPtr(title),
		Subject: &datatypes.Reference{
			Identifier: &datatypes.Identifier{System: fhir.ToUriPtr(types.BsnSystem), Value: fhir.ToStringPtr(*patient.Ssn)},
		},
	}
	if err := s.FHIRClient.Create(ctx, carePlan, &carePlan); err != nil {
		return nil, err
	}
	// Create SharedCarePlan record
	reference := s.FHIRClient.BuildRequestURI(fmt.Sprintf("CarePlan/%s", *carePlan.ID)).String()
	if err := s.Repository.Create(ctx, customerID, dossierID, reference); err != nil {
		return nil, err
	}
	return &types.SharedCarePlan{
		FHIRCarePlan: carePlan,
	}, nil
}

func (s Service) AllForPatient(ctx context.Context, customerID int, patientID string) ([]types.SharedCarePlan, error) {
	dossiers, err := s.DossierRepository.AllByPatient(ctx, customerID, patientID)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch dossiers for patient %s: %w", patientID, err)
	}
	var result []types.SharedCarePlan
	for _, current := range dossiers {
		sharedCarePlan, err := s.FindByID(ctx, customerID, current.Id)
		if err != nil {
			return nil, fmt.Errorf("unable to get shared care plan for dossier %s: %w", current.Id, err)
		}
		result = append(result, *sharedCarePlan)
	}
	return result, nil
}

func (s Service) FindByID(ctx context.Context, customerID int, dossierID string) (*types.SharedCarePlan, error) {
	sharedCarePlan, err := s.Repository.FindByDossierID(ctx, customerID, dossierID)
	if err != nil {
		return nil, err
	}
	carePlan := resources.CarePlan{}
	if err := s.FHIRClient.ReadOne(ctx, sharedCarePlan.Reference, &carePlan); err != nil {
		return nil, err
	}

	// Lookup CareTeam for filling overview of participants
	careTeams := []resources.CareTeam{}
	for _, careTeamRef := range carePlan.CareTeam {
		if careTeamRef.Reference == nil {
			logrus.Infof("CareTeam reference is nil, skipping")
			continue
		}
		var careTeam resources.CareTeam
		if err := s.FHIRClient.ReadOne(ctx, string(*careTeamRef.Reference), &careTeam); err != nil {
			return nil, err
		}
		careTeams = append(careTeams, careTeam)
	}
	organizationMap := make(map[string]types.Organization)
	for _, careTeam := range careTeams {
		for _, participant := range careTeam.Participant {
			// Prevent nil deref
			if participant.OnBehalfOf == nil ||
				participant.OnBehalfOf.Identifier == nil ||
				participant.OnBehalfOf.Display == nil {
				logrus.Infof("CareTeam/%s participant OnBehalfOf is invalid (missing identifier or display), skipping", *careTeam.ID)
				continue
			}
			if participant.OnBehalfOf.Display == nil {
				logrus.Infof("CareTeam/%s participant has no OnBehalfOf, skipping", *careTeam.ID)
				continue
			}
			if err != nil {
				return nil, err
			}
			code := fmt.Sprintf("%s|%s", *participant.OnBehalfOf.Identifier.System, *participant.OnBehalfOf.Identifier.Value)
			organizationMap[code] = types.Organization{
				Name: fmt.Sprintf("%s ( %s)", participant.OnBehalfOf.Display, code),
			}
		}
	}
	result := types.SharedCarePlan{
		DossierID:    dossierID,
		FHIRCarePlan: carePlan,
		Participants: make([]types.Organization, 0),
	}
	for _, org := range organizationMap {
		result.Participants = append(result.Participants, org)
	}
	return &result, nil
}
