package eoverdracht

import (
	"context"
	"errors"
	"fmt"

	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
)

type TransferService interface {
	GetTask(ctx context.Context, taskID string) (*TransferTask, error)
	UpdateTaskStatus(ctx context.Context, fhirTaskID string, newState string) error
	GetNursingHandoff(ctx context.Context, fhirCompositionPath string) (NursingHandoff, error)
	GetAdvanceNotice(ctx context.Context, fhirCompositionPath string) (AdvanceNotice, error)
}

func NewReceiverFHIRTransferService(repo fhir.Repository) TransferService {
	return &receiverTransferService{fhirRepo: repo}
}

type receiverTransferService struct {
	fhirRepo   fhir.Repository
	fhirClient fhir.Client
}

func (s receiverTransferService) GetTask(ctx context.Context, taskID string) (*TransferTask, error) {
	fhirTask, err := s.fhirRepo.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	task := &TransferTask{
		ID:     fhir.FromIDPtr(fhirTask.ID),
		Status: fhir.FromCodePtr(fhirTask.Status),
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

func (s receiverTransferService) UpdateTaskStatus(ctx context.Context, fhirTaskID string, newStatus string) error {
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
func (s receiverTransferService) GetAdvanceNotice(ctx context.Context, fhirCompositionPath string) (AdvanceNotice, error) {

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
func (s receiverTransferService) GetNursingHandoff(ctx context.Context, fhirCompositionPath string) (NursingHandoff, error) {

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

func (s receiverTransferService) findTaskInputOutputByCode(ios []resources.TaskInputOutput, code datatypes.Code) *resources.TaskInputOutput {
	for _, io := range ios {
		if fhir.FromCodePtr(io.Type.Coding[0].Code) == string(code) {
			return &io
		}
	}
	return nil
}
