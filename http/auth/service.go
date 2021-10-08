package auth

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/gommon/log"
	nutsAuthClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client/auth"

	"github.com/nuts-foundation/go-did/vc"
)

var (
	bearerAuthSchema = "Bearer"

	ErrEmptyBearerToken = errors.New("empty access token")
)

type Service interface {
	RequestAccessToken(ctx context.Context, requester, authorizer, service string, vcs []vc.VerifiableCredential, identity *string) (*nutsAuthClient.AccessTokenResponse, error)
	IntrospectAccessToken(ctx context.Context, accessToken string) (*nutsAuthClient.TokenIntrospectionResponse, error)
	ParseBearerToken(request *http.Request) (string, error)
}

type authService struct {
	client *nutsAuthClient.ClientWithResponses
}

func fromStringPtr(ptr *string) (output string) {
	if ptr != nil {
		output = *ptr
	}

	return
}

func NewService(server string) (Service, error) {
	authClient, err := nutsAuthClient.NewClientWithResponses(server)
	if err != nil {
		return nil, err
	}

	return &authService{
		client: authClient,
	}, nil
}

func (s *authService) RequestAccessToken(ctx context.Context, requester, authorizer, service string, vcs []vc.VerifiableCredential, identity *string) (*nutsAuthClient.AccessTokenResponse, error) {
	httpResponse, err := s.client.RequestAccessToken(ctx, nutsAuthClient.RequestAccessTokenJSONRequestBody{
		Requester:   requester,
		Authorizer:  authorizer,
		Service:     service,
		Credentials: vcs,
		Identity:    fromStringPtr(identity),
	})
	if err != nil {
		return nil, err
	}

	response, err := nutsAuthClient.ParseRequestAccessTokenResponse(httpResponse)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK {
		log.Warnf("Server response: %s", string(response.Body))
		return nil, fmt.Errorf("invalid status code when requesting token: %d", response.StatusCode())
	}

	return response.JSON200, nil
}

func (s *authService) IntrospectAccessToken(ctx context.Context, accessToken string) (*nutsAuthClient.TokenIntrospectionResponse, error) {
	values := &url.Values{}
	values.Set("token", accessToken)

	httpResponse, err := s.client.IntrospectAccessTokenWithBody(ctx, "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(values.Encode())))
	if err != nil {
		return nil, err
	}

	response, err := nutsAuthClient.ParseIntrospectAccessTokenResponse(httpResponse)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK {
		log.Warnf("Server response: %s", string(response.Body))
		return nil, fmt.Errorf("invalid status code when requesting token: %d", response.StatusCode())
	}

	return response.JSON200, nil
}

func (s *authService) ParseBearerToken(request *http.Request) (string, error) {
	authorizationHeader := request.Header.Get("Authorization")
	if authorizationHeader == "" {
		return "", ErrEmptyBearerToken
	}

	bearerToken := authorizationHeader[len(bearerAuthSchema)+1:]
	if bearerToken == "" {
		return "", ErrEmptyBearerToken
	}

	return bearerToken, nil
}
