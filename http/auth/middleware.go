package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	nutsAuthClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client/auth"
)

const AccessToken = "accessToken"

type ErrorFunc func(c echo.Context, err error) error

func DefaultErrorFunc(c echo.Context, err error) error {
	return err
}

type AccessFunc func(ctx echo.Context, request *http.Request, token *nutsAuthClient.TokenIntrospectionResponse) error

func DefaultAccessFunc(ctx echo.Context, request *http.Request, token *nutsAuthClient.TokenIntrospectionResponse) error {
	return nil
}

type Config struct {
	Skipper middleware.Skipper
	ErrorF  ErrorFunc
	AccessF AccessFunc
}

type SecurityFilter struct {
	Auth Service
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

			if internalValue := c.Get("internal"); internalValue != nil {
				if value, ok := internalValue.(bool); ok && value {
					// this call has already been validated, it must be a Task update
					return next(c)
				}
			}

			c.Logger().Debugf("Checking access token on %s %s", c.Request().Method, c.Request().RequestURI)
			token, err := filter.parseAccessToken(c)
			if err != nil {
				return config.ErrorF(c, errors.New("invalid access-token"))
			}

			c.Set(AccessToken, *token)

			if err := config.AccessF(c, c.Request(), token); err != nil {
				return config.ErrorF(c, errors.New("not authorized"))
			}

			return next(c)
		}
	}
}

func (filter SecurityFilter) parseAccessToken(c echo.Context) (*nutsAuthClient.TokenIntrospectionResponse, error) {
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
