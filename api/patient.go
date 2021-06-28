package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

func (w Wrapper) GetPatients(ctx echo.Context) error {
	patients, err := w.PatientRepository.All(ctx.Request().Context(), "c1")
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, patients)
}


func (w Wrapper) NewPatient(ctx echo.Context) error {
	patientProperties := domain.PatientProperties{}
	if err := ctx.Bind(&patientProperties); err != nil {
		return err
	}

	patient, err := w.PatientRepository.NewPatient(ctx.Request().Context(), "c1", patientProperties)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, patient)
}

