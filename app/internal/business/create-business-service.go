package business

import (
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
)

type CreateBusinessService struct {
	repository types.Repository[types.Business]
}

func NewCreateBusinessService(repository types.Repository[types.Business]) *CreateBusinessService {
	return &CreateBusinessService{
		repository: repository,
	}
}

func (s *CreateBusinessService) Execute(request types.Business) error {
	if err := s.repository.Save(request); err != nil {
		return exception.Wrap("s.repository.Save", "create-business-service.go", err)
	}

	return nil
}
