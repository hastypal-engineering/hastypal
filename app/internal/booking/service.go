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
}

type Service struct {
	logger      *slog.Logger
	sessionRepo SessionRepository
}

func NewService(logger *slog.Logger, sessionRepo SessionRepository) *Service {
	return &Service{
		logger:      logger,
		sessionRepo: sessionRepo,
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
		DateAdd:  time.Now().UTC().Format(time.DateTime),
		DateUpd:  time.Now().UTC().Format(time.DateTime),
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
