package eoverdracht

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer"
)

type TransferService interface {
	GetTask(ctx context.Context, taskID string) (*TransferTask, error)
	CreateTask(ctx context.Context, domainTask TransferTask) (TransferTask, error)
	UpdateTaskStatus(ctx context.Context, fhirTaskID string, newState string) error
	UpdateTask(ctx context.Context, fhirTaskID string, callbackFn func(domainTask TransferTask) TransferTask) error

	CreateAdvanceNotice(ctx context.Context, advanceNotice AdvanceNotice) error
	CreateNursingHandoff(ctx context.Context, nursingHandoff NursingHandoff) error

	GetAdvanceNotice(ctx context.Context, fhirCompositionPath string) (AdvanceNotice, error)
	GetNursingHandoff(ctx context.Context, fhirCompositionPath string) (NursingHandoff, error)
}

func NewFHIRTransferService(client fhir.Client) TransferService {
	return &transferService{fhirClient: client, resourceBuilder: NewFHIRBuilder()}
}

type transferService struct {
	fhirClient fhir.Client
	resourceBuilder TransferFHIRBuilder
}

func (s transferService) CreateTask(ctx context.Context, domainTask TransferTask) (TransferTask, error) {
	transferTask := s.resourceBuilder.BuildTask(fhir.TaskProperties{
		RequesterID: domainTask.SenderDID,
		OwnerID:     domainTask.ReceiverDID,
		Status:      transfer.RequestedState,
	})

	if domainTask.AdvanceNoticeID != nil {
		transferTask.Input = append(transferTask.Input, resources.TaskInputOutput{
			Type:           &fhir.LoincAdvanceNoticeType,
			ValueReference: &datatypes.Reference{Reference: fhir.ToStringPtr("/Composition/" + *domainTask.AdvanceNoticeID)},
		})
	}
	if domainTask.NursingHandoffID != nil {
		transferTask.Input = append(transferTask.Input, resources.TaskInputOutput{
			Type:           &fhir.SnomedNursingHandoffType,
			ValueReference: &datatypes.Reference{Reference: fhir.ToStringPtr("/Composition/" + *domainTask.NursingHandoffID)},
		})
	}

	err := s.fhirClient.CreateOrUpdate(ctx, transferTask)
	if err != nil {
		return domainTask, fmt.Errorf("could not create FHIR Task: %w", err)
	}
	domainTask.ID = fhir.FromIDPtr(transferTask.ID)
	return domainTask, nil
}

func (s transferService) UpdateTask(ctx context.Context, fhirTaskID string, callbackFn func(domainTask TransferTask) TransferTask) error {
	task, err := s.GetTask(ctx, fhirTaskID)
	if err != nil {
		return err
	}

	domainTask := callbackFn(*task)

	transferTask := s.resourceBuilder.BuildTask(fhir.TaskProperties{
		ID:          &domainTask.ID,
		RequesterID: domainTask.SenderDID,
		OwnerID:     domainTask.ReceiverDID,
		Status:      transfer.RequestedState,
	})

	if domainTask.AdvanceNoticeID != nil {
		transferTask.Input = append(transferTask.Input, resources.TaskInputOutput{
			Type:           &fhir.LoincAdvanceNoticeType,
			ValueReference: &datatypes.Reference{Reference: fhir.ToStringPtr("/Composition/" + *domainTask.AdvanceNoticeID)},
		})
	}
	if domainTask.NursingHandoffID != nil {
		transferTask.Input = append(transferTask.Input, resources.TaskInputOutput{
			Type:           &fhir.SnomedNursingHandoffType,
			ValueReference: &datatypes.Reference{Reference: fhir.ToStringPtr("/Composition/" + *domainTask.NursingHandoffID)},
		})
	}

	err = s.fhirClient.CreateOrUpdate(ctx, transferTask)
	if err != nil {
		return fmt.Errorf("could not update FHIR Task: %w", err)
	}
	return nil
}

func (s transferService) CreateAdvanceNotice(ctx context.Context, advanceNotice AdvanceNotice) error {
	// Save the Patient
	err := s.fhirClient.CreateOrUpdate(ctx, advanceNotice.Patient)
	if err != nil {
		return err
	}
	// Save the all the problems
	for _, problem := range advanceNotice.Problems {
		err = s.fhirClient.CreateOrUpdate(ctx, problem)
		if err != nil {
			return err
		}
	}
	// Save all the interventions
	for _, intervention := range advanceNotice.Interventions {
		err = s.fhirClient.CreateOrUpdate(ctx, intervention)
		if err != nil {
			return err
		}
	}
	// At least save the composition
	err = s.fhirClient.CreateOrUpdate(ctx, advanceNotice.Composition)
	if err != nil {
		return err
	}
	return nil
}

func (s transferService) CreateNursingHandoff(ctx context.Context, nursingHandoff NursingHandoff) error {
	panic("implement me")
}

func (s transferService) GetTask(ctx context.Context, taskID string) (*TransferTask, error) {
	fhirTask := resources.Task{}
	err := s.fhirClient.ReadOne(ctx, "/Task/"+taskID, &fhirTask)
	if err != nil {
		return nil, fmt.Errorf("error while fetching task (task-id=%s): %w", taskID, err)
	}

	task := &TransferTask{
		ID:          fhir.FromIDPtr(fhirTask.ID),
		Status:      fhir.FromCodePtr(fhirTask.Status),
		SenderDID:   fhir.FromStringPtr(fhirTask.Requester.Agent.Identifier.Value),
		ReceiverDID: fhir.FromStringPtr(fhirTask.Owner.Identifier.Value),
	}

	if input := s.findTaskInputOutputByCode(fhirTask.Input, fhir.LoincAdvanceNoticeCode); input != nil {
		ref := fhir.FromStringPtr(input.ValueReference.Reference)
		ref = strings.Split(ref, "/Composition/")[1]
		task.AdvanceNoticeID = &ref
	}
	if input := s.findTaskInputOutputByCode(fhirTask.Input, fhir.SnomedNursingHandoffCode); input != nil {
		ref := fhir.FromStringPtr(input.ValueReference.Reference)
		ref = strings.Split(ref, "/Composition/")[1]
		task.NursingHandoffID = &ref
	}

	return task, nil
}

