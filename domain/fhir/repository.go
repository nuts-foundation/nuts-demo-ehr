package fhir

import (
	"context"

	"github.com/google/uuid"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

type TaskState string

// The following states:

// CreatedState : Task created, not yet announced to filler
const CreatedState = TaskState("created")

// RequestedState : placer has made the registration available
const RequestedState = TaskState("requested")

// ReceivedState : filler has received the request and is judging the request
const ReceivedState = TaskState("received")

// AcceptedState : filler accepts the registration, and can provide the care asked for
const AcceptedState = TaskState("accepted")

// OnHoldState : filler proposes a different date
const OnHoldState = TaskState("on-hold")

// InProgressState : placer confirms to the filler it proceeds with the transfer
const InProgressState = TaskState("in-progress")

// RejectedState : filler rejects the registration
const RejectedState = TaskState("rejected")

// CancelledState : placer or filler has cancelled the transfer
const CancelledState = TaskState("cancelled")

// CompletedState : filler received the nursing handoff
const CompletedState = TaskState("completed")

type Repository interface {
	CreateTask(ctx context.Context, taskProperties domain.TaskProperties) (*domain.Task, error)
}

type TaskFactory struct{}

func (TaskFactory) New(taskProperties domain.TaskProperties) *domain.Task {
	return &domain.Task{
		ID:             uuid.New().String(),
		TaskProperties: taskProperties,
	}
}
