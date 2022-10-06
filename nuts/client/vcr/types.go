package vcr

import (
	"github.com/nuts-foundation/go-did/did"
	"github.com/nuts-foundation/go-did/vc"
	v2 "github.com/nuts-foundation/nuts-node/vcr/api/v2"
)

// SearchVCRequest is the request body for searching VCs
type SearchVCRequest v2.SearchVCRequest

type VerifiableCredential = vc.VerifiableCredential
type VerifiablePresentation = vc.VerifiablePresentation
type DID = did.DID
type CredentialSubject = interface{}
