package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
)

// GetTransferRequest handles requests to receive a transfer request.
func (w Wrapper) GetTransferRequest(ctx echo.Context, requestorDID string, fhirTaskID string) error {
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}

	var sessionWithUserContext *Session

	for _, session := range w.APIAuth.GetSessions() {
		if session.CustomerID != cid || !session.UserContext {
			continue
		}

		sessionWithUserContext = &session
		break
	}

	if sessionWithUserContext == nil {
		return errors.New("unable to get transfer request without elevation")
	}

	transferRequest, err := w.TransferReceiverService.GetTransferRequest(
		ctx.Request().Context(),
		cid,
		requestorDID,
		sessionWithUserContext.Presentation,
		fhirTaskID,
	)
	if err != nil {
		return fmt.Errorf("unable to get transferRequest: %w", err)
	}

	return ctx.JSON(http.StatusOK, transferRequest)
}

func (w Wrapper) GetInboxInfo(ctx echo.Context) error {
	customer := w.getCustomer(ctx)

	count, err := w.TransferReceiverRepo.GetNotCompletedCount(ctx.Request().Context(), customer.Id)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, types.InboxInfo{MessageCount: count})
}

func (w Wrapper) GetInbox(ctx echo.Context) error {
	customer := w.getCustomer(ctx)

	transfers, err := w.TransferReceiverRepo.GetAll(ctx.Request().Context(), customer.Id)
	if err != nil {
		return err
	}

	var entries []types.InboxEntry

	for _, transfer := range transfers {
		sender := transfer.Sender

		// We need to fetch the organization as we only have it's ID
		organization, err := w.OrganizationRegistry.Get(ctx.Request().Context(), transfer.Sender.Did)
		if err != nil {
			logrus.Errorf("failed to get organization: %s", err.Error())
		}

		if organization != nil {
			sender = types.FromNutsOrganization(*organization)
		}

		entries = append(entries, types.InboxEntry{
			Date:              transfer.CreatedAt.Format("02-01-2006 15:04:05"),
			RequiresAttention: true,
			ResourceID:        transfer.FhirTaskID,
			Sender:            sender,
			Title:             "Overdracht van zorg",
			Type:              "incomingTransfer",
			Status:            transfer.Status,
		})
	}

	return ctx.JSON(http.StatusOK, entries)
}
