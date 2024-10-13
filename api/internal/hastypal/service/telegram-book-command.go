package service

import (
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/constants"
	"github.com/adriein/hastypal/internal/hastypal/types"
	"math"
	"strings"
	"time"
)

type TelegramBookCommandService struct {
	bot *TelegramBot
}

func NewTelegramBookCommandService(
	bot *TelegramBot,
) *TelegramBookCommandService {
	return &TelegramBookCommandService{
		bot: bot,
	}
}

func (s *TelegramBookCommandService) Execute(business types.Business, update types.TelegramUpdate) error {
	var markdownText strings.Builder

	welcome := fmt.Sprintf(
		"*%s, estas son las fechas ![📅](tg://emoji?id=5368324170671202286) ",
		update.Message.From.FirstName,
	)

	commandInformation := fmt.Sprintf(
		"que %s tiene disponibles para el servicio %s*\n\n",
		"Hastypal Business Test",
		"Corte de pelo y barba express 18€",
	)

	processInstructions := "*Selecciona un día y te respondere con las horas disponibles:*\n\n"

	markdownText.WriteString(welcome)
	markdownText.WriteString(commandInformation)
	markdownText.WriteString(processInstructions)

	today := time.Now()

	buttons := make([]types.KeyboardButton, 15)

	for i := 0; i < 15; i++ {
		newDate := today.AddDate(0, 0, i)

		buttons[i] = types.KeyboardButton{
			Text: fmt.Sprintf("%s",
				newDate.Weekday().String()),
			CallbackData: fmt.Sprintf("inspect %d", i),
		}
	}

	inlineKeyboard := s.chunkKeyboardButtons(buttons, 5)

	message := types.SendTelegramMessage{
		ChatId:         update.Message.Chat.Id,
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

func (s *TelegramBookCommandService) chunkKeyboardButtons(
	buttons []types.KeyboardButton,
	chunkSize int,
) [][]types.KeyboardButton {
	var chunked [][]types.KeyboardButton

	for i := 0; i < len(buttons); i += chunkSize {
		end := int(math.Min(float64(i+chunkSize), float64(len(buttons))))
		chunked = append(chunked, buttons[i:end])
	}

	return chunked
}
