package service

import (
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
	return nil
}

func (s *TelegramWebhookService) serviceFactory(update types.TelegramUpdate) types.WebhookService {
	structType := reflect.TypeOf(update)

	structVal := reflect.ValueOf(update)
	fieldNum := structVal.NumField()

	for i := 0; i < fieldNum; i++ {
		// Field(i) returns i'th value of the struct
		field := structVal.Field(i)
		fieldName := structType.Field(i).Name

		// CAREFUL! IsZero interprets empty strings and int equal 0 as a zero value.
		// To check only if the pointers have been initialized,
		// you can check the kind of the field:
		// if field.Kind() == reflect.Pointer { // check }

		// IsZero panics if the value is invalid.
		// Most functions and methods never return an invalid Value.
	}

	return nil
}
