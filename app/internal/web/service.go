package web

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/adriein/hastypal/internal/booking"
	"github.com/adriein/hastypal/internal/business"
	"github.com/adriein/hastypal/pkg/helper/conversion"
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
			Duration:    conversion.BeautifyDuration(service.Duration),
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

/*
================================================================================
WEB PATCH BOOKING
==============================================================================
*/

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

/*
================================================================================
WEB SHOW DATES
==============================================================================
*/

func (s *Service) ShowDates(ctx context.Context, req GetDatesReq) (*BookingDTO, error) {
	business, err := s.business.GetBusinessByPublicID(ctx, req.BusinessPublicID)
	if err != nil {
		return nil, eris.Wrap(err, "Error showing dates trying to fetch business by public ID")
	}

	session, err := s.booking.GetCurrentSession(ctx, req.SessionID)
	if err != nil {
		return nil, eris.Wrapf(err, "Error showing dates while fetching session, to book on bussiness %d", business.ID)
	}

	if err := session.EnsureIsValid(); err != nil {
		return nil, eris.Wrap(err, "Error showing dates, session has expired")
	}

	schedule, err := s.business.GetBusinessSchedule(ctx, business.ID)
	if err != nil {
		return nil, eris.Wrapf(err, "Error showing dates while fetching business schedule on business %d", business.ID)
	}

	selectedService, err := s.business.GetServiceByID(ctx, session.ServiceID)
	if err != nil {
		return nil, eris.Wrapf(err, "Error showing dates while fetching service with id %d", session.ServiceID)
	}
	fmt.Println(selectedService)

	reqDay := req.Day

	if req.Day.IsZero() {
		reqDay = time.Now().Add(-24 * time.Hour)
	}

	// We initialize a time table with 1440 positions, every position is a minute inside a day
	timeTable := TimeTable{Data: make([]bool, 1440)}

	otherUserSessions, err := s.booking.GetSessionsOnDateByBusiness(ctx, business.ID, reqDay)
	if err != nil {
		return nil, eris.Wrap(err, "Error showing dates while fetching other users sessions")
	}

	var slots []*SlotDTO

	for _, day := range schedule.WeeklySchedule {
		if reqDay.Weekday() != day.DayOfWeek {
			continue
		}

		for _, timeSlot := range day.TimeSlots {
			openHour, err := conversion.StringToTime(timeSlot.OpenTime, time.TimeOnly)
			if err != nil {
				return nil, eris.Wrapf(err, "Error showing dates while converting openHour %s to time.Time", openHour)
			}

			closeHour, err := conversion.StringToTime(timeSlot.CloseTime, time.TimeOnly)
			if err != nil {
				return nil, eris.Wrapf(err, "Error showing dates while converting closeHour %s to time.Time", closeHour)
			}

			openTime := conversion.TimeToMinFromMidnight(openHour)
			closeTime := conversion.TimeToMinFromMidnight(closeHour)

			timeTable.MarkTimeSlot(openTime, closeTime, true)
		}

		for _, otherUserSession := range otherUserSessions {
			serviceStartInMin := int(otherUserSession.Interval.Start)
			serviceEndInMin := int(otherUserSession.Interval.End)

			timeTable.MarkTimeSlot(serviceStartInMin, serviceEndInMin, false)
		}

		stepInterval := 30

		for _, timeSlot := range day.TimeSlots {
			openHour, err := conversion.StringToTime(timeSlot.OpenTime, time.TimeOnly)
			if err != nil {
				return nil, eris.Wrapf(err, "Error showing dates while converting openHour %s to time.Time", openHour)
			}

			closeHour, err := conversion.StringToTime(timeSlot.CloseTime, time.TimeOnly)
			if err != nil {
				return nil, eris.Wrapf(err, "Error showing dates while converting closeHour %s to time.Time", closeHour)
			}

			openTime := conversion.TimeToMinFromMidnight(openHour)
			closeTime := conversion.TimeToMinFromMidnight(closeHour)

			for i := openTime; i+stepInterval <= closeTime; i += stepInterval {
				isAvailable := timeTable.IsChunkAllTrue(i, stepInterval)

				slotTime, err := conversion.TimeFromInt(i)
				if err != nil {
					return nil, eris.Wrap(err, "Error showing dates while converting timeTable int to time")
				}

				hourStr := fmt.Sprintf("%02d:%02d", slotTime.Hour(), slotTime.Minute())

				slots = append(slots, &SlotDTO{
					Hour:        hourStr,
					IsAvailable: isAvailable,
				})
			}
		}
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
		Slots: &BookingDatesDTO{
			Day:   reqDay.String(),
			Slots: slots,
		},
	}

	return dto, nil
}
