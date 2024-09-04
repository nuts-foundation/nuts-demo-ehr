package api

import (
	"errors"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
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
	if request.Sender.Did == "" {
		return errors.New("DID missing for other party")
	}

	dossier, err := w.DossierRepository.FindByID(ctx.Request().Context(), customer.Id, dossierID)
	if err != nil {
		return err
	}

	patient, err := w.PatientRepository.FindByID(ctx.Request().Context(), customer.Id, dossier.PatientID)
	if err != nil {
		return err
	}

	if patient.Ssn == nil {
		return errors.New("no SSN registered for patient")
	}

	if err := w.EpisodeService.CreateCollaboration(
		ctx.Request().Context(),
		customer.Id,
		dossierID,
		*patient.Ssn,
		request.Sender.Did,
		w.FHIRService.ClientFactory(fhir.WithTenant(customer.Id)),
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
	// TODO: We changed this API to showing organizations we shared this episode with, need to
	//       add another method that shows organizations that shared with us
	collaborations, err := w.EpisodeService.GetCollaborations(ctx.Request().Context(), customer.Id, dossierID, *patient.Ssn, w.FHIRService.ClientFactory(fhir.WithTenant(customer.Id)))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, collaborations)
}
