package api

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
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

func (w Wrapper) UpdateTransfer(ctx echo.Context, transferID string) error {
	updateRequest := &domain.TransferProperties{}
	err := ctx.Bind(updateRequest)
	if err != nil {
		return err
	}
	transfer, err := w.TransferRepository.Update(ctx.Request().Context(), w.getCustomerID(), transferID, func(t domain.Transfer) (*domain.Transfer, error) {
		t.Description = updateRequest.Description
		t.TransferDate = updateRequest.TransferDate
		return &t, nil
	})
	return ctx.JSON(http.StatusOK, transfer)
}

func (w Wrapper) StartTransferNegotiation(ctx echo.Context, transferID string, organizationDID string) error {
	var negotiation *domain.TransferNegotiation
	_, err := w.TransferRepository.Update(ctx.Request().Context(), w.getCustomerID(), transferID, func(transfer domain.Transfer) (*domain.Transfer, error) {
		// Validate transfer
		if transfer.Status == domain.TransferStatusCancelled || transfer.Status == domain.TransferStatusCompleted || transfer.Status == domain.TransferStatusAssigned {
			return nil, errors.New("can't start new transfer negotiation when status is 'cancelled', 'assigned' or 'completed'")
		}
		senderDID := w.getCustomerDID()
		if senderDID == nil {
			return nil, errors.New("transferring care organization isn't registered on Nuts Network")
		}
		// Create negotiation and share it to the other party
		// TODO: Share transaction to this repository call as well
		var err error
		negotiation, err = w.TransferRepository.CreateNegotiation(ctx.Request().Context(), w.getCustomerID(), transferID, organizationDID, transfer.TransferDate.Time)
		if err != nil {
			return nil, err
		}
		// Send FHIR task notification
		task := transfer2.EOverdrachtTask{
			SenderNutsDID:   *senderDID,
			ReceiverNutsDID: organizationDID,
			Status:          domain.TransferNegotiationStatusRequested,
		}
		err = w.FHIRGateway.CreateTask(task)
		if err != nil {
			return nil, err
		}

		// Update transfer.Status = requested
		transfer.Status = domain.TransferStatusRequested
		return &transfer, nil
	})
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, *negotiation)
}

func (w Wrapper) AssignTransferNegotiation(ctx echo.Context, transferID string, organizationDID string) error {
	var negotiation *domain.TransferNegotiation
	_, err := w.TransferRepository.Update(ctx.Request().Context(), w.getCustomerID(), transferID, func(transfer domain.Transfer) (*domain.Transfer, error) {
		// Validate transfer
		if transfer.Status == domain.TransferStatusRequested {
			return nil, errors.New("can't assign transfer to care organization when status is not 'requested'")
		}
		senderDID := w.getCustomerDID()
		if senderDID == nil {
			return nil, errors.New("transferring care organization isn't registered on Nuts Network")
		}
		// Make sure the negotiation is accepted by the receiving care organization
		var err error
		negotiation, err = w.findNegotiation(ctx.Request().Context(), transferID, organizationDID, err)
		if err != nil {
			return nil, err
		}
		if negotiation.Status != domain.TransferNegotiationStatusAccepted {
			return nil, errors.New("can't assign transfer to care organization when it hasn't accepted the transfer")
		}
		// All is fine, update task
		task := transfer2.EOverdrachtTask{
			SenderNutsDID:   *senderDID,
			ReceiverNutsDID: organizationDID,
			Status:          domain.TransferNegotiationStatusInProgress,
		}
		err = w.FHIRGateway.CreateTask(task)
		if err != nil {
			return nil, err
		}
		// Update transfer.Status = assigned
		transfer.Status = domain.TransferStatusAssigned
		return &transfer, nil
	})
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
	// Enrich with organization info
	for i, negotiation := range negotiations {
		organization, err := w.OrganizationRegistry.Get(ctx.Request().Context(), negotiation.OrganizationDID)
		if err != nil {
			logrus.Warnf("Error while fetching organization info for negotiation (DID=%s): %v", negotiation.OrganizationDID, err)
			continue
		}
		negotiations[i].Organization = organization
	}
	return ctx.JSON(http.StatusOK, negotiations)
}

func (w Wrapper) findNegotiation(ctx context.Context, transferID string, organizationDID string, err error) (*domain.TransferNegotiation, error) {
	negotiations, err := w.TransferRepository.ListNegotiations(ctx, transferID, transferID)
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
