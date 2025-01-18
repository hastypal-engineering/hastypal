package business

import (
	"github.com/adriein/hastypal/internal/hastypal/shared/constants"
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
)

type UpdateBusinessService struct {
	repository               types.Repository[types.Business]
	serviceCatalogRepository types.Repository[types.ServiceCatalog]
}

func NewUpdateBusinessService(repository types.Repository[types.Business], serviceCatalogRepository types.Repository[types.ServiceCatalog]) *UpdateBusinessService {
	return &UpdateBusinessService{
		repository:               repository,
		serviceCatalogRepository: serviceCatalogRepository,
	}
}

func (s *UpdateBusinessService) Execute(request types.Business) error {
	s.ensureBusinessExists(request)

	if err := s.repository.Update(request); err != nil {
		return exception.Wrap(
			"s.repository.Update", "update-business-service.go", err,
		)
	}

	if len(request.ServiceCatalog) != 0 {
		for i := 0; i < len(request.ServiceCatalog); i++ {
			currentServiceCatalog := types.ServiceCatalog{
				Id:         request.ServiceCatalog[i].Id,
				Name:       request.ServiceCatalog[i].Name,
				Price:      request.ServiceCatalog[i].Price,
				Currency:   request.ServiceCatalog[i].Currency,
				Duration:   request.ServiceCatalog[i].Duration,
				BusinessId: request.Id,
			}

			filter := types.Filter{
				Name:    "id",
				Operand: constants.Equal,
				Value:   currentServiceCatalog.Id,
			}

			criteria := types.Criteria{Filters: []types.Filter{filter}}

			_, err := s.serviceCatalogRepository.FindOne(criteria)

			if err != nil {
				if saveErr := s.serviceCatalogRepository.Save(currentServiceCatalog); saveErr != nil {
					return exception.Wrap("s.serviceCatalogRepository.Save", "update-business-service.go", saveErr)
				}

				continue
			}

			if updateErr := s.serviceCatalogRepository.Update(currentServiceCatalog); updateErr != nil {
				return exception.Wrap("s.serviceCatalogRepository.Update", "update-business-service.go", updateErr)
			}
		}
	}

	// TODO: Check if business has any service catalog associated and if the request is missing some of the services delete them

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

	return nil
}
