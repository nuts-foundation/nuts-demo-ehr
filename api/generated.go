// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.7.1 DO NOT EDIT.
package api

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /auth)
	SetCustomer(ctx echo.Context) error

	// (POST /auth/irma/session)
	AuthenticateWithIRMA(ctx echo.Context) error

	// (GET /auth/irma/session/{sessionToken}/result)
	GetIRMAAuthenticationResult(ctx echo.Context, sessionToken string) error

	// (POST /auth/passwd)
	AuthenticateWithPassword(ctx echo.Context) error

	// (GET /customers)
	ListCustomers(ctx echo.Context) error

	// (GET /private)
	CheckSession(ctx echo.Context) error

	// (GET /private/customer)
	GetCustomer(ctx echo.Context) error

	// (GET /private/dossier)
	GetDossier(ctx echo.Context, params GetDossierParams) error

	// (POST /private/dossier)
	CreateDossier(ctx echo.Context) error

	// (GET /private/network/organizations)
	SearchOrganizations(ctx echo.Context, params SearchOrganizationsParams) error

	// (GET /private/patient/{patientID})
	GetPatient(ctx echo.Context, patientID string) error

	// (PUT /private/patient/{patientID})
	UpdatePatient(ctx echo.Context, patientID string) error

	// (GET /private/patients)
	GetPatients(ctx echo.Context) error

	// (POST /private/patients)
	NewPatient(ctx echo.Context) error

	// (GET /private/transfer)
	GetPatientTransfers(ctx echo.Context, params GetPatientTransfersParams) error

	// (POST /private/transfer)
	CreateTransfer(ctx echo.Context) error

	// (GET /private/transfer/{transferID})
	GetTransfer(ctx echo.Context, transferID string) error

	// (GET /private/transfer/{transferID}/negotiation)
	ListTransferNegotiations(ctx echo.Context, transferID string) error

	// (POST /private/transfer/{transferID}/negotiation/{organizationDID})
	StartTransferNegotiation(ctx echo.Context, transferID string, organizationDID string) error

	// (POST /private/transfer/{transferID}/negotiation/{organizationDID}/assign)
	AssignTransferNegotiation(ctx echo.Context, transferID string, organizationDID string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// SetCustomer converts echo context to params.
func (w *ServerInterfaceWrapper) SetCustomer(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.SetCustomer(ctx)
	return err
}

// AuthenticateWithIRMA converts echo context to params.
func (w *ServerInterfaceWrapper) AuthenticateWithIRMA(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.AuthenticateWithIRMA(ctx)
	return err
}

// GetIRMAAuthenticationResult converts echo context to params.
func (w *ServerInterfaceWrapper) GetIRMAAuthenticationResult(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "sessionToken" -------------
	var sessionToken string

	err = runtime.BindStyledParameterWithLocation("simple", false, "sessionToken", runtime.ParamLocationPath, ctx.Param("sessionToken"), &sessionToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter sessionToken: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetIRMAAuthenticationResult(ctx, sessionToken)
	return err
}

// AuthenticateWithPassword converts echo context to params.
func (w *ServerInterfaceWrapper) AuthenticateWithPassword(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.AuthenticateWithPassword(ctx)
	return err
}

// ListCustomers converts echo context to params.
func (w *ServerInterfaceWrapper) ListCustomers(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ListCustomers(ctx)
	return err
}

// CheckSession converts echo context to params.
func (w *ServerInterfaceWrapper) CheckSession(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CheckSession(ctx)
	return err
}

// GetCustomer converts echo context to params.
func (w *ServerInterfaceWrapper) GetCustomer(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetCustomer(ctx)
	return err
}

// GetDossier converts echo context to params.
func (w *ServerInterfaceWrapper) GetDossier(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetDossierParams
	// ------------- Required query parameter "patientID" -------------

	err = runtime.BindQueryParameter("form", true, true, "patientID", ctx.QueryParams(), &params.PatientID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter patientID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetDossier(ctx, params)
	return err
}

// CreateDossier converts echo context to params.
func (w *ServerInterfaceWrapper) CreateDossier(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateDossier(ctx)
	return err
}

// SearchOrganizations converts echo context to params.
func (w *ServerInterfaceWrapper) SearchOrganizations(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params SearchOrganizationsParams
	// ------------- Required query parameter "query" -------------

	err = runtime.BindQueryParameter("form", true, true, "query", ctx.QueryParams(), &params.Query)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter query: %s", err))
	}

	// ------------- Optional query parameter "didServiceType" -------------

	err = runtime.BindQueryParameter("form", true, false, "didServiceType", ctx.QueryParams(), &params.DidServiceType)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter didServiceType: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.SearchOrganizations(ctx, params)
	return err
}

// GetPatient converts echo context to params.
func (w *ServerInterfaceWrapper) GetPatient(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "patientID" -------------
	var patientID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "patientID", runtime.ParamLocationPath, ctx.Param("patientID"), &patientID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter patientID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetPatient(ctx, patientID)
	return err
}

// UpdatePatient converts echo context to params.
func (w *ServerInterfaceWrapper) UpdatePatient(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "patientID" -------------
	var patientID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "patientID", runtime.ParamLocationPath, ctx.Param("patientID"), &patientID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter patientID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UpdatePatient(ctx, patientID)
	return err
}

// GetPatients converts echo context to params.
func (w *ServerInterfaceWrapper) GetPatients(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetPatients(ctx)
	return err
}

// NewPatient converts echo context to params.
func (w *ServerInterfaceWrapper) NewPatient(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.NewPatient(ctx)
	return err
}

// GetPatientTransfers converts echo context to params.
func (w *ServerInterfaceWrapper) GetPatientTransfers(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetPatientTransfersParams
	// ------------- Required query parameter "patientID" -------------

	err = runtime.BindQueryParameter("form", true, true, "patientID", ctx.QueryParams(), &params.PatientID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter patientID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetPatientTransfers(ctx, params)
	return err
}

// CreateTransfer converts echo context to params.
func (w *ServerInterfaceWrapper) CreateTransfer(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateTransfer(ctx)
	return err
}

// GetTransfer converts echo context to params.
func (w *ServerInterfaceWrapper) GetTransfer(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "transferID" -------------
	var transferID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "transferID", runtime.ParamLocationPath, ctx.Param("transferID"), &transferID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter transferID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetTransfer(ctx, transferID)
	return err
}

// ListTransferNegotiations converts echo context to params.
func (w *ServerInterfaceWrapper) ListTransferNegotiations(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "transferID" -------------
	var transferID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "transferID", runtime.ParamLocationPath, ctx.Param("transferID"), &transferID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter transferID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ListTransferNegotiations(ctx, transferID)
	return err
}

// StartTransferNegotiation converts echo context to params.
func (w *ServerInterfaceWrapper) StartTransferNegotiation(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "transferID" -------------
	var transferID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "transferID", runtime.ParamLocationPath, ctx.Param("transferID"), &transferID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter transferID: %s", err))
	}

	// ------------- Path parameter "organizationDID" -------------
	var organizationDID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "organizationDID", runtime.ParamLocationPath, ctx.Param("organizationDID"), &organizationDID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter organizationDID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.StartTransferNegotiation(ctx, transferID, organizationDID)
	return err
}

// AssignTransferNegotiation converts echo context to params.
func (w *ServerInterfaceWrapper) AssignTransferNegotiation(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "transferID" -------------
	var transferID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "transferID", runtime.ParamLocationPath, ctx.Param("transferID"), &transferID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter transferID: %s", err))
	}

	// ------------- Path parameter "organizationDID" -------------
	var organizationDID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "organizationDID", runtime.ParamLocationPath, ctx.Param("organizationDID"), &organizationDID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter organizationDID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.AssignTransferNegotiation(ctx, transferID, organizationDID)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/auth", wrapper.SetCustomer)
	router.POST(baseURL+"/auth/irma/session", wrapper.AuthenticateWithIRMA)
	router.GET(baseURL+"/auth/irma/session/:sessionToken/result", wrapper.GetIRMAAuthenticationResult)
	router.POST(baseURL+"/auth/passwd", wrapper.AuthenticateWithPassword)
	router.GET(baseURL+"/customers", wrapper.ListCustomers)
	router.GET(baseURL+"/private", wrapper.CheckSession)
	router.GET(baseURL+"/private/customer", wrapper.GetCustomer)
	router.GET(baseURL+"/private/dossier", wrapper.GetDossier)
	router.POST(baseURL+"/private/dossier", wrapper.CreateDossier)
	router.GET(baseURL+"/private/network/organizations", wrapper.SearchOrganizations)
	router.GET(baseURL+"/private/patient/:patientID", wrapper.GetPatient)
	router.PUT(baseURL+"/private/patient/:patientID", wrapper.UpdatePatient)
	router.GET(baseURL+"/private/patients", wrapper.GetPatients)
	router.POST(baseURL+"/private/patients", wrapper.NewPatient)
	router.GET(baseURL+"/private/transfer", wrapper.GetPatientTransfers)
	router.POST(baseURL+"/private/transfer", wrapper.CreateTransfer)
	router.GET(baseURL+"/private/transfer/:transferID", wrapper.GetTransfer)
	router.GET(baseURL+"/private/transfer/:transferID/negotiation", wrapper.ListTransferNegotiations)
	router.POST(baseURL+"/private/transfer/:transferID/negotiation/:organizationDID", wrapper.StartTransferNegotiation)
	router.POST(baseURL+"/private/transfer/:transferID/negotiation/:organizationDID/assign", wrapper.AssignTransferNegotiation)

}
