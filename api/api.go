package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
)

type Wrapper struct {
	Auth       auth
	Repository customers.Repository
}

func (w Wrapper) CheckSession(ctx echo.Context) error {
	// If this function is reached, it means the session is still valid
	return ctx.NoContent(http.StatusNoContent)
}

func (w Wrapper) CreateSession(ctx echo.Context) error {
	sessionRequest := domain.CreateSessionRequest{}
	if err := ctx.Bind(&sessionRequest); err != nil {
		return err
	}

	if !w.Auth.CheckCredentials(sessionRequest.Username, sessionRequest.Password) {
		return echo.NewHTTPError(http.StatusForbidden, "invalid credentials")
	}

	token, err := w.Auth.CreateJWT(sessionRequest.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(200, domain.CreateSessionResponse{Token: string(token)})
}

func (w Wrapper) ListCustomers(ctx echo.Context) error {
	customers, err := w.Repository.All()
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return ctx.JSON(http.StatusOK, customers)
}
