package service

import (
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/constants"
	"github.com/adriein/hastypal/internal/hastypal/helper"
	"github.com/adriein/hastypal/internal/hastypal/types"
	"net/url"
	"strings"
	"time"
)

type TelegramHoursCommandService struct {
	bot        *TelegramBot
	repository types.Repository[types.BookingSession]
}

func NewTelegramHoursCommandService(
	bot *TelegramBot,
	repository types.Repository[types.BookingSession],
) *TelegramHoursCommandService {
	return &TelegramHoursCommandService{
		bot:        bot,
		repository: repository,
	}
}

func (s *TelegramHoursCommandService) Execute(business types.Business, update types.TelegramUpdate) error {
	if ackErr := s.ackToTelegramClient(update.CallbackQuery.Id); ackErr != nil {
		return ackErr
	}

	var markdownText strings.Builder

	location, loadLocationErr := time.LoadLocation("Europe/Madrid")

	if loadLocationErr != nil {
		return types.ApiError{
			Msg:      loadLocationErr.Error(),
			Function: "Execute -> time.LoadLocation()",
			File:     "telegram-dates-command.go",
			Values:   []string{"Europe/Madrid"},
		}
	}

	time.Local = location

	parsedUrl, parseUrlErr := url.Parse(update.CallbackQuery.Data)

	if parseUrlErr != nil {
		return types.ApiError{
			Msg:      parseUrlErr.Error(),
			Function: "Execute -> url.Parse()",
			File:     "telegram-hours-command.go",
			Values:   []string{update.CallbackQuery.Data},
		}
	}

	queryParams := parsedUrl.Query()

	stringSelectedDate := fmt.Sprintf("%s %s", queryParams.Get("date"), "07:00:00")
	sessionId := queryParams.Get("session")

	selectedDate, timeParseErr := time.Parse(time.DateTime, stringSelectedDate)

	if timeParseErr != nil {
		return types.ApiError{
			Msg:      timeParseErr.Error(),
			Function: "Execute -> time.Parse()",
			File:     "telegram-hours-command.go",
			Values:   []string{stringSelectedDate},
		}
	}

	session, getSessionErr := s.getCurrentSession(sessionId)

	if getSessionErr != nil {
		return getSessionErr
	}

	if invalidSession := session.EnsureIsValid(); invalidSession != nil {
		return invalidSession
	}

	if updateSessionErr := s.updateSession(session, stringSelectedDate); updateSessionErr != nil {
		return updateSessionErr
	}

	dateParts := strings.Split(selectedDate.Format(time.RFC822), " ")

	day := dateParts[0]
	month := dateParts[1]

	welcome := fmt.Sprintf(
		"*![⌚️](tg://emoji?id=5368324170671202286) Las horas disponibles para el servicio %s el día %s son:*\n\n",
		"Corte de pelo y barba express 18€",
		fmt.Sprintf("%s %s", day, month),
	)

	processInstructions := "*Selecciona una hora y te escribiré un resumen para que puedas confirmar la reserva*\n\n"

	markdownText.WriteString(welcome)
	markdownText.WriteString(processInstructions)

	buttons := make([]types.KeyboardButton, 12)

	for i := 8; i <= len(buttons)+7; i++ {
		hour := fmt.Sprintf("%02d:00", i)

		buttons[i-8] = types.KeyboardButton{
			Text: fmt.Sprintf("%s", hour),
			CallbackData: fmt.Sprintf(
				"/confirmation?session=%s&hour=%s",
				sessionId,
				hour,
			),
		}
	}

	array := helper.NewArrayHelper[types.KeyboardButton]()

	inlineKeyboard := array.Chunk(buttons, 3)

	message := types.SendTelegramMessage{
		ChatId:         update.CallbackQuery.From.Id,
		Text:           markdownText.String(),
		ParseMode:      constants.TelegramMarkdown,
		ProtectContent: true,
		ReplyMarkup:    types.ReplyMarkup{InlineKeyboard: inlineKeyboard},
	}

	if botSendMsgErr := s.bot.SendMsg(message); botSendMsgErr != nil {
		return botSendMsgErr
	}

	return nil
}

func (s *TelegramHoursCommandService) ackToTelegramClient(callbackQueryId string) error {
	return s.bot.AnswerCallbackQuery(types.AnswerCallbackQuery{CallbackQueryId: callbackQueryId})
}

func (s *TelegramHoursCommandService) getCurrentSession(sessionId string) (types.BookingSession, error) {
	filter := types.Filter{
		Name:    "id",
		Operand: constants.Equal,
		Value:   sessionId,
	}

	criteria := types.Criteria{Filters: []types.Filter{filter}}

	session, findOneErr := s.repository.FindOne(criteria)

	if findOneErr != nil {
		return types.BookingSession{}, findOneErr
	}

	return session, nil
}

func (s *TelegramHoursCommandService) updateSession(actualSession types.BookingSession, date string) error {
	updatedSession := types.BookingSession{
		Id:         actualSession.Id,
		BusinessId: actualSession.BusinessId,
		ChatId:     actualSession.ChatId,
		ServiceId:  actualSession.ServiceId,
		Date:       date,
		Hour:       "",
		CreatedAt:  time.Now().Format(time.DateTime), //We refresh the created at on purpose
		Ttl:        actualSession.Ttl,
	}

	reflection := helper.NewReflectionHelper[types.BookingSession]()

	mergedSession := reflection.Merge(actualSession, updatedSession)

	if err := s.repository.Update(mergedSession); err != nil {
		return err
	}

	return nil
}
