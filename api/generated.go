// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /auth)
	SetCustomer(ctx echo.Context) error

	// (POST /auth/dummy)
	AuthenticateWithDummy(ctx echo.Context) error

	// (GET /auth/dummy/session/{sessionToken}/result)
	GetDummyAuthenticationResult(ctx echo.Context, sessionToken string) error

	// (POST /auth/employeeid/session)
	AuthenticateWithEmployeeID(ctx echo.Context) error

	// (GET /auth/employeeid/session/{sessionToken}/result)
	GetEmployeeIDAuthenticationResult(ctx echo.Context, sessionToken string) error

	// (POST /auth/irma/session)
	AuthenticateWithIRMA(ctx echo.Context) error

	// (GET /auth/irma/session/{sessionToken}/result)
	GetIRMAAuthenticationResult(ctx echo.Context, sessionToken string) error

	// (POST /auth/openid4vp)
	CreateAuthorizationRequest(ctx echo.Context) error

	// (GET /auth/openid4vp/{token})
	GetOpenID4VPAuthenticationResult(ctx echo.Context, token string) error

	// (POST /auth/passwd)
	AuthenticateWithPassword(ctx echo.Context) error

	// (GET /customers)
	ListCustomers(ctx echo.Context) error

	// (POST /external/transfer/notify/{taskID})
	NotifyTransferUpdate(ctx echo.Context, taskID string) error

	// (PUT /internal/customer/{customerID}/task/{taskID})
	TaskUpdate(ctx echo.Context, customerID int, taskID string) error

	// (GET /private)
	CheckSession(ctx echo.Context) error

	// (GET /private/customer)
	GetCustomer(ctx echo.Context) error

	// (POST /private/dossier)
	CreateDossier(ctx echo.Context) error

	// (GET /private/dossier/{patientID})
	GetDossier(ctx echo.Context, patientID string) error

	// (POST /private/episode)
	CreateEpisode(ctx echo.Context) error

	// (GET /private/episode/{episodeID})
	GetEpisode(ctx echo.Context, episodeID string) error

	// (GET /private/episode/{episodeID}/collaboration)
	GetCollaboration(ctx echo.Context, episodeID string) error

	// (POST /private/episode/{episodeID}/collaboration)
	CreateCollaboration(ctx echo.Context, episodeID string) error

	// (GET /private/network/inbox)
	GetInbox(ctx echo.Context) error

	// (GET /private/network/inbox/info)
	GetInboxInfo(ctx echo.Context) error

	// (GET /private/network/organizations)
	SearchOrganizations(ctx echo.Context, params SearchOrganizationsParams) error

	// (GET /private/patient/{patientID})
	GetPatient(ctx echo.Context, patientID string) error

	// (PUT /private/patient/{patientID})
	UpdatePatient(ctx echo.Context, patientID string) error

	// (GET /private/patients)
	GetPatients(ctx echo.Context, params GetPatientsParams) error

	// (POST /private/patients)
	NewPatient(ctx echo.Context) error

	// (GET /private/reports/{patientID})
	GetReports(ctx echo.Context, patientID string, params GetReportsParams) error

	// (POST /private/reports/{patientID})
	CreateReport(ctx echo.Context, patientID string) error

	// (GET /private/transfer)
	GetPatientTransfers(ctx echo.Context, params GetPatientTransfersParams) error

	// (POST /private/transfer)
	CreateTransfer(ctx echo.Context) error

	// (GET /private/transfer-request/{requestorDID}/{fhirTaskID})
	GetTransferRequest(ctx echo.Context, requestorDID string, fhirTaskID string) error

	// (POST /private/transfer-request/{requestorDID}/{fhirTaskID})
	ChangeTransferRequestState(ctx echo.Context, requestorDID string, fhirTaskID string) error

	// (DELETE /private/transfer/{transferID})
	CancelTransfer(ctx echo.Context, transferID string) error

	// (GET /private/transfer/{transferID})
	GetTransfer(ctx echo.Context, transferID string) error

	// (PUT /private/transfer/{transferID})
	UpdateTransfer(ctx echo.Context, transferID string) error

	// (PUT /private/transfer/{transferID}/assign)
	AssignTransferDirect(ctx echo.Context, transferID string) error

	// (GET /private/transfer/{transferID}/negotiation)
	ListTransferNegotiations(ctx echo.Context, transferID string) error

	// (POST /private/transfer/{transferID}/negotiation)
	StartTransferNegotiation(ctx echo.Context, transferID string) error

	// (PUT /private/transfer/{transferID}/negotiation/{negotiationID})
	UpdateTransferNegotiationStatus(ctx echo.Context, transferID string, negotiationID string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// SetCustomer converts echo context to params.
func (w *ServerInterfaceWrapper) SetCustomer(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.SetCustomer(ctx)
	return err
}

// AuthenticateWithDummy converts echo context to params.
func (w *ServerInterfaceWrapper) AuthenticateWithDummy(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.AuthenticateWithDummy(ctx)
	return err
}

// GetDummyAuthenticationResult converts echo context to params.
func (w *ServerInterfaceWrapper) GetDummyAuthenticationResult(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "sessionToken" -------------
	var sessionToken string

	err = runtime.BindStyledParameterWithOptions("simple", "sessionToken", ctx.Param("sessionToken"), &sessionToken, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter sessionToken: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetDummyAuthenticationResult(ctx, sessionToken)
	return err
}

// AuthenticateWithEmployeeID converts echo context to params.
func (w *ServerInterfaceWrapper) AuthenticateWithEmployeeID(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.AuthenticateWithEmployeeID(ctx)
	return err
}

// GetEmployeeIDAuthenticationResult converts echo context to params.
func (w *ServerInterfaceWrapper) GetEmployeeIDAuthenticationResult(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "sessionToken" -------------
	var sessionToken string

	err = runtime.BindStyledParameterWithOptions("simple", "sessionToken", ctx.Param("sessionToken"), &sessionToken, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter sessionToken: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetEmployeeIDAuthenticationResult(ctx, sessionToken)
	return err
}

// AuthenticateWithIRMA converts echo context to params.
func (w *ServerInterfaceWrapper) AuthenticateWithIRMA(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.AuthenticateWithIRMA(ctx)
	return err
}

// GetIRMAAuthenticationResult converts echo context to params.
func (w *ServerInterfaceWrapper) GetIRMAAuthenticationResult(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "sessionToken" -------------
	var sessionToken string

	err = runtime.BindStyledParameterWithOptions("simple", "sessionToken", ctx.Param("sessionToken"), &sessionToken, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter sessionToken: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetIRMAAuthenticationResult(ctx, sessionToken)
	return err
}

// CreateAuthorizationRequest converts echo context to params.
func (w *ServerInterfaceWrapper) CreateAuthorizationRequest(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateAuthorizationRequest(ctx)
	return err
}

// GetOpenID4VPAuthenticationResult converts echo context to params.
func (w *ServerInterfaceWrapper) GetOpenID4VPAuthenticationResult(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "token" -------------
	var token string

	err = runtime.BindStyledParameterWithOptions("simple", "token", ctx.Param("token"), &token, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter token: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetOpenID4VPAuthenticationResult(ctx, token)
	return err
}

// AuthenticateWithPassword converts echo context to params.
func (w *ServerInterfaceWrapper) AuthenticateWithPassword(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.AuthenticateWithPassword(ctx)
	return err
}

// ListCustomers converts echo context to params.
func (w *ServerInterfaceWrapper) ListCustomers(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListCustomers(ctx)
	return err
}

// NotifyTransferUpdate converts echo context to params.
func (w *ServerInterfaceWrapper) NotifyTransferUpdate(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "taskID" -------------
	var taskID string

	err = runtime.BindStyledParameterWithOptions("simple", "taskID", ctx.Param("taskID"), &taskID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter taskID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.NotifyTransferUpdate(ctx, taskID)
	return err
}

// TaskUpdate converts echo context to params.
func (w *ServerInterfaceWrapper) TaskUpdate(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "customerID" -------------
	var customerID int

	err = runtime.BindStyledParameterWithOptions("simple", "customerID", ctx.Param("customerID"), &customerID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter customerID: %s", err))
	}

	// ------------- Path parameter "taskID" -------------
	var taskID string

	err = runtime.BindStyledParameterWithOptions("simple", "taskID", ctx.Param("taskID"), &taskID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter taskID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.TaskUpdate(ctx, customerID, taskID)
	return err
}

// CheckSession converts echo context to params.
func (w *ServerInterfaceWrapper) CheckSession(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CheckSession(ctx)
	return err
}

// GetCustomer converts echo context to params.
func (w *ServerInterfaceWrapper) GetCustomer(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetCustomer(ctx)
	return err
}

// CreateDossier converts echo context to params.
func (w *ServerInterfaceWrapper) CreateDossier(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateDossier(ctx)
	return err
}

// GetDossier converts echo context to params.
func (w *ServerInterfaceWrapper) GetDossier(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "patientID" -------------
	var patientID string

	err = runtime.BindStyledParameterWithOptions("simple", "patientID", ctx.Param("patientID"), &patientID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter patientID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetDossier(ctx, patientID)
	return err
}

// CreateEpisode converts echo context to params.
func (w *ServerInterfaceWrapper) CreateEpisode(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateEpisode(ctx)
	return err
}

// GetEpisode converts echo context to params.
func (w *ServerInterfaceWrapper) GetEpisode(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "episodeID" -------------
	var episodeID string

	err = runtime.BindStyledParameterWithOptions("simple", "episodeID", ctx.Param("episodeID"), &episodeID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter episodeID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetEpisode(ctx, episodeID)
	return err
}

// GetCollaboration converts echo context to params.
func (w *ServerInterfaceWrapper) GetCollaboration(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "episodeID" -------------
	var episodeID string

	err = runtime.BindStyledParameterWithOptions("simple", "episodeID", ctx.Param("episodeID"), &episodeID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter episodeID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetCollaboration(ctx, episodeID)
	return err
}

// CreateCollaboration converts echo context to params.
func (w *ServerInterfaceWrapper) CreateCollaboration(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "episodeID" -------------
	var episodeID string

	err = runtime.BindStyledParameterWithOptions("simple", "episodeID", ctx.Param("episodeID"), &episodeID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter episodeID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateCollaboration(ctx, episodeID)
	return err
}

// GetInbox converts echo context to params.
func (w *ServerInterfaceWrapper) GetInbox(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetInbox(ctx)
	return err
}

// GetInboxInfo converts echo context to params.
func (w *ServerInterfaceWrapper) GetInboxInfo(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetInboxInfo(ctx)
	return err
}

// SearchOrganizations converts echo context to params.
func (w *ServerInterfaceWrapper) SearchOrganizations(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

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

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.SearchOrganizations(ctx, params)
	return err
}

// GetPatient converts echo context to params.
func (w *ServerInterfaceWrapper) GetPatient(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "patientID" -------------
	var patientID string

	err = runtime.BindStyledParameterWithOptions("simple", "patientID", ctx.Param("patientID"), &patientID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter patientID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetPatient(ctx, patientID)
	return err
}

// UpdatePatient converts echo context to params.
func (w *ServerInterfaceWrapper) UpdatePatient(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "patientID" -------------
	var patientID string

	err = runtime.BindStyledParameterWithOptions("simple", "patientID", ctx.Param("patientID"), &patientID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter patientID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.UpdatePatient(ctx, patientID)
	return err
}

// GetPatients converts echo context to params.
func (w *ServerInterfaceWrapper) GetPatients(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetPatientsParams
	// ------------- Optional query parameter "name" -------------

	err = runtime.BindQueryParameter("form", true, false, "name", ctx.QueryParams(), &params.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter name: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetPatients(ctx, params)
	return err
}

// NewPatient converts echo context to params.
func (w *ServerInterfaceWrapper) NewPatient(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.NewPatient(ctx)
	return err
}

// GetReports converts echo context to params.
func (w *ServerInterfaceWrapper) GetReports(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "patientID" -------------
	var patientID string

	err = runtime.BindStyledParameterWithOptions("simple", "patientID", ctx.Param("patientID"), &patientID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter patientID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetReportsParams
	// ------------- Optional query parameter "episodeID" -------------

	err = runtime.BindQueryParameter("form", true, false, "episodeID", ctx.QueryParams(), &params.EpisodeID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter episodeID: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetReports(ctx, patientID, params)
	return err
}

// CreateReport converts echo context to params.
func (w *ServerInterfaceWrapper) CreateReport(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "patientID" -------------
	var patientID string

	err = runtime.BindStyledParameterWithOptions("simple", "patientID", ctx.Param("patientID"), &patientID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter patientID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateReport(ctx, patientID)
	return err
}

// GetPatientTransfers converts echo context to params.
func (w *ServerInterfaceWrapper) GetPatientTransfers(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetPatientTransfersParams
	// ------------- Required query parameter "patientID" -------------

	err = runtime.BindQueryParameter("form", true, true, "patientID", ctx.QueryParams(), &params.PatientID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter patientID: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetPatientTransfers(ctx, params)
	return err
}

// CreateTransfer converts echo context to params.
func (w *ServerInterfaceWrapper) CreateTransfer(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateTransfer(ctx)
	return err
}

// GetTransferRequest converts echo context to params.
func (w *ServerInterfaceWrapper) GetTransferRequest(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "requestorDID" -------------
	var requestorDID string

	err = runtime.BindStyledParameterWithOptions("simple", "requestorDID", ctx.Param("requestorDID"), &requestorDID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter requestorDID: %s", err))
	}

	// ------------- Path parameter "fhirTaskID" -------------
	var fhirTaskID string

	err = runtime.BindStyledParameterWithOptions("simple", "fhirTaskID", ctx.Param("fhirTaskID"), &fhirTaskID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter fhirTaskID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetTransferRequest(ctx, requestorDID, fhirTaskID)
	return err
}

// ChangeTransferRequestState converts echo context to params.
func (w *ServerInterfaceWrapper) ChangeTransferRequestState(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "requestorDID" -------------
	var requestorDID string

	err = runtime.BindStyledParameterWithOptions("simple", "requestorDID", ctx.Param("requestorDID"), &requestorDID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter requestorDID: %s", err))
	}

	// ------------- Path parameter "fhirTaskID" -------------
	var fhirTaskID string

	err = runtime.BindStyledParameterWithOptions("simple", "fhirTaskID", ctx.Param("fhirTaskID"), &fhirTaskID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter fhirTaskID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ChangeTransferRequestState(ctx, requestorDID, fhirTaskID)
	return err
}

// CancelTransfer converts echo context to params.
func (w *ServerInterfaceWrapper) CancelTransfer(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "transferID" -------------
	var transferID string

	err = runtime.BindStyledParameterWithOptions("simple", "transferID", ctx.Param("transferID"), &transferID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter transferID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CancelTransfer(ctx, transferID)
	return err
}

// GetTransfer converts echo context to params.
func (w *ServerInterfaceWrapper) GetTransfer(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "transferID" -------------
	var transferID string

	err = runtime.BindStyledParameterWithOptions("simple", "transferID", ctx.Param("transferID"), &transferID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter transferID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetTransfer(ctx, transferID)
	return err
}

// UpdateTransfer converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateTransfer(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "transferID" -------------
	var transferID string

	err = runtime.BindStyledParameterWithOptions("simple", "transferID", ctx.Param("transferID"), &transferID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter transferID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.UpdateTransfer(ctx, transferID)
	return err
}

// AssignTransferDirect converts echo context to params.
func (w *ServerInterfaceWrapper) AssignTransferDirect(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "transferID" -------------
	var transferID string

	err = runtime.BindStyledParameterWithOptions("simple", "transferID", ctx.Param("transferID"), &transferID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter transferID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.AssignTransferDirect(ctx, transferID)
	return err
}

// ListTransferNegotiations converts echo context to params.
func (w *ServerInterfaceWrapper) ListTransferNegotiations(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "transferID" -------------
	var transferID string

	err = runtime.BindStyledParameterWithOptions("simple", "transferID", ctx.Param("transferID"), &transferID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter transferID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListTransferNegotiations(ctx, transferID)
	return err
}

// StartTransferNegotiation converts echo context to params.
func (w *ServerInterfaceWrapper) StartTransferNegotiation(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "transferID" -------------
	var transferID string

	err = runtime.BindStyledParameterWithOptions("simple", "transferID", ctx.Param("transferID"), &transferID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter transferID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.StartTransferNegotiation(ctx, transferID)
	return err
}

// UpdateTransferNegotiationStatus converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateTransferNegotiationStatus(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "transferID" -------------
	var transferID string

	err = runtime.BindStyledParameterWithOptions("simple", "transferID", ctx.Param("transferID"), &transferID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter transferID: %s", err))
	}

	// ------------- Path parameter "negotiationID" -------------
	var negotiationID string

	err = runtime.BindStyledParameterWithOptions("simple", "negotiationID", ctx.Param("negotiationID"), &negotiationID, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter negotiationID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.UpdateTransferNegotiationStatus(ctx, transferID, negotiationID)
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
	router.POST(baseURL+"/auth/dummy", wrapper.AuthenticateWithDummy)
	router.GET(baseURL+"/auth/dummy/session/:sessionToken/result", wrapper.GetDummyAuthenticationResult)
	router.POST(baseURL+"/auth/employeeid/session", wrapper.AuthenticateWithEmployeeID)
	router.GET(baseURL+"/auth/employeeid/session/:sessionToken/result", wrapper.GetEmployeeIDAuthenticationResult)
	router.POST(baseURL+"/auth/irma/session", wrapper.AuthenticateWithIRMA)
	router.GET(baseURL+"/auth/irma/session/:sessionToken/result", wrapper.GetIRMAAuthenticationResult)
	router.POST(baseURL+"/auth/openid4vp", wrapper.CreateAuthorizationRequest)
	router.GET(baseURL+"/auth/openid4vp/:token", wrapper.GetOpenID4VPAuthenticationResult)
	router.POST(baseURL+"/auth/passwd", wrapper.AuthenticateWithPassword)
	router.GET(baseURL+"/customers", wrapper.ListCustomers)
	router.POST(baseURL+"/external/transfer/notify/:taskID", wrapper.NotifyTransferUpdate)
	router.PUT(baseURL+"/internal/customer/:customerID/task/:taskID", wrapper.TaskUpdate)
	router.GET(baseURL+"/private", wrapper.CheckSession)
	router.GET(baseURL+"/private/customer", wrapper.GetCustomer)
	router.POST(baseURL+"/private/dossier", wrapper.CreateDossier)
	router.GET(baseURL+"/private/dossier/:patientID", wrapper.GetDossier)
	router.POST(baseURL+"/private/episode", wrapper.CreateEpisode)
	router.GET(baseURL+"/private/episode/:episodeID", wrapper.GetEpisode)
	router.GET(baseURL+"/private/episode/:episodeID/collaboration", wrapper.GetCollaboration)
	router.POST(baseURL+"/private/episode/:episodeID/collaboration", wrapper.CreateCollaboration)
	router.GET(baseURL+"/private/network/inbox", wrapper.GetInbox)
	router.GET(baseURL+"/private/network/inbox/info", wrapper.GetInboxInfo)
	router.GET(baseURL+"/private/network/organizations", wrapper.SearchOrganizations)
	router.GET(baseURL+"/private/patient/:patientID", wrapper.GetPatient)
	router.PUT(baseURL+"/private/patient/:patientID", wrapper.UpdatePatient)
	router.GET(baseURL+"/private/patients", wrapper.GetPatients)
	router.POST(baseURL+"/private/patients", wrapper.NewPatient)
	router.GET(baseURL+"/private/reports/:patientID", wrapper.GetReports)
	router.POST(baseURL+"/private/reports/:patientID", wrapper.CreateReport)
	router.GET(baseURL+"/private/transfer", wrapper.GetPatientTransfers)
	router.POST(baseURL+"/private/transfer", wrapper.CreateTransfer)
	router.GET(baseURL+"/private/transfer-request/:requestorDID/:fhirTaskID", wrapper.GetTransferRequest)
	router.POST(baseURL+"/private/transfer-request/:requestorDID/:fhirTaskID", wrapper.ChangeTransferRequestState)
	router.DELETE(baseURL+"/private/transfer/:transferID", wrapper.CancelTransfer)
	router.GET(baseURL+"/private/transfer/:transferID", wrapper.GetTransfer)
	router.PUT(baseURL+"/private/transfer/:transferID", wrapper.UpdateTransfer)
	router.PUT(baseURL+"/private/transfer/:transferID/assign", wrapper.AssignTransferDirect)
	router.GET(baseURL+"/private/transfer/:transferID/negotiation", wrapper.ListTransferNegotiations)
	router.POST(baseURL+"/private/transfer/:transferID/negotiation", wrapper.StartTransferNegotiation)
	router.PUT(baseURL+"/private/transfer/:transferID/negotiation/:negotiationID", wrapper.UpdateTransferNegotiationStatus)

}
