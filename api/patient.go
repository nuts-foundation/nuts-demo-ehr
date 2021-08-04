package api

import (
	"net/http"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
)

// GetPatientsParams defines parameters for GetPatients.
type GetPatientsParams struct {

	// Search patients by name
	Name *string `json:"name,omitempty"`
}

func (w Wrapper) GetPatients(ctx echo.Context, params GetPatientsParams) error {
	customerID := w.getCustomerID(ctx)
	patients, err := w.PatientRepository.All(ctx.Request().Context(), customerID, params.Name)
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

	patient, err := w.PatientRepository.NewPatient(ctx.Request().Context(), w.getCustomerID(ctx), patientProperties)
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
	patient, err := w.PatientRepository.Update(ctx.Request().Context(), w.getCustomerID(ctx), patientID, func(c domain.Patient) (*domain.Patient, error) {
		c.PatientProperties = patientProps
		return &c, nil
	})
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, patient)
}

func (w Wrapper) GetPatient(ctx echo.Context, patientID string) error {
	patient, err := w.PatientRepository.FindByID(ctx.Request().Context(), w.getCustomerID(ctx), patientID)
	if err != nil {
		return err
	}
	if patient == nil {
		return ctx.NoContent(http.StatusNotFound)
	}
	return ctx.JSON(http.StatusOK, patient)

}

func (w Wrapper) getCustomerID(ctx echo.Context) string {
	cid, ok := ctx.Get(CustomerID).(string)
	if !ok {
		return ""
	}
	customer, _ := w.CustomerRepository.FindByID(cid)
	if customer.Id != cid {
		return ""
	}

	return customer.Id
}

func (w Wrapper) getCustomerDID(ctx echo.Context) *string {
	cid, ok := ctx.Get(CustomerID).(string)
	if !ok {
		return nil
	}
	customer, _ := w.CustomerRepository.FindByID(cid)
	if customer.Id != cid {
		return nil
	}

	return customer.Did
}
