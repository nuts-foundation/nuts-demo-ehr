package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/dossier"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/inbox"
	nutsClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/registry"

	"github.com/lestrrat-go/jwx/jwt"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer"

	"github.com/labstack/echo/v4"
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
	APIAuth              *Auth
	NutsAuth             nutsClient.Auth
	CustomerRepository   customers.Repository
	PatientRepository    patients.Repository
	DossierRepository    dossier.Repository
	TransferRepository   transfer.Repository
	OrganizationRegistry registry.OrganizationRegistry
	TransferService      transfer.Service
	Inbox                *inbox.Service
	TenantInitializer    func(tenant int) error
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

	token, err := w.APIAuth.CreateCustomerJWT(customer.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := w.TenantInitializer(customer.Id); err != nil {
		return fmt.Errorf("unable to initialize tenant: %w", err)
	}

	return ctx.JSON(200, domain.SessionToken{Token: string(token)})
}

func (w Wrapper) AuthenticateWithPassword(ctx echo.Context) error {
	req := domain.PasswordAuthenticateRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{err})
	}
	sessionId, err := w.APIAuth.AuthenticatePassword(req.CustomerID, req.Password)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, errorResponse{err})
	}
	token, err := w.APIAuth.CreateSessionJWT(req.CustomerID, sessionId, false)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(200, domain.SessionToken{Token: string(token)})
}

func (w Wrapper) getCustomerIDFromHeader(ctx echo.Context) (int, error) {
	token, err := w.APIAuth.extractJWTFromHeader(ctx)
	if err != nil {
		ctx.Echo().Logger.Error(err)
		return 0, echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	rawID, ok := token.Get(CustomerID)
	if !ok {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "missing customerID in token")
	}
	return int(rawID.(float64)), nil
}

func (w Wrapper) AuthenticateWithIRMA(ctx echo.Context) error {
	customerID, err := w.getCustomerIDFromHeader(ctx)
	if err != nil {
		return err
	}
	customer, err := w.CustomerRepository.FindByID(customerID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{err})
	}

	// forward to node
	bytes, err := w.NutsAuth.CreateIrmaSession(*customer)
	if err != nil {
		return err
	}

	// convert to map so echo rendering doesn't escape double quotes
	j := map[string]interface{}{}
	json.Unmarshal(bytes, &j)
	return ctx.JSON(http.StatusOK, j)
}

func (w Wrapper) GetIRMAAuthenticationResult(ctx echo.Context, sessionToken string) error {
	customerID, err := w.getCustomerIDFromHeader(ctx)
	if err != nil {
		return err
	}

	// forward to node
	bytes, err := w.NutsAuth.GetIrmaSessionResult(sessionToken)
	if err != nil {
		return err
	}

	base64String := base64.StdEncoding.EncodeToString(bytes)
	sessionID := w.APIAuth.StoreVP(customerID, base64String)

	newToken, err := w.APIAuth.CreateSessionJWT(customerID, sessionID, true)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(200, domain.SessionToken{Token: string(newToken)})
}

func (w Wrapper) GetCustomer(ctx echo.Context) error {
	customerID := ctx.Get(CustomerID)

	customer, err := w.CustomerRepository.FindByID(customerID.(int))
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

// customerIDFromToken gets the customerID from the jwt
// given some libs it can happen the customerID is returned as a float64...
func customerIDFromToken(token jwt.Token) (int, bool) {
	rawCustomerID, ok := token.Get(CustomerID)
	if !ok {
		return 0, false
	}
	switch customerID := rawCustomerID.(type) {
	case float64:
		return int(customerID), true
	case int:
		return customerID, true
	case string:
		i, err := strconv.Atoi(customerID)
		if err != nil {
			return 0, false
		}
		return i, true
	}

	return 0, false
}
