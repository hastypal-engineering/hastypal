package telegram

import (
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/shared/constants"
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"net/url"
	"strings"
)

type NotificationWebhookTelegramService struct {
	startCommandHandler        types.TelegramCommandHandler
	datesCommandHandler        types.TelegramCommandHandler
	hoursCommandHandler        types.TelegramCommandHandler
	confirmationCommandHandler types.TelegramCommandHandler
	finishCommandHandler       types.TelegramCommandHandler
}

func NewNotificationWebhookTelegramService(
	startCommandHandler types.TelegramCommandHandler,
	datesCommandHandler types.TelegramCommandHandler,
	hoursCommandHandler types.TelegramCommandHandler,
	confirmationCommandHandler types.TelegramCommandHandler,
	finishCommandHandler types.TelegramCommandHandler,
) *NotificationWebhookTelegramService {
	return &NotificationWebhookTelegramService{
		startCommandHandler:        startCommandHandler,
		datesCommandHandler:        datesCommandHandler,
		hoursCommandHandler:        hoursCommandHandler,
		confirmationCommandHandler: confirmationCommandHandler,
		finishCommandHandler:       finishCommandHandler,
	}
}

func (s *NotificationWebhookTelegramService) Execute(update types.TelegramUpdate) error {
	pipe := [2]types.ResolveTelegramUpdate{s.resolveBotCommand, s.resolveCallbackQueryCommand}

	for i := 0; i < len(pipe); i++ {
		parseFunc := pipe[i]

		if err := parseFunc(update); err != nil {
			switch i {
			case 0:
				return exception.Wrap(
					"resolveBotCommand",
					"notification-webhook-telegram-service",
					err,
				)
			case 1:
				return exception.Wrap(
					"resolveCallbackQueryCommand",
					"notification-webhook-telegram-service",
					err,
				)
			}
		}
	}

	return nil
}

func (s *NotificationWebhookTelegramService) resolveBotCommand(update types.TelegramUpdate) error {
	reflection := helper.NewReflectionHelper[types.TelegramUpdate]()

	if !reflection.HasField(update, constants.TelegramMessageField) {
		return nil
	}

	text := strings.Split(update.Message.Text, " ")

	handler, err := s.resolveHandler(text[0])

	if err != nil {
		return exception.Wrap("s.resolveHandler", "notification-webhook-telegram-service", err).
			WithValues([]string{text[0]})
	}

	if handlerErr := handler.Execute(update); handlerErr != nil {
		return exception.Wrap("handler.Execute", "notification-webhook-telegram-service", handlerErr)
	}

	return nil
}

func (s *NotificationWebhookTelegramService) resolveCallbackQueryCommand(update types.TelegramUpdate) error {
	reflection := helper.NewReflectionHelper[types.TelegramUpdate]()

	if !reflection.HasField(update, constants.TelegramCallbackQueryField) {
		return nil
	}

	parsedUrl, parseUrlErr := url.Parse(update.CallbackQuery.Data)

	if parseUrlErr != nil {
		return exception.New(parseUrlErr.Error()).
			Trace("url.Parse", "notification-webhook-telegram-service").
			WithValues([]string{update.CallbackQuery.Data})
	}

	handler, err := s.resolveHandler(parsedUrl.Path)

	if err != nil {
		return exception.Wrap("s.resolveHandler", "notification-webhook-telegram-service", err)
	}

	if handlerErr := handler.Execute(update); handlerErr != nil {
		return exception.Wrap("handler.Execute", "notification-webhook-telegram-service", handlerErr)
	}

	return nil
}

func (s *NotificationWebhookTelegramService) resolveHandler(command string) (types.TelegramCommandHandler, error) {
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

	return nil, exception.New(fmt.Sprintf("Hanlder not found for command: %s", command)).
		Trace("resolveHandler", "notification-webhook-telegram-service")
}
