package task

import (
	"context"
)

type TaskState string

// The following states:

// RequestedState : placer has made the registration available
const RequestedState = TaskState("requested")

// ReceivedState : filler has received the request and is juding the request
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

const SnomedCodingSystem = "http://snomed.info/sct"
const SnomedTransferCode = "308292007"


type Repository interface {
	Create(ctx context.Context, patientID, customerID string) ()
}