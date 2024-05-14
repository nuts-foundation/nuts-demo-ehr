.PHONY: dev

run-generators: gen-api

install-tools:
	go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.0.0

docker:
	docker build -t nutsfoundation/nuts-demo-ehr:main .

gen-api:
	npm run gen-api

	oapi-codegen -generate types -package types -exclude-schemas SharedCarePlan -o domain/types/generated_types.go api/api.yaml
	oapi-codegen -generate server -package api -o api/generated.go api/api.yaml
	oapi-codegen -generate client,types -package auth \
	  -import-mapping='../common/ssi_types.yaml:github.com/nuts-foundation/nuts-demo-ehr/nuts/client/common' \
	  -o nuts/client/auth/generated.go https://nuts-node.readthedocs.io/en/latest/_static/auth/v1.yaml
	oapi-codegen -generate client,types -package common -exclude-schemas VerifiableCredential,VerifiablePresentation,DID,DIDDocument -generate types,skip-prune -o nuts/client/common/generated.go https://nuts-node.readthedocs.io/en/latest/_static/common/ssi_types.yaml
	oapi-codegen -generate client,types -package discovery \
	   -import-mapping='../common/ssi_types.yaml:github.com/nuts-foundation/nuts-demo-ehr/nuts/client/common' \
	   -exclude-schemas SearchVCRequest,CredentialSubject \
	   -o nuts/client/discovery/generated.go https://nuts-node.readthedocs.io/en/latest/_static/discovery/v1.yaml
	oapi-codegen -generate client,types -package vcr \
	   -import-mapping='../common/ssi_types.yaml:github.com/nuts-foundation/nuts-demo-ehr/nuts/client/common' \
	   -exclude-schemas SearchVCRequest,CredentialSubject \
	   -o nuts/client/vcr/generated.go https://nuts-node.readthedocs.io/en/latest/_static/vcr/vcr_v2.yaml
	oapi-codegen -generate client,types -package didman \
	  -import-mapping='../common/ssi_types.yaml:github.com/nuts-foundation/nuts-demo-ehr/nuts/client/common' \
	  -o nuts/client/didman/generated.go -exclude-schemas OrganizationSearchResult https://nuts-node.readthedocs.io/en/latest/_static/didman/v1.yaml
	oapi-codegen -generate client,types -package vdr \
	  -import-mapping='../common/ssi_types.yaml:github.com/nuts-foundation/nuts-demo-ehr/nuts/client/common' \
	  -o nuts/client/vdr/generated.go https://nuts-node.readthedocs.io/en/latest/_static/vdr/v1.yaml
	oapi-codegen -generate client,types -package vdr_v2 \
	  -import-mapping='../common/ssi_types.yaml:github.com/nuts-foundation/nuts-demo-ehr/nuts/client/common' \
	  -o nuts/client/vdr_v2/generated.go https://nuts-node.readthedocs.io/en/latest/_static/vdr/v2.yaml
	oapi-codegen -generate client,types -package iam \
	  -import-mapping='../common/ssi_types.yaml:github.com/nuts-foundation/nuts-demo-ehr/nuts/client/common' \
	  -o nuts/client/iam/generated.go https://nuts-node.readthedocs.io/en/latest/_static/auth/v2.yaml
