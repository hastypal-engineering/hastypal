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

type TelegramDatesCommandService struct {
	bot        *TelegramBot
	repository types.Repository[types.BookingSession]
}

func NewTelegramDatesCommandService(
	bot *TelegramBot,
	repository types.Repository[types.BookingSession],
) *TelegramDatesCommandService {
	return &TelegramDatesCommandService{
		bot:        bot,
		repository: repository,
	}
}

func (s *TelegramDatesCommandService) Execute(business types.Business, update types.TelegramUpdate) error {
	answerCbErr := s.bot.AnswerCallbackQuery(types.AnswerCallbackQuery{CallbackQueryId: update.CallbackQuery.Id})

	if answerCbErr != nil {
		return answerCbErr
	}

	var markdownText strings.Builder

	parsedUrl, parseUrlErr := url.Parse(update.CallbackQuery.Data)

	if parseUrlErr != nil {
		return types.ApiError{
			Msg:      parseUrlErr.Error(),
			Function: "Execute -> url.Parse()",
			File:     "telegram-dates-command.go",
			Values:   []string{update.CallbackQuery.Data},
		}
	}

	queryParams := parsedUrl.Query()

	sessionId := queryParams.Get("session")

	filter := types.Filter{
		Name:    "id",
		Operand: constants.Equal,
		Value:   sessionId,
	}

	criteria := types.Criteria{Filters: []types.Filter{filter}}

	session, findOneErr := s.repository.FindOne(criteria)

	if findOneErr != nil {
		return findOneErr
	}

	updatedSession := types.BookingSession{
		Id:         sessionId,
		BusinessId: session.BusinessId,
		ChatId:     session.ChatId,
		ServiceId:  session.ServiceId,
		Date:       "",
		Hour:       "",
		CreatedAt:  "",
		Ttl:        0,
	}

	reflection := helper.NewReflectionHelper[types.BookingSession]()

	mergedSession := reflection.Merge(session, updatedSession)

	if err := s.repository.Update(mergedSession); err != nil {
		return err
	}

	commandInformation := fmt.Sprintf(
		"%s tiene disponibles para:*\n\n![🔸](tg://emoji?id=5368324170671202286) %s\n\n",
		"Hastypal Business Test",
		"Corte de pelo y barba express 18€",
	)

	processInstructions := "*Selecciona un día para ver las horas disponibles:*\n\n"

	markdownText.WriteString("![📅](tg://emoji?id=5368324170671202286) *A continuación puedes ver las fechas que ")
	markdownText.WriteString(commandInformation)
	markdownText.WriteString(processInstructions)

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

	today := time.Now()

	buttons := make([]types.KeyboardButton, 15)

	for i := 0; i < 15; i++ {
		newDate := today.AddDate(0, 0, i)

		dateParts := strings.Split(newDate.Format(time.RFC822), " ")

		day := dateParts[0]
		month := dateParts[1]

		buttons[i] = types.KeyboardButton{
			Text:         fmt.Sprintf("%s %s", day, month),
			CallbackData: fmt.Sprintf("/hours?session=%s&date=%s", sessionId, newDate.Format(time.DateOnly)),
		}
	}

	array := helper.NewArrayHelper[types.KeyboardButton]()

	inlineKeyboard := array.Chunk(buttons, 5)

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
