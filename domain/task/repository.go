package task

import (
	"context"

	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

type TaskState string

// The following states:

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

// Coding systems:
const SnomedCodingSystem = "http://snomed.info/sct"
const LoincCodingSystem = "http://loinc.org"

// Codes:
const SnomedTransferCode = "308292007"
const TransferDisplay = "Overdracht van zorg"
const LoincAdvanceNoticeCode = "57830-2"
const SnomedAlternaticeDateCode = "146851000146105"
const SnomedNursingHandoffCode = "371535009"



type Repository interface {
	Create(ctx context.Context, task domain.Task) (*domain.Task, error)
}

type Factory struct {}