package google

import (
	"github.com/adriein/hastypal/internal/hastypal/service"
	"github.com/adriein/hastypal/internal/hastypal/types"
	"net/url"
)

type AuthCallbackGoogleService struct {
	repository types.Repository[types.GoogleToken]
	googleApi  *service.GoogleApi
}

func NewAuthCallbackGoogleService(
	repository types.Repository[types.GoogleToken],
	googleApi *service.GoogleApi,
) *AuthCallbackGoogleService {
	return &AuthCallbackGoogleService{
		repository: repository,
		googleApi:  googleApi,
	}
}

func (s *AuthCallbackGoogleService) Execute(request string) error {
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

	googleToken, exchangeTokenErr := s.googleApi.ExchangeToken(state, code)

	if exchangeTokenErr != nil {
		return exchangeTokenErr
	}

	if err := s.repository.Save(googleToken); err != nil {
		return err
	}

	return nil
}
