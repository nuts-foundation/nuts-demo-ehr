package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/monarko/fhirgo/STU3/resources"
)

func (w Wrapper) TaskUpdate(ctx echo.Context, customerID int, taskID string) error {
	// get customer
	customer, err := w.CustomerRepository.FindByID(customerID)
	if err != nil {
		return err
	}
	if customer == nil {
		// shouldn't happen since this is an internal call
		return echo.NewHTTPError(http.StatusNotFound, "customer unknown")
	}

	// marshal body into fhir task
	task := resources.Task{}
	if err = ctx.Bind(&task); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if task.Status == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Task.Status not found")
	}
	status := *task.Status

	// update existing task
	err = w.TransferSenderService.UpdateTask(ctx.Request().Context(), *customer, taskID, string(status))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return err
	}

	return ctx.NoContent(http.StatusAccepted)
}
