package api

import (
	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"net/http"
)

func (w Wrapper) GetInboxInfo(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, domain.InboxInfo{MessageCount: 10})
}

func (w Wrapper) GetInbox(ctx echo.Context) error {
	entries, err := w.InboxRepository.List(ctx.Request().Context())
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, entries)
}
