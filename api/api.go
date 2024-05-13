package api

import (
	"encoding/json"
	"fmt"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/acl"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/sharedcareplan"
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

var _ ServerInterface = (*Wrapper)(nil)

type Wrapper struct {
	APIAuth                 *Auth
	ACL                     *acl.Repository
	NutsClient              *nutsClient.HTTPClient
	CustomerRepository      customers.Repository
	PatientRepository       patients.Repository
	ReportRepository        reports.Repository
	DossierRepository       dossier.Repository
	OrganizationRegistry    registry.OrganizationRegistry
	TransferSenderRepo      sender.TransferRepository
	TransferSenderService   sender.TransferService
	TransferReceiverService receiver.TransferService
	TransferReceiverRepo    receiver.TransferRepository
	ZorginzageService       domain.ZorginzageService
	SharedCarePlanService   *sharedcareplan.Service
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

	token, err := w.APIAuth.CreateSessionJWT(customer.Name, userInfo.Identifier, req.CustomerID, sessionId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(200, types.SessionToken{Token: string(token)})
}

func (w Wrapper) CreateAuthorizationRequest(ctx echo.Context) error {
	customerID, err := w.APIAuth.GetCustomerIDFromHeader(ctx)
	if err != nil {
		return err
	}
	customer, err := w.CustomerRepository.FindByID(customerID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{err})
	}
	response, err := w.NutsClient.CreateAuthenticationRequest(*customer.Did)
	if err != nil {
		return err
	}

	// convert to map so echo rendering doesn't escape double quotes
	return ctx.JSON(http.StatusOK, response)
}

func (w Wrapper) GetOpenID4VPAuthenticationResult(ctx echo.Context, token string) error {
	//customerID, err := w.APIAuth.GetCustomerIDFromHeader(ctx)
	//if err != nil {
	//	return err
	//}
	//customer, err := w.CustomerRepository.FindByID(customerID)
	//if err != nil {
	//	return ctx.JSON(http.StatusInternalServerError, errorResponse{err})
	//}
	response, err := w.NutsClient.GetAuthenticationResult(token)
	if err != nil {
		return err
	}

	// convert to map so echo rendering doesn't escape double quotes
	return ctx.JSON(http.StatusOK, response)
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
