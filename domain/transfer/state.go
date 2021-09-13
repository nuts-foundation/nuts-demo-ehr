package transfer

// All possible states as described by the Nictiz eOverdracht v4.0:
// https://informatiestandaarden.nictiz.nl/wiki/vpk:V4.0_FHIR_eOverdracht#Using_Task_to_manage_the_workflow

const RequestedState = "requested"
const AcceptedState = "accepted"
const RejectedState = "rejected"
const OnHoldState = "on-hold"
const CancelledState = "cancelled"
const InProgressState = "in-progress"
const CompletedState = "completed"
