package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

func (w Wrapper) GetPatients(ctx echo.Context) error {
	customerID := w.getCustomerID()
	patients, err := w.PatientRepository.All(ctx.Request().Context(), customerID)
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

	patient, err := w.PatientRepository.NewPatient(ctx.Request().Context(), w.getCustomerID(), patientProperties)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, patient)
}


func (w Wrapper) getCustomerID() string {
	var customerID string
	// TODO: Determine customer ID from auth token
	customers, _ := w.CustomerRepository.All()
	if len(customers) > 0 {
		customerID = customers[0].Id
	}
	return customerID
}

