package business

import (
	"context"
	"log/slog"

	"github.com/rotisserie/eris"
)

type BusinessService interface {
	GetBusinessByID(ctx context.Context, ID int) (*Business, error)
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
