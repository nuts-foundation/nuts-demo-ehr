package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
)

type SearchOrganizationsParams = types.SearchOrganizationsParams

func (w Wrapper) SearchOrganizations(ctx echo.Context, params SearchOrganizationsParams) error {
	customer, err := w.getCustomer(ctx)
	if err != nil {
		return err
	}

	organizations, err := w.OrganizationRegistry.Search(ctx.Request().Context(), params.Query, params.DidServiceType)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var results = make([]types.Organization, 0)

	for _, organization := range organizations {
		// Hide our own organization
		if organization.ID == *customer.Did {
			continue
		}

		results = append(results, types.FromNutsOrganization(organization))
	}

	return ctx.JSON(http.StatusOK, results)
}
