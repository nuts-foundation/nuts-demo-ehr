package api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
)

type SearchOrganizationsParams = types.SearchOrganizationsParams

func (w Wrapper) SearchOrganizations(ctx echo.Context, params SearchOrganizationsParams) error {
	customer := w.getCustomer(ctx)
	if customer == nil {
		return errors.New("customer not found")
	}

	organizations, err := w.OrganizationRegistry.Search(ctx.Request().Context(), params.Query, params.DidServiceType)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var results []types.Organization

	for _, organization := range organizations {
		// Hide our own organization
		if organization.Did == *customer.Did {
			continue
		}

		results = append(results, organization)
	}

	return ctx.JSON(http.StatusOK, results)
}
