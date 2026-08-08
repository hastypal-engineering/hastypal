package telegram

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/adriein/hastypal/internal/booking"
	"github.com/adriein/hastypal/internal/business"
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"github.com/adriein/hastypal/pkg/constants"
	"github.com/adriein/hastypal/pkg/helper/array"
	"github.com/adriein/hastypal/pkg/helper/reflection"
	"github.com/rotisserie/eris"
)

type TelegramService interface {
	HandleMessage(ctx context.Context, update TelegramUpdate) error
}

type Service struct {
	business business.BusinessService
	booking  booking.BookingService
	bot      *TelegramBot
}

func NewService(business business.BusinessService) *Service {
	return &Service{
		business: business,
	}
}

/*
================================================================================
TELEGRAM UPDATE HANLDER
================================================================================
*/

func (s *Service) HandleMessage(ctx context.Context, update TelegramUpdate) error {
	if reflection.HasField(update, constants.TelegramMessageField) {
		if err := s.resolveBotCommand(ctx, update); err != nil {
			eris.Wrap(err, "Error resolving bot command")
		}

		return nil
	}

	if err := s.resolveCallbackQueryCommand(ctx, update); err != nil {
		eris.Wrap(err, "Error resolving callback query command")
	}

	return nil
}

func (s *Service) resolveBotCommand(ctx context.Context, update TelegramUpdate) error {
	if !reflection.HasField(update, constants.TelegramMessageField) {
		return nil
	}

	txtArr := strings.Split(update.Message.Text, " ")
	command := txtArr[0]

	switch command {
	case constants.StartCommand:
		return s.startConversation(ctx, update)
	}

	return nil
}

func (s *Service) resolveCallbackQueryCommand(ctx context.Context, update TelegramUpdate) error {
	if !reflection.HasField(update, constants.TelegramCallbackQueryField) {
		return nil
	}

	url, err := url.Parse(update.CallbackQuery.Data)

	if err != nil {
		return eris.Wrap(err, "Error parsing the callback query url")
	}

	switch url.Path {
	case constants.ServiceCommand:
		return nil
	case constants.DatesCommand:
		return nil
	case constants.HoursCommand:
		return nil
	case constants.ConfirmationCommand:
		return nil
	case constants.FinishCommand:
		return nil
	}

	return nil
}

/*
================================================================================
TELEGRAM START CONVERSATION COMMAND
================================================================================
*/

func (s *Service) startConversation(ctx context.Context, update TelegramUpdate) error {
	var markdownText strings.Builder

	businessIdRaw := strings.ReplaceAll(update.Message.Text, "/start ", "")

	businessID, err := strconv.Atoi(businessIdRaw)

	if err != nil {
		return eris.Wrap(err, "Error converting business ID to int")
	}

	business, err := s.business.GetBusinessByID(ctx, businessID)

	if err != nil {
		return eris.Wrap(err, "Error fetching business")
	}

	sessionID, err := s.booking.InitSession(ctx, businessID, update.Message.Chat.Id)

	if err != nil {
		return eris.Wrap(err, "Error creating a session for this conversation")
	}

	welcome := fmt.Sprintf(
		"Hola %s ![👋](tg://emoji?id=5368324170671202286), soy HastypalBot el ayudante de %s\\.\n\n",
		update.Message.From.FirstName,
		"Hastypal Business Test",
	)

	services := []string{
		"Corte de pelo y barba express 18€",
		"Corte de pelo y barba premium 22€",
	}

	emoji := "![🔸](tg://emoji?id=5368324170671202286)"

	markdownText.WriteString(welcome)
	markdownText.WriteString("*Te muestro a continuación los servicios que ofrecemos:*\n\n")

	buttons := make([]KeyboardButton, len(services))

	for i, serv := range services {
		markdownText.WriteString(fmt.Sprintf("%s %s\n\n", emoji, serv))

		buttons[i] = KeyboardButton{
			Text:         fmt.Sprintf("%s 📅", services[i]),
			CallbackData: fmt.Sprintf("/dates?session=%s&service=%s&page=0", sessionID, "test-short"),
		}
	}

	inlineKeyboard := array.Chunk(buttons, 1)

	message := TelegramMessage{
		ChatId:         update.Message.Chat.Id,
		Text:           markdownText.String(),
		ParseMode:      constants.TelegramMarkdown,
		ProtectContent: true,
		ReplyMarkup:    ReplyMarkup{InlineKeyboard: inlineKeyboard},
	}

	bookingMessage := BookingTelegramMessage{
		BusinessName:     business.Name,
		BookingSessionId: sessionID,
		Message:          message,
	}

	if err := s.bot.SendMsg(bookingMessage); err != nil {
		return eris.Wrap(err, "Error sending message to telegram")
	}

	return nil
}

