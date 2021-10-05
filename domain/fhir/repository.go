package fhir

import (
	"context"
	"fmt"

	"github.com/monarko/fhirgo/STU3/resources"
)

type Repository interface {
	GetTask(ctx context.Context, fhirTaskID string) (resources.Task, error)
	UpdateTaskStatus(ctx context.Context, fhirTaskID string, newState string) error
}

type fhirRepository struct {
	client Client
}

func NewFHIRRepository(client Client) Repository {
	return &fhirRepository{client: client}
}

func (r fhirRepository) GetTask(ctx context.Context, fhirTaskID string) (resources.Task, error) {
	task := resources.Task{}
	err := r.client.ReadOne(ctx, "/Task/"+fhirTaskID, &task)
	if err != nil {
		return resources.Task{}, fmt.Errorf("error while fetching task (task-id=%s): %w", fhirTaskID, err)
	}
	return task, nil
}

func (r fhirRepository) UpdateTaskStatus(ctx context.Context, fhirTaskID string, newStatus string) error {
	const updateErr = "could not update task state: %w"
	task, err := r.GetTask(ctx, fhirTaskID)
	if err != nil {
		return fmt.Errorf(updateErr, err)
	}
	task.Status = ToCodePtr(newStatus)
	err = r.client.CreateOrUpdate(ctx, task)
	if err != nil {
		return fmt.Errorf(updateErr, err)
	}
	return nil
}
