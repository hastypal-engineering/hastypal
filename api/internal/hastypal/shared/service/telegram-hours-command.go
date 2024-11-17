package service

import (
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/shared/constants"
	helper2 "github.com/adriein/hastypal/internal/hastypal/shared/helper"
	types2 "github.com/adriein/hastypal/internal/hastypal/shared/types"
	"net/url"
	"strings"
	"time"
)

type TelegramHoursCommandService struct {
	bot        *TelegramBot
	repository types2.Repository[types2.BookingSession]
}

func NewTelegramHoursCommandService(
	bot *TelegramBot,
	repository types2.Repository[types2.BookingSession],
) *TelegramHoursCommandService {
	return &TelegramHoursCommandService{
		bot:        bot,
		repository: repository,
	}
}

func (s *TelegramHoursCommandService) Execute(business types2.Business, update types2.TelegramUpdate) error {
	if ackErr := s.ackToTelegramClient(update.CallbackQuery.Id); ackErr != nil {
		return ackErr
	}

	var markdownText strings.Builder

	location, loadLocationErr := time.LoadLocation("Europe/Madrid")

	if loadLocationErr != nil {
		return types2.ApiError{
			Msg:      loadLocationErr.Error(),
			Function: "Execute -> time.LoadLocation()",
			File:     "telegram-hours-command.go",
			Values:   []string{"Europe/Madrid"},
		}
	}

	time.Local = location

	parsedUrl, parseUrlErr := url.Parse(update.CallbackQuery.Data)

	if parseUrlErr != nil {
		return types2.ApiError{
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
		return types2.ApiError{
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
		"![⌚️](tg://emoji?id=5368324170671202286) Las horas disponibles para:\n\n",
	)

	service := fmt.Sprintf(
		"![🔸](tg://emoji?id=5368324170671202286) %s\n\n",
		"Corte de pelo y barba express 18€",
	)

	date := fmt.Sprintf(
		"![📅](tg://emoji?id=5368324170671202286) %s\n\n",
		fmt.Sprintf("%s %s", day, month),
	)

	processInstructions := "*Selecciona una hora y te escribiré un resumen para que puedas confirmar la reserva*\n\n"

	markdownText.WriteString(welcome)
	markdownText.WriteString(service)
	markdownText.WriteString(date)
	markdownText.WriteString(processInstructions)

	buttons := make([]types2.KeyboardButton, 12)

	for i := 8; i <= len(buttons)+7; i++ {
		hour := fmt.Sprintf("%02d:00", i)

		buttons[i-8] = types2.KeyboardButton{
			Text: fmt.Sprintf("%s", hour),
			CallbackData: fmt.Sprintf(
				"/confirmation?session=%s&hour=%s",
				sessionId,
				hour,
			),
		}
	}

	array := helper2.NewArrayHelper[types2.KeyboardButton]()

	inlineKeyboard := array.Chunk(buttons, 3)

	message := types2.SendTelegramMessage{
		ChatId:         update.CallbackQuery.From.Id,
		Text:           markdownText.String(),
		ParseMode:      constants.TelegramMarkdown,
		ProtectContent: true,
		ReplyMarkup:    types2.ReplyMarkup{InlineKeyboard: inlineKeyboard},
	}

	if botSendMsgErr := s.bot.SendMsg(message); botSendMsgErr != nil {
		return botSendMsgErr
	}

	return nil
}

func (s *TelegramHoursCommandService) ackToTelegramClient(callbackQueryId string) error {
	return s.bot.AnswerCallbackQuery(types2.AnswerCallbackQuery{CallbackQueryId: callbackQueryId})
}

func (s *TelegramHoursCommandService) getCurrentSession(sessionId string) (types2.BookingSession, error) {
	filter := types2.Filter{
		Name:    "id",
		Operand: constants.Equal,
		Value:   sessionId,
	}

	criteria := types2.Criteria{Filters: []types2.Filter{filter}}

	session, findOneErr := s.repository.FindOne(criteria)

	if findOneErr != nil {
		return types2.BookingSession{}, findOneErr
	}

	return session, nil
}

func (s *TelegramHoursCommandService) updateSession(actualSession types2.BookingSession, date string) error {
	updatedSession := types2.BookingSession{
		Id:         actualSession.Id,
		BusinessId: actualSession.BusinessId,
		ChatId:     actualSession.ChatId,
		ServiceId:  actualSession.ServiceId,
		Date:       date,
		Hour:       "",
		CreatedAt:  actualSession.CreatedAt,
		UpdatedAt:  time.Now().Format(time.DateTime), //We refresh the created at on purpose
		Ttl:        actualSession.Ttl,
	}

	reflection := helper2.NewReflectionHelper[types2.BookingSession]()

	mergedSession := reflection.Merge(actualSession, updatedSession)

	if err := s.repository.Update(mergedSession); err != nil {
		return err
	}

	return nil
}
