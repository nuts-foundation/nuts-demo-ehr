package api

import (
	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/sharedcareplan"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	"net/http"
)

type GetPatientCarePlansParams = types.GetPatientCarePlansParams
type FHIRCodeableConcept = types.FHIRCodeableConcept
type FHIRIdentifier = types.FHIRIdentifier

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

	carePlan, err := w.SharedCarePlanService.FindByID(ctx.Request().Context(), cid, dossierID, true)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, carePlan)
}

func (w Wrapper) CreateCarePlanActivity(ctx echo.Context, dossierID string) error {
	if w.SharedCarePlanService == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Shared Care Planning is not enabled")
	}
	customer, err := w.getCustomer(ctx)
	if err != nil {
		return err
	}
	if customer.URA == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Customer has no URA")
	}
	request := types.CreateCarePlanActivityRequest{}
	if err := ctx.Bind(&request); err != nil {
		return err
	}
	if len(request.Code.Coding) == 0 || request.Code.Coding[0].Code == nil ||
		request.Code.Coding[0].System == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid code")
	}

	requestor := sharedcareplan.MakeIdentifier("http://fhir.nl/fhir/NamingSystem/ura", *customer.URA)
	err = w.SharedCarePlanService.CreateActivity(ctx.Request().Context(), customer.Id, dossierID, request.Code, *requestor, request.Owner)
	if err != nil {
		return err
	}
	result, err := w.SharedCarePlanService.FindByID(ctx.Request().Context(), customer.Id, dossierID, true)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, result)
}
