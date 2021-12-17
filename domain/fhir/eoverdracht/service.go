package eoverdracht

import (
	"context"
	"fmt"
	"reflect"
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

	GetAdvanceNotice(ctx context.Context, fhirCompositionID string) (AdvanceNotice, error)
	GetNursingHandoff(ctx context.Context, fhirCompositionID string) (NursingHandoff, error)
}

func NewFHIRTransferService(client fhir.Client) TransferService {
	return &transferService{fhirClient: client, resourceBuilder: NewFHIRBuilder()}
}

type transferService struct {
	fhirClient      fhir.Client
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
			Type:           &LoincAdvanceNoticeType,
			ValueReference: &datatypes.Reference{Reference: fhir.ToStringPtr("/Composition/" + *domainTask.AdvanceNoticeID)},
		})
	}
	if domainTask.NursingHandoffID != nil {
		transferTask.Input = append(transferTask.Input, resources.TaskInputOutput{
			Type:           &SnomedNursingHandoffType,
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
			Type:           &LoincAdvanceNoticeType,
			ValueReference: &datatypes.Reference{Reference: fhir.ToStringPtr("/Composition/" + *domainTask.AdvanceNoticeID)},
		})
	}
	if domainTask.NursingHandoffID != nil {
		transferTask.Input = append(transferTask.Input, resources.TaskInputOutput{
			Type:           &SnomedNursingHandoffType,
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
	err := s.fhirClient.ReadOne(ctx, "Task/"+taskID, &fhirTask)
	if err != nil {
		return nil, fmt.Errorf("error while fetching task (task-id=%s): %w", taskID, err)
	}

	task := &TransferTask{
		ID:          fhir.FromIDPtr(fhirTask.ID),
		Status:      fhir.FromCodePtr(fhirTask.Status),
		SenderDID:   fhir.FromStringPtr(fhirTask.Requester.Agent.Identifier.Value),
		ReceiverDID: fhir.FromStringPtr(fhirTask.Owner.Identifier.Value),
	}

	if input := s.findTaskInputOutputByCode(fhirTask.Input, LoincAdvanceNoticeCode); input != nil {
		ref := fhir.FromStringPtr(input.ValueReference.Reference)
		ref = strings.Split(ref, "Composition/")[1]
		task.AdvanceNoticeID = &ref
	}
	if input := s.findTaskInputOutputByCode(fhirTask.Input, SnomedNursingHandoffCode); input != nil {
		ref := fhir.FromStringPtr(input.ValueReference.Reference)
		ref = strings.Split(ref, "Composition/")[1]
		task.NursingHandoffID = &ref
	}

	return task, nil
}

func (s transferService) UpdateTaskStatus(ctx context.Context, fhirTaskID string, newStatus string) error {
	// TODO: check for valid state changes
	const updateErr = "could not update task state: %w"

	task := &resources.Task{}

	if err := s.fhirClient.ReadOne(ctx, "Task/"+fhirTaskID, &task); err != nil {
		return err
	}

	task.Status = fhir.ToCodePtr(newStatus)

	if err := s.fhirClient.CreateOrUpdate(ctx, task); err != nil {
		return fmt.Errorf(updateErr, err)
	}

	return nil
}

// GetAdvanceNotice converts a resolved composition into a AdvanceNotice
func (s transferService) GetAdvanceNotice(ctx context.Context, fhirCompositionID string) (AdvanceNotice, error) {
	composition, patient, err := s.ResolveComposition(ctx, "Composition/"+fhirCompositionID)
	if err != nil {
		return AdvanceNotice{}, fmt.Errorf("could not get nursing handoff composition: %w", err)
	}
	advanceNotice := AdvanceNotice{
		Composition: *composition,
		Patient:     *patient,
	}

	carePlanCompositions, err := s.resolveCompositionSections(composition.Section, CarePlanConcept)
	if err != nil {
		return advanceNotice, err
	}

	for _, handoffComposition := range carePlanCompositions {
		conditions, err := s.resolveCompositionEntry(ctx, handoffComposition, resources.Condition{})
		if err != nil {
			return AdvanceNotice{}, err
		}
		for _, condition := range conditions {
			advanceNotice.Problems = append(advanceNotice.Problems, *condition.(*resources.Condition))
		}

		procedures, err := s.resolveCompositionEntry(ctx, handoffComposition, fhir.Procedure{})
		if err != nil {
			return AdvanceNotice{}, err
		}
		for _, procedure := range procedures {
			advanceNotice.Interventions = append(advanceNotice.Interventions, *procedure.(*fhir.Procedure))
		}
	}

	return advanceNotice, nil
}

// GetNursingHandoff converts a resolved composition into a NursingHandoff
func (s transferService) GetNursingHandoff(ctx context.Context, fhirCompositionID string) (NursingHandoff, error) {
	composition, patient, err := s.ResolveComposition(ctx, "Composition/"+fhirCompositionID)
	if err != nil {
		return NursingHandoff{}, fmt.Errorf("could not get nursing handoff composition: %w", err)
	}
	nursingHandoff := NursingHandoff{
		Composition: *composition,
		Patient:     *patient,
	}

	carePlanCompositions, err := s.resolveCompositionSections(composition.Section, CarePlanConcept)
	for _, handoffComposition := range carePlanCompositions {
		conditions, err := s.resolveCompositionEntry(ctx, handoffComposition, resources.Condition{})
		if err != nil {
			return NursingHandoff{}, err
		}
		for _, condition := range conditions {
			nursingHandoff.Problems = append(nursingHandoff.Problems, *condition.(*resources.Condition))
		}

		procedures, err := s.resolveCompositionEntry(ctx, handoffComposition, fhir.Procedure{})
		if err != nil {
			return NursingHandoff{}, err
		}
		for _, procedure := range procedures {
			nursingHandoff.Interventions = append(nursingHandoff.Interventions, *procedure.(*fhir.Procedure))
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

func (s transferService) resolveCompositionSections(sections []fhir.CompositionSection, code datatypes.CodeableConcept) ([]fhir.CompositionSection, error) {
	for _, section := range sections {
		if fhir.FromCodePtr(section.Code.Coding[0].Code) == fhir.FromCodePtr(code.Coding[0].Code) {
			return section.Section, nil
		}
	}
	return nil, fmt.Errorf("sections not found")
}

func (s transferService) resolveCompositionEntry(ctx context.Context, section fhir.CompositionSection, resource interface{}) ([]interface{}, error) {
	resourceType := reflect.TypeOf(resource)
	res := []interface{}{}

	typeStr := resourceType.Name()
	for _, entry := range section.Entry {
		entryPath := fhir.FromStringPtr(entry.Reference)
		resourceTypeStr := strings.Split(entryPath, "/")[0]
		if resourceTypeStr != typeStr {
			continue
		}
		newResourceType := reflect.New(resourceType)
		newResource := newResourceType.Interface()
		err := s.fhirClient.ReadOne(ctx, entryPath, &newResource)
		if err != nil {
			return nil, err
		}
		res = append(res, newResource)
	}
	return res, nil
}

func (s transferService) ResolveComposition(ctx context.Context, compositionPath string) (*fhir.Composition, *resources.Patient, error) {
	composition := fhir.Composition{}
	patient := resources.Patient{}

	err := s.fhirClient.ReadOne(ctx, "/"+compositionPath, &composition)
	if err != nil {
		return nil, nil, fmt.Errorf("error while fetching the composition(id=%s): %w", compositionPath, err)
	}

	// Fetch the Patient
	err = s.fhirClient.ReadOne(ctx, "/"+fhir.FromStringPtr(composition.Subject.Reference), &patient)
	if err != nil {
		return nil, nil, fmt.Errorf("error while fetching the transfer subject (patient): %w", err)
	}

	return &composition, &patient, nil
}
