package service

import (
	"github.com/adriein/hastypal/internal/hastypal/types"
)

type TelegramWebhookService struct {
	bot *TelegramBot
}

func NewTelegramWebhookService(
	bot *TelegramBot,
) *TelegramWebhookService {
	return &TelegramWebhookService{
		bot: bot,
	}
}

func (s *TelegramWebhookService) Execute(update types.TelegramUpdate) error {

	return nil
}
