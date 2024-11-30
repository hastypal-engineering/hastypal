package business

import (
	"github.com/adriein/hastypal/internal/hastypal/shared/constants"
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
		return types.Business{}, getBusinessErr
	}

	if comparePasswordsError := s.comparePasswords(business.Password, request.Password); comparePasswordsError != nil {
		return types.Business{}, comparePasswordsError
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
		return types.Business{}, findOneErr
	}

	return business, nil
}

func (s *LoginBusinessService) comparePasswords(storedPassword string, givenPassword string) error {
	if storedPassword != givenPassword {
		return types.ApiError{
			Msg:      "Passwords don't match",
			Function: "Execute -> comparePasswords()",
			File:     "login-business-service.go",
			Domain:   true,
		}
	}

	return nil
}
