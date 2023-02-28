package client

import (
	"context"
	"net/http"

	"github.com/nuts-foundation/nuts-demo-ehr/keyring"
)

// RequestEditorFn defines the type of function used to modify outgoing HTTP requests
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Authorizer authorizes requests to nuts-node endpoints using the given private key
type Authorizer struct {
	Key keyring.Key
}

// RequestEditorFn returns a RequestEditorFn suitable for use in WithRequestEditorFn
// of HTTP client constructors.
func (a *Authorizer) RequestEditorFn() RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "foo")
		return nil
	}
}
