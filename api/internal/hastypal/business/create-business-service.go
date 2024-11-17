package business

import (
	types2 "github.com/adriein/hastypal/internal/hastypal/shared/types"
)

type CreateBusinessService struct {
	repository types2.Repository[types2.Business]
}

func NewCreateBusinessService(repository types2.Repository[types2.Business]) *CreateBusinessService {
	return &CreateBusinessService{
		repository: repository,
	}
}

func (s *CreateBusinessService) Execute(request types2.Business) error {
	if err := s.repository.Save(request); err != nil {
		return err
	}

	return nil
}
