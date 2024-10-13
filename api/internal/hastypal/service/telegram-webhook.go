package service

import (
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/constants"
	"github.com/adriein/hastypal/internal/hastypal/helper"
	"github.com/adriein/hastypal/internal/hastypal/types"
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
	pipe := [2]types.ResolveTelegramUpdate{s.resolveBotCommand, s.resolveCallbackQueryCommand}

	for i := 0; i < len(pipe); i++ {
		parseFunc := pipe[i]

		if err := parseFunc(update); err != nil {
			return err
		}
	}

	return nil
}

func (s *TelegramWebhookService) resolveBotCommand(update types.TelegramUpdate) error {
	if !helper.HasField(update, constants.TelegramMessageField) {
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

	handler, err := s.resolveHandler(text[0])

	if err != nil {
		return err
	}

	if handlerErr := handler.Execute(types.Business{}, update); handlerErr != nil {
		return handlerErr
	}

	return nil
}

func (s *TelegramWebhookService) resolveCallbackQueryCommand(update types.TelegramUpdate) error {
	if !helper.HasField(update, constants.TelegramCallbackQueryField) {
		return nil
	}

	text := strings.Split(update.CallbackQuery.Data, " ")

	/*filters := make([]types.Filter, 1)

	filters[0] = types.Filter{Name: "diffusion_channel", Value: text[1]}

	criteria := types.Criteria{Filters: filters}

	business, err := s.repository.FindOne(criteria)

	if err != nil {
		return err
	}*/

	handler, err := s.resolveHandler(text[0])

	if err != nil {
		return err
	}

	if handlerErr := handler.Execute(types.Business{}, update); handlerErr != nil {
		return handlerErr
	}

	return nil
}

func (s *TelegramWebhookService) resolveHandler(command string) (types.TelegramCommandHandler, error) {
	switch command {
	case constants.StartCommand:
		return s.startCommandHandler, nil
	case constants.BookCommand:
		return s.bookCommandHandler, nil
	}

	return nil, types.ApiError{
		Msg:      fmt.Sprintf("Hanlder not found for command: %s", command),
		Function: "Execute -> resolveHandler()",
		File:     "service/telegram-webhook.go",
	}
}
