package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
)

type CreateDossierRequest = types.CreateDossierRequest

func (w Wrapper) GetDossier(ctx echo.Context, patientID string) error {
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}

	dossiers, err := w.DossierRepository.AllByPatient(ctx.Request().Context(), cid, patientID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, dossiers)
}

func (w Wrapper) CreateDossier(ctx echo.Context) error {
	request := types.CreateDossierRequest{}
	if err := ctx.Bind(&request); err != nil {
		return err
	}
	logrus.Infof("Creating dossier (name=%s, patientID=%s)", request.Name, request.PatientID)
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
	dossier, err := w.DossierRepository.Create(ctx.Request().Context(), cid, request.Name, string(request.PatientID))
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, dossier)
}
