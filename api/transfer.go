package api

import (
	"github.com/labstack/echo/v4"
)

func (w Wrapper) CreateTransfer(ctx echo.Context, patientID string) error {
	panic("implement me")
}

func (w Wrapper) GetTransfer(ctx echo.Context, patientID string, transferID string) error {
	panic("implement me")
}

func (w Wrapper) ListTransferNegotiations(ctx echo.Context, patientID string, transferID string) error {
	panic("implement me")
}

func (w Wrapper) StartTransferNegotiation(ctx echo.Context, patientID string, transferID string) error {
	panic("implement me")
}

func (w Wrapper) AcceptTransferNegotiation(ctx echo.Context, patientID string, transferID string, negotiationID string) error {
	panic("implement me")
}
