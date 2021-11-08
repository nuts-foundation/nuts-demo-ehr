package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
)

func (w Wrapper) GetReports(ctx echo.Context, patientID string) error {
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}

	reports, err := w.ReportRepository.AllByPatient(ctx.Request().Context(), cid, patientID)
	if err != nil {
		return err
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
