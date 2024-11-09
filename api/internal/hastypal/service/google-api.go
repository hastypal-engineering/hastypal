package service

import (
	"context"
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/constants"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"log"
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

	verifier := oauth2.GenerateVerifier()

	return config.AuthCodeURL(verifier, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
}

func (g *GoogleApi) Client() *calendar.Service {
	ctx := context.Background()

	config := &oauth2.Config{
		ClientID:     os.Getenv(constants.GoogleClientId),
		ClientSecret: os.Getenv(constants.GoogleClientSecret),
		RedirectURL:  "http://localhost:4000/",
		Endpoint:     google.Endpoint,
		Scopes:       []string{calendar.CalendarEventsScope},
	}

	verifier := oauth2.GenerateVerifier()

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := config.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
	fmt.Printf("Visit the URL for the auth dialog: %v", url)

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	token, err := config.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		log.Fatal(err)
	}

	client, newServiceErr := calendar.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))

	if newServiceErr != nil {
		log.Fatal(newServiceErr)
	}

	return client
}
