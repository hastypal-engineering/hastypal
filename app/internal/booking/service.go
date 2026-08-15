package booking

import (
	"context"
	"log/slog"
	"time"

	"github.com/adriein/hastypal/pkg/helper"
	"github.com/rotisserie/eris"
)

type BookingService interface {
	InitSession(ctx context.Context, businessID int, chatID int) (string, error)
	GetCurrentSession(ctx context.Context, sessionID string) (*Session, error)
	RefreshSession(ctx context.Context, session *Session) error
	GetSessionsOnDate(ctx context.Context, date time.Time) ([]*Session, error)
	GetSessionOnHour(ctx context.Context, date time.Time) (*Session, error)
	RegisterBooking(ctx context.Context, sessionID string, businessID int, serviceID string, date time.Time) error
}

type Service struct {
	logger      *slog.Logger
	sessionRepo SessionRepository
	bookingRepo BookingRepository
}

func NewService(logger *slog.Logger, sessionRepo SessionRepository, bookingRepo BookingRepository) *Service {
	return &Service{
		logger:      logger,
		sessionRepo: sessionRepo,
		bookingRepo: bookingRepo,
	}
}

func (s *Service) InitSession(ctx context.Context, businessID int, chatID int) (string, error) {
	sessionId := helper.ShortUuid()

	session := &Session{
		Id:         sessionId,
		BusinessId: businessID,
		ChatId:     chatID,
		ServiceId:  "",
		Date:       "",
		Hour:       "",
		DateAdd:    time.Now().UTC(),
		DateUpd:    time.Now().UTC(),
		Ttl:        time.Minute.Milliseconds() * 5,
	}

	if err := s.sessionRepo.Save(ctx, session); err != nil {
		return "", eris.Wrap(err, "Error storing the current session")
	}

	return sessionId, nil
}

func (s *Service) GetCurrentSession(ctx context.Context, sessionID string) (*Session, error) {
	session, err := s.sessionRepo.GetByID(ctx, sessionID)

	if err != nil {
		return nil, eris.Wrap(err, "Error fetching session by ID")
	}

	return session, nil
}

func (s *Service) RefreshSession(ctx context.Context, session *Session) error {
	session.Refresh()

	if err := s.sessionRepo.Update(ctx, session); err != nil {
		return eris.Wrap(err, "Error patching the current session")
	}

	return nil
}

func (s *Service) GetSessionsOnDate(ctx context.Context, date time.Time) ([]*Session, error) {
	sessions, err := s.sessionRepo.GetByDate(ctx, date)

	if err != nil {
		return nil, eris.Wrap(err, "Error fetching sessions on specific date")
	}

	return sessions, nil
}

func (s *Service) GetSessionOnHour(ctx context.Context, date time.Time) (*Session, error) {
	sessions, err := s.sessionRepo.GetByHour(ctx, date)

	if err != nil {
		return nil, eris.Wrap(err, "Error fetching sessions on specific date")
	}

	return sessions, nil
}

func (s *Service) RegisterBooking(ctx context.Context, sessionID string, businessID int, serviceID string, date time.Time) error {
	booking := &Booking{
		SessionID:  sessionID,
		BusinessID: businessID,
		ServiceID:  serviceID,
		Date:       date,
	}

	if err := s.bookingRepo.Save(ctx, booking); err != nil {
		return eris.Wrap(err, "Error saving the booking")
	}

	return nil
}
