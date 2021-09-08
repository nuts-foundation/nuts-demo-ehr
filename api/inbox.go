package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

func (w Wrapper) GetInboxInfo(ctx echo.Context) error {
	customer := w.getCustomer(ctx)

	count, err := w.TransferReceiverRepo.GetCount(ctx.Request().Context(), customer.Id)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, domain.InboxInfo{MessageCount: count})
}

func (w Wrapper) GetInbox(ctx echo.Context) error {
	customer := w.getCustomer(ctx)

	transfers, err := w.TransferReceiverRepo.GetAll(ctx.Request().Context(), customer.Id)
	if err != nil {
		return err
	}

	var entries []domain.InboxEntry

	for _, transfer := range transfers {
		entries = append(entries, domain.InboxEntry{
			Date:              "TODO",
			RequiresAttention: true,
			ResourceID:        transfer.FhirTaskID,
			Sender:            transfer.Sender,
			Title:             "Overdracht van zorg",
			Type:              "incomingTransfer",
		})
	}

	return ctx.JSON(http.StatusOK, entries)
}
