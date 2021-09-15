package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/notification"
	httpAuth "github.com/nuts-foundation/nuts-demo-ehr/http/auth"
	nutsAuthClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client/auth"
	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

type GetPatientTransfersParams = domain.GetPatientTransfersParams

func (w Wrapper) CreateTransfer(ctx echo.Context) error {
	request := domain.CreateTransferRequest{}
	if err := ctx.Bind(&request); err != nil {
		return err
	}
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
	transfer, err := w.TransferSenderService.Create(ctx.Request().Context(), cid, request)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, transfer)
}

func (w Wrapper) GetPatientTransfers(ctx echo.Context, params GetPatientTransfersParams) error {
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
	transfers, err := w.TransferSenderRepo.FindByPatientID(ctx.Request().Context(), cid, params.PatientID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, transfers)
}

func (w Wrapper) GetTransfer(ctx echo.Context, transferID string) error {
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
	transfer, err := w.TransferSenderRepo.FindByID(ctx.Request().Context(), cid, transferID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, transfer)
}

func (w Wrapper) GetTransferRequest(ctx echo.Context, requestorDID string, fhirTaskID string) error {
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
	transferRequest, err := w.TransferSenderService.GetTransferRequest(ctx.Request().Context(), cid, requestorDID, fhirTaskID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, transferRequest)
}


func (w Wrapper) ChangeTransferRequestState(ctx echo.Context, requestorDID string, fhirTaskID string) error {
	updateRequest := &domain.TransferNegotiationStatus{}
	err := ctx.Bind(updateRequest)
	if err != nil {
		return err
	}
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}

	err = w.TransferReceiverService.UpdateTransferRequestState(ctx.Request().Context(), cid, requestorDID, fhirTaskID, string(updateRequest.Status))
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (w Wrapper) UpdateTransfer(ctx echo.Context, transferID string) error {
	updateRequest := &domain.TransferProperties{}
	err := ctx.Bind(updateRequest)
	if err != nil {
		return err
	}
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
	transfer, err := w.TransferSenderRepo.Update(ctx.Request().Context(), cid, transferID, func(t *domain.Transfer) (*domain.Transfer, error) {
		//t.Description = updateRequest.Description
		t.TransferDate = updateRequest.TransferDate
		return t, nil
	})
	return ctx.JSON(http.StatusOK, transfer)
}

func (w Wrapper) CancelTransfer(ctx echo.Context, transferID string) error {
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
	transfer, err := w.TransferSenderRepo.Cancel(ctx.Request().Context(), cid, transferID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, transfer)
}

func (w Wrapper) StartTransferNegotiation(ctx echo.Context, transferID string) error {
	request := domain.CreateTransferNegotiationRequest{}
	if err := ctx.Bind(&request); err != nil {
		return err
	}
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
	negotiation, err := w.TransferSenderService.CreateNegotiation(ctx.Request().Context(), cid, transferID, request.OrganizationDID, request.TransferDate.Time)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, *negotiation)
}

func (w Wrapper) AssignTransferDirect(ctx echo.Context, transferID string) error {
	request := domain.CreateTransferNegotiationRequest{}
	if err := ctx.Bind(&request); err != nil {
		return err
	}
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
	negotiation, err := w.TransferSenderService.CreateNegotiation(ctx.Request().Context(), cid, transferID, request.OrganizationDID, request.TransferDate.Time)
	if err != nil {
		return err
	}
	negotiation, err = w.TransferSenderService.ConfirmNegotiation(ctx.Request().Context(), cid, transferID, string(negotiation.Id))
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (w Wrapper) ListTransferNegotiations(ctx echo.Context, transferID string) error {
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
	negotiations, err := w.TransferSenderRepo.ListNegotiations(ctx.Request().Context(), cid, transferID)
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
		negotiations[i].Organization = *organization
	}
	return ctx.JSON(http.StatusOK, negotiations)
}

func (w Wrapper) UpdateTransferNegotiationStatus(ctx echo.Context, transferID string, negotiationID string) error {
	request := domain.TransferNegotiationStatus{}
	if err := ctx.Bind(&request); err != nil {
		return err
	}
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
	negotiation, err := w.TransferSenderRepo.UpdateNegotiationState(ctx.Request().Context(), cid, negotiationID, request.Status)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, negotiation)
}

func (w Wrapper) NotifyTransferUpdate(ctx echo.Context) error {
	// This gets called by a transfer sending XIS to inform the local node there's FHIR tasks to be retrieved.
	rawToken := ctx.Get(httpAuth.AccessToken)
	if rawToken == nil {
		// should have been caught by security filter
		return errors.New("missing access-token")
	}
	token, ok := rawToken.(nutsAuthClient.TokenIntrospectionResponse)
	if !ok {
		// should have been caught by security filter
		return errors.New("missing access-token")
	}

	senderDID := token.Sub
	if senderDID == nil {
		return errors.New("missing 'sub' in access-token")
	}
	customerDID := token.Iss
	if customerDID == nil {
		return errors.New("missing 'Iss' in access-token")
	}

	customer, err := w.CustomerRepository.FindByDID(*customerDID)
	if err != nil {
		return err
	}

	if customer == nil {
		logrus.Warnf("Received transfer notification for unknown customer DID: %s", *senderDID)
		return echo.NewHTTPError(http.StatusNotFound, "taskOwner unknown on this server")
	}

	if err := w.NotificationHandler.Handle(ctx.Request().Context(), notification.Notification{
		SenderDID:   *senderDID,
		CustomerDID: *customerDID,
		CustomerID:  customer.Id,
	}); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (w Wrapper) findNegotiation(ctx context.Context, customerID int, transferID, negotiationID string) (*domain.TransferNegotiation, error) {
	negotiations, err := w.TransferSenderRepo.ListNegotiations(ctx, customerID, transferID)
	if err != nil {
		return nil, err
	}
	for _, curr := range negotiations {
		if string(curr.Id) == negotiationID {
			return &curr, nil
		}
	}
	return nil, errors.New("transfer negotiation not found")
}
