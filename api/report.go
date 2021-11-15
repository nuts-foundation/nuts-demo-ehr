package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
)

// GetReportsParams defines parameters for GetReports.
type GetReportsParams struct {
	// The identifier of episode the report must be part of.
	EpisodeID *string `json:"episodeID,omitempty"`
}

func (w Wrapper) GetReports(ctx echo.Context, patientID string, params GetReportsParams) error {
	customer := w.getCustomer(ctx)

	// Get the local reports for the patient
	reports, err := w.ReportRepository.AllByPatient(ctx.Request().Context(), customer.Id, patientID, params.EpisodeID)
	if err != nil {
		return err
	}

	// Get the remote reports for the patient
	patient, err := w.PatientRepository.FindByID(ctx.Request().Context(), customer.Id, patientID)
	if err != nil {
		return err
	}

	if patient.Ssn != nil {
		remoteReports, err := w.EpisodeService.GetReports(ctx.Request().Context(), *customer.Did, *patient.Ssn)
		if err != nil {
			return err
		}

		reports = append(reports, remoteReports...)
	}

	return ctx.JSON(http.StatusOK, reports)
}

func (w Wrapper) CreateReport(ctx echo.Context, patientID string) error {
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}

	reportToCreate := types.Report{}
	if err := ctx.Bind(&reportToCreate); err != nil {
		return err
	}

	if err = w.ReportRepository.Create(ctx.Request().Context(), cid, patientID, reportToCreate); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}
