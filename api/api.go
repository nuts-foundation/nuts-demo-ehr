package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/notification"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer/receiver"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer/sender"

	"github.com/lestrrat-go/jwx/jwt"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/dossier"

	nutsClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/registry"

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
	APIAuth                 *Auth
	NutsAuth                nutsClient.Auth
	CustomerRepository      customers.Repository
	PatientRepository       patients.Repository
	DossierRepository       dossier.Repository
	OrganizationRegistry    registry.OrganizationRegistry
	TransferSenderRepo      sender.TransferRepository
	TransferSenderService   sender.TransferService
	TransferReceiverService receiver.TransferService
	TransferReceiverRepo    receiver.TransferRepository
	NotificationHandler     notification.Handler
	TenantInitializer       func(tenant int) error
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

func (w Wrapper) AuthenticateWithIRMA(ctx echo.Context) error {
	customerID, err := w.APIAuth.GetCustomerIDFromHeader(ctx)
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
	customerID, err := w.APIAuth.GetCustomerIDFromHeader(ctx)
	if err != nil {
		return err
	}

	// forward to node
	sessionStatus, err := w.NutsAuth.GetIrmaSessionResult(sessionToken)
	if err != nil {
		return err
	}

	if sessionStatus.Status != "DONE" {
		return echo.NewHTTPError(http.StatusNotFound, "signing session not completed")
	}

	sessionBytes, _ := json.Marshal(sessionStatus)
	base64String := base64.StdEncoding.EncodeToString(sessionBytes)
	sessionID := w.APIAuth.StoreVP(customerID, base64String)

	newToken, err := w.APIAuth.CreateSessionJWT(customerID, sessionID, true)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(200, domain.SessionToken{Token: string(newToken)})
}

func (w Wrapper) AuthenticateWithDummy(ctx echo.Context) error {
	customerID, err := w.APIAuth.GetCustomerIDFromHeader(ctx)
	if err != nil {
		return err
	}
	customer, err := w.CustomerRepository.FindByID(customerID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{err})
	}
	bytes, err := w.NutsAuth.CreateDummySession(*customer)
	if err != nil {
		return err
	}

	// convert to map so echo rendering doesn't escape double quotes
	j := map[string]interface{}{}
	json.Unmarshal(bytes, &j)
	return ctx.JSON(http.StatusOK, j)
}

func (w Wrapper) GetDummyAuthenticationResult(ctx echo.Context, sessionToken string) error {
	customerID, err := w.APIAuth.GetCustomerIDFromHeader(ctx)
	if err != nil {
		return err
	}

	var (
		sessionStatus string
	)

	// for dummy, it takes a few request to get to status completed.
	for i := 0; i < 4 && sessionStatus != "completed"; i++ {
		// forward to node
		sessionResult, err := w.NutsAuth.GetDummySessionResult(sessionToken)
		if err != nil {
			return err
		}
		sessionStatus = sessionResult.Status
	}

	if sessionStatus != "completed" {
		return echo.NewHTTPError(http.StatusNotFound, "signing session not completed")
	}
	sessionBytes, _ := json.Marshal(sessionStatus)
	base64String := base64.StdEncoding.EncodeToString(sessionBytes)

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
