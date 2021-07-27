package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

type SearchOrganizationsParams = domain.SearchOrganizationsParams

func (w Wrapper) SearchOrganizations(ctx echo.Context, params SearchOrganizationsParams) error {
	organizations, err := w.OrganizationRegistry.Search(ctx.Request().Context(), params.Query, params.DidServiceType)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, organizations)
}
