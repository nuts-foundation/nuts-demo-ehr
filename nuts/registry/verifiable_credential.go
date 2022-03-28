package registry

import (
	"context"
	"encoding/json"
	"fmt"
	ssi "github.com/nuts-foundation/go-did"
	"github.com/nuts-foundation/go-did/vc"

	nutsClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client"
	"github.com/nuts-foundation/nuts-node/vcr/credential"
)

type VCRSearchParams struct {
	PurposeOfUse string
	SubjectID    string
	Issuer       string
	Subject      string
	ResourcePath string
}

type VerifiableCredentialRegistry interface {
	// CreateAuthorizationCredential creates a NutsAuthorizationCredential on the nuts node
	CreateAuthorizationCredential(ctx context.Context, issuer string, subject *credential.NutsAuthorizationCredentialSubject) error
	// RevokeAuthorizationCredential revokes a credential based on the resourcePath contained in the credential
	RevokeAuthorizationCredential(ctx context.Context, purposeOfUse, subjectID, resourcePath string) error
	// ResolveVerifiableCredential from the nuts node. It also returns untrusted credentials
	ResolveVerifiableCredential(ctx context.Context, credentialID string) (*vc.VerifiableCredential, error)
	// FindAuthorizationCredentials returns the NutsAuthorizationCredential for the given params
	FindAuthorizationCredentials(ctx context.Context, params *VCRSearchParams) ([]vc.VerifiableCredential, error)
}

type httpVerifiableCredentialRegistry struct {
	nutsClient nutsClient.VCRClient
}

func NewVerifiableCredentialRegistry(client nutsClient.VCRClient) VerifiableCredentialRegistry {
	return &httpVerifiableCredentialRegistry{
		nutsClient: client,
	}
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

	return registry.nutsClient.CreateVC(ctx, credential.NutsAuthorizationCredentialType, issuer, subjectMap, nil, false)
}

func (registry *httpVerifiableCredentialRegistry) FindAuthorizationCredentials(ctx context.Context, params *VCRSearchParams) ([]vc.VerifiableCredential, error) {
	query := nutsClient.GetNutsCredentialTemplate(*credential.NutsAuthorizationCredentialTypeURI)
	credentialSubject := make(map[string]interface{}, 0)
	query.CredentialSubject = []interface{}{credentialSubject}

	// may be extended by issuanceDate for even faster results.
	credentialSubject["purposeOfUse"] = params.PurposeOfUse
	credentialSubject["resources"] = map[string]string{"path": params.ResourcePath}

	if params.SubjectID != "" {
		credentialSubject["id"] = params.SubjectID
	}

	if params.Subject != "" {
		credentialSubject["subject"] = params.Subject
	}

	if params.Issuer != "" {
		query.Issuer = ssi.MustParseURI(params.Issuer)
	}

	return registry.nutsClient.FindCredentials(ctx, query, false)
}

func (registry *httpVerifiableCredentialRegistry) RevokeAuthorizationCredential(ctx context.Context, purposeOfUse, subjectID, resourcePath string) error {
	// may be extended by issuanceDate for even faster results.
	query := nutsClient.GetNutsCredentialTemplate(*credential.NutsAuthorizationCredentialTypeURI)
	query.CredentialSubject = []interface{}{
		map[string]interface{}{
			"id":           subjectID,
			"purposeOfUse": purposeOfUse,
			"resources": map[string]interface{}{
				"path": resourcePath,
			},
		},
	}
	credentialIDs, err := registry.nutsClient.FindCredentialIDs(ctx, query, false)
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
	return registry.nutsClient.ResolveCredential(ctx, credentialID)
}
