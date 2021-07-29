package domain

import "time"

type TaskProperties struct {
	Status      string
	PatientID   string
	RequesterID string
	OwnerID     string
}

type Task struct {
	ID string
	TaskProperties
	FHIRAdvanceNoticeID  *string
	FHIRNursingHandoffID *string
	AlternativeDate      *time.Time
}
