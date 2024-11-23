package business

import (
	types "github.com/adriein/hastypal/internal/hastypal/shared/types"
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
		return err
	}

	return nil
}
