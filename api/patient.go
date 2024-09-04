package api

import (
	"errors"
	"fmt"
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

type GetRemotePatientParams = types.GetRemotePatientParams

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
		c.FirstName = patientProps.FirstName
		c.Surname = patientProps.Surname
		c.Ssn = patientProps.Ssn
		c.Dob = patientProps.Dob
		c.Zipcode = patientProps.Zipcode
		c.AvatarUrl = patientProps.AvatarUrl
		c.Email = patientProps.Email
		c.Gender = patientProps.Gender
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

func (w Wrapper) GetRemotePatient(ctx echo.Context, params GetRemotePatientParams) error {
	customer, err := w.getCustomer(ctx)
	if err != nil {
		return err
	}
	patient, err := w.ZorginzageService.RemotePatient(ctx.Request().Context(), customer.Id, params.RemotePartyDID, params.PatientSSN)
	if err != nil {
		return fmt.Errorf("unable to load remote patient: %w", err)
	}
	return ctx.JSON(http.StatusOK, patient)
}

func (w Wrapper) getSessionID(ctx echo.Context) (string, error) {
	sessionID, ok := ctx.Get(SessionID).(string)
	if !ok {
		return "", errors.New("no active session")
	}
	return sessionID, nil
}

func (w Wrapper) getSession(ctx echo.Context) (*Session, error) {
	sessionID, err := w.getSessionID(ctx)
	if err != nil {
		return nil, err
	}
	session := w.APIAuth.GetSession(sessionID)
	if session == nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "unknown session ID")
	}
	sessionCopy := *session // Return copy
	return &sessionCopy, nil
}

func (w Wrapper) getCustomerID(ctx echo.Context) (string, error) {
	session, err := w.getSession(ctx)
	if err != nil {
		return "", err
	}
	return session.CustomerID, nil
}

func (w Wrapper) getCustomer(ctx echo.Context) (*types.Customer, error) {
	customerID, err := w.getCustomerID(ctx)
	if err != nil {
		return nil, err
	}
	result, err := w.CustomerRepository.FindByID(customerID)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("customer not found")
	}
	return result, nil
}
