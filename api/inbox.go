package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

func (w Wrapper) GetInboxInfo(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, domain.InboxInfo{MessageCount: 10})
}

func (w Wrapper) GetInbox(ctx echo.Context) error {
	customer := w.getCustomer(ctx)
	entries, err := w.Inbox.List(ctx.Request().Context(), customer)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, entries)
}
