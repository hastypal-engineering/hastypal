package google

import (
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/service"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
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
		return exception.New(parseUrlErr.Error()).
			Trace("url.Parse", "auth-callback-google-service.go").
			WithValues([]string{request})
	}

	code := parsedUrl.Query().Get("code")
	state := parsedUrl.Query().Get("state")

	googleToken, exchangeTokenErr := s.googleApi.ExchangeToken(state, code)

	if exchangeTokenErr != nil {
		return exception.Wrap(
			"s.googleApi.ExchangeToken",
			"auth-callback-google-service.go",
			exchangeTokenErr,
		)
	}

	if err := s.repository.Save(googleToken); err != nil {
		return exception.Wrap(
			"s.repository.Save",
			"auth-callback-google-service.go",
			err,
		)
	}

	return nil
}
