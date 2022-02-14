package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nuts-foundation/go-did/vc"

	nutsClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/client/vcr"
	"github.com/nuts-foundation/nuts-node/vcr/credential"
)

type VCRSearchParams struct {
	PurposeOfUse string
	SubjectID    string
	Issuer       string
	Subject      string
	ResourcePath string
}

func convertCredential(input *vcr.VerifiableCredential) (*vc.VerifiableCredential, error) {
	result := &vc.VerifiableCredential{}

	bytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, result); err != nil {
		return nil, err
	}

	return result, nil
}

type VerifiableCredentialRegistry interface {
	// CreateKVKCredential creates a new KVKCredential in the registry
	CreateKVKCredential(ctx context.Context, issuer string, proof string) error
	// CreateAuthorizationCredential creates a NutsAuthorizationCredential on the nuts node
	CreateAuthorizationCredential(ctx context.Context, issuer string, subject *credential.NutsAuthorizationCredentialSubject) error
	// RevokeAuthorizationCredential revokes a credential based on the resourcePath contained in the credential
	RevokeAuthorizationCredential(ctx context.Context, purposeOfUse, subjectID, resourcePath string) error
	// ResolveVerifiableCredential from the nuts node. It also returns untrusted credentials
	ResolveVerifiableCredential(ctx context.Context, credentialID string) (*vc.VerifiableCredential, error)
	// FindAuthorizationCredentials returns the NutsAuthorizationCredential for the given params
	FindAuthorizationCredentials(ctx context.Context, params *VCRSearchParams) ([]vc.VerifiableCredential, error)
	// FindKVKCredential returns the KVKCredential for the given issuer
	FindKVKCredential(ctx context.Context, issuer string) (*vc.VerifiableCredential, error)
}

type httpVerifiableCredentialRegistry struct {
	nutsClient nutsClient.VCRClient
}

func NewVerifiableCredentialRegistry(client nutsClient.VCRClient) VerifiableCredentialRegistry {
	return &httpVerifiableCredentialRegistry{
		nutsClient: client,
	}
}

func (registry *httpVerifiableCredentialRegistry) CreateKVKCredential(ctx context.Context, issuer, proof string) error {
	return registry.nutsClient.CreateVC(ctx, "NutsKVKCredential", issuer, nil, nil, []interface{}{map[string]interface{}{
		"type":       "IRMASignatureProof",
		"proofValue": proof,
	}})
}

func (registry *httpVerifiableCredentialRegistry) FindKVKCredential(ctx context.Context, issuer string) (*vc.VerifiableCredential, error) {
	result, err := registry.nutsClient.FindKVKCredential(ctx, issuer)
	if err != nil {
		return nil, err
	}

	return convertCredential(result)
}

func (registry *httpVerifiableCredentialRegistry) CreateAuthorizationCredential(ctx context.Context, issuer string, subject *credential.NutsAuthorizationCredentialSubject) error {
	subjectMap := map[string]interface{}{}

	data, err := json.Marshal(subject)
	if err != nil {
		return fmt.Errorf("invalid subject: %w", err)
	}

	if err := json.Unmarshal(data, &subjectMap); err != nil {
		return fmt.Errorf("invalid subject: %w", err)
	}

	return registry.nutsClient.CreateVC(ctx, credential.NutsAuthorizationCredentialType, issuer, subjectMap, nil, nil)
}

func (registry *httpVerifiableCredentialRegistry) FindAuthorizationCredentials(ctx context.Context, params *VCRSearchParams) ([]vc.VerifiableCredential, error) {
	// may be extended by issuanceDate for even faster results.
	searchParams := map[string]string{
		"credentialSubject.purposeOfUse":     params.PurposeOfUse,
		"credentialSubject.resources.#.path": params.ResourcePath,
	}

	if params.SubjectID != "" {
		searchParams["credentialSubject.id"] = params.SubjectID
	}

	if params.Subject != "" {
		searchParams["credentialSubject.subject"] = params.Subject
	}

	if params.Issuer != "" {
		searchParams["issuer"] = params.Issuer
	}

	credentials, err := registry.nutsClient.FindAuthorizationCredentials(ctx, searchParams)
	if err != nil {
		return nil, err
	}

	results := make([]vc.VerifiableCredential, len(credentials))

	for i, authCredential := range credentials {
		result, err := convertCredential(&authCredential) //nolint:gosec
		if err != nil {
			return nil, err
		}

		results[i] = *result
	}

	return results, nil
}

func (registry *httpVerifiableCredentialRegistry) RevokeAuthorizationCredential(ctx context.Context, purposeOfUse, subjectID, resourcePath string) error {
	// may be extended by issuanceDate for even faster results.
	params := map[string]string{
		"credentialSubject.id":               subjectID,
		"credentialSubject.purposeOfUse":     purposeOfUse,
		"credentialSubject.resources.#.path": resourcePath,
	}
	credentialIDs, err := registry.nutsClient.FindAuthorizationCredentialIDs(ctx, params)
	if err != nil {
		return err
	}
	for _, ID := range credentialIDs {
		if err = registry.nutsClient.RevokeCredential(ctx, ID); err != nil {
			return err
		}
	}

	return nil
}

func (registry *httpVerifiableCredentialRegistry) ResolveVerifiableCredential(ctx context.Context, credentialID string) (*vc.VerifiableCredential, error) {
	authCredential, err := registry.nutsClient.ResolveVerifiableCredential(ctx, credentialID, true)
	if err != nil {
		return nil, err
	}

	return convertCredential(authCredential)
}
