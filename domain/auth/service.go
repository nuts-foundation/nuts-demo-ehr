package auth

import (
	"bytes"
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"net/http"
	"net/url"

	"github.com/nuts-foundation/go-did/vc"
	client "github.com/nuts-foundation/nuts-demo-ehr/client/auth"
)

type Service interface {
	RequestAccessToken(ctx context.Context, actor, custodian, service string, vcs []vc.VerifiableCredential) (*client.AccessTokenResponse, error)
	IntrospectAccessToken(ctx context.Context, accessToken string) (*client.TokenIntrospectionResponse, error)
}

type authService struct {
	client *client.ClientWithResponses
}

func NewService(server string) (Service, error) {
	authClient, err := client.NewClientWithResponses(server)
	if err != nil {
		return nil, err
	}

	return &authService{
		client: authClient,
	}, nil
}

func (s *authService) RequestAccessToken(ctx context.Context, actor, custodian, service string, vcs []vc.VerifiableCredential) (*client.AccessTokenResponse, error) {
	httpResponse, err := s.client.RequestAccessToken(ctx, client.RequestAccessTokenJSONRequestBody{
		Actor:       actor,
		Custodian:   custodian,
		Service:     service,
		Credentials: vcs,
	})
	if err != nil {
		return nil, err
	}

	response, err := client.ParseRequestAccessTokenResponse(httpResponse)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("invalid status code when requesting token: %d", response.StatusCode())
	} else {
		log.Warnf("Server response: %s", string(response.Body))
	}

	return response.JSON200, nil
}

func (s *authService) IntrospectAccessToken(ctx context.Context, accessToken string) (*client.TokenIntrospectionResponse, error) {
	values := &url.Values{}
	values.Set("token", accessToken)

	httpResponse, err := s.client.IntrospectAccessTokenWithBody(ctx, "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(values.Encode())))
	if err != nil {
		return nil, err
	}

	response, err := client.ParseIntrospectAccessTokenResponse(httpResponse)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("invalid status code when requesting token: %d", response.StatusCode())
	}

	return response.JSON200, nil
}