func (s transferService) UpdateTaskStatus(ctx context.Context, fhirTaskID string, newStatus string) error {
	// TODO: check for valid state changes
	const updateErr = "could not update task state: %w"
	task, err := s.GetTask(ctx, fhirTaskID)
	if err != nil {
		return fmt.Errorf(updateErr, err)
	}
	task.Status = newStatus
	err = s.fhirClient.CreateOrUpdate(ctx, task)
	if err != nil {
		return fmt.Errorf(updateErr, err)
	}
	return nil
}

// GetAdvanceNotice converts a resolved composition into a AdvanceNotice
func (s transferService) GetAdvanceNotice(ctx context.Context, fhirCompositionPath string) (AdvanceNotice, error) {

	composition, sections, patient, err := s.ResolveComposition(ctx, fhirCompositionPath)
	if err != nil {
		return AdvanceNotice{}, fmt.Errorf("could not get nursing handoff composition: %w", err)
	}
	advanceNotice := AdvanceNotice{
		Composition: *composition,
		Patient:     *patient,
	}

	nursingDiagnosisEntries, ok := sections[CarePlanCode+"/"+NursingDiagnosisCode]
	if !ok {
		return AdvanceNotice{}, errors.New("NursingDiagnosis not set in composition")
	}

	// the nursing diagnosis contains both conditions and procedures
	for _, entry := range nursingDiagnosisEntries {
		if condition, ok := entry.(resources.Condition); ok {
			advanceNotice.Problems = append(advanceNotice.Problems, condition)
		}
		if procedure, ok := entry.(fhir.Procedure); ok {
			advanceNotice.Interventions = append(advanceNotice.Interventions, procedure)
		}
	}

	return advanceNotice, nil
}

// GetNursingHandoff converts a resolved composition into a NursingHandoff
func (s transferService) GetNursingHandoff(ctx context.Context, fhirCompositionPath string) (NursingHandoff, error) {

	composition, sections, patient, err := s.ResolveComposition(ctx, fhirCompositionPath)
	if err != nil {
		return NursingHandoff{}, fmt.Errorf("could not get nursing handoff composition: %w", err)
	}
	nursingHandoff := NursingHandoff{
		Composition: *composition,
		Patient:     *patient,
	}

	nursingDiagnosisEntries, ok := sections[CarePlanCode+"/"+NursingDiagnosisCode]
	if !ok {
		return NursingHandoff{}, errors.New("NursingDiagnosis not set in composition")
	}

	// the nursing diagnosis contains both conditions and procedures
	for _, entry := range nursingDiagnosisEntries {
		if condition, ok := entry.(resources.Condition); ok {
			nursingHandoff.Problems = append(nursingHandoff.Problems, condition)
		}
		if procedure, ok := entry.(fhir.Procedure); ok {
			nursingHandoff.Interventions = append(nursingHandoff.Interventions, procedure)
		}
	}

	return nursingHandoff, nil
}

func (s transferService) findTaskInputOutputByCode(ios []resources.TaskInputOutput, code datatypes.Code) *resources.TaskInputOutput {
	for _, io := range ios {
		if fhir.FromCodePtr(io.Type.Coding[0].Code) == string(code) {
			return &io
		}
	}
	return nil
}

func (s transferService) ResolveComposition(ctx context.Context, compositionPath string) (*fhir.Composition, map[string][]interface{}, *resources.Patient, error) {
	composition := fhir.Composition{}
	patient := resources.Patient{}
	sections := map[string][]interface{}{}

	err := s.fhirClient.ReadOne(ctx, "/"+compositionPath, &composition)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error while fetching the composition(id=%s): %w", compositionPath, err)
	}

	// Fetch the Patient
	err = s.fhirClient.ReadOne(ctx, "/"+fhir.FromStringPtr(composition.Subject.Reference), &patient)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error while fetching the transfer subject (patient): %w", err)
	}

	// Fill a map with sections by code
	for _, l1 := range composition.Section {
		l1code := fhir.FromCodePtr(l1.Code.Coding[0].Code)
		if l1code == "" {
			continue
		}
		for _, l2 := range l1.Section {
			l2code := fhir.FromCodePtr(l2.Code.Coding[0].Code)
			if l2code == "" {
				continue
			}
			var entryResources []interface{}
			for _, l2entry := range l2.Entry {
				entryPath := fhir.FromStringPtr(l2entry.Reference)
				resourceTypeStr := strings.Split(entryPath, "/")[0]
				switch resourceTypeStr {
				case "Condition":
					resource := resources.Condition{}
					err := s.fhirClient.ReadOne(ctx, entryPath, &resource)
					if err != nil {
						return nil, nil, nil, err
					}
					entryResources = append(entryResources, resource)
				case "Procedure":
					resource := fhir.Procedure{}
					err := s.fhirClient.ReadOne(ctx, entryPath, &resource)
					if err != nil {
						return nil, nil, nil, err
					}
					entryResources = append(entryResources, resource)
				}

			}
			sections[l1code+"/"+l2code] = entryResources
		}
	}

	return &composition, sections, &patient, nil
}
