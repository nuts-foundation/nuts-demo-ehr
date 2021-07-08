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

	// (POST /web/auth)
	SetCustomer(ctx echo.Context) error

	// (POST /web/auth/irma/session)
	AuthenticateWithIRMA(ctx echo.Context) error

	// (GET /web/auth/irma/session/{sessionToken}/result)
	GetIRMAAuthenticationResult(ctx echo.Context, sessionToken string) error

	// (POST /web/auth/passwd)
	AuthenticateWithPassword(ctx echo.Context) error

	// (GET /web/customers)
	ListCustomers(ctx echo.Context) error

	// (GET /web/private)
	CheckSession(ctx echo.Context) error

	// (GET /web/private/customer)
	GetCustomer(ctx echo.Context) error

	// (GET /web/private/network/organizations)
	SearchOrganizations(ctx echo.Context, params SearchOrganizationsParams) error

	// (GET /web/private/patient/{patientID})
	GetPatient(ctx echo.Context, patientID string) error

	// (PUT /web/private/patient/{patientID})
	UpdatePatient(ctx echo.Context, patientID string) error

	// (POST /web/private/patient/{patientID}/transfer)
	CreateTransfer(ctx echo.Context, patientID string) error

	// (GET /web/private/patient/{patientID}/transfer/{transferID})
	GetTransfer(ctx echo.Context, patientID string, transferID string) error

	// (GET /web/private/patient/{patientID}/transfer/{transferID}/negotiation)
	ListTransferNegotiations(ctx echo.Context, patientID string, transferID string) error

	// (POST /web/private/patient/{patientID}/transfer/{transferID}/negotiation)
	StartTransferNegotiation(ctx echo.Context, patientID string, transferID string) error

	// (PUT /web/private/patient/{patientID}/transfer/{transferID}/negotiation/{negotiationID})
	AcceptTransferNegotiation(ctx echo.Context, patientID string, transferID string, negotiationID string) error

	// (GET /web/private/patients)
	GetPatients(ctx echo.Context) error

	// (POST /web/private/patients)
	NewPatient(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// SetCustomer converts echo context to params.
func (w *ServerInterfaceWrapper) SetCustomer(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.SetCustomer(ctx)
	return err
}

// AuthenticateWithIRMA converts echo context to params.
func (w *ServerInterfaceWrapper) AuthenticateWithIRMA(ctx echo.Context) error {
	var err error

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

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetIRMAAuthenticationResult(ctx, sessionToken)
	return err
}

// AuthenticateWithPassword converts echo context to params.
func (w *ServerInterfaceWrapper) AuthenticateWithPassword(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.AuthenticateWithPassword(ctx)
	return err
}

// ListCustomers converts echo context to params.
func (w *ServerInterfaceWrapper) ListCustomers(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ListCustomers(ctx)
	return err
}

// CheckSession converts echo context to params.
func (w *ServerInterfaceWrapper) CheckSession(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CheckSession(ctx)
	return err
}

// GetCustomer converts echo context to params.
func (w *ServerInterfaceWrapper) GetCustomer(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetCustomer(ctx)
	return err
}

// SearchOrganizations converts echo context to params.
func (w *ServerInterfaceWrapper) SearchOrganizations(ctx echo.Context) error {
	var err error

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

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UpdatePatient(ctx, patientID)
	return err
}

// CreateTransfer converts echo context to params.
func (w *ServerInterfaceWrapper) CreateTransfer(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "patientID" -------------
	var patientID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "patientID", runtime.ParamLocationPath, ctx.Param("patientID"), &patientID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter patientID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateTransfer(ctx, patientID)
	return err
}

// GetTransfer converts echo context to params.
func (w *ServerInterfaceWrapper) GetTransfer(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "patientID" -------------
	var patientID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "patientID", runtime.ParamLocationPath, ctx.Param("patientID"), &patientID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter patientID: %s", err))
	}

	// ------------- Path parameter "transferID" -------------
	var transferID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "transferID", runtime.ParamLocationPath, ctx.Param("transferID"), &transferID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter transferID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetTransfer(ctx, patientID, transferID)
	return err
}

