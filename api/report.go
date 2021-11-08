package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
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
