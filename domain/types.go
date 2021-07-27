package domain

import "time"

type Task struct {
	ID                   string
	Status               string
	PatientID            string
	RequesterID          string
	OwnerID              string
	FHIRAdvanceNoticeID  *string
	FHIRNursingHandoffID *string
	AlternativeDate      *time.Time
}
