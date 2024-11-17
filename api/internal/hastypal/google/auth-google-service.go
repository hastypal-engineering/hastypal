package google

import (
	"github.com/adriein/hastypal/internal/hastypal/service"
)

type AuthGoogleService struct {
	googleApi *service.GoogleApi
}

func NewAuthGoogleService(
	googleApi *service.GoogleApi,
) *AuthGoogleService {
	return &AuthGoogleService{
		googleApi: googleApi,
	}
}

func (s *AuthGoogleService) Execute() string {
	googleAuthUrl := s.googleApi.GetAuthCodeUrl()

	return googleAuthUrl
}
