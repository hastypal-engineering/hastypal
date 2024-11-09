package service

import (
	"context"
	"github.com/adriein/hastypal/internal/hastypal/types"
	"golang.org/x/oauth2"
	"net/url"
	"time"
)

type GoogleAuthCallbackService struct {
	repository types.Repository[types.GoogleToken]
	googleApi  *GoogleApi
}

func NewGoogleAuthCallbackService(
	repository types.Repository[types.GoogleToken],
	googleApi *GoogleApi,
) *GoogleAuthCallbackService {
	return &GoogleAuthCallbackService{
		repository: repository,
		googleApi:  googleApi,
	}
}

func (s *GoogleAuthCallbackService) Execute(request string) error {
	ctx := context.Background()
	config := s.googleApi.GetOauth2Config()

	parsedUrl, parseUrlErr := url.Parse(request)

	if parseUrlErr != nil {
		return types.ApiError{
			Msg:      parseUrlErr.Error(),
			Function: "Handler -> url.Parse()",
			File:     "handler/google-auth-callback.go",
			Values:   []string{request},
		}
	}

	code := parsedUrl.Query().Get("code")
	state := parsedUrl.Query().Get("state")

	token, exchangeErr := config.Exchange(ctx, code, oauth2.VerifierOption(state))

	if exchangeErr != nil {
		return types.ApiError{
			Msg:      exchangeErr.Error(),
			Function: "Handler -> config.Exchange()",
			File:     "handler/google-auth-callback.go",
		}
	}

	googleToken := types.GoogleToken{
		BusinessId:   state,
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		CreatedAt:    time.Now().Format(time.DateTime),
		UpdatedAt:    time.Now().Format(time.DateTime),
	}

	if err := s.repository.Save(googleToken); err != nil {
		return err
	}

	return nil
}
