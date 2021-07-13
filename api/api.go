package api

import (
	"encoding/base64"
	"encoding/json"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/dossier"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/client"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/patients"
)

const BearerAuthScopes = domain.BearerAuthScopes

type errorResponse struct {
	Error error
}

func (e errorResponse) MarshalJSON() ([]byte, error) {
	asMap := make(map[string]string, 1)
	asMap["error"] = e.Error.Error()
	return json.Marshal(asMap)
}

type Wrapper struct {
	Auth               *Auth
	Client             client.HTTPClient
	CustomerRepository customers.Repository
	PatientRepository  patients.Repository
	DossierRepository  dossier.Repository
	TransferRepository transfer.Repository
	FHIRGateway        fhir.Gateway
}

func (w Wrapper) CheckSession(ctx echo.Context) error {
	// If this function is reached, it means the session is still valid
	return ctx.NoContent(http.StatusNoContent)
}

func (w Wrapper) SetCustomer(ctx echo.Context) error {
	customer := domain.Customer{}
	if err := ctx.Bind(&customer); err != nil {
		return err
	}

	token, err := w.Auth.CreateCustomerJWT(customer.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(200, domain.SessionToken{Token: string(token)})
}

func (w Wrapper) AuthenticateWithPassword(ctx echo.Context) error {
	req := domain.PasswordAuthenticateRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{err})
	}
	sessionId, err := w.Auth.AuthenticatePassword(req.CustomerID, req.Password)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, errorResponse{err})
	}
	token, err := w.Auth.CreateSessionJWT(req.CustomerID, sessionId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(200, domain.SessionToken{Token: string(token)})
}

func (w Wrapper) AuthenticateWithIRMA(ctx echo.Context) error {
	req := domain.IRMAAuthenticationRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{err})
	}

	customer, err := w.CustomerRepository.FindByID(req.CustomerID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{err})
	}

	// forward to node
	bytes, err := w.Client.CreateIrmaSession(*customer)
	if err != nil {
		return err
	}

	// convert to map so echo rendering doesn't escape double quotes
	j := map[string]interface{}{}
	json.Unmarshal(bytes, &j)
	return ctx.JSON(http.StatusOK, j)
}

func (w Wrapper) GetIRMAAuthenticationResult(ctx echo.Context, sessionToken string) error {
	// current customerID
	token, err := w.Auth.extractJWTFromHeader(ctx)
	if err != nil {
		ctx.Echo().Logger.Error(err)
		return ctx.NoContent(http.StatusUnauthorized)
	}
	customerID, ok := token.Get(CustomerID)
	if ok {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	// forward to node
	bytes, err := w.Client.GetIrmaSessionResult(sessionToken)
	if err != nil {
		return err
	}

	base64String := base64.StdEncoding.EncodeToString(bytes)
	sessionID := w.Auth.StoreVP(customerID.(string), base64String)

	newToken, err := w.Auth.CreateSessionJWT(customerID.(string), sessionID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(200, domain.SessionToken{Token: string(newToken)})
}

func (w Wrapper) GetCustomer(ctx echo.Context) error {
	customerID := ctx.Get(CustomerID)

	customer, err := w.CustomerRepository.FindByID(customerID.(string))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{err})
	}
	return ctx.JSON(http.StatusOK, customer)
}

func (w Wrapper) ListCustomers(ctx echo.Context) error {
	customers, err := w.CustomerRepository.All()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{err})
	}
	return ctx.JSON(http.StatusOK, customers)
}
