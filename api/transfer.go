package api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	transfer2 "github.com/nuts-foundation/nuts-demo-ehr/domain/transfer"
)

type GetPatientTransfersParams = domain.GetPatientTransfersParams

func (w Wrapper) CreateTransfer(ctx echo.Context) error {
	request := domain.CreateTransferRequest{}
	if err := ctx.Bind(&request); err != nil {
		return err
	}
	transfer, err := w.TransferRepository.Create(ctx.Request().Context(), w.getCustomerID(), string(request.DossierID), request.Description, request.TransferDate.Time)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, transfer)
}

func (w Wrapper) GetPatientTransfers(ctx echo.Context, params GetPatientTransfersParams) error {
	transfers, err := w.TransferRepository.FindByPatientID(ctx.Request().Context(), w.getCustomerID(), params.PatientID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, transfers)
}

func (w Wrapper) GetTransfer(ctx echo.Context, transferID string) error {
	transfer, err := w.TransferRepository.FindByID(ctx.Request().Context(), w.getCustomerID(), transferID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, transfer)
}

func (w Wrapper) StartTransferNegotiation(ctx echo.Context, transferID string, organizationDID string) error {
	transfer, err := w.TransferRepository.FindByID(ctx.Request().Context(), w.getCustomerID(), transferID)
	if err != nil {
		return err
	}
	// Validate transfer
	if transfer.Status == domain.TransferStatusCancelled || transfer.Status == domain.TransferStatusCompleted || transfer.Status == domain.TransferStatusAssigned {
		return errors.New("can't start new transfer negotiation when status is 'cancelled', 'assigned' or 'completed'")
	}
	senderDID := w.getCustomerDID()
	if senderDID == nil {
		return errors.New("transferring care organization isn't registered on Nuts Network")
	}
	// Create negotiation and share it to the other party
	negotiation, err := w.TransferRepository.CreateNegotiation(ctx.Request().Context(), w.getCustomerID(), transferID, organizationDID, transfer.TransferDate.Time)
	if err != nil {
		return err
	}
	task := transfer2.EOverdrachtTask{
		SenderNutsDID:   *senderDID,
		ReceiverNutsDID: organizationDID,
		Status:          domain.TransferNegotiationStatusRequested,
	}
	err = w.FHIRGateway.CreateTask(task)
	if err != nil {
		return err
	}
	// Update transfer.Status = assigned
	err = w.TransferRepository.Update(ctx.Request().Context(), w.getCustomerID(), transfer.Description, transfer.TransferDate.Time, domain.TransferStatusRequested)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, negotiation)
}

func (w Wrapper) AssignTransferNegotiation(ctx echo.Context, transferID string, organizationDID string) error {
	transfer, err := w.TransferRepository.FindByID(ctx.Request().Context(), w.getCustomerID(), transferID)
	if err != nil {
		return err
	}
	// Validate transfer
	if transfer.Status == domain.TransferStatusRequested {
		return errors.New("can't assign transfer to care organization when status is not 'requested'")
	}
	senderDID := w.getCustomerDID()
	if senderDID == nil {
		return errors.New("transferring care organization isn't registered on Nuts Network")
	}
	// Make sure the negotation is accepted by the receiving care organization
	negotiation, err := w.findNegotiation(ctx, transferID, organizationDID, err)
	if err != nil {
		return err
	}
	if negotiation.Status != domain.TransferNegotiationStatusAccepted {
		return errors.New("can't assign transfer to care organization when it hasn't accepted the transfer")
	}
	// All is fine, update task
	task := transfer2.EOverdrachtTask{
		SenderNutsDID:   *senderDID,
		ReceiverNutsDID: organizationDID,
		Status:          domain.TransferNegotiationStatusInProgress,
	}
	err = w.FHIRGateway.CreateTask(task)
	if err != nil {
		return err
	}
	// Update transfer.Status = assigned
	err = w.TransferRepository.Update(ctx.Request().Context(), w.getCustomerID(), transfer.Description, transfer.TransferDate.Time, domain.TransferStatusRequested)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, negotiation)
}

func (w Wrapper) ListTransferNegotiations(ctx echo.Context, transferID string) error {
	negotiations, err := w.TransferRepository.ListNegotiations(ctx.Request().Context(), w.getCustomerID(), transferID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, negotiations)
}

func (w Wrapper) findNegotiation(ctx echo.Context, transferID string, organizationDID string, err error) (*domain.TransferNegotiation, error) {
	negotiations, err := w.TransferRepository.ListNegotiations(ctx.Request().Context(), transferID, transferID)
	if err != nil {
		return nil, err
	}
	for _, curr := range negotiations {
		if curr.OrganizationDID == organizationDID {
			return &curr, nil
		}
	}
	return nil, errors.New("transfer negotiation not found")
}
