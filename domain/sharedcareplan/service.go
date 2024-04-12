package sharedcareplan

import (
	"context"
	"fmt"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/dossier"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/patients"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/client"
	r4 "github.com/samply/golang-fhir-models/fhir-models/fhir"
	"github.com/sirupsen/logrus"
)

type Service struct {
	DossierRepository dossier.Repository
	PatientRepository patients.Repository
	Repository        Repository
	FHIRClient        fhir.Client
	NutsClient        *client.HTTPClient
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
	carePlan := r4.CarePlan{
		Status: r4.RequestStatusActive,
		Intent: r4.CarePlanIntentProposal,
		Title:  &title,
		Subject: r4.Reference{
			Identifier: MakeIdentifier(types.BsnSystem, *patient.Ssn),
		},
	}
	if err := s.FHIRClient.Create(ctx, carePlan, &carePlan); err != nil {
		return nil, err
	}
	// Create SharedCarePlan record
	reference := s.FHIRClient.BuildRequestURI(fmt.Sprintf("CarePlan/%s", *carePlan.Id)).String()
	if err := s.Repository.Create(ctx, customerID, dossierID, reference); err != nil {
		return nil, err
	}
	return &types.SharedCarePlan{
		DossierID:    dossierID,
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
		sharedCarePlan, err := s.FindByID(ctx, customerID, current.Id, false)
		if err != nil {
			return nil, fmt.Errorf("unable to get shared care plan for dossier %s: %w", current.Id, err)
		}
		result = append(result, *sharedCarePlan)
	}
	return result, nil
}

func (s Service) FindByID(ctx context.Context, customerID int, dossierID string, resolvefhirReferences bool) (*types.SharedCarePlan, error) {
	sharedCarePlan, err := s.find(ctx, customerID, dossierID)
	if err != nil {
		return nil, err
	}

	// Lookup CareTeam for filling overview of participants
	careTeams := []r4.CareTeam{}
	for _, careTeamRef := range sharedCarePlan.FHIRCarePlan.CareTeam {
		if careTeamRef.Reference == nil {
			logrus.Infof("CareTeam reference is nil, skipping")
			continue
		}
		var careTeam r4.CareTeam
		if err := s.FHIRClient.ReadOne(ctx, *careTeamRef.Reference, &careTeam); err != nil {
			return nil, err
		}
		careTeams = append(careTeams, careTeam)
	}
	organizationMap := make(map[string]types.Organization)
	for _, careTeam := range careTeams {
		for _, participant := range careTeam.Participant {
			// Prevent nil deref
			if participant.OnBehalfOf == nil ||
				participant.OnBehalfOf.Identifier == nil {
				logrus.Infof("CareTeam/%s participant OnBehalfOf is invalid (missing identifier or display), skipping", *careTeam.Id)
				continue
			}
			if err != nil {
				return nil, err
			}
			code := fmt.Sprintf("%s|%s", *participant.OnBehalfOf.Identifier.System, *participant.OnBehalfOf.Identifier.Value)
			organizationName, err := s.organizationDisplayName(*participant.OnBehalfOf.Identifier.Value)
			if err != nil {
				logrus.Warnf("unable to resolve organization name for URA %s: %v", *participant.OnBehalfOf.Identifier.Value, err)
				organizationName = *participant.OnBehalfOf.Identifier.Value
			}
			organizationMap[code] = types.Organization{
				Name: organizationName,
			}
		}
	}
	if resolvefhirReferences {
		// Resolve activities to actual tasks
		for i, activity := range sharedCarePlan.FHIRCarePlan.Activity {
			if activity.Reference == nil || activity.Reference.Reference == nil {
				logrus.Infof("CarePlanActivity/%d has no reference, skipping", i)
				continue
			}
			var task r4.Task
			if err := s.FHIRClient.ReadOne(ctx, *activity.Reference.Reference, &task); err != nil {
				return nil, err
			}
			if task.Code == nil {
				logrus.Infof("Task/%s has no code, skipping", *task.Id)
				continue
			}
			sharedCarePlan.FHIRActivityTasks[*activity.Reference.Reference] = task
		}
	}
	sharedCarePlan.Participants = make([]types.Organization, 0)
	for _, org := range organizationMap {
		sharedCarePlan.Participants = append(sharedCarePlan.Participants, org)
	}
	return sharedCarePlan, nil
}

func (s Service) organizationDisplayName(uraNummer string) (string, error) {
	searchResults, err := s.NutsClient.SearchDiscoveryService(context.Background(), map[string]string{
		"credentialSubject.ura": uraNummer,
	}, nil, nil)
	if err != nil {
		return "", err
	}
	if len(searchResults) == 0 {
		return "", fmt.Errorf("no organization found with URA number %s", uraNummer)
	}
	name, ok := searchResults[0].Fields["organization_name"].(string)
	if !ok {
		return "", fmt.Errorf("organization with URA number %s has no organization_name", uraNummer)
	}
	return name + " (URA: " + uraNummer + ")", nil
}

func (s Service) CreateActivity(ctx context.Context, customerID int, dossierID string, code types.FHIRCodeableConcept, requester types.FHIRIdentifier, owner types.FHIRIdentifier) error {
	sharedCarePlan, err := s.find(ctx, customerID, dossierID)
	if err != nil {
		return err
	}
	// Create Task, then update care plan
	// TODO: Should be in 1 bundle with PATCH
	carePlanRef := "CarePlan/" + *sharedCarePlan.FHIRCarePlan.Id
	task := r4.Task{
		Requester: &r4.Reference{
			Identifier: &requester,
		},
		Owner: &r4.Reference{
			Identifier: &owner,
		},
		BasedOn: []r4.Reference{
			{
				Reference: &carePlanRef,
			},
		},
		Code:   &code,
		Status: r4.TaskStatusAccepted,
		Intent: "order",
	}
	if err := s.FHIRClient.Create(ctx, task, &task); err != nil {
		return fmt.Errorf("unable to create Task: %w", err)
	}
	taskRef := "Task/" + *task.Id
	carePlan := sharedCarePlan.FHIRCarePlan
	carePlan.Activity = append(carePlan.Activity, r4.CarePlanActivity{
		Reference: &r4.Reference{Reference: &taskRef},
	})
	// This is racy! (if other parties update the plan at the same time)
	if err := s.FHIRClient.CreateOrUpdate(ctx, carePlan, &carePlan); err != nil {
		return fmt.Errorf("unable to update CarePlan with new activity: %w", err)
	}
	return nil
}

func (s Service) find(ctx context.Context, customerID int, dossierID string) (*types.SharedCarePlan, error) {
	sharedCarePlan, err := s.Repository.FindByDossierID(ctx, customerID, dossierID)
	if err != nil {
		return nil, err
	}
	carePlan := r4.CarePlan{}
	if err := s.FHIRClient.ReadOne(ctx, sharedCarePlan.Reference, &carePlan); err != nil {
		return nil, err
	}
	return &types.SharedCarePlan{
		DossierID:         dossierID,
		FHIRCarePlan:      carePlan,
		FHIRActivityTasks: map[string]r4.Task{},
	}, nil
}

func MakeIdentifier(system, value string) *r4.Identifier {
	return &r4.Identifier{System: &system, Value: &value}
}
