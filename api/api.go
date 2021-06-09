package api

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/client"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
)

type Wrapper struct {
	Auth       *Auth
	Client     client.HTTPClient
	Repository customers.Repository
}

func (w Wrapper) CheckSession(ctx echo.Context) error {
	// If this function is reached, it means the session is still valid
	return ctx.NoContent(http.StatusNoContent)
}

func (w Wrapper) CreateSession(ctx echo.Context) error {
	customer := domain.Customer{}
	if err := ctx.Bind(&customer); err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	// forward to node
	bytes, err := w.Client.CreateIrmaSession(customer)
	if err != nil {
		return err
	}

	// convert to map so echo rendering doesn't escape double quotes
	j := map[string]interface{}{}
	json.Unmarshal(bytes, &j)
	return ctx.JSON(200, j)
}

func (w Wrapper) SessionResult(ctx echo.Context, sessionToken string) error {
	// forward to node
	bytes, err := w.Client.GetIrmaSessionResult(sessionToken)
	if err != nil {
		return err
	}

	base64String := base64.StdEncoding.EncodeToString(bytes)
	token := w.Auth.StoreVP(base64String)
	writeSession(ctx, token)
	return ctx.JSON(200, domain.SessionToken{
		Token: token,
	})
}

func (w Wrapper) ListCustomers(ctx echo.Context) error {
	customers, err := w.Repository.All()
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return ctx.JSON(http.StatusOK, customers)
}

// writeSession expect the VP as base64encoded json
func writeSession(ctx echo.Context, vp string) {
	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = vp
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(24 * time.Hour)
	ctx.SetCookie(cookie)
}
