package api

import (
	"github.com/sirupsen/logrus"
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
		sender := transfer.Sender

		// We need to fetch the organization as we only have it's ID
		organization, err := w.OrganizationRegistry.Get(ctx.Request().Context(), transfer.Sender.Did)
		if err != nil {
			logrus.Errorf("failed to get organization: %s", err.Error())
		}

		if organization != nil {
			sender = *organization
		}

		entries = append(entries, domain.InboxEntry{
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
