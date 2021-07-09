package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

type SearchOrganizationsParams = domain.SearchOrganizationsParams

func (w Wrapper) SearchOrganizations(ctx echo.Context, params SearchOrganizationsParams) error {
	results := []domain.Organization{
		{
			City: "Hengelo",
			Did:  "did:nuts:123",
			Name: "Hengelzorg BV",
		},
	}
	return ctx.JSON(http.StatusOK, results)

}
