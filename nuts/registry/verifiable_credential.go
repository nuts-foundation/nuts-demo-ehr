package registry

import (
	"context"
	"encoding/json"
	"fmt"

	nutsClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client"
	"github.com/nuts-foundation/nuts-node/vcr/credential"
)

type VerifiableCredentialRegistry interface {
	// CreateAuthorizationCredential creates a NutsAuthorizationCredential on the nuts node
	CreateAuthorizationCredential(ctx context.Context, purposeOfUse, issuer, subjectID string, resources []credential.Resource) error
	// RevokeAuthorizationCredential revokes a credential based on the resourcePath contained in the credential
	RevokeAuthorizationCredential(ctx context.Context, purposeOfUse, subjectID, resourcePath string) error
}

type httpVerifiableCredentialRegistry struct {
	nutsClient nutsClient.VCRClient
}

func NewVerifiableCredentialRegistry(client *nutsClient.HTTPClient) VerifiableCredentialRegistry {
	return &httpVerifiableCredentialRegistry{
		nutsClient: client,
	}
}

func (registry *httpVerifiableCredentialRegistry) CreateAuthorizationCredential(ctx context.Context, purposeOfUse, issuer, subjectID string, resources []credential.Resource) error {
	subject := credential.NutsAuthorizationCredentialSubject{
		ID: subjectID,
		LegalBase: credential.LegalBase{
			ConsentType: "implied",
		},
		PurposeOfUse: purposeOfUse,
		Resources:    resources,
	}
	subjectMap := map[string]interface{}{}

	data, err := json.Marshal(subject)
	if err != nil {
		return fmt.Errorf("invalid subject: %w", err)
	}

	if err := json.Unmarshal(data, &subjectMap); err != nil {
		return fmt.Errorf("invalid subject: %w", err)
	}

	return registry.nutsClient.CreateVC(ctx, credential.NutsAuthorizationCredentialType, issuer, subjectMap, nil)
}

func (registry *httpVerifiableCredentialRegistry) RevokeAuthorizationCredential(ctx context.Context, purposeOfUse, subjectID, resourcePath string) error {
	// may be extended by issuanceDate for even faster results.
	params := map[string]string{
		"credentialSubject.id": subjectID,
		"credentialSubject.purposeOfUse": purposeOfUse,
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
