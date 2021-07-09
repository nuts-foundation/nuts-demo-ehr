package api

import (
	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

type GetDossierParams = domain.GetDossierParams

func (w Wrapper) GetDossier(ctx echo.Context, params GetDossierParams) error {
	panic("not implemented")
}
