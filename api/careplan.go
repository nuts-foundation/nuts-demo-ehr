package api

import (
	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	"net/http"
)

type GetPatientCarePlansParams = types.GetPatientCarePlansParams

func (w Wrapper) CreateCarePlan(ctx echo.Context) error {
	if w.SharedCarePlanService == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Shared Care Planning is not enabled")
	}
	request := types.CreateCarePlanRequest{}
	if err := ctx.Bind(&request); err != nil {
		return err
	}

	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}

	carePlan, err := w.SharedCarePlanService.Create(ctx.Request().Context(), cid, request.DossierID, request.Title)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, *carePlan)
}

func (w Wrapper) GetPatientCarePlans(ctx echo.Context, params types.GetPatientCarePlansParams) error {
	if w.SharedCarePlanService == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Shared Care Planning is not enabled")
	}
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}

	carePlans, err := w.SharedCarePlanService.AllForPatient(ctx.Request().Context(), cid, params.PatientID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, carePlans)
}

func (w Wrapper) GetCarePlan(ctx echo.Context, dossierID string) error {
	if w.SharedCarePlanService == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Shared Care Planning is not enabled")
	}
	cid, err := w.getCustomerID(ctx)
	if err != nil {
		return err
	}

	carePlan, err := w.SharedCarePlanService.FindByID(ctx.Request().Context(), cid, dossierID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, carePlan)
}
