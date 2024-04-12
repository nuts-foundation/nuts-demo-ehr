package sharedcareplan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/dossier"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/patients"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/client"
	openapi_types "github.com/oapi-codegen/runtime/types"
	r4 "github.com/samply/golang-fhir-models/fhir-models/fhir"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"time"
)

type Service struct {
	DossierRepository dossier.Repository
	PatientRepository patients.Repository
	Repository        Repository
	SCPFHIRClient     fhir.Client
	NutsClient        *client.HTTPClient
	LocalFHIRClient   fhir.Factory
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
	if err := s.SCPFHIRClient.Create(ctx, carePlan, &carePlan); err != nil {
		return nil, err
	}
	// Create SharedCarePlan record
	reference := s.SCPFHIRClient.BuildRequestURI(fmt.Sprintf("CarePlan/%s", *carePlan.Id)).String()
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
		if err := s.SCPFHIRClient.ReadOne(ctx, *careTeamRef.Reference, &careTeam); err != nil {
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
			if err := s.SCPFHIRClient.ReadOne(ctx, *activity.Reference.Reference, &task); err != nil {
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
	if err := s.SCPFHIRClient.Create(ctx, task, &task); err != nil {
		return fmt.Errorf("unable to create Task: %w", err)
	}
	taskRef := "Task/" + *task.Id
	carePlan := sharedCarePlan.FHIRCarePlan
	carePlan.Activity = append(carePlan.Activity, r4.CarePlanActivity{
		Reference: &r4.Reference{Reference: &taskRef},
	})
	// This is racy! (if other parties update the plan at the same time)
	if err := s.SCPFHIRClient.CreateOrUpdate(ctx, carePlan, &carePlan); err != nil {
		return fmt.Errorf("unable to update CarePlan with new activity: %w", err)
	}

	// Notify Task owner
	if *owner.Value != *requester.Value {
		go func() {
			logrus.Infof("Notifying Task owner %s", *owner.Value)
			dossier, err := s.DossierRepository.FindByID(ctx, customerID, dossierID)
			if err != nil {
				logrus.Errorf("unable to find dossier %s: %v", dossierID, err)
				return
			}
			var patient r4.Patient
			if err = s.LocalFHIRClient(fhir.WithTenant(customerID)).ReadOne(ctx, "Patient/"+dossier.PatientID, &patient); err != nil {
				logrus.Errorf("unable to find patient %s: %v", dossier.PatientID, err)
				return
			}
			requestBody := types.NotifyCarePlanUpdateJSONRequestBody{
				CarePlanURL: sharedCarePlan.FHIRCarePlanURL,
				Patient:     patient,
				Task:        task,
			}
			requestData, _ := json.Marshal(requestBody)
			httpRequest, err := http.NewRequest("POST", "http://localhost:1304/web/external/careplan/notify", bytes.NewReader(requestData))
			if err != nil {
				logrus.Errorf("unable to create HTTP request: %v", err)
				return
			}
			httpRequest.Header.Set("Content-Type", "application/json")
			httpResponse, err := http.DefaultClient.Do(httpRequest)
			if err != nil {
				logrus.Errorf("unable to send HTTP request: %v", err)
				return
			}
			responseData, _ := io.ReadAll(httpResponse.Body)
			logrus.Infof("Notified Task owner %s. Status: %s, response: %s", *owner.Value, httpResponse.Status, string(responseData))
		}()
	}

	return nil
}

func (s Service) find(ctx context.Context, customerID int, dossierID string) (*types.SharedCarePlan, error) {
	sharedCarePlan, err := s.Repository.FindByDossierID(ctx, customerID, dossierID)
	if err != nil {
		return nil, err
	}
	carePlan := r4.CarePlan{}
	if err := s.SCPFHIRClient.ReadOne(ctx, sharedCarePlan.Reference, &carePlan); err != nil {
		return nil, err
	}
	return &types.SharedCarePlan{
		DossierID:         dossierID,
		FHIRCarePlan:      carePlan,
		FHIRCarePlanURL:   sharedCarePlan.Reference,
		FHIRActivityTasks: map[string]r4.Task{},
	}, nil
}

func (s Service) HandleNotify(ctx context.Context, customerID int, patientResource r4.Patient, task r4.Task, carePlanURL string) error {
	// create patient (make sure it doesn't already exist)
	existingPatients, err := s.PatientRepository.All(ctx, customerID, nil)
	if err != nil {
		return err
	}
	var patient *types.Patient
	for _, p := range existingPatients {
		if *p.Ssn == *patientResource.Identifier[0].Value {
			patient = &p
			break
		}
	}
	if patient == nil {
		patientProperties := types.PatientProperties{
			Ssn: patientResource.Identifier[0].Value,
			Dob: &openapi_types.Date{Time: time.Date(1990, 1, 12, 0, 0, 0, 0, time.Local)},
		}
		if patientResource.Gender != nil {
			patientProperties.Gender = types.Gender(patientResource.Gender.String())
		}
		if len(patientResource.Name) > 0 {
			patientProperties.FirstName = strings.Join(patientResource.Name[0].Given, " ")
			if patientResource.Name[0].Family != nil {
				patientProperties.Surname = *patientResource.Name[0].Family
			}
		}
		patient, err = s.PatientRepository.NewPatient(ctx, customerID, patientProperties)
		if err != nil {
			return fmt.Errorf("error creating patient: %w", err)
		}
	}

	// Create dossier, if it doesn't exist
	existingCarePlan, _ := s.Repository.FindByCarePlanURL(ctx, customerID, carePlanURL)
	if err != nil {
		return fmt.Errorf("error finding shared care plan by URL: %w", err)
	}
	if existingCarePlan == nil {
		// Create a dossier
		var carePlan r4.CarePlan
		if err := s.SCPFHIRClient.ReadOne(ctx, carePlanURL, &carePlan); err != nil {
			return fmt.Errorf("error reading care plan: %w", err)
		}
		title := "<missing title>"
		if carePlan.Title != nil {
			title = *carePlan.Title
		}
		newDossier, err := s.DossierRepository.Create(ctx, customerID, title, patient.ObjectID)
		if err != nil {
			return fmt.Errorf("error creating dossier: %w", err)
		}
		// Now create a CarePlan record
		if err := s.Repository.Create(ctx, customerID, newDossier.Id, carePlanURL); err != nil {
			return err
		}
	}
	return nil
}

func MakeIdentifier(system, value string) *r4.Identifier {
	return &r4.Identifier{System: &system, Value: &value}
}
