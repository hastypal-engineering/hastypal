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

type TelegramConfirmationCommandService struct {
	bot        *TelegramBot
	repository types2.Repository[types2.BookingSession]
}

func NewTelegramConfirmationCommandService(
	bot *TelegramBot,
	repository types2.Repository[types2.BookingSession],
) *TelegramConfirmationCommandService {
	return &TelegramConfirmationCommandService{
		bot:        bot,
		repository: repository,
	}
}

func (s *TelegramConfirmationCommandService) Execute(business types2.Business, update types2.TelegramUpdate) error {
	if ackErr := s.ackToTelegramClient(update.CallbackQuery.Id); ackErr != nil {
		return ackErr
	}

	var markdownText strings.Builder

	location, loadLocationErr := time.LoadLocation("Europe/Madrid")

	if loadLocationErr != nil {
		return types2.ApiError{
			Msg:      loadLocationErr.Error(),
			Function: "Execute -> time.LoadLocation()",
			File:     "telegram-confirmation-command.go",
			Values:   []string{"Europe/Madrid"},
		}
	}

	time.Local = location

	parsedUrl, parseUrlErr := url.Parse(update.CallbackQuery.Data)

	if parseUrlErr != nil {
		return types2.ApiError{
			Msg:      parseUrlErr.Error(),
			Function: "Execute -> url.Parse()",
			File:     "telegram-confirmation-command.go",
			Values:   []string{update.CallbackQuery.Data},
		}
	}

	queryParams := parsedUrl.Query()
	sessionId := queryParams.Get("session")
	hour := queryParams.Get("hour")

	session, getSessionErr := s.getCurrentSession(sessionId)

	if getSessionErr != nil {
		return getSessionErr
	}

	if invalidSession := session.EnsureIsValid(); invalidSession != nil {
		return invalidSession
	}

	if updateSessionErr := s.updateSession(session, hour); updateSessionErr != nil {
		return updateSessionErr
	}

	selectedDate, timeParseErr := time.Parse(time.DateTime, session.Date)

	if timeParseErr != nil {
		return types2.ApiError{
			Msg:      timeParseErr.Error(),
			Function: "Execute -> time.Parse()",
			File:     "telegram-confirmation-command.go",
			Values:   []string{session.Date},
		}
	}

	dateParts := strings.Split(selectedDate.Format(time.RFC822), " ")

	day := dateParts[0]
	month := dateParts[1]

	welcome := "![🙂](tg://emoji?id=5368324170671202286) Último paso te lo prometo\\! Confirma que todo esta correcto:\n\n"

	service := fmt.Sprintf(
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
	markdownText.WriteString(service)
	markdownText.WriteString(date)
	markdownText.WriteString(hourMarkdown)
	markdownText.WriteString(processInstructions)

	buttons := make([]types2.KeyboardButton, 2)

	availableButtons := [2]string{"Confirmar", "Cancelar"}

	for i, text := range availableButtons {
		buttons[i] = types2.KeyboardButton{
			Text: text,
			CallbackData: fmt.Sprintf(
				"/book?session=%s",
				sessionId,
			),
		}
	}

	array := helper2.NewArrayHelper[types2.KeyboardButton]()

	inlineKeyboard := array.Chunk(buttons, 1)

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

func (s *TelegramConfirmationCommandService) ackToTelegramClient(callbackQueryId string) error {
	return s.bot.AnswerCallbackQuery(types2.AnswerCallbackQuery{CallbackQueryId: callbackQueryId})
}

func (s *TelegramConfirmationCommandService) getCurrentSession(sessionId string) (types2.BookingSession, error) {
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

func (s *TelegramConfirmationCommandService) updateSession(actualSession types2.BookingSession, hour string) error {
	updatedSession := types2.BookingSession{
		Id:         actualSession.Id,
		BusinessId: actualSession.BusinessId,
		ChatId:     actualSession.ChatId,
		ServiceId:  actualSession.ServiceId,
		Date:       actualSession.Date,
		Hour:       hour,
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
