package proxy

import (
	"errors"
	"fmt"
	"github.com/nuts-foundation/go-did/vc"
	client "github.com/nuts-foundation/nuts-demo-ehr/client/auth"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/auth"
	"github.com/nuts-foundation/nuts-node/vcr/credential"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

type Server struct {
	proxy *httputil.ReverseProxy
	auth  auth.Service
	path  string
}

func NewServer(authService auth.Service, targetURL url.URL, path string) *Server {
	// Does not support query parameters in targetURL
	proxyDirector := func(req *http.Request) {
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.URL.RawPath = "" // Not required?
		req.URL.Path = targetURL.Path + req.URL.Path[len(path):]
		logrus.Infof("Rewritten to: %s", req.URL.Path)
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	return &Server{
		path:  path,
		proxy: &httputil.ReverseProxy{Director: proxyDirector},
		auth:  authService,
	}
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

func (server *Server) parseAccessToken(c echo.Context) (*client.TokenIntrospectionResponse, error) {
	bearerToken, err := server.auth.ParseBearerToken(c.Request())
	if err != nil {
		return nil, fmt.Errorf("failed to parse the bearer token: %w", err)
	}

	token, err := server.auth.IntrospectAccessToken(c.Request().Context(), bearerToken)
	if err != nil {
		return nil, fmt.Errorf("failed to introspect token: %w", err)
	}

	if !token.Active {
		return nil, fmt.Errorf("access-token is not active")
	}

	return token, nil
}

func (server *Server) verifyAccess(route *fhirRoute, token *client.TokenIntrospectionResponse) error {
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

func (server *Server) Handler(_ echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Logger().Debugf("FHIR Proxy: proxying %s %s", c.Request().Method, c.Request().RequestURI)
		token, err := server.parseAccessToken(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, NewOperationOutcome(err, "Invalid access-token", CodeSecurity, SeverityError))
		}

		route := parseRoute(c.Request())

		if err := server.verifyAccess(route, token); err != nil {
			return c.JSON(http.StatusUnauthorized, NewOperationOutcome(err, "Not authorized", CodeSecurity, SeverityError))
		}

		server.proxy.ServeHTTP(c.Response().Writer, c.Request())

		return nil
	}
}
