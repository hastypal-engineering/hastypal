package seed

import (
	"context"
	"log/slog"

	"github.com/adriein/hastypal/internal/business"
	"github.com/rotisserie/eris"
)

type SeedService interface {
	Run(ctx context.Context) error
}

type Service struct {
	logger   *slog.Logger
	business business.BusinessService
}

func (s *Service) Run(ctx context.Context) error {
	s.logger.Debug("Starting seed")
	for _, fakeBusiness := range BusinessSeed {
		if err := s.business.CreateBusiness(ctx, fakeBusiness); err != nil {
			return eris.Wrapf(err, "Error creating fake business %s with ID: %d", fakeBusiness.Name, fakeBusiness.ID)
		}
	}

	return nil
}
