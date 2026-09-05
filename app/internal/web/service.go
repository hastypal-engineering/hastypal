package web

import (
	"context"
	"log/slog"

	"github.com/adriein/hastypal/internal/booking"
	"github.com/adriein/hastypal/internal/business"
	"github.com/adriein/hastypal/internal/google"
	"github.com/adriein/hastypal/internal/reminder"
	"github.com/adriein/hastypal/internal/translation"
	"github.com/rotisserie/eris"
)

type WebService interface {
	ShowServices(ctx context.Context, req GetServicesReq) ([]*ServiceDTO, error)
}

type Service struct {
	logger   slog.Logger
	business business.BusinessService
	booking  booking.BookingService
	lang     translation.TranslationService
	reminder reminder.ReminderService
	google   google.GoogleService
}

func NewService(
	business business.BusinessService,
	booking booking.BookingService,
	lang translation.TranslationService,
	reminder reminder.ReminderService,
	google google.GoogleService,
) *Service {
	return &Service{
		business: business,
		booking:  booking,
		lang:     lang,
		reminder: reminder,
		google:   google,
	}
}

/*
================================================================================
WEB SHOW SERVICES
==============================================================================
*/

func (s *Service) ShowServices(ctx context.Context, req GetServicesReq) ([]*ServiceDTO, error) {
	business, err := s.business.GetBusinessByID(ctx, req.BusinessID)
	if err != nil {
		return nil, eris.Wrapf(err, "Error fetching business with ID %d", req.BusinessID)
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
