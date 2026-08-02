package google

import (
	"github.com/adriein/hastypal/internal/hastypal/shared/service"
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

func (s *AuthGoogleService) Execute(businessId string) string {
	googleAuthUrl := s.googleApi.GetAuthCodeUrlForBusiness(businessId)

	return googleAuthUrl
}
