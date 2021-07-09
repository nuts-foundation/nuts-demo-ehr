package transfer

import "github.com/nuts-foundation/nuts-demo-ehr/domain"

type EOverdrachtTask struct {
	SenderNutsDID   string
	ReceiverNutsDID string
	Status          domain.TransferNegotiationStatus
}
