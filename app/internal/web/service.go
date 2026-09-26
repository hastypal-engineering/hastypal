package web

import (
	"context"
	"log/slog"

	"github.com/adriein/hastypal/internal/booking"
	"github.com/adriein/hastypal/internal/business"
	"github.com/adriein/hastypal/pkg/middleware"
	"github.com/rotisserie/eris"
)

type WebService interface {
	ShowServices(ctx context.Context, req GetServicesReq) (*BookingDTO, error)
	StoreService(ctx context.Context, dto *BookingPatchDTO) error
	ShowDates(ctx context.Context, req GetDatesReq) (*BookingDTO, error)
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

func (s *Service) ShowServices(ctx context.Context, req GetServicesReq) (*BookingDTO, error) {
	business, err := s.business.GetBusinessByPublicID(ctx, req.BusinessPublicID)
	if err != nil {
		return nil, eris.Wrap(err, "Error showing services trying to fetch business by public ID")
	}

	sessionID, err := s.booking.InitSession(ctx, business.ID)
	if err != nil {
		return nil, eris.Wrapf(err, "Error creating session, to book on bussiness %d", business.ID)
	}

	catalog, err := s.business.GetServiceCatalog(ctx, business.ID)
	if err != nil {
		return nil, eris.Wrapf(err, "Error showing services trying to fetch the catalog of business %d", business.ID)
	}

	var services []*ServiceDTO

	for _, service := range catalog {
		dto := &ServiceDTO{
			ID:          service.ID,
			Name:        service.Name,
			Price:       service.Price,
			Currency:    service.Currency,
			Duration:    service.Duration,
			Description: service.Description,
		}

		services = append(services, dto)
	}

	dto := &BookingDTO{
		SessionID: sessionID,
		Step:      1,
		Services:  services,
		Business: &BusinessDTO{
			PublicID:    business.PublicID,
			Name:        business.Name,
			Email:       business.Email,
			Address:     business.Address,
			Phone:       business.ContactPhone,
			Description: "the better business",
		},
	}

	return dto, nil
}

func (s *Service) StoreService(ctx context.Context, dto *BookingPatchDTO) error {
	traceID := ctx.Value(middleware.TraceIDKey)

	session, err := s.booking.GetCurrentSession(ctx, dto.SessionID)
	if err != nil {
		s.logger.Error("Error retrieving the session while storing the selected service", "trace_id", traceID, "error", eris.ToString(err, true), "session_id", session.ID, "business_id", session.BusinessID)
		return eris.Wrap(err, "Failed retrieving session to patch the booking")
	}

	if err := session.EnsureIsValid(); err != nil {
		return eris.Wrap(err, "Session has expired")
	}

	session.ServiceID = dto.ServiceID

	if err := s.booking.PatchSession(ctx, session); err != nil {
		s.logger.Error("Error patching the session while storing the selected service", "trace_id", traceID, "error", eris.ToString(err, true), "session_id", session.ID, "business_id", session.BusinessID)
		return eris.Wrap(err, "Error patching the session with the serviceID")
	}

	return nil
}

func (s *Service) ShowDates(ctx context.Context, req GetDatesReq) (*BookingDTO, error) {
	business, err := s.business.GetBusinessByPublicID(ctx, req.BusinessPublicID)
	if err != nil {
		return nil, eris.Wrap(err, "Error showing services trying to fetch business by public ID")
	}

	session, err := s.booking.GetCurrentSession(ctx, req.SessionID)
	if err != nil {
		return nil, eris.Wrapf(err, "Error creating session, to book on bussiness %d", business.ID)
	}

	dto := &BookingDTO{
		SessionID: session.ID,
		Step:      2,
		Business: &BusinessDTO{
			PublicID:    business.PublicID,
			Name:        business.Name,
			Email:       business.Email,
			Address:     business.Address,
			Phone:       business.ContactPhone,
			Description: "the better business",
		},
	}

	return dto, nil
}
