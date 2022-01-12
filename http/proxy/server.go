package proxy

import (
	"context"
	"errors"
	"fmt"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir/zorginzage"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/nuts-foundation/go-did/vc"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer"
	"github.com/nuts-foundation/nuts-demo-ehr/http/auth"
	nutsAuthClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client/auth"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/registry"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
	"github.com/nuts-foundation/nuts-node/vcr/credential"
	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
)

var fhirServerTenant = struct{}{}

type Server struct {
	proxy               *httputil.ReverseProxy
	auth                auth.Service
	path                string
	customerRepository  customers.Repository
	vcRegistry          registry.VerifiableCredentialRegistry
	multiTenancyEnabled bool
}

func NewServer(authService auth.Service, customerRepository customers.Repository, vcRegistry registry.VerifiableCredentialRegistry, targetURL url.URL, path string, multiTenancyEnabled bool) *Server {
	server := &Server{
		path:                path,
		auth:                authService,
		customerRepository:  customerRepository,
		vcRegistry:          vcRegistry,
		multiTenancyEnabled: multiTenancyEnabled,
	}

	server.proxy = &httputil.ReverseProxy{
		// Does not support query parameters in targetURL
		Director: func(req *http.Request) {
			requestURL := &url.URL{}
			*requestURL = *req.URL
			requestURL.Scheme = targetURL.Scheme
			requestURL.Host = targetURL.Host
			requestURL.RawPath = "" // Not required?

			if server.multiTenancyEnabled {
				tenant := req.Context().Value(fhirServerTenant).(int) // this shouldn't/can't fail, because the middleware handler should've set it.
				requestURL.Path = fmt.Sprintf("%s/%d%s", targetURL.Path, tenant, req.URL.Path[len(path):])
			} else {
				requestURL.Path = targetURL.Path + req.URL.Path[len(path):]
			}

			req.URL = requestURL
			req.Host = requestURL.Host

			logrus.Debugf("Rewritten to: %s", req.URL.Path)

			if _, ok := req.Header["User-Agent"]; !ok {
				// explicitly disable User-Agent so it's not set to default value
				req.Header.Set("User-Agent", "")
			}
		},
	}

	return server
}

func (server *Server) AuthMiddleware() echo.MiddlewareFunc {
	config := auth.Config{
		Skipper: server.skipper,
		ErrorF:  errorFunc,
		AccessF: server.verifyAccess,
	}

	return auth.SecurityFilter{Auth: server.auth}.AuthWithConfig(config)
}

func (server *Server) skipper(ctx echo.Context) bool {
	requestURI := ctx.Request().RequestURI
	// everything with /web is handled
	return !strings.HasPrefix(requestURI, server.path)
}

func errorFunc(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusUnauthorized, NewOperationOutcome(err, "access denied", CodeSecurity, SeverityError))
}

func (server *Server) Handler(other echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Task update
		if intervalValue := c.Get("internal"); intervalValue != nil {
			if value, ok := intervalValue.(bool); ok && value {
				c.Logger().Debugf("routing internally to %s", c.Request().URL.Path)

				return other(c)
			}
		}

		// non task FHIR resources
		c.Logger().Debugf("FHIR Proxy: proxying %s %s", c.Request().Method, c.Request().RequestURI)
		accessToken := c.Get(auth.AccessToken).(nutsAuthClient.TokenIntrospectionResponse)

		if server.multiTenancyEnabled {
			// Enrich request with resource owner's FHIR server tenant, which is the customer's ID
			tenant, err := server.getTenant(*accessToken.Iss)
			if err != nil {
				return c.JSON(http.StatusBadRequest, NewOperationOutcome(err, err.Error(), CodeSecurity, SeverityError))
			}

			c.SetRequest(c.Request().WithContext(context.WithValue(
				c.Request().Context(),
				fhirServerTenant,
				tenant,
			)))
		}

		// proxy handling
		server.proxy.ServeHTTP(c.Response(), c.Request())

		return nil
	}
}

