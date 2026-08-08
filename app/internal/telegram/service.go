package telegram

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/adriein/hastypal/pkg/constants"
	"github.com/adriein/hastypal/pkg/helper/array"
	"github.com/adriein/hastypal/pkg/helper/reflection"
	"github.com/rotisserie/eris"
)

type TelegramService interface {
	HandleMessage(update TelegramUpdate) error
}

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

// TELEGRAM UPDATE HANDLING

func (s *Service) HandleMessage(update TelegramUpdate) error {
	if reflection.HasField(update, constants.TelegramMessageField) {
		if err := s.resolveBotCommand(update); err != nil {
			eris.Wrap(err, "Error resolving bot command")
		}

		return nil
	}

	if err := s.resolveCallbackQueryCommand(update); err != nil {
		eris.Wrap(err, "Error resolving callback query command")
	}

	return nil
}

func (s *Service) resolveBotCommand(update TelegramUpdate) error {
	if !reflection.HasField(update, constants.TelegramMessageField) {
		return nil
	}

	txtArr := strings.Split(update.Message.Text, " ")

	switch txtArr[0] {
	case constants.StartCommand:
		return nil
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

func (s *Service) resolveCallbackQueryCommand(update TelegramUpdate) error {
	if !reflection.HasField(update, constants.TelegramCallbackQueryField) {
		return nil
	}

	url, err := url.Parse(update.CallbackQuery.Data)

	if err != nil {
		return eris.Wrap(err, "Error parsing the callback query url")
	}

	switch url.Path {
	case constants.StartCommand:
		return nil
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

// START TELEGRAM CONVERSATION

func (s *Service) startConversation(update TelegramUpdate) error {
	var markdownText strings.Builder

	businessId := strings.ReplaceAll(update.Message.Text, "/start ", "")

	business, err := s.getBusiness(businessId)

	if err != nil {
		return eris.Wrap(err, "Error fetching business")
	}

	session, err := s.createSession(business.Id, update.Message.Chat.Id)

	if err != nil {
		return eris.Wrap(err, "Error creating session")
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
			CallbackData: fmt.Sprintf("/dates?session=%s&service=%s&page=0", session.Id, "test-short"),
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
		BookingSessionId: session.Id,
		Message:          message,
	}

	if err := s.bot.SendMsg(bookingMessage); err != nil {
		return eris.Wrap(err, "Error sending message to telegram")
	}

	return nil
}
