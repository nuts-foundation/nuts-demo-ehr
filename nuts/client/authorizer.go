package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/nuts-foundation/nuts-demo-ehr/internal/keyring"

	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/google/uuid"
)

// RequestEditorFn defines the type of function used to modify outgoing HTTP requests
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Authorizer authorizes requests to nuts-node endpoints using the given private key
type Authorizer struct {
	Key      keyring.Key
	Audience string
}

// RequestEditorFn returns a RequestEditorFn suitable for use in WithRequestEditorFn
// of HTTP client constructors.
func (a *Authorizer) RequestEditorFn() RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		// Determine the period of validity for the JWT
		issuedAt := time.Now()
		notBefore := issuedAt
		expires := notBefore.Add(time.Second * time.Duration(60))

		// Build the JWT
		token, err := jwt.NewBuilder().
			Issuer("nuts-demo-ehr").
			Subject("nuts-demo-ehr").
			Audience([]string{a.Audience}).
			IssuedAt(issuedAt).
			NotBefore(notBefore).
			Expiration(expires).
			JwtID(uuid.NewString()).
			Build()

		// Ensure the JWT was successfully built
		if err != nil {
			return fmt.Errorf("failed to build JWT: %w", err)
		}

		// Sign the JWT using the configured key
		signed, err := a.Key.SignJWT(token)
		if err != nil {
			return fmt.Errorf("failed to sign JWT: %w", err)
		}

		// Put the signed JWT in the Authorization header
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", signed))
		return nil
	}
}
