package service

import (
	"context"
	"github.com/adriein/hastypal/internal/hastypal/constants"
	"github.com/adriein/hastypal/internal/hastypal/types"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"os"
)

type GoogleApi struct{}

func NewGoogleApi() *GoogleApi {
	return &GoogleApi{}
}

func (g *GoogleApi) GetOauth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv(constants.GoogleClientId),
		ClientSecret: os.Getenv(constants.GoogleClientSecret),
		RedirectURL:  "http://localhost:4000/api/v1/business/google-auth-callback",
		Endpoint:     google.Endpoint,
		Scopes:       []string{calendar.CalendarEventsScope},
	}
}

func (g *GoogleApi) GetAuthCodeUrl() string {
	config := g.GetOauth2Config()

	verifier := "testId"

	return config.AuthCodeURL(verifier, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
}

func (g *GoogleApi) Client(businessToken types.GoogleToken) (*calendar.Service, error) {
	ctx := context.Background()

	config := g.GetOauth2Config()

	token := &oauth2.Token{
		AccessToken:  businessToken.AccessToken,
		TokenType:    businessToken.TokenType,
		RefreshToken: businessToken.RefreshToken,
	}

	client, newServiceErr := calendar.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))

	if newServiceErr != nil {
		return nil, types.ApiError{
			Msg:      newServiceErr.Error(),
			Function: "Client -> calendar.NewService()",
			File:     "service/google-api.go",
		}
	}

	return client, nil
}