// verifyAccess checks the access policy rules. The token has already been checked and the introspected token is used.
func (server *Server) verifyAccess(ctx echo.Context, request *http.Request, token *nutsAuthClient.TokenIntrospectionResponse) error {
	route := parseRoute(request)

	// check purposeOfUse/service according to §6.2 eOverdracht-sender policy
	service := token.Service
	if service == nil {
		return errors.New("access-token doesn't contain 'service' claim")
	}

	switch *service {
	case zorginzage.ServiceName:
		subjects, err := server.parseNutsAuthorizationCredentials(request.Context(), token)
		if err != nil {
			return err
		}

		// observation specific access
		observationPath := server.path + "/Observation"

		if route.path() == observationPath {
			if route.operation != "read" {
				return fmt.Errorf("incorrect operation %s on: %s, must be read", route.operation, observationPath)
			}

			var episodeOfCareID string

			for _, subject := range subjects {
				for _, resource := range subject.Resources {
					if strings.HasPrefix(resource.Path, "/EpisodeOfCare/") {
						episodeOfCareID = resource.Path[len("/EpisodeOfCare/"):]
						break
					}
				}
			}

			if episodeOfCareID == "" {
				return fmt.Errorf("unable to find context for route: %s", route.path())
			}

			if route.url.Query().Get("context") != fmt.Sprintf("EpisodeOfCare/%s", episodeOfCareID) {
				return fmt.Errorf("access denied for episode %s in route: %s", episodeOfCareID, route.path())
			}

			return nil
		}

		if err := server.validateWithNutsAuthorizationCredential(token, subjects, *route); err != nil {
			return fmt.Errorf("access denied for %s on %s: %w", route.operation, route.path(), err)
		}
	case transfer.SenderServiceName:
		// task specific access
		serverTaskPath := server.path + "/Task"

		if route.path() == serverTaskPath {
			// §6.2.1.1 retrieving tasks via a search operation
			// validate GET [base]/Task?code=http://snomed.info/sct|308292007&_lastUpdated=[time of last request]
			if route.operation != "read" {
				return fmt.Errorf("incorrect operation %s on: %s, must be read", route.operation, serverTaskPath)
			}

			// query params(code and _lastUpdated) are optional
			// ok for Task search
			return nil
		}

		subjects, err := server.parseNutsAuthorizationCredentials(request.Context(), token)
		if err != nil {
			return err
		}

		// §6.2.1.2 Updating the Task
		// and
		// §6.2.2 other resources that require a credential and a user contract
		// the existence of the user contract is validated by validateWithNutsAuthorizationCredential
		if err := server.validateWithNutsAuthorizationCredential(token, subjects, *route); err != nil {
			return fmt.Errorf("access denied for %s on %s: %w", route.operation, route.path(), err)
		}

		// Task updates must be routed internally
		if route.operation == "update" && strings.HasPrefix(route.path(), serverTaskPath) {
			tenant, err := server.getTenant(*token.Iss)
			if err != nil {
				return fmt.Errorf("access denied for %s on %s, tenant %s: %w", route.operation, route.path(), *token.Iss, err)
			}

			// task handling
			req := ctx.Request()
			path := fmt.Sprintf("/web/internal/customer/%d/task/%s", tenant, route.resourceID)
			req.URL.Path = path
			req.URL.RawPath = path
			req.RequestURI = path
			ctx.SetRequest(req)

			ctx.Set("internal", true)
		}
	default:
		return fmt.Errorf("access-token contains incorrect 'service' claim: %s, must be %s", *service, transfer.SenderServiceName)
	}

	return nil
}

func (server *Server) parseNutsAuthorizationCredentials(ctx context.Context, token *nutsAuthClient.TokenIntrospectionResponse) ([]credential.NutsAuthorizationCredentialSubject, error) {
	var subjects []credential.NutsAuthorizationCredentialSubject

	if token.Vcs == nil {
		return subjects, nil
	}

	for _, credentialID := range *token.Vcs {
		// resolve credential. NutsAuthCredential must be resolved with the untrusted flag
		authCredential, err := server.vcRegistry.ResolveVerifiableCredential(ctx, credentialID)
		if err != nil {
			return nil, fmt.Errorf("invalid credential: %w", err)
		}

		if !validCredentialType(*authCredential) {
			continue
		}

		subject := make([]credential.NutsAuthorizationCredentialSubject, 0)

		if err := authCredential.UnmarshalCredentialSubject(&subject); err != nil {
			return nil, fmt.Errorf("invalid content for NutsAuthorizationCredential credentialSubject: %w", err)
		}

		subjects = append(subjects, subject...)
	}

	return subjects, nil
}

func (server *Server) validateWithNutsAuthorizationCredential(token *nutsAuthClient.TokenIntrospectionResponse, subjects []credential.NutsAuthorizationCredentialSubject, route fhirRoute) error {
	hasUser := token.Email != nil

	if token.Vcs != nil {
		for _, subject := range subjects {
			for _, resource := range subject.Resources {
				// path should match
				if !strings.Contains(route.path(), resource.Path) {
					continue
				}

				// usi must be present when resource requires user context
				if resource.UserContext && !hasUser {
					continue
				}

				// operation must match
				for _, operation := range resource.Operations {
					if operation == route.operation {
						// all is ok, no need to continue after a match
						return nil
					}
				}
			}
		}

		return errors.New("no matching NutsAuthorizationCredential found in access-token")
	}

	return errors.New("no NutsAuthorizationCredential in access-token")
}

func (server *Server) getTenant(requesterDID string) (int, error) {
	customer, err := server.customerRepository.FindByDID(requesterDID)
	if err != nil {
		return 0, err
	}
	if customer == nil {
		return 0, errors.New("unknown tenant")
	}
	return customer.Id, nil
}

func validCredentialType(verifiableCredential vc.VerifiableCredential) bool {
	return verifiableCredential.IsType(*credential.NutsAuthorizationCredentialTypeURI)
}
