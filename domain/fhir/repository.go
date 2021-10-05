package fhir

import (
	"context"
	"fmt"

	"github.com/monarko/fhirgo/STU3/resources"
)

type Repository interface {
	GetTask(ctx context.Context, fhirTaskID string) (resources.Task, error)
}

type fhirRepository struct {
	client Factory
}

func NewFHIRRepository(client Factory) Repository {
	return &fhirRepository{client: client}
}

func (r fhirRepository) GetTask(ctx context.Context, fhirTaskID string) (resources.Task, error) {
	// TODO: Read AdvanceNotification here instead of the transfer task
	task := resources.Task{}
	err := r.client().ReadOne(ctx, "/Task/"+fhirTaskID, &task)
	if err != nil {
		return resources.Task{}, fmt.Errorf("error while looking up transfer task remotely(task-id=%s): %w", fhirTaskID, err)
	}
	return task, nil
}
