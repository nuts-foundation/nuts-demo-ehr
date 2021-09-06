package didman

import "github.com/nuts-foundation/go-did/did"

type OrganizationSearchResult struct {
	DIDDocument  did.Document           `json:"didDocument"`
	Organization map[string]interface{} `json:"organization"`
}
