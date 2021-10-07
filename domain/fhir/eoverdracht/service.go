package eoverdracht

import (
	"context"
	"errors"
	"fmt"

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

	CreateAdvanceNotice(ctx context.Context, nursingHandoff NursingHandoff) (AdvanceNotice, error)
	CreateNursingHandoff(ctx context.Context, nursingHandoff NursingHandoff) (NursingHandoff, error)

	GetAdvanceNotice(ctx context.Context, fhirCompositionPath string) (AdvanceNotice, error)
	GetNursingHandoff(ctx context.Context, fhirCompositionPath string) (NursingHandoff, error)
}

func NewFHIRTransferService(repo fhir.Repository) TransferService {
	return &transferService{fhirRepo: repo}
}

type transferService struct {
	fhirRepo   fhir.Repository
	fhirClient fhir.Client
}

func (s transferService) CreateTask(ctx context.Context, domainTask TransferTask) (TransferTask, error) {
	transferTask := NewFHIRBuilder().BuildTask(fhir.TaskProperties{
		RequesterID: domainTask.SenderDID,
		OwnerID:     domainTask.ReceiverDID,
		Status:      transfer.RequestedState,
		Input: []resources.TaskInputOutput{
			{
				Type:           &fhir.LoincAdvanceNoticeType,
				ValueReference: &datatypes.Reference{Reference: fhir.ToStringPtr("/Composition/" + *domainTask.AdvanceNoticeID)},
			},
		},
	})

	err := s.fhirRepo.CreateOrUpdateResource(ctx, transferTask)
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

	transferTask := NewFHIRBuilder().BuildTask(fhir.TaskProperties{
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

	err = s.fhirRepo.CreateOrUpdateResource(ctx, transferTask)
	if err != nil {
		return fmt.Errorf("could not update FHIR Task: %w", err)
	}
	return nil
}

func (s transferService) CreateAdvanceNotice(ctx context.Context, nursingHandoff NursingHandoff) (AdvanceNotice, error) {
	panic("implement me")
}

func (s transferService) CreateNursingHandoff(ctx context.Context, nursingHandoff NursingHandoff) (NursingHandoff, error) {
	panic("implement me")
}

func (s transferService) GetTask(ctx context.Context, taskID string) (*TransferTask, error) {
	fhirTask, err := s.fhirRepo.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	task := &TransferTask{
		ID:          fhir.FromIDPtr(fhirTask.ID),
		Status:      fhir.FromCodePtr(fhirTask.Status),
		SenderDID:   fhir.FromStringPtr(fhirTask.Requester.Agent.Identifier.Value),
		ReceiverDID: fhir.FromStringPtr(fhirTask.Owner.Identifier.Value),
	}

	if input := s.findTaskInputOutputByCode(fhirTask.Input, fhir.LoincAdvanceNoticeCode); input != nil {
		ref := fhir.FromStringPtr(input.ValueReference.Reference)
		task.AdvanceNoticeID = &ref
	}
	if input := s.findTaskInputOutputByCode(fhirTask.Input, fhir.SnomedNursingHandoffCode); input != nil {
		ref := fhir.FromStringPtr(input.ValueReference.Reference)
		task.NursingHandoffID = &ref
	}

	return task, nil
}

func (s transferService) UpdateTaskStatus(ctx context.Context, fhirTaskID string, newStatus string) error {
	// TODO: check for valid state changes
	const updateErr = "could not update task state: %w"
	task, err := s.fhirRepo.GetTask(ctx, fhirTaskID)
	if err != nil {
		return fmt.Errorf(updateErr, err)
	}
	task.Status = fhir.ToCodePtr(newStatus)
	err = s.fhirRepo.UpdateTask(ctx, task)
	if err != nil {
		return fmt.Errorf(updateErr, err)
	}
	return nil
}

// GetAdvanceNotice converts a resolved composition into a AdvanceNotice
func (s transferService) GetAdvanceNotice(ctx context.Context, fhirCompositionPath string) (AdvanceNotice, error) {

	composition, sections, patient, err := s.fhirRepo.ResolveComposition(ctx, fhirCompositionPath)
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

	composition, sections, patient, err := s.fhirRepo.ResolveComposition(ctx, fhirCompositionPath)
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
