package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/nuts-foundation/nuts-demo-ehr/http/proxy"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/notification"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	httpAuth "github.com/nuts-foundation/nuts-demo-ehr/http/auth"
	nutsAuthClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client/auth"
	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
)

type GetPatientTransfersParams = types.GetPatientTransfersParams

func (w Wrapper) CreateTransfer(ctx echo.Context) error {
	request := types.CreateTransferRequest{}
	if err := ctx.Bind(&request); err != nil {
		return err
	}
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
	transfer, err := w.TransferSenderService.CreateTransfer(ctx.Request().Context(), cid, request)
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

	transfer, err := w.TransferSenderService.GetTransferByID(ctx.Request().Context(), cid, transferID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, transfer)
}

func (w Wrapper) ChangeTransferRequestState(ctx echo.Context, requesterDID string, fhirTaskID string) error {
	updateRequest := &types.TransferNegotiationStatus{}
	err := ctx.Bind(updateRequest)
	if err != nil {
		return err
	}
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}

	err = w.TransferReceiverService.UpdateTransferRequestState(ctx.Request().Context(), cid, requesterDID, fhirTaskID, string(updateRequest.Status))
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (w Wrapper) UpdateTransfer(ctx echo.Context, transferID string) error {
	updateRequest := &types.TransferProperties{}
	err := ctx.Bind(updateRequest)
	if err != nil {
		return err
	}
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}

	_, err = w.TransferSenderRepo.Update(ctx.Request().Context(), cid, transferID, func(t *types.Transfer) (*types.Transfer, error) {
		//t.Description = updateRequest.Description
		t.TransferDate = updateRequest.TransferDate
		return t, nil
	})
	if err != nil {
		return err
	}

	transfer, err := w.TransferSenderService.GetTransferByID(ctx.Request().Context(), cid, transferID)
	if err != nil {
		return err
	}

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
	request := types.CreateTransferNegotiationRequest{}
	if err := ctx.Bind(&request); err != nil {
		return err
	}
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
	negotiation, err := w.TransferSenderService.CreateNegotiation(ctx.Request().Context(), cid, transferID, request.OrganizationDID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, *negotiation)
}

func (w Wrapper) AssignTransferDirect(ctx echo.Context, transferID string) error {
	request := types.CreateTransferNegotiationRequest{}
	if err := ctx.Bind(&request); err != nil {
		return err
	}
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
	_, err = w.TransferSenderService.AssignTransfer(ctx.Request().Context(), cid, transferID, request.OrganizationDID)
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
		negotiations[i].Organization = types.FromNutsOrganization(*organization)
	}
	return ctx.JSON(http.StatusOK, negotiations)
}

func (w Wrapper) UpdateTransferNegotiationStatus(ctx echo.Context, transferID string, negotiationID string) error {
	request := types.TransferNegotiationStatus{}
	if err := ctx.Bind(&request); err != nil {
		return err
	}
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
	newState := request.Status
	if newState == transfer.InProgressState {
		_, err = w.TransferSenderService.ConfirmNegotiation(ctx.Request().Context(), cid, transferID, negotiationID)
	} else if newState == transfer.CancelledState {
		_, err = w.TransferSenderService.CancelNegotiation(ctx.Request().Context(), cid, transferID, negotiationID)
	}
	if err != nil {
		return fmt.Errorf("unable to update transfer negotiation state: %w", err)
	}
	negotiation, err := w.TransferSenderRepo.UpdateNegotiationState(ctx.Request().Context(), cid, negotiationID, request.Status)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, negotiation)
}

func (w Wrapper) NotifyTransferUpdate(ctx echo.Context, taskID string) error {
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
		return ctx.JSON(http.StatusInternalServerError, &proxy.OperationOutcome{
			Text: "an error occurred",
			Issue: &proxy.Issue{
				Code:     "error",
				Severity: "error",
				Details: &proxy.IssueDetails{
					Text: err.Error(),
				},
			},
		})
	}

	if customer == nil {
		logrus.Warnf("Received transfer notification for unknown customer DID: %s", *senderDID)

		return ctx.JSON(http.StatusNotFound, &proxy.OperationOutcome{
			Text: "taskOwner unknown on this server",
			Issue: &proxy.Issue{
				Code:     "invalid",
				Severity: "error",
				Details: &proxy.IssueDetails{
					Text: fmt.Sprintf("received transfer notification for unknown taskOwner with DID: %s", *senderDID),
				},
			},
		})
	}

	if err := w.NotificationHandler.Handle(ctx.Request().Context(), notification.Notification{
		TaskID:      taskID,
		SenderDID:   *senderDID,
		CustomerDID: *customerDID,
		CustomerID:  customer.Id,
	}); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusAccepted)
}
