package eoverdracht

import (
	"context"

	"github.com/monarko/fhirgo/STU3/datatypes"
	"github.com/monarko/fhirgo/STU3/resources"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
)

type TransferService interface {
	GetTask(ctx context.Context, taskID string) (*TransferTask, error)
}

func NewReceiverFHIRTransferService(repo fhir.Repository) TransferService {
	return &receiverTransferService{fhirRepo: repo}
}

type receiverTransferService struct {
	fhirRepo fhir.Repository
}

func (s receiverTransferService) GetTask(ctx context.Context, taskID string) (*TransferTask, error) {
	fhirTask, err := s.fhirRepo.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	task := &TransferTask{
		ID:               fhir.FromIDPtr(fhirTask.ID),
		Status:           fhir.FromCodePtr(fhirTask.Status),
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

func (s receiverTransferService) findTaskInputOutputByCode(ios []resources.TaskInputOutput, code datatypes.Code) *resources.TaskInputOutput {
	for _, io := range ios {
		if fhir.FromCodePtr(io.Type.Coding[0].Code) == string(code) {
			return &io
		}
	}
	return nil
}