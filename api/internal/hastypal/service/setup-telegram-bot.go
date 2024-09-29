package service

import (
	"github.com/adriein/hastypal/internal/hastypal/types"
)

type SetupTelegramBotService struct {
	bot        *TelegramBot
	repository types.Repository[types.Business]
}

func NewSetupTelegramBotService(
	bot *TelegramBot,
	repository types.Repository[types.Business],
) *SetupTelegramBotService {
	return &SetupTelegramBotService{
		bot:        bot,
		repository: repository,
	}
}

func (s *SetupTelegramBotService) Execute(communicationPhoneNumber string, commands []types.TelegramBotCommand) error {
	criteria := s.buildCommunicationPhoneNumberCriteria(communicationPhoneNumber)

	_, findOneErr := s.repository.FindOne(criteria)

	if findOneErr != nil {
		return findOneErr
	}

	return nil
}

func (s *SetupTelegramBotService) buildCommunicationPhoneNumberCriteria(phoneNumber string) types.Criteria {
	var filters []types.Filter

	filters = make([]types.Filter, 1)

	filters[0] = types.Filter{Name: "communication_phone_number", Value: phoneNumber}

	return types.Criteria{Filters: filters}
}
