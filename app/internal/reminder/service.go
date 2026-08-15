package reminder

import (
	"context"
	"log/slog"
	"time"
)

type ReminderService interface {
	NewReminder(ctx context.Context, bookingID string, scheduledAt time.Time) error
}

type Service struct {
	logger *slog.Logger
	repo   ReminderRepository
}

func NewService(logger *slog.Logger, repo ReminderRepository) *Service {
	return &Service{
		logger: logger,
		repo:   repo,
	}
}

func (s *Service) NewReminder(ctx context.Context, bookingID string, scheduledAt time.Time) error {
	reminder := &Reminder{
		BookingID:   bookingID,
		ScheduledAt: scheduledAt,
	}

	s.repo.Save(ctx, reminder)

	return nil
}
