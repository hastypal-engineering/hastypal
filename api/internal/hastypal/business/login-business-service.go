package business

import (
	"github.com/adriein/hastypal/internal/hastypal/constants"
	"github.com/adriein/hastypal/internal/hastypal/types"
)

type LoginBusiness struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginBusinessService struct {
	repository types.Repository[LoginBusiness]
}

func NewLoginBusinessService(repository types.Repository[LoginBusiness]) *LoginBusinessService {
	return &LoginBusinessService{
		repository: repository,
	}
}

func (s *LoginBusinessService) Execute(request LoginBusiness) error {
	// Get business with email from request
	_, getBusinessErr := s.getBusiness(request.Email)

	if getBusinessErr != nil {
		return getBusinessErr
	}
	// If not exists return error

	// Compare given password with stored password

	// If passwords don't match return error

	return nil
}

func (s *LoginBusinessService) getBusiness(email string) (LoginBusiness, error) {
	filter := types.Filter{
		Name:    "email",
		Operand: constants.Equal,
		Value:   email,
	}

	criteria := types.Criteria{Filters: []types.Filter{filter}}

	business, findOneErr := s.repository.FindOne(criteria)

	if findOneErr != nil {
		return LoginBusiness{}, findOneErr
	}

	return business, nil
}
