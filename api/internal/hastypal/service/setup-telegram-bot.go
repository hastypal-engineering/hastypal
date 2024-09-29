package service

import (
	"github.com/adriein/hastypal/internal/hastypal/types"
)

type SetupTelegramBotService struct {
	bot *TelegramBot
}

func NewSetupTelegramBotService(
	bot *TelegramBot,
) *SetupTelegramBotService {
	return &SetupTelegramBotService{
		bot: bot,
	}
}

func (s *SetupTelegramBotService) Execute(setup types.AdminTelegramBotSetup) error {
	if setCommandsErr := s.bot.SetCommands(setup.Commands); setCommandsErr != nil {
		return setCommandsErr
	}

	if setWebhookErr := s.bot.SetWebhook(setup.Webhook); setWebhookErr != nil {
		return setWebhookErr
	}

	return nil
}

func (s *SetupTelegramBotService) buildCommunicationPhoneNumberCriteria(phoneNumber string) types.Criteria {
	var filters []types.Filter

	filters = make([]types.Filter, 1)

	filters[0] = types.Filter{Name: "communication_phone_number", Value: phoneNumber}

	return types.Criteria{Filters: filters}
}
