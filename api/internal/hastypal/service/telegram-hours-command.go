package service

import (
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/constants"
	"github.com/adriein/hastypal/internal/hastypal/helper"
	"github.com/adriein/hastypal/internal/hastypal/types"
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

	welcome := fmt.Sprintf(
		"*![⌚️](tg://emoji?id=5368324170671202286) Las horas disponibles para el servicio %s en día %s son:*\n\n",
		"Corte de pelo y barba express 18€",
		"aaaa",
	)

	processInstructions := "*Selecciona una hora y te escribiré un resumen para que puedas confirmar*\n\n"

	markdownText.WriteString(welcome)
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
			CallbackData: fmt.Sprintf("/hours %s", newDate.Format(time.DateTime)),
		}
	}

	inlineKeyboard := helper.Chunk[types.KeyboardButton](buttons, 5)

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
