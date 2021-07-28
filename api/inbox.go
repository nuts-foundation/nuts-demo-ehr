package api

import (
	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"net/http"
)

func (w Wrapper) GetInboxInfo(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, domain.InboxInfo{MessageCount: 10})
}
