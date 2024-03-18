package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (w Wrapper) GetACL(ctx echo.Context, tenantDID string, authorizedDID string) error {
	authorizedResources, err := w.ACL.AuthorizedResources(ctx.Request().Context(), tenantDID, authorizedDID)
	if err != nil {
		return err
	}

	result := make(map[string][]string)
	for _, resource := range authorizedResources {
		result[resource.Resource] = append(result[resource.Resource], resource.Operation)
	}

	return ctx.JSON(http.StatusOK, result)
}
