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

func (s *TelegramWebhookService) Execute(setup types.AdminTelegramBotSetup) error {

	return nil
}

func (s *TelegramWebhookService) buildCommunicationPhoneNumberCriteria(phoneNumber string) types.Criteria {
	var filters []types.Filter

	filters = make([]types.Filter, 1)

	filters[0] = types.Filter{Name: "communication_phone_number", Value: phoneNumber}

	return types.Criteria{Filters: filters}
}
