package google

import (
	"github.com/adriein/hastypal/internal/hastypal/shared/service"
	types2 "github.com/adriein/hastypal/internal/hastypal/shared/types"
	"net/url"
)

type AuthCallbackGoogleService struct {
	repository types2.Repository[types2.GoogleToken]
	googleApi  *service.GoogleApi
}

func NewAuthCallbackGoogleService(
	repository types2.Repository[types2.GoogleToken],
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
		return types2.ApiError{
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
