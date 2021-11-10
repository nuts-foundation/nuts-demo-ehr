package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
)

func (w Wrapper) CreateEpisode(ctx echo.Context) error {
	request := types.CreateEpisodeRequest{}

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

	episode, err := w.EpisodeService.Create(ctx.Request().Context(), cid, string(dossier.PatientID), request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, episode)
}

func (w Wrapper) getEpisode(ctx echo.Context, dossierID string) (*types.Episode, error) {
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return nil, err
	}

	dossier, err := w.DossierRepository.FindByID(ctx.Request().Context(), cid, dossierID)
	if err != nil {
		return nil, err
	}
	if dossier == nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "dossier not found")
	}

	episode, err := w.EpisodeService.Get(ctx.Request().Context(), cid, string(dossier.Id))
	if err != nil {
		return nil, err
	}

	return episode, nil
}

func (w Wrapper) GetEpisode(ctx echo.Context, dossierID string) error {
	episode, err := w.getEpisode(ctx, dossierID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, episode)
}
