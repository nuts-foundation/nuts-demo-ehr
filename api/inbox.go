package api

import (
	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"net/http"
)

type NotifyTransferUpdateParams = domain.NotifyTransferUpdateParams

func (w Wrapper) GetInboxInfo(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, domain.InboxInfo{MessageCount: 10})
}

func (w Wrapper) GetInbox(ctx echo.Context) error {
	customerID := w.getCustomerID(ctx)
	entries, err := w.Inbox.List(ctx.Request().Context(), customerID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, entries)
}
