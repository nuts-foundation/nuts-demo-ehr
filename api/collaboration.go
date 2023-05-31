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

	customer, err := w.getCustomer(ctx)
	if err != nil {
		return err
	}

	dossier, err := w.DossierRepository.FindByID(ctx.Request().Context(), customer.Id, dossierID)
	if err != nil {
		return err
	}

	patient, err := w.PatientRepository.FindByID(ctx.Request().Context(), customer.Id, string(dossier.PatientID))
	if err != nil {
		return err
	}

	if patient.Ssn == nil {
		return errors.New("no SSN registered for patient")
	}

	if err := w.EpisodeService.CreateCollaboration(
		ctx.Request().Context(),
		*customer.Did,
		dossierID,
		*patient.Ssn,
		request.Sender.Did,
	); err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, nil)
}

func (w Wrapper) GetCollaboration(ctx echo.Context, dossierID string) error {
	customer, err := w.getCustomer(ctx)
	if err != nil {
		return err
	}
	if customer.Did == nil {
		return errors.New("DID missing for customer")
	}

	dossier, err := w.DossierRepository.FindByID(ctx.Request().Context(), customer.Id, dossierID)
	if err != nil {
		return err
	}

	patient, err := w.PatientRepository.FindByID(ctx.Request().Context(), customer.Id, string(dossier.PatientID))
	if err != nil {
		return err
	}

	if patient.Ssn == nil {
		return errors.New("no SSN registered for patient")
	}

	// We want to find collaborations pointing to us, so we don't want to search on the customer DID
	collaborations, err := w.EpisodeService.GetCollaborations(ctx.Request().Context(), "", dossierID, *patient.Ssn)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, collaborations)
}
