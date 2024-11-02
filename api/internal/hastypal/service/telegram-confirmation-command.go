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

type TelegramConfirmationCommandService struct {
	bot        *TelegramBot
	repository types.Repository[types.BookingSession]
}

func NewTelegramConfirmationCommandService(
	bot *TelegramBot,
	repository types.Repository[types.BookingSession],
) *TelegramConfirmationCommandService {
	return &TelegramConfirmationCommandService{
		bot:        bot,
		repository: repository,
	}
}

func (s *TelegramConfirmationCommandService) Execute(business types.Business, update types.TelegramUpdate) error {
	if ackErr := s.ackToTelegramClient(update.CallbackQuery.Id); ackErr != nil {
		return ackErr
	}

	var markdownText strings.Builder

	location, loadLocationErr := time.LoadLocation("Europe/Madrid")

	if loadLocationErr != nil {
		return types.ApiError{
			Msg:      loadLocationErr.Error(),
			Function: "Execute -> time.LoadLocation()",
			File:     "telegram-confirmation-command.go",
			Values:   []string{"Europe/Madrid"},
		}
	}

	time.Local = location

	parsedUrl, parseUrlErr := url.Parse(update.CallbackQuery.Data)

	if parseUrlErr != nil {
		return types.ApiError{
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
		return types.ApiError{
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

	buttons := make([]types.KeyboardButton, 2)

	availableButtons := [2]string{"Confirmar", "Cancelar"}

	for i, text := range availableButtons {
		buttons[i] = types.KeyboardButton{
			Text: text,
			CallbackData: fmt.Sprintf(
				"/book?session=%s",
				sessionId,
			),
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
		return botSendMsgErr
	}

	return nil
}

func (s *TelegramConfirmationCommandService) ackToTelegramClient(callbackQueryId string) error {
	return s.bot.AnswerCallbackQuery(types.AnswerCallbackQuery{CallbackQueryId: callbackQueryId})
}

func (s *TelegramConfirmationCommandService) getCurrentSession(sessionId string) (types.BookingSession, error) {
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

func (s *TelegramConfirmationCommandService) updateSession(actualSession types.BookingSession, hour string) error {
	updatedSession := types.BookingSession{
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

	reflection := helper.NewReflectionHelper[types.BookingSession]()

	mergedSession := reflection.Merge(actualSession, updatedSession)

	if err := s.repository.Update(mergedSession); err != nil {
		return err
	}

	return nil
}
