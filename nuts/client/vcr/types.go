package vcr

import "github.com/nuts-foundation/go-did/vc"

// SearchVCRequest is the request body for searching VCs
type SearchVCRequest struct {
	// A partial VerifiableCredential in JSON-LD format. Each field will be used to match credentials against. All fields MUST be present.
	Query         vc.VerifiableCredential `json:"query"`
	SearchOptions *SearchOptions          `json:"searchOptions,omitempty"`
}
