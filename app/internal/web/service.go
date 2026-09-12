package web

import (
	"context"
	"log/slog"

	"github.com/adriein/hastypal/internal/booking"
	"github.com/adriein/hastypal/internal/business"
	"github.com/rotisserie/eris"
)

type WebService interface {
	ShowServices(ctx context.Context, req GetServicesReq) ([]*ServiceDTO, error)
}

type Service struct {
	logger   slog.Logger
	business business.BusinessService
	booking  booking.BookingService
}

func NewService(
	logger slog.Logger,
	business business.BusinessService,
	booking booking.BookingService,
) *Service {
	return &Service{
		logger:   logger,
		business: business,
		booking:  booking,
	}
}

/*
================================================================================
WEB SHOW SERVICES
==============================================================================
*/

func (s *Service) ShowServices(ctx context.Context, req GetServicesReq) ([]*ServiceDTO, error) {
	business, err := s.business.GetBusinessByPublicID(ctx, req.BusinessPublicID)
	if err != nil {
		return nil, eris.Wrap(err, "Error showing services trying to fetch business by public ID")
	}

	sessionID, err := s.booking.InitSession(ctx, business.ID)
	if err != nil {
		return nil, eris.Wrapf(err, "Error creating session, to book on bussiness %d", business.ID)
	}

	var dtos []*ServiceDTO

	for _, service := range business.ServiceCatalog {
		dto := &ServiceDTO{
			ID:          service.ID,
			SessionID:   sessionID,
			Name:        service.Name,
			Price:       service.Price,
			Currency:    service.Currency,
			Duration:    service.Duration,
			Description: service.Description,
		}

		dtos = append(dtos, dto)
	}

	return dtos, nil
}
