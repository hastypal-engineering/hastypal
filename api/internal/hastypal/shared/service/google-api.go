package service

import (
	"context"
	"github.com/adriein/hastypal/internal/hastypal/shared/constants"
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"os"
	"time"
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

func (g *GoogleApi) GetAuthCodeUrlForBusiness(businessId string) string {
	config := g.GetOauth2Config()

	return config.AuthCodeURL(businessId, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(businessId))
}

func (g *GoogleApi) ExchangeToken(state string, code string) (types.GoogleToken, error) {
	ctx := context.Background()
	config := g.GetOauth2Config()

	token, exchangeErr := config.Exchange(ctx, code, oauth2.VerifierOption(state))

	if exchangeErr != nil {
		return types.GoogleToken{}, exception.
			New(exchangeErr.Error()).
			Trace("config.Exchange", "google-api.go")
	}

	googleToken := types.GoogleToken{
		BusinessId:   state,
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		CreatedAt:    time.Now().Format(time.DateTime),
		UpdatedAt:    time.Now().Format(time.DateTime),
	}

	return googleToken, nil
}

func (g *GoogleApi) CalendarClient(businessToken types.GoogleToken) (*calendar.Service, error) {
	ctx := context.Background()

	config := g.GetOauth2Config()

	token := &oauth2.Token{
		AccessToken:  businessToken.AccessToken,
		TokenType:    businessToken.TokenType,
		RefreshToken: businessToken.RefreshToken,
	}

	client, newServiceErr := calendar.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))

	if newServiceErr != nil {
		return nil, exception.
			New(newServiceErr.Error()).
			Trace("calendar.NewService", "google-api.go")
	}

	return client, nil
}
