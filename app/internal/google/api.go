package google

import (
	"context"
	"os"
	"time"

	"github.com/adriein/hastypal/pkg/constants"
	"github.com/rotisserie/eris"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type CalendarApi interface {
	GetAuthenticationURL(businessID string) string
	ExchangeToken(ctx context.Context, state string, code string) (*GoogleToken, error)
	Client(ctx context.Context, businessToken *GoogleToken) (*calendar.Service, error)
}

type GoogleCalendarApi struct{}

func NewGoogleApi() *GoogleCalendarApi {
	return &GoogleCalendarApi{}
}

func (g *GoogleCalendarApi) getOauth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv(constants.GoogleClientId),
		ClientSecret: os.Getenv(constants.GoogleClientSecret),
		RedirectURL:  "http://localhost:4000/api/v1/business/google-auth-callback",
		Endpoint:     google.Endpoint,
		Scopes:       []string{calendar.CalendarEventsScope},
	}
}

func (g *GoogleCalendarApi) GetAuthenticationURL(businessID string) string {
	config := g.getOauth2Config()

	return config.AuthCodeURL(businessID, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(businessID))
}

func (g *GoogleCalendarApi) ExchangeToken(ctx context.Context, state string, code string) (*GoogleToken, error) {
	config := g.getOauth2Config()

	token, err := config.Exchange(ctx, code, oauth2.VerifierOption(state))

	if err != nil {
		return nil, eris.Wrap(err, "Error exchanging token")
	}

	googleToken := &GoogleToken{
		BusinessID:   state,
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		DateAdd:      time.Now(),
		DateUpd:      time.Now(),
	}

	return googleToken, nil
}

func (g *GoogleCalendarApi) Client(ctx context.Context, businessToken *GoogleToken) (*calendar.Service, error) {
	config := g.getOauth2Config()

	token := &oauth2.Token{
		AccessToken:  businessToken.AccessToken,
		TokenType:    businessToken.TokenType,
		RefreshToken: businessToken.RefreshToken,
	}

	client, err := calendar.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))

	if err != nil {
		return nil, eris.Wrap(err, "Error creating the calendar service")
	}

	return client, nil
}
