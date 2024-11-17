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

type TelegramDatesCommandService struct {
	bot        *TelegramBot
	repository types2.Repository[types2.BookingSession]
}

func NewTelegramDatesCommandService(
	bot *TelegramBot,
	repository types2.Repository[types2.BookingSession],
) *TelegramDatesCommandService {
	return &TelegramDatesCommandService{
		bot:        bot,
		repository: repository,
	}
}

func (s *TelegramDatesCommandService) Execute(business types2.Business, update types2.TelegramUpdate) error {
	if ackErr := s.ackToTelegramClient(update.CallbackQuery.Id); ackErr != nil {
		return ackErr
	}

	var markdownText strings.Builder

	parsedUrl, parseUrlErr := url.Parse(update.CallbackQuery.Data)

	if parseUrlErr != nil {
		return types2.ApiError{
			Msg:      parseUrlErr.Error(),
			Function: "Execute -> url.Parse()",
			File:     "telegram-dates-command.go",
			Values:   []string{update.CallbackQuery.Data},
		}
	}

	queryParams := parsedUrl.Query()

	sessionId := queryParams.Get("session")
	serviceId := queryParams.Get("service")

	session, getSessionErr := s.getCurrentSession(sessionId)

	if getSessionErr != nil {
		return getSessionErr
	}

	if invalidSession := session.EnsureIsValid(); invalidSession != nil {
		return invalidSession
	}

	if updateSessionErr := s.updateSession(session, serviceId); updateSessionErr != nil {
		return updateSessionErr
	}

	commandInformation := fmt.Sprintf(
		"%s tiene disponibles para:\n\n![🔸](tg://emoji?id=5368324170671202286) %s\n\n",
		"Hastypal Business Test",
		"Corte de pelo y barba express 18€",
	)

	processInstructions := "*Selecciona un día para ver las horas disponibles:*\n\n"

	markdownText.WriteString("![📅](tg://emoji?id=5368324170671202286) A continuación puedes ver las fechas que ")
	markdownText.WriteString(commandInformation)
	markdownText.WriteString(processInstructions)

	location, loadLocationErr := time.LoadLocation("Europe/Madrid")

	if loadLocationErr != nil {
		return types2.ApiError{
			Msg:      loadLocationErr.Error(),
			Function: "Execute -> time.LoadLocation()",
			File:     "telegram-dates-command.go",
			Values:   []string{"Europe/Madrid"},
		}
	}

	time.Local = location

	today := time.Now()

	buttons := make([]types2.KeyboardButton, 15)

	for i := 0; i < 15; i++ {
		newDate := today.AddDate(0, 0, i)

		dateParts := strings.Split(newDate.Format(time.RFC822), " ")

		day := dateParts[0]
		month := dateParts[1]

		buttons[i] = types2.KeyboardButton{
			Text:         fmt.Sprintf("%s %s", day, month),
			CallbackData: fmt.Sprintf("/hours?session=%s&date=%s", sessionId, newDate.Format(time.DateOnly)),
		}
	}

	array := helper2.NewArrayHelper[types2.KeyboardButton]()

	inlineKeyboard := array.Chunk(buttons, 5)

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

func (s *TelegramDatesCommandService) getCurrentSession(sessionId string) (types2.BookingSession, error) {
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

func (s *TelegramDatesCommandService) ackToTelegramClient(callbackQueryId string) error {
	return s.bot.AnswerCallbackQuery(types2.AnswerCallbackQuery{CallbackQueryId: callbackQueryId})
}

func (s *TelegramDatesCommandService) updateSession(actualSession types2.BookingSession, serviceId string) error {
	updatedSession := types2.BookingSession{
		Id:         actualSession.Id,
		BusinessId: actualSession.BusinessId,
		ChatId:     actualSession.ChatId,
		ServiceId:  serviceId,
		Date:       "",
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
