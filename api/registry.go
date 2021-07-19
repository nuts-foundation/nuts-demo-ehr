package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

type SearchOrganizationsParams = domain.SearchOrganizationsParams

func (w Wrapper) SearchOrganizations(ctx echo.Context, params SearchOrganizationsParams) error {
	organizations, err := w.Client.SearchOrganizations(ctx.Request().Context(), params.Query)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	// TODO: Filter on params.didServiceType
	var results []domain.Organization
	for _, curr := range organizations {
		concept := curr["organization"].(map[string]interface{})
		results = append(results, domain.Organization{
			Did:  curr["subject"].(string),
			City: concept["city"].(string),
			Name: concept["name"].(string),
		})
	}
	return ctx.JSON(http.StatusOK, results)
}
