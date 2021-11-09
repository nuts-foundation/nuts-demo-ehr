package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
)

func (w Wrapper) CreateCollaboration(ctx echo.Context) error {
	request := types.CreateCollaborationRequest{}

	if err := ctx.Bind(&request); err != nil {
		return err
	}

	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}

	dossier, err := w.DossierRepository.FindByID(ctx.Request().Context(), cid, string(request.DossierID))
	if err != nil {
		return err
	}

	collaboration, err := w.CollaborationService.Create(ctx.Request().Context(), cid, string(dossier.Id), string(dossier.PatientID))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, collaboration)
}

func (w Wrapper) GetCollaboration(ctx echo.Context, dossierID string) error {
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}

	dossier, err := w.DossierRepository.FindByID(ctx.Request().Context(), cid, dossierID)
	if err != nil {
		return err
	}

	collaboration, err := w.CollaborationService.Get(ctx.Request().Context(), cid, string(dossier.Id))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, collaboration)
}
