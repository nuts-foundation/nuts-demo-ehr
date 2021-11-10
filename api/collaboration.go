package api

import (
	"errors"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CreateCollaborationRequest = types.CreateCollaborationRequest

func (w Wrapper) CreateCollaboration(ctx echo.Context, dossierID string) error {
	request := &CreateCollaborationRequest{}

	if err := ctx.Bind(request); err != nil {
		return err
	}

	customerDID := w.getCustomerDID(ctx)

	if customerDID == nil {
		return errors.New("DID missing for customer")
	}

	episode, err := w.getEpisode(ctx, dossierID)
	if err != nil {
		return err
	}

	if err := w.EpisodeService.CreateCollaboration(ctx.Request().Context(), *customerDID, string(episode.Id), request.Sender.Did); err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, nil)
}