// ListTransferNegotiations converts echo context to params.
func (w *ServerInterfaceWrapper) ListTransferNegotiations(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "patientID" -------------
	var patientID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "patientID", runtime.ParamLocationPath, ctx.Param("patientID"), &patientID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter patientID: %s", err))
	}

	// ------------- Path parameter "transferID" -------------
	var transferID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "transferID", runtime.ParamLocationPath, ctx.Param("transferID"), &transferID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter transferID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ListTransferNegotiations(ctx, patientID, transferID)
	return err
}

// StartTransferNegotiation converts echo context to params.
func (w *ServerInterfaceWrapper) StartTransferNegotiation(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "patientID" -------------
	var patientID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "patientID", runtime.ParamLocationPath, ctx.Param("patientID"), &patientID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter patientID: %s", err))
	}

	// ------------- Path parameter "transferID" -------------
	var transferID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "transferID", runtime.ParamLocationPath, ctx.Param("transferID"), &transferID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter transferID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.StartTransferNegotiation(ctx, patientID, transferID)
	return err
}

// AcceptTransferNegotiation converts echo context to params.
func (w *ServerInterfaceWrapper) AcceptTransferNegotiation(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "patientID" -------------
	var patientID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "patientID", runtime.ParamLocationPath, ctx.Param("patientID"), &patientID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter patientID: %s", err))
	}

	// ------------- Path parameter "transferID" -------------
	var transferID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "transferID", runtime.ParamLocationPath, ctx.Param("transferID"), &transferID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter transferID: %s", err))
	}

	// ------------- Path parameter "negotiationID" -------------
	var negotiationID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "negotiationID", runtime.ParamLocationPath, ctx.Param("negotiationID"), &negotiationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter negotiationID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.AcceptTransferNegotiation(ctx, patientID, transferID, negotiationID)
	return err
}

// GetPatients converts echo context to params.
func (w *ServerInterfaceWrapper) GetPatients(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetPatients(ctx)
	return err
}

// NewPatient converts echo context to params.
func (w *ServerInterfaceWrapper) NewPatient(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.NewPatient(ctx)
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

	router.POST(baseURL+"/web/auth", wrapper.SetCustomer)
	router.POST(baseURL+"/web/auth/irma/session", wrapper.AuthenticateWithIRMA)
	router.GET(baseURL+"/web/auth/irma/session/:sessionToken/result", wrapper.GetIRMAAuthenticationResult)
	router.POST(baseURL+"/web/auth/passwd", wrapper.AuthenticateWithPassword)
	router.GET(baseURL+"/web/customers", wrapper.ListCustomers)
	router.GET(baseURL+"/web/private", wrapper.CheckSession)
	router.GET(baseURL+"/web/private/customer", wrapper.GetCustomer)
	router.GET(baseURL+"/web/private/network/organizations", wrapper.SearchOrganizations)
	router.GET(baseURL+"/web/private/patient/:patientID", wrapper.GetPatient)
	router.PUT(baseURL+"/web/private/patient/:patientID", wrapper.UpdatePatient)
	router.POST(baseURL+"/web/private/patient/:patientID/transfer", wrapper.CreateTransfer)
	router.GET(baseURL+"/web/private/patient/:patientID/transfer/:transferID", wrapper.GetTransfer)
	router.GET(baseURL+"/web/private/patient/:patientID/transfer/:transferID/negotiation", wrapper.ListTransferNegotiations)
	router.POST(baseURL+"/web/private/patient/:patientID/transfer/:transferID/negotiation", wrapper.StartTransferNegotiation)
	router.PUT(baseURL+"/web/private/patient/:patientID/transfer/:transferID/negotiation/:negotiationID", wrapper.AcceptTransferNegotiation)
	router.GET(baseURL+"/web/private/patients", wrapper.GetPatients)
	router.POST(baseURL+"/web/private/patients", wrapper.NewPatient)

}

