package api

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/client"
)

type Auth struct {
	// sessions maps session cookies to base64 encoded VPs
	sessions   map[string]string
	nodeClient client.HTTPClient
}

func NewAuth(nodeClient client.HTTPClient) *Auth {
	return &Auth {
		sessions: map[string]string{},
		nodeClient: nodeClient,
	}
}

// StoreVP stores the given VP under a new identifier or existing identifier
func (auth *Auth) StoreVP(VP string) string {
	for k, v := range auth.sessions {
		if v == VP {
			return k
		}
	}
	randomBytes := make([]byte, 64)
	_, _ = rand.Read(randomBytes)

	token := hex.EncodeToString(randomBytes)
	auth.sessions[token] = VP
	return token
}

func (auth *Auth) VPHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		protectedPaths := []string{
			"/web/private",
		}
		for _, path := range protectedPaths {
			if strings.HasPrefix(ctx.Request().RequestURI, path) {
				// check cookie
				authorized := false
				token := ""
				for _, c := range ctx.Cookies() {
					if c.Name == "session" {
						token = c.Value
						_, ok := auth.sessions[token]
						if !ok {
							return ctx.NoContent(http.StatusForbidden)
						}
						authorized = true
					}
				}
				if !authorized {
					return ctx.NoContent(http.StatusForbidden)
				}
			}
		}
		return next(ctx)
	}
}
