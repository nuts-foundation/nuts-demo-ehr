package api

import (
	"encoding/json"
	"fmt"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"net/http"
	"strconv"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/dossier"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/episode"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/notification"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/patients"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/reports"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer/receiver"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer/sender"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"

	"github.com/lestrrat-go/jwx/jwt"
	nutsClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client"
	nutsAuth "github.com/nuts-foundation/nuts-demo-ehr/nuts/client/auth"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/registry"

	"github.com/labstack/echo/v4"
)

const BearerAuthScopes = types.BearerAuthScopes

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
	ReportRepository        reports.Repository
	DossierRepository       dossier.Repository
	OrganizationRegistry    registry.OrganizationRegistry
	TransferSenderRepo      sender.TransferRepository
	TransferSenderService   sender.TransferService
	TransferReceiverService receiver.TransferService
	TransferReceiverRepo    receiver.TransferRepository
	FHIRService             fhir.Service
	EpisodeService          episode.Service
	NotificationHandler     notification.Handler
	TenantInitializer       func(tenant int) error
}

func (w Wrapper) CheckSession(ctx echo.Context) error {
	// If this function is reached, it means the session is still valid
	return ctx.NoContent(http.StatusNoContent)
}

func (w Wrapper) SetCustomer(ctx echo.Context) error {
	customer := types.Customer{}
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

	return ctx.JSON(200, types.SessionToken{Token: string(token)})
}

func (w Wrapper) AuthenticateWithPassword(ctx echo.Context) error {
	req := types.PasswordAuthenticateRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{err})
	}

	customer, err := w.CustomerRepository.FindByID(req.CustomerID)
	if err != nil {
		return err
	}

	sessionId, userInfo, err := w.APIAuth.AuthenticatePassword(req.CustomerID, req.Password)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, errorResponse{err})
	}

	token, err := w.APIAuth.CreateSessionJWT(customer.Name, userInfo.Identifier, req.CustomerID, sessionId, false)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(200, types.SessionToken{Token: string(token)})
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
	bytes, err := w.NutsAuth.CreateIrmaSession(*customer.Did)
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

	authSessionID, err := w.getSessionID(ctx)
	if err != nil {
		// No current session, create a new one. Introspect IRMA VP and extract properties for UserInfo.
		userPresentation, err := w.NutsAuth.VerifyPresentation(*sessionStatus.VerifiablePresentation)
		if err != nil {
			return fmt.Errorf("unable to verify presentation: %w", err)
		}
		attrs := *userPresentation.IssuerAttributes
		userInfo := UserInfo{
			Identifier: fmt.Sprintf("%v", attrs["sidn-pbdf.email.email"]),
			Initials:   fmt.Sprintf("%v", attrs["gemeente.personalData.initials"]),
			FamilyName: fmt.Sprintf("%v", attrs["gemeente.personalData.familyname"]),
		}
		authSessionID = w.APIAuth.createSession(customerID, userInfo)
	}

	err = w.APIAuth.Elevate(authSessionID, *sessionStatus.VerifiablePresentation)
	if err != nil {
		return fmt.Errorf("unable to elevate session: %w", err)
	}

	session := w.APIAuth.GetSession(authSessionID)

	customer, err := w.CustomerRepository.FindByID(customerID)
	if err != nil {
		return err
	}

	newToken, err := w.APIAuth.CreateSessionJWT(customer.Name, session.UserInfo.Identifier, customerID, authSessionID, true)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(200, types.SessionToken{Token: string(newToken)})
}

func (w Wrapper) AuthenticateWithEmployeeID(ctx echo.Context) error {
	// The method is called "Authenticate" but it is actually elevation,
	// since it requires an existing session from with employee info.
	var sessionID string
	if sid, ok := ctx.Get(SessionID).(string); ok {
		sessionID = sid
	} else {
		return echo.NewHTTPError(http.StatusUnauthorized, "existing session is required for EmployeeID means (missing token)")
	}

	session := w.APIAuth.GetSession(sessionID)
	if session == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "existing session is required for EmployeeID means (unknown session)")
	}

	customer, _ := w.CustomerRepository.FindByID(session.CustomerID)
	if customer == nil || customer.Did == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "customer with DID required for EmployeeID means")
	}

	params := map[string]interface{}{
		"employer": *customer.Did,
		"employee": map[string]interface{}{
			"identifier": session.UserInfo.Identifier,
			"roleName":   session.UserInfo.RoleName,
			"initials":   session.UserInfo.Initials,
			"familyName": session.UserInfo.FamilyName,
		},
	}
	bytes, err := w.NutsAuth.CreateEmployeeIDSession(params)
	if err != nil {
		return err
	}

	// convert to map so echo rendering doesn't escape double quotes
	j := map[string]interface{}{}

	if err := json.Unmarshal(bytes, &j); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, j)
}

func (w Wrapper) GetEmployeeIDAuthenticationResult(ctx echo.Context, sessionToken string) error {
	authSession, err := w.getSession(ctx)
	if err != nil {
		return err
	}
	authSessionID, _ := w.getSessionID(ctx) // can't fail

	// forward to node
	sessionStatus, err := w.NutsAuth.GetEmployeeIDSessionResult(sessionToken)
	if err != nil {
		return err
	}

	if sessionStatus.Status != "completed" {
		return echo.NewHTTPError(http.StatusNotFound, sessionStatus.Status)
	}

	err = w.APIAuth.Elevate(authSessionID, *sessionStatus.VerifiablePresentation)
	if err != nil {
		return fmt.Errorf("unable to elevate session: %w", err)
	}

	customer, err := w.CustomerRepository.FindByID(authSession.CustomerID)
	if err != nil {
		return err
	}

	newToken, err := w.APIAuth.CreateSessionJWT(customer.Name, authSession.UserInfo.Identifier, authSession.CustomerID, authSessionID, true)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(200, types.SessionToken{Token: string(newToken)})
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

	bytes, err := w.NutsAuth.CreateDummySession(*customer.Did)
	if err != nil {
		return err
	}

	// convert to map so echo rendering doesn't escape double quotes
	j := map[string]interface{}{}

	if err := json.Unmarshal(bytes, &j); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, j)
}

func (w Wrapper) GetDummyAuthenticationResult(ctx echo.Context, sessionToken string) error {
	authSession, err := w.getSession(ctx)
	if err != nil {
		return err
	}
	authSessionID, _ := w.getSessionID(ctx) // can't fail

	customerID, err := w.APIAuth.GetCustomerIDFromHeader(ctx)
	if err != nil {
		return err
	}

	var sessionResult *nutsAuth.SignSessionStatusResponse

	// for dummy, it takes a few request to get to status completed.
	for i := 0; i < 4; i++ {
		// forward to node
		sessionResult, err = w.NutsAuth.GetDummySessionResult(sessionToken)
		if err != nil {
			return err
		}

		if sessionResult.Status == "completed" {
			break
		}
	}

	if sessionResult.Status != "completed" {
		return echo.NewHTTPError(http.StatusNotFound, "signing session not completed")
	}

	err = w.APIAuth.Elevate(authSessionID, *sessionResult.VerifiablePresentation)
	if err != nil {
		return fmt.Errorf("failed to elevate session: %w", err)
	}

	customer, err := w.CustomerRepository.FindByID(customerID)
	if err != nil {
		return err
	}

	newToken, err := w.APIAuth.CreateSessionJWT(customer.Name, authSession.UserInfo.Identifier, customerID, authSessionID, true)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(200, types.SessionToken{Token: string(newToken)})
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
