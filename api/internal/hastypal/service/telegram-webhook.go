package service

import (
	"fmt"
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
	s.serviceFactory(update)
	return nil
}

func (s *TelegramWebhookService) serviceFactory(update types.TelegramUpdate) types.WebhookService {
	structType := reflect.TypeOf(update)

	structVal := reflect.ValueOf(update)
	fieldNum := structVal.NumField()

	for i := 0; i < fieldNum; i++ {
		field := structVal.Field(i)
		fieldName := structType.Field(i).Name

		fmt.Println(field.IsZero())
		fmt.Println(fieldName)
	}

	return nil
}
