package telegram

import (
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/shared/constants"
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	"github.com/adriein/hastypal/internal/hastypal/shared/service"
	"github.com/adriein/hastypal/internal/hastypal/shared/translation"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"net/url"
	"strings"
	"time"
)

type PickServiceCommandTelegramService struct {
	bot               *service.TelegramBot
	sessionRepository types.Repository[types.BookingSession]
	bookingRepository types.Repository[types.Booking]
	translation       *translation.Translation
}

func NewPickServiceCommandTelegramService(
	bot *service.TelegramBot,
	sessionRepository types.Repository[types.BookingSession],
) *PickServiceCommandTelegramService {
	return &PickServiceCommandTelegramService{
		bot:               bot,
		sessionRepository: sessionRepository,
	}
}

func (s *PickServiceCommandTelegramService) Execute(update types.TelegramUpdate) error {
	if ackErr := s.ackToTelegramClient(update.CallbackQuery.Id); ackErr != nil {
		return exception.Wrap(
			"s.ackToTelegramClient",
			"pick-service-command-telegram-service.go",
			ackErr,
		)
	}

	var markdownText strings.Builder

	parsedUrl, parseUrlErr := url.Parse(update.CallbackQuery.Data)

	if parseUrlErr != nil {
		return exception.New(parseUrlErr.Error()).
			Trace("url.Parse(update.CallbackQuery.Data)", "pick-service-command-telegram-service.go").
			WithValues([]string{update.CallbackQuery.Data})
	}

	queryParams := parsedUrl.Query()

	sessionId := queryParams.Get("sessionId")

	session, getSessionErr := s.getCurrentSession(sessionId)

	if getSessionErr != nil {
		return exception.Wrap(
			"s.getCurrentSession",
			"pick-service-command-telegram-service.go",
			getSessionErr,
		)
	}

	if invalidSession := session.EnsureIsValid(); invalidSession != nil {
		return exception.Wrap(
			"session.EnsureIsValid",
			"pick-service-command-telegram-service.go",
			invalidSession,
		)
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
			CallbackData: fmt.Sprintf("/dates?session=%s&service=%s", session.Id, "test-short"),
		}
	}

	array := helper.NewArrayHelper[types.KeyboardButton]()

	inlineKeyboard := array.Chunk(buttons, 1)

	message := types.SendTelegramMessage{
		ChatId:         update.CallbackQuery.From.Id,
		Text:           markdownText.String(),
		ParseMode:      constants.TelegramMarkdown,
		ProtectContent: true,
		ReplyMarkup:    types.ReplyMarkup{InlineKeyboard: inlineKeyboard},
	}

	if botSendMsgErr := s.bot.SendMsg(message); botSendMsgErr != nil {
		return exception.Wrap(
			"s.bot.SendMsg",
			"pick-service-command-telegram-service.go",
			botSendMsgErr,
		)
	}

	return nil
}

func (s *PickServiceCommandTelegramService) getCurrentSession(sessionId string) (types.BookingSession, error) {
	filter := types.Filter{
		Name:    "id",
		Operand: constants.Equal,
		Value:   sessionId,
	}

	criteria := types.Criteria{Filters: []types.Filter{filter}}

	session, findOneErr := s.sessionRepository.FindOne(criteria)

	if findOneErr != nil {
		return types.BookingSession{}, exception.Wrap(
			"s.repository.FindOne",
			"pick-service-command-telegram-service.go",
			findOneErr,
		)
	}

	return session, nil
}

func (s *PickServiceCommandTelegramService) ackToTelegramClient(callbackQueryId string) error {
	return s.bot.AnswerCallbackQuery(types.AnswerCallbackQuery{CallbackQueryId: callbackQueryId})
}

func (s *PickServiceCommandTelegramService) updateSession(actualSession types.BookingSession) error {
	updatedSession := types.BookingSession{
		Id:         actualSession.Id,
		BusinessId: actualSession.BusinessId,
		ChatId:     actualSession.ChatId,
		ServiceId:  "",
		Date:       "",
		Hour:       "",
		CreatedAt:  actualSession.CreatedAt,
		UpdatedAt:  time.Now().UTC().Format(time.DateTime), //We refresh the created at on purpose
		Ttl:        actualSession.Ttl,
	}

	if err := s.sessionRepository.Update(updatedSession); err != nil {
		return exception.Wrap(
			"s.repository.Update",
			"pick-service-command-telegram-service.go",
			err,
		)
	}

	return nil
}
