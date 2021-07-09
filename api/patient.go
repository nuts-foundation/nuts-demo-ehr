package api

import (
	"net/http"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

func (w Wrapper) GetPatients(ctx echo.Context) error {
	customerID := w.getCustomerID()
	patients, err := w.PatientRepository.All(ctx.Request().Context(), customerID)
	if err != nil {
		return err
	}
	// Sort patients by surname
	sort.Slice(patients, func(i, j int) bool {
		s1 := patients[i].Surname
		s2 := patients[j].Surname
		if s1 == "" {
			return true
		}
		if s2 == "" {
			return false
		}
		return strings.Compare(s1, s2) < 0
	})
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

func (w Wrapper) UpdatePatient(ctx echo.Context, patientID string) error {
	patientProps := domain.PatientProperties{}
	if err := ctx.Bind(&patientProps); err != nil {
		return err
	}
	patient, err := w.PatientRepository.Update(ctx.Request().Context(), w.getCustomerID(), patientID, func(c domain.Patient) (*domain.Patient, error) {
		c.PatientProperties = patientProps
		return &c, nil
	})
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, patient)
}

func (w Wrapper) GetPatient(ctx echo.Context, patientID string) error {
	patient, err := w.PatientRepository.FindByID(ctx.Request().Context(), w.getCustomerID(), patientID)
	if err != nil {
		return err
	}
	if patient == nil {
		return ctx.NoContent(http.StatusNotFound)
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

func (w Wrapper) getCustomerDID() *string {
	customer, _ := w.CustomerRepository.FindByID(w.getCustomerID())
	return customer.Did
}
