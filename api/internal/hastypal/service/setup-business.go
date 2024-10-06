package service

import (
	"github.com/adriein/hastypal/internal/hastypal/types"
)

type SetupBusinessService struct {
	bot *TelegramBot
}

func NewSetupBusinessService(
	bot *TelegramBot,
) *SetupBusinessService {
	return &SetupBusinessService{
		bot: bot,
	}
}

func (s *SetupBusinessService) Execute(update types.AdminTelegramBotSetup) error {
	return nil
}

func (s *SetupBusinessService) buildCommunicationPhoneNumberCriteria(phoneNumber string) types.Criteria {
	filters := make([]types.Filter, 1)

	filters[0] = types.Filter{Name: "communication_phone_number", Value: phoneNumber}

	return types.Criteria{Filters: filters}
}
