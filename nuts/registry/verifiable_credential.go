package registry

import (
	"context"
	"encoding/json"
	"fmt"

	nutsClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client"
	"github.com/nuts-foundation/nuts-node/vcr/credential"
)

type VerifiableCredentialRegistry interface {
	CreateAuthorizationCredential(ctx context.Context, purposeOfUse, issuer, subjectID string, resources []credential.Resource) error
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
