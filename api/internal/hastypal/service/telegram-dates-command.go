package service

import (
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/constants"
	"github.com/adriein/hastypal/internal/hastypal/helper"
	"github.com/adriein/hastypal/internal/hastypal/types"
	"strings"
	"time"
)

type TelegramDatesCommandService struct {
	bot *TelegramBot
}

func NewTelegramDatesCommandService(
	bot *TelegramBot,
) *TelegramDatesCommandService {
	return &TelegramDatesCommandService{
		bot: bot,
	}
}

func (s *TelegramDatesCommandService) Execute(business types.Business, update types.TelegramUpdate) error {
	var markdownText strings.Builder

	welcome := fmt.Sprintf(
		"*%s, las ![📅](tg://emoji?id=5368324170671202286) fechas ",
		update.CallbackQuery.From.FirstName,
	)

	commandInformation := fmt.Sprintf(
		"que %s tiene disponibles para el servicio:*\n\n![🔸](tg://emoji?id=5368324170671202286) %s",
		"Hastypal Business Test",
		"Corte de pelo y barba express 18€",
	)

	processInstructions := "*Selecciona un día y te respondere con las horas disponibles:*\n\n"

	markdownText.WriteString(welcome)
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
			CallbackData: fmt.Sprintf("hours %d", i),
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
