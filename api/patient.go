package api

import (
	"errors"
	"net/http"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
)

// GetPatientsParams defines parameters for GetPatients.
type GetPatientsParams struct {

	// Search patients by name
	Name *string `json:"name,omitempty"`
}

func (w Wrapper) GetPatients(ctx echo.Context, params GetPatientsParams) error {
	customerID, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
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
	patientProperties := types.PatientProperties{}
	if err := ctx.Bind(&patientProperties); err != nil {
		return err
	}

	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
	patient, err := w.PatientRepository.NewPatient(ctx.Request().Context(), cid, patientProperties)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, patient)
}

func (w Wrapper) UpdatePatient(ctx echo.Context, patientID string) error {
	patientProps := types.PatientProperties{}
	if err := ctx.Bind(&patientProps); err != nil {
		return err
	}
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
	patient, err := w.PatientRepository.Update(ctx.Request().Context(), cid, patientID, func(c types.Patient) (*types.Patient, error) {
		c.PatientProperties = patientProps
		return &c, nil
	})
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, patient)
}

func (w Wrapper) GetPatient(ctx echo.Context, patientID string) error {
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}
	patient, err := w.PatientRepository.FindByID(ctx.Request().Context(), cid, patientID)
	if err != nil {
		return err
	}
	if patient == nil {
		return ctx.NoContent(http.StatusNotFound)
	}
	return ctx.JSON(http.StatusOK, patient)

}

func (w Wrapper) getCustomerID(ctx echo.Context) (int, error) {
	customer := w.getCustomer(ctx)
	if customer == nil {
		return 0, errors.New("not found")
	}
	return customer.Id, nil
}

func (w Wrapper) getCustomer(ctx echo.Context) *types.Customer {
	cid, ok := ctx.Get(CustomerID).(int)
	if !ok {
		return nil
	}
	customer, _ := w.CustomerRepository.FindByID(cid)
	if customer.Id != cid {
		return nil
	}
	return customer
}

func (w Wrapper) getCustomerDID(ctx echo.Context) *string {
	cid, ok := ctx.Get(CustomerID).(int)
	if !ok {
		return nil
	}
	customer, _ := w.CustomerRepository.FindByID(cid)
	if customer.Id != cid {
		return nil
	}

	return customer.Did
}
