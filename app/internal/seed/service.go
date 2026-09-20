package seed

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/adriein/hastypal/internal/business"
	"github.com/adriein/hastypal/pkg/helper"
	"github.com/rotisserie/eris"
)

type SeedService interface {
	Run(ctx context.Context) error
}

type Service struct {
	logger   *slog.Logger
	business business.BusinessService
}

func NewService(logger *slog.Logger, business business.BusinessService) *Service {
	return &Service{
		logger:   logger,
		business: business,
	}
}

func (s *Service) Run(ctx context.Context) error {
	s.logger.Debug("Starting seed")

	for _, fakeBusiness := range BusinessSeed {
		s.logger.Debug(fmt.Sprintf("Creating fake business %s", fakeBusiness.Name))
		publicID, err := helper.GenerateUniqueSlug(fakeBusiness.Name, false)
		if err != nil {
			return eris.Wrapf(err, "Failed creating slug for business %s", fakeBusiness.Name)
		}

		fakeBusiness.PublicID = publicID

		businessID, err := s.business.CreateBusiness(ctx, fakeBusiness)
		if err != nil {
			return eris.Wrapf(err, "Error creating fake business %s", fakeBusiness.Name)
		}

		for _, fakeService := range fakeBusiness.ServiceCatalog {
			s.logger.Debug(fmt.Sprintf("Creating fake service %s", fakeService.Name))

			fakeService.BusinessID = businessID

			if _, err := s.business.CreateServiceCatalog(ctx, fakeService); err != nil {
				return eris.Wrapf(err, "Error creating fake service %s", fakeService.Name)
			}
		}
	}

	return nil
}
