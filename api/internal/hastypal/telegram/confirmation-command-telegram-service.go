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

type ConfirmationCommandTelegramService struct {
	bot         *service.TelegramBot
	repository  types.Repository[types.BookingSession]
	translation *translation.Translation
}

func NewConfirmationCommandTelegramService(
	bot *service.TelegramBot,
	repository types.Repository[types.BookingSession],
	translation *translation.Translation,
) *ConfirmationCommandTelegramService {
	return &ConfirmationCommandTelegramService{
		bot:         bot,
		repository:  repository,
		translation: translation,
	}
}

func (s *ConfirmationCommandTelegramService) Execute(update types.TelegramUpdate) error {
	if ackErr := s.ackToTelegramClient(update.CallbackQuery.Id); ackErr != nil {
		return exception.Wrap(
			"s.ackToTelegramClient",
			"confirmation-command-telegram-service.go",
			ackErr,
		)
	}

	var markdownText strings.Builder

	location, loadLocationErr := time.LoadLocation("Europe/Madrid")

	if loadLocationErr != nil {
		return exception.New(loadLocationErr.Error()).
			Trace("time.LoadLocation", "confirmation-command-telegram-service.go").
			WithValues([]string{"Europe/Madrid"})
	}

	time.Local = location

	parsedUrl, parseUrlErr := url.Parse(update.CallbackQuery.Data)

	if parseUrlErr != nil {
		return exception.New(parseUrlErr.Error()).
			Trace("url.Parse(update.CallbackQuery.Data)", "confirmation-command-telegram-service.go").
			WithValues([]string{update.CallbackQuery.Data})
	}

	queryParams := parsedUrl.Query()
	sessionId := queryParams.Get("session")
	hour := queryParams.Get("hour")

	session, getSessionErr := s.getCurrentSession(sessionId)

	if getSessionErr != nil {
		return exception.Wrap(
			"s.getCurrentSession",
			"confirmation-command-telegram-service.go",
			getSessionErr,
		)
	}

	if invalidSession := session.EnsureIsValid(); invalidSession != nil {
		return exception.Wrap(
			"session.EnsureIsValid",
			"confirmation-command-telegram-service.go",
			invalidSession,
		)
	}

	if updateSessionErr := s.updateSession(session, hour); updateSessionErr != nil {
		return exception.Wrap(
			"s.updateSession",
			"confirmation-command-telegram-service.go",
			updateSessionErr,
		)
	}

	selectedDate, timeParseErr := time.Parse(time.DateTime, session.Date)

	if timeParseErr != nil {
		return exception.New(timeParseErr.Error()).
			Trace("time.Parse(time.DateTime, session.Date)", "confirmation-command-telegram-service.go").
			WithValues([]string{session.Date})
	}

	dateParts := strings.Split(selectedDate.Format(time.RFC822), " ")

	day := dateParts[0]
	month := s.translation.GetSpanishMonthShortForm(selectedDate.Month())

	welcome := "![🙂](tg://emoji?id=5368324170671202286) Último paso te lo prometo\\! Confirma que todo esta correcto:\n\n"

	bookedService := fmt.Sprintf(
		"![🟢](tg://emoji?id=5368324170671202286) %s\n\n",
		"Corte de pelo y barba express 18€",
	)

	date := fmt.Sprintf("![📅](tg://emoji?id=5368324170671202286) %s %s\n\n", day, month)

	hourMarkdown := fmt.Sprintf(
		"![⌚️](tg://emoji?id=5368324170671202286) %s\n\n",
		hour,
	)

	processInstructions := "*Pulsa confirmar si todo es correcto o cancelar de lo contrario*\n\n"

	markdownText.WriteString(welcome)
	markdownText.WriteString(bookedService)
	markdownText.WriteString(date)
	markdownText.WriteString(hourMarkdown)
	markdownText.WriteString(processInstructions)

	buttons := make([]types.KeyboardButton, 2)

	confirmButton := types.KeyboardButton{
		Text:         "Confirmar",
		CallbackData: fmt.Sprintf("/book?session=%s", sessionId),
	}

	buttons = append(buttons, confirmButton)

	backButton := types.KeyboardButton{
		Text:         "Atrás",
		CallbackData: fmt.Sprintf("/hours?session=%s&date=%s", sessionId, selectedDate.Format(time.DateOnly)),
	}

	buttons = append(buttons, backButton)

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
			"confirmation-command-telegram-service.go",
			botSendMsgErr,
		)
	}

	return nil
}

func (s *ConfirmationCommandTelegramService) ackToTelegramClient(callbackQueryId string) error {
	return s.bot.AnswerCallbackQuery(types.AnswerCallbackQuery{CallbackQueryId: callbackQueryId})
}

func (s *ConfirmationCommandTelegramService) getCurrentSession(sessionId string) (types.BookingSession, error) {
	filter := types.Filter{
		Name:    "id",
		Operand: constants.Equal,
		Value:   sessionId,
	}

	criteria := types.Criteria{Filters: []types.Filter{filter}}

	session, findOneErr := s.repository.FindOne(criteria)

	if findOneErr != nil {
		return types.BookingSession{}, exception.Wrap(
			"s.repository.FindOne",
			"confirmation-command-telegram-service.go",
			findOneErr,
		)
	}

	return session, nil
}

func (s *ConfirmationCommandTelegramService) updateSession(actualSession types.BookingSession, hour string) error {
	updatedSession := types.BookingSession{
		Id:         actualSession.Id,
		BusinessId: actualSession.BusinessId,
		ChatId:     actualSession.ChatId,
		ServiceId:  actualSession.ServiceId,
		Date:       actualSession.Date,
		Hour:       hour,
		CreatedAt:  actualSession.CreatedAt,
		UpdatedAt:  time.Now().UTC().Format(time.DateTime), //We refresh the created at on purpose
		Ttl:        actualSession.Ttl,
	}

	if err := s.repository.Update(updatedSession); err != nil {
		return exception.Wrap(
			"s.repository.Update",
			"confirmation-command-telegram-service.go",
			err,
		)
	}

	return nil
}
