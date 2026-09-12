package business

import (
	"context"
	"log/slog"

	"github.com/rotisserie/eris"
)

type BusinessService interface {
	GetBusinessByPublicID(ctx context.Context, ID string) (*Business, error)
	CreateBusiness(ctx context.Context, business *Business) (int, error)
	CreateServiceCatalog(ctx context.Context, service *ServiceCatalog) (int, error)
}

type Service struct {
	logger *slog.Logger
	repo   BusinessRepository
}

func NewService(logger *slog.Logger, repo BusinessRepository) *Service {
	return &Service{
		logger: logger,
		repo:   repo,
	}
}

func (s *Service) GetBusinessByPublicID(ctx context.Context, ID string) (*Business, error) {
	business, err := s.repo.GetByPublicID(ctx, ID)
	if err != nil {
		return nil, eris.Wrapf(err, "Error fetching business by public ID %s", ID)
	}

	return business, nil
}

func (s *Service) CreateBusiness(ctx context.Context, business *Business) (int, error) {
	ID, err := s.repo.Create(ctx, business)
	if err != nil {
		return 0, eris.Wrapf(err, "Error creating business %s", business.Name)
	}

	return ID, nil
}

func (s *Service) CreateServiceCatalog(ctx context.Context, service *ServiceCatalog) (int, error) {
	ID, err := s.repo.CreateService(ctx, service)
	if err != nil {
		return 0, eris.Wrapf(err, "Error creating service %s, for business with ID %d", service.Name, service.BusinessID)
	}

	return ID, nil
}
