package service

import (
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/constants"
	"github.com/adriein/hastypal/internal/hastypal/helper"
	"github.com/adriein/hastypal/internal/hastypal/types"
	"net/url"
	"strings"
)

type TelegramWebhookService struct {
	repository                 types.Repository[types.Business]
	startCommandHandler        types.TelegramCommandHandler
	datesCommandHandler        types.TelegramCommandHandler
	hoursCommandHandler        types.TelegramCommandHandler
	confirmationCommandHandler types.TelegramCommandHandler
	finishCommandHandler       types.TelegramCommandHandler
}

func NewTelegramWebhookService(
	repository types.Repository[types.Business],
	startCommandHandler types.TelegramCommandHandler,
	datesCommandHandler types.TelegramCommandHandler,
	hoursCommandHandler types.TelegramCommandHandler,
	confirmationCommandHandler types.TelegramCommandHandler,
	finishCommandHandler types.TelegramCommandHandler,
) *TelegramWebhookService {
	return &TelegramWebhookService{
		repository:                 repository,
		startCommandHandler:        startCommandHandler,
		datesCommandHandler:        datesCommandHandler,
		hoursCommandHandler:        hoursCommandHandler,
		confirmationCommandHandler: confirmationCommandHandler,
		finishCommandHandler:       finishCommandHandler,
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
	reflection := helper.NewReflectionHelper[types.TelegramUpdate]()

	if !reflection.HasField(update, constants.TelegramMessageField) {
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
	reflection := helper.NewReflectionHelper[types.TelegramUpdate]()

	if !reflection.HasField(update, constants.TelegramCallbackQueryField) {
		return nil
	}

	parsedUrl, parseUrlErr := url.Parse(update.CallbackQuery.Data)

	if parseUrlErr != nil {
		return types.ApiError{
			Msg:      parseUrlErr.Error(),
			Function: "Execute -> url.Parse()",
			File:     "telegram-webhook.go",
			Values:   []string{update.CallbackQuery.Data},
		}
	}

	/*filters := make([]types.Filter, 1)

	filters[0] = types.Filter{Name: "diffusion_channel", Value: text[1]}

	criteria := types.Criteria{Filters: filters}

	business, err := s.repository.FindOne(criteria)

	if err != nil {
		return err
	}*/

	handler, err := s.resolveHandler(parsedUrl.Path)

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
	case constants.DatesCommand:
		return s.datesCommandHandler, nil
	case constants.HoursCommand:
		return s.hoursCommandHandler, nil
	case constants.ConfirmationCommand:
		return s.confirmationCommandHandler, nil
	case constants.FinishCommand:
		return s.finishCommandHandler, nil
	}

	return nil, types.ApiError{
		Msg:      fmt.Sprintf("Hanlder not found for command: %s", command),
		Function: "Execute -> resolveHandler()",
		File:     "service/telegram-webhook.go",
	}
}
