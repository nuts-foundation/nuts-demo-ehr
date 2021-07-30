package domain

import "time"

type TaskProperties struct {
	Status      string
	PatientID   string
	// nuts DID of the placer
	RequesterID string
	// nuts DID of the filler
	OwnerID     string
}

type Task struct {
	ID string
	TaskProperties
	FHIRAdvanceNoticeID  *string
	FHIRNursingHandoffID *string
	AlternativeDate      *time.Time
}
