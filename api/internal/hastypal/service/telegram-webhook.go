package service

import (
	"github.com/adriein/hastypal/internal/hastypal/constants"
	"github.com/adriein/hastypal/internal/hastypal/types"
	"reflect"
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
	pipe := [1]types.ParseTelegramUpdate{s.parseBotCommand}

	for i := 0; i < len(pipe); i++ {
		parseFunc := pipe[i]

		if err := parseFunc(update); err != nil {
			return err
		}
	}

	return nil
}

func (s *TelegramWebhookService) isChatMessage(update types.TelegramUpdate) bool {
	isChatMessage := false

	structType := reflect.TypeOf(update)

	structVal := reflect.ValueOf(update)
	fieldNum := structVal.NumField()

	for i := 0; i < fieldNum; i++ {
		field := structVal.Field(i)
		fieldName := structType.Field(i).Name

		if fieldName == constants.TelegramMessageField && field.IsZero() {
			isChatMessage = false

			break
		}

		isChatMessage = true
	}

	return isChatMessage
}

func (s *TelegramWebhookService) parseBotCommand(update types.TelegramUpdate) error {
	if !s.isChatMessage(update) {
		return nil
	}

	return nil
}
