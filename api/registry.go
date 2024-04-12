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
	searchResults, err := w.NutsClient.SearchDiscoveryService(ctx.Request().Context(), request.Query, request.DiscoveryServiceID, request.DidServiceType)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var results = make(map[string]types.Organization, 0)
	for _, searchResult := range searchResults {
		// Hide our own organization
		if request.ExcludeOwn != nil && *request.ExcludeOwn && searchResult.ID == *customer.Did {
			continue
		}
		current, exists := results[searchResult.ID]
		if !exists {
			current = types.FromNutsOrganization(searchResult.NutsOrganization)
		}
		if ura, ok := searchResult.Fields["organization_ura"].(string); ok {
			current.Identifiers["ura"] = ura
		}
		current.DiscoveryServices = append(current.DiscoveryServices, searchResult.ServiceID)
		results[searchResult.ID] = current
	}

	return ctx.JSON(http.StatusOK, results)
}
