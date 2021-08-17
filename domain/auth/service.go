package auth

import (
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"net/http"

	client "github.com/nuts-foundation/nuts-demo-ehr/client/auth"
)

type Service interface {
	RequestAccessToken(ctx context.Context, actor, custodian, service string) (*client.AccessTokenResponse, error)
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

func (s *authService) RequestAccessToken(ctx context.Context, actor, custodian, service string) (*client.AccessTokenResponse, error) {
	httpResponse, err := s.client.RequestAccessToken(ctx, client.RequestAccessTokenJSONRequestBody{
		Actor:     actor,
		Custodian: custodian,
		Service:   service,
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
