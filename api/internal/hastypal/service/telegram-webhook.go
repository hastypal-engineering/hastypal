package service

import (
	"github.com/adriein/hastypal/internal/hastypal/constants"
	"github.com/adriein/hastypal/internal/hastypal/types"
	"reflect"
	"strings"
)

type TelegramWebhookService struct {
	repository          types.Repository[types.Business]
	startCommandHandler types.TelegramCommandHandler
	bookCommandHandler  types.TelegramCommandHandler
}

func NewTelegramWebhookService(
	repository types.Repository[types.Business],
	startCommandHandler types.TelegramCommandHandler,
	bookCommandHandler types.TelegramCommandHandler,
) *TelegramWebhookService {
	return &TelegramWebhookService{
		repository:          repository,
		startCommandHandler: startCommandHandler,
		bookCommandHandler:  bookCommandHandler,
	}
}

func (s *TelegramWebhookService) Execute(update types.TelegramUpdate) error {
	pipe := [1]types.ResolveTelegramUpdate{s.resolveBotCommand}

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

func (s *TelegramWebhookService) resolveBotCommand(update types.TelegramUpdate) error {
	if !s.isChatMessage(update) {
		return nil
	}

	text := strings.Split(update.Message.Text, " ")

	/*filters := make([]types.Filter, 1)

	filters[0] = types.Filter{Name: "diffusion_channel", Value: text[1]}

	criteria := types.Criteria{Filters: filters}

	business, err := s.repository.FindOne(criteria)

	if err != nil {
		return err
	}*/

	switch text[0] {
	case constants.StartCommand:
		handlerErr := s.startCommandHandler.Execute(types.Business{}, update)

		if handlerErr != nil {
			return handlerErr
		}
	case constants.BookCommand:
		handlerErr := s.bookCommandHandler.Execute(types.Business{}, update)

		if handlerErr != nil {
			return handlerErr
		}
	}

	return nil
}
