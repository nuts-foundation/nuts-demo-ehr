package proxy

import (
	"context"
	"errors"
	"fmt"
	"github.com/nuts-foundation/go-did/vc"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	client "github.com/nuts-foundation/nuts-demo-ehr/client/auth"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/auth"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
	http2 "github.com/nuts-foundation/nuts-demo-ehr/http"
	"github.com/nuts-foundation/nuts-node/vcr/credential"
	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
)

var fhirServerTenant = struct{}{}

type Server struct {
	proxy              *httputil.ReverseProxy
	auth               auth.Service
	path               string
	customerRepository customers.Repository
}

func NewServer(authService auth.Service, customerRepository customers.Repository, targetURL url.URL, path string) *Server {
	server := &Server{
		path:               path,
		auth:               authService,
		customerRepository: customerRepository,
	}
	server.proxy = &httputil.ReverseProxy{
		// Does not support query parameters in targetURL
		Director: func(req *http.Request) {
			tenant := req.Context().Value(fhirServerTenant).(string) // this shouldn't/can't fail, because the middleware handler should've set it.

			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host
			req.URL.RawPath = "" // Not required?
			req.URL.Path = targetURL.Path + "/" + tenant + req.URL.Path[len(path):]
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
	config := http2.Config{
		Skipper: server.skipper,
		ErrorF:  errorFunc,
		AccessF: server.verifyAccess,
	}
	return http2.SecurityFilter{Auth: server.auth}.AuthWithConfig(config)
}

func (server *Server) skipper(ctx echo.Context) bool {
	requestURI := ctx.Request().RequestURI
	return !strings.HasPrefix(requestURI, server.path)
}

func errorFunc(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusUnauthorized, NewOperationOutcome(err, "access denied", CodeSecurity, SeverityError))
}

func (server *Server) Handler(_ echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Logger().Debugf("FHIR Proxy: proxying %s %s", c.Request().Method, c.Request().RequestURI)
		accessToken := c.Get(http2.AccessToken).(client.TokenIntrospectionResponse)
		// Enrich request with resource owner's FHIR server tenant, which is the customer's ID
		tenant, err := server.getTenant(*accessToken.Iss)
		if err != nil {
			return c.JSON(http.StatusBadRequest, NewOperationOutcome(err, err.Error(), CodeSecurity, SeverityError))
		}
		c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), fhirServerTenant, tenant)))

		server.proxy.ServeHTTP(c.Response().Writer, c.Request())

		return nil
	}
}

func (server *Server) verifyAccess(request *http.Request, token *client.TokenIntrospectionResponse) error {
	route := parseRoute(request)

	// check purposeOfUse/service
	service := token.Service
	if service == nil {
		return errors.New("access-token doesn't contain 'service' claim")
	}
	if *service != "eOverdracht-sender" {
		return fmt.Errorf("access-token contains incorrect 'service' claim: %s, must be eOverdracht-sender", *service)
	}

	// TODO: Assert that token.subject equals the requester?
	// NutsAuthorizationCredential is only required for:
	// 1. Retrieving FHIR resources that contain personal information (for the sake of simplicity; everything other than the Task for now)
	// 2. Updating a task resource (so everything other than a HTTP GET/read)
	if route.operation == "read" && strings.HasPrefix(route.path, server.path+"/Task") {
		return nil
	}

	authCredential, err := findNutsAuthorizationCredential(token)
	if err != nil {
		return err
	}

	subject := &credential.NutsAuthorizationCredentialSubject{}

	if err := authCredential.UnmarshalCredentialSubject(subject); err != nil {
		return fmt.Errorf("invalid type for NutsAuthorizationCredential subject: %w", err)
	}

	allowed := false

	for _, resource := range subject.Resources {
		if route.path != resource.Path {
			continue
		}

		for _, operation := range resource.Operations {
			if operation == route.operation {
				allowed = true
				break
			}
		}
	}

	if !allowed {
		return fmt.Errorf("access denied for path '%s' with operation: %s", route.path, route.operation)
	}

	return nil
}

func findNutsAuthorizationCredential(token *client.TokenIntrospectionResponse) (*vc.VerifiableCredential, error) {
	if token.Vcs != nil {
		for _, verifiableCredential := range *token.Vcs {
			types := credential.ExtractTypes(verifiableCredential)

			for _, typeName := range types {
				if typeName == credential.NutsAuthorizationCredentialType {
					return &verifiableCredential, nil
				}
			}
		}
	}

	return nil, errors.New("NutsAuthorizationCredential was not found in the access-token")
}

func (server *Server) getTenant(requesterDID string) (string, error) {
	customer, err := server.customerRepository.FindByDID(requesterDID)
	if err != nil {
		return "", err
	}
	if customer == nil {
		return "", errors.New("unknown tenant")
	}
	return customer.Id, nil
}
