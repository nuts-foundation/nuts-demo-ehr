package api

import (
	"encoding/json"
	"errors"
	"fmt"
	ssi "github.com/nuts-foundation/go-did"
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

	sessionId, err := w.APIAuth.AuthenticatePassword(req.CustomerID, req.Password)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, errorResponse{err})
	}

	customer, err := w.CustomerRepository.FindByID(req.CustomerID)
	if err != nil {
		return err
	}

	session := w.APIAuth.GetSession(sessionId)
	employeeIdentifier, err := extractEmployeeIdentifier(session.Credential)
	if err != nil {
		return err
	}

	token, err := w.APIAuth.CreateSessionJWT(customer.Name, employeeIdentifier, req.CustomerID, sessionId, false)
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

	sessionID := w.APIAuth.StoreVP(customerID, *sessionStatus.VerifiablePresentation)

	customer, err := w.CustomerRepository.FindByID(customerID)
	if err != nil {
		return err
	}

	employeeIdentifier, err := extractEmployeeIdentifier(*sessionStatus.VerifiablePresentation)
	if err != nil {
		return err
	}

	newToken, err := w.APIAuth.CreateSessionJWT(customer.Name, employeeIdentifier, customerID, sessionID, true)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(200, types.SessionToken{Token: string(newToken)})
}

func (w Wrapper) AuthenticateWithSelfSigned(ctx echo.Context) error {
	// The method is called "Authenticate" but it is actually elevation,
	// since it requires an existing session from with employee info.
	var sessionID string
	if sid, ok := ctx.Get(SessionID).(string); ok {
		sessionID = sid
	} else {
		return echo.NewHTTPError(http.StatusUnauthorized, "existing session is required for self-signed means (missing token)")
	}

	session := w.APIAuth.GetSession(sessionID)
	if session == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "existing session is required for self-signed means (unknown session)")
	}
	params, err := extractEmployeeCredentialSubject(session.Credential)
	if err != nil {
		return err
	}

	bytes, err := w.NutsAuth.CreateSelfSignedSession(params)
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

func (w Wrapper) GetSelfSignedAuthenticationResult(ctx echo.Context, sessionToken string) error {
	customerID, err := w.APIAuth.GetCustomerIDFromHeader(ctx)
	if err != nil {
		return err
	}

	// forward to node
	sessionStatus, err := w.NutsAuth.GetSelfSignedSessionResult(sessionToken)
	if err != nil {
		return err
	}

	if sessionStatus.Status != "completed" {
		return echo.NewHTTPError(http.StatusNotFound, sessionStatus.Status)
	}

	sessionID := w.APIAuth.StoreVP(customerID, *sessionStatus.VerifiablePresentation)

	customer, err := w.CustomerRepository.FindByID(customerID)
	if err != nil {
		return err
	}

	userIdentifier, err := extractEmployeeIdentifier(*sessionStatus.VerifiablePresentation)
	if err != nil {
		return err
	}

	newToken, err := w.APIAuth.CreateSessionJWT(customer.Name, userIdentifier, customerID, sessionID, true)
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

	sessionID := w.APIAuth.StoreVP(customerID, *sessionResult.VerifiablePresentation)

	customer, err := w.CustomerRepository.FindByID(customerID)
	if err != nil {
		return err
	}

	// TODO: Extract user name from Dummy VP
	newToken, err := w.APIAuth.CreateSessionJWT(customer.Name, "(dummy)", customerID, sessionID, true)
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

func extractEmployeeCredentialSubject(presentation nutsAuth.VerifiablePresentation) (map[string]interface{}, error) {
	// TODO: Support IRMA VP
	// EmployeeCredential is created by password authentication
	// NutsEmployeeCredential is created by Nuts SelfSigned auth means
	// They share (almost?) the same credential subject structure (maybe we should make them the same credential)
	for _, credential := range presentation.VerifiableCredential {
		if credential.IsType(ssi.MustParseURI("EmployeeCredential")) ||
			credential.IsType(ssi.MustParseURI("NutsEmployeeCredential")) {
			return credential.CredentialSubject[0].(map[string]interface{}), nil
		}
	}

	return nil, errors.New("no EmployeeCredential found in VP")
}

func extractEmployeeIdentifier(presentation nutsAuth.VerifiablePresentation) (string, error) {
	for _, credential := range presentation.VerifiableCredential {
		if credential.IsType(ssi.MustParseURI("EmployeeCredential")) {
			subject := credential.CredentialSubject[0].(map[string]interface{})
			employee := subject["employee"].(map[string]interface{})
			return employee["identifier"].(string), nil
		}
		if credential.IsType(ssi.MustParseURI("NutsEmployeeCredential")) {
			subject := credential.CredentialSubject[0].(map[string]interface{})
			member := subject["member"].(map[string]interface{})
			return member["identifier"].(string), nil
		}
	}

	return "", errors.New("no EmployeeCredential or NutsEmployeeCredential found in VP")
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
