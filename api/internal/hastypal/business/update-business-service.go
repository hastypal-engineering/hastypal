package business

import (
	"github.com/adriein/hastypal/internal/hastypal/shared/constants"
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
)

type UpdateBusinessService struct {
	repository types.Repository[types.Business]
}

func NewUpdateBusinessService(repository types.Repository[types.Business]) *UpdateBusinessService {
	return &UpdateBusinessService{
		repository: repository,
	}
}

func (s *UpdateBusinessService) Execute(request types.Business) error {
	if err := s.repository.Update(request); err != nil {
		return exception.Wrap(
			"s.repository.Update", "update-business-service.go", err,
		)
	}

	return nil
}

func (s *UpdateBusinessService) ensureBusinessExists(business types.Business) error {
	filter := types.Filter{
		Name:    "id",
		Operand: constants.Equal,
		Value:   business.Id,
	}

	criteria := types.Criteria{Filters: []types.Filter{filter}}

	_, err := s.repository.FindOne(criteria)

	if err != nil {
		return exception.Wrap("s.repository.FindOne", "update-business-service.go", err)
	}

	// if result {

	// }

	return nil
}
