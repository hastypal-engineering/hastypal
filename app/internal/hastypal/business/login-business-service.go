package business

import (
	"github.com/adriein/hastypal/internal/hastypal/shared/constants"
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
)

type LoginBusiness struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginBusinessService struct {
	repository types.Repository[types.Business]
}

func NewLoginBusinessService(repository types.Repository[types.Business]) *LoginBusinessService {
	return &LoginBusinessService{
		repository: repository,
	}
}

func (s *LoginBusinessService) Execute(request LoginBusiness) (types.Business, error) {
	business, getBusinessErr := s.getBusiness(request.Email)

	if getBusinessErr != nil {
		return types.Business{}, exception.Wrap(
			"s.getBusiness",
			"login-business-service.go",
			getBusinessErr,
		)
	}

	if comparePasswordsError := s.comparePasswords(business.Password, request.Password); comparePasswordsError != nil {
		return types.Business{}, exception.Wrap(
			"s.comparePasswords",
			"login-business-service.go",
			comparePasswordsError,
		)
	}

	return business, nil
}

func (s *LoginBusinessService) getBusiness(email string) (types.Business, error) {
	filter := types.Filter{
		Name:    "email",
		Operand: constants.Equal,
		Value:   email,
	}

	criteria := types.Criteria{Filters: []types.Filter{filter}}

	business, findOneErr := s.repository.FindOne(criteria)

	if findOneErr != nil {
		return types.Business{}, exception.Wrap(
			"s.repository.FindOne",
			"login-business-service.go",
			findOneErr,
		)
	}

	return business, nil
}

func (s *LoginBusinessService) comparePasswords(storedPassword string, givenPassword string) error {
	if storedPassword != givenPassword {
		return exception.
			New("Passwords don't match").
			Trace("comparePasswords", "login-business-service.go").
			Domain()
	}

	return nil
}
