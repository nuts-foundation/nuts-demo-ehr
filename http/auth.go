package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	client "github.com/nuts-foundation/nuts-demo-ehr/client/auth"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/auth"
)

const AccessToken = "accessToken"

type ErrorFunc func(c echo.Context, err error) error
func DefaultErrorFunc(c echo.Context, err error) error {
	return err
}

type AccessFunc func(request *http.Request, token *client.TokenIntrospectionResponse) error
func DefaultAccessFunc(request *http.Request, token *client.TokenIntrospectionResponse) error {
	return nil
}

type Config struct {
	Skipper middleware.Skipper
	ErrorF  ErrorFunc
	AccessF AccessFunc
}

type SecurityFilter struct {
	Auth  auth.Service
}

func (filter SecurityFilter) AuthWithConfig(config Config) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}
	if config.ErrorF == nil {
		config.ErrorF = DefaultErrorFunc
	}
	if config.AccessF == nil {
		config.AccessF = DefaultAccessFunc
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}

			c.Logger().Debugf("Checking access token on %s %s", c.Request().Method, c.Request().RequestURI)
			token, err := filter.parseAccessToken(c)
			if err != nil {
				return config.ErrorF(c, errors.New("invalid access-token"))
			}

			c.Set(AccessToken, *token)

			if err := config.AccessF(c.Request(), token); err != nil {
				return config.ErrorF(c, errors.New("not authorized"))
			}

			return next(c)
		}
	}
}

func (filter SecurityFilter) parseAccessToken(c echo.Context) (*client.TokenIntrospectionResponse, error) {
	bearerToken, err := filter.Auth.ParseBearerToken(c.Request())
	if err != nil {
		return nil, fmt.Errorf("failed to parse the bearer token: %w", err)
	}

	token, err := filter.Auth.IntrospectAccessToken(c.Request().Context(), bearerToken)
	if err != nil {
		return nil, fmt.Errorf("failed to introspect token: %w", err)
	}

	if !token.Active {
		return nil, fmt.Errorf("access-token is not active")
	}

	return token, nil
}
