package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/sirupsen/logrus"
)

type GetDossierParams = domain.GetDossierParams
type CreateDossierRequest = domain.CreateDossierRequest

func (w Wrapper) GetDossier(ctx echo.Context, params GetDossierParams) error {
	if params.PatientID != "" {
		dossiers, err := w.DossierRepository.AllByPatient(ctx.Request().Context(), w.getCustomerID(), params.PatientID)
		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusOK, dossiers)
	}
	return echo.NewHTTPError(http.StatusBadRequest, "missing patient_id query param")
}

func (w Wrapper) CreateDossier(ctx echo.Context) error {
	request := domain.CreateDossierRequest{}
	if err := ctx.Bind(&request); err != nil {
		return err
	}
	logrus.Infof("Creating dossier (name=%s, patientID=%s)", request.Name, request.PatientID)
	dossier, err := w.DossierRepository.Create(ctx.Request().Context(), w.getCustomerID(), request.Name, string(request.PatientID))
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, dossier)
}
