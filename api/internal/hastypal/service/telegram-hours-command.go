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
	bot *TelegramBot
}

func NewTelegramHoursCommandService(
	bot *TelegramBot,
) *TelegramHoursCommandService {
	return &TelegramHoursCommandService{
		bot: bot,
	}
}

func (s *TelegramHoursCommandService) Execute(business types.Business, update types.TelegramUpdate) error {
	answerCbErr := s.bot.AnswerCallbackQuery(types.AnswerCallbackQuery{CallbackQueryId: update.CallbackQuery.Id})

	if answerCbErr != nil {
		return answerCbErr
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
	service := queryParams.Get("service")

	selectedDate, timeParseErr := time.Parse(time.DateTime, stringSelectedDate)

	if timeParseErr != nil {
		return types.ApiError{
			Msg:      timeParseErr.Error(),
			Function: "Execute -> time.Parse()",
			File:     "telegram-hours-command.go",
			Values:   []string{stringSelectedDate},
		}
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
				"/confirmation?service=%s&date=%s&hour=%s",
				service,
				selectedDate.Format(time.DateOnly),
				hour,
			),
		}
	}

	inlineKeyboard := helper.Chunk[types.KeyboardButton](buttons, 3)

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
