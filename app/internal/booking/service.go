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

	session := Session{
		Id:         sessionId,
		BusinessId: businessID,
		ChatId:     chatID,
		ServiceId:  "",
		Date:       "",
		Hour:       "",
		CreatedAt:  time.Now().UTC().Format(time.DateTime),
		UpdatedAt:  time.Now().UTC().Format(time.DateTime),
		Ttl:        time.Minute.Milliseconds() * 5,
	}

	if err := s.sessionRepo.Save(ctx, session); err != nil {
		return "", eris.Wrap(err, "Error storing the current session")
	}

	return sessionId, nil
}

func (s *Service) GetCurrentSession(ctx context.Context, sessionID string) (*Session, error) {
	return nil, nil
}
