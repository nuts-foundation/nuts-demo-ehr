package api

import (
	"net/http"
	"slices"

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
	organizations, err := w.NutsClient.SearchDiscoveryService(ctx.Request().Context(), request.Query, request.DiscoveryServiceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// get all customer.Id DIDs
	dids, err := w.NutsClient.ListSubjectDIDs(ctx.Request().Context(), customer.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var results = make(map[string]types.Organization, 0)
	for _, organization := range organizations {
		// Hide our own organization
		if request.ExcludeOwn != nil && *request.ExcludeOwn && slices.Contains(dids, organization.ID) {
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
