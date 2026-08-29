package whatsapp

import (
	"context"
	"log/slog"
)

type WhatsappService interface {
	HandleMessage(ctx context.Context, update WhatsappUpdate) error
}

type Service struct {
	logger *slog.Logger
}

func NewService(logger *slog.Logger) *Service {
	return &Service{
		logger: logger,
	}
}