/*
================================================================================
TELEGRAM SELECT SERVICE COMMAND
================================================================================
*/

func (s *Service) selectService(ctx context.Context, update TelegramUpdate) error {
	ack := AnswerCallbackQuery{CallbackQueryId: update.CallbackQuery.Id}

	if err := s.bot.AnswerCallbackQuery(ack); err != nil {
		return eris.Wrap(err, "Error acking telegram conversation")
	}

	var markdownText strings.Builder

	parsedUrl, err := url.Parse(update.CallbackQuery.Data)

	if err != nil {
		return eris.Wrap(err, "Error parsing url")
	}

	queryParams := parsedUrl.Query()

	sessionID := queryParams.Get("sessionId")

	session, err := s.booking.GetCurrentSession(ctx, sessionID)

	if err != nil {
		return eris.Wrap(err, "Error fetching current booking session")
	}

	business, err := s.business.GetBusinessByID(ctx, session.BusinessId)

	if err != nil {
		return eris.Wrap(err, "Error fetching business")
	}

	if err := session.EnsureIsValid(); err != nil {
		message := TelegramMessage{ChatId: update.CallbackQuery.From.Id}

		expiredSessionMessage := message.SessionExpired()

		bookingExpiredSessionMessage := BookingTelegramMessage{
			BusinessName:     business.Name,
			BookingSessionId: session.Id,
			Message:          expiredSessionMessage,
		}

		if err := s.bot.SendMsg(bookingExpiredSessionMessage); err != nil {
			return eris.Wrap(err, "Error sending message to telegram")
		}

		return nil
	}

	if updateSessionErr := s.updateSession(session); updateSessionErr != nil {
		return exception.Wrap("s.updateSession", "pick-service-command-telegram-service.go", updateSessionErr)
	}

	services := []string{
		"Corte de pelo y barba express 18€",
		"Corte de pelo y barba premium 22€",
	}

	emoji := "![🔸](tg://emoji?id=5368324170671202286)"

	markdownText.WriteString("*Te muestro a continuación los servicios que ofrecemos:*\n\n")

	buttons := make([]types.KeyboardButton, len(services))

	for i, serv := range services {
		markdownText.WriteString(fmt.Sprintf("%s %s\n\n", emoji, serv))

		buttons[i] = types.KeyboardButton{
			Text:         fmt.Sprintf("%s 📅", services[i]),
			CallbackData: fmt.Sprintf("/dates?session=%s&service=%s&page=0", session.Id, "test-short"),
		}
	}

	array := helper.NewArrayHelper[types.KeyboardButton]()

	inlineKeyboard := array.Chunk(buttons, 1)

	message := types.TelegramMessage{
		ChatId:         update.CallbackQuery.From.Id,
		Text:           markdownText.String(),
		ParseMode:      constants.TelegramMarkdown,
		ProtectContent: true,
		ReplyMarkup:    types.ReplyMarkup{InlineKeyboard: inlineKeyboard},
	}

	bookingMessage := types.BookingTelegramMessage{
		BusinessName:     business.Name,
		BookingSessionId: session.Id,
		Message:          message,
	}

	if botSendMsgErr := s.bot.SendMsg(bookingMessage); botSendMsgErr != nil {
		return exception.Wrap(
			"s.bot.SendMsg",
			"pick-service-command-telegram-service.go",
			botSendMsgErr,
		)
	}

	return nil
}
