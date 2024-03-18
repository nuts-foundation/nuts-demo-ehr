package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
)

func (w Wrapper) SearchOrganizations(ctx echo.Context) error {
	customer, err := w.getCustomer(ctx)
	if err != nil {
		return err
	}

	var request types.SearchOrganizationsJSONRequestBody
	if err := ctx.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	organizations, err := w.NutsClient.SearchDiscoveryService(ctx.Request().Context(), request.Query, request.DiscoveryServiceID, request.DidServiceType)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var results = make(map[string]types.Organization, 0)
	for _, organization := range organizations {
		// Hide our own organization
		if request.ExcludeOwn != nil && *request.ExcludeOwn && organization.ID == *customer.Did {
			continue
		}
		current, exists := results[organization.ID]
		if !exists {
			current = types.FromNutsOrganization(organization.NutsOrganization)
		}
		current.DiscoveryServices = append(current.DiscoveryServices, organization.ServiceID)
		results[organization.ID] = current
	}

	return ctx.JSON(http.StatusOK, results)
}
