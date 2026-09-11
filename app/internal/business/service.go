package business

import (
	"context"
	"log/slog"

	"github.com/rotisserie/eris"
)

type BusinessService interface {
	GetBusinessByID(ctx context.Context, ID int) (*Business, error)
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

func (s *Service) GetBusinessByID(ctx context.Context, ID int) (*Business, error) {
	business, err := s.repo.GetByID(ctx, ID)
	if err != nil {
		return nil, eris.Wrap(err, "Error fetching business by ID")
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
