package service

import (
	"github.com/adriein/hastypal/internal/hastypal/constants"
	"github.com/adriein/hastypal/internal/hastypal/types"
)

type TelegramStartCommandService struct {
	bot *TelegramBot
}

func NewTelegramStartCommandService(
	bot *TelegramBot,
) *TelegramStartCommandService {
	return &TelegramStartCommandService{
		bot: bot,
	}
}

func (s *TelegramStartCommandService) Execute(business types.Business, update types.TelegramUpdate) error {
	markdownText := "*Los servicios ofrecidos por Hastypal Business test son:*\n\n" +
		"![🔸](tg://emoji?id=5368324170671202286) Corte de pelo y barba express 18€\n\n" +
		"![🔸](tg://emoji?id=5368324170671202286) Corte de pelo y barba premium 22€\n\n"

	service1 := make([]types.KeyboardButton, 1)
	service1[0] = types.KeyboardButton{Text: "Corte de pelo y barba express 18€", CallbackData: "/book 1"}
	service2 := make([]types.KeyboardButton, 1)
	service2[0] = types.KeyboardButton{Text: "Corte de pelo y barba premium 22€", CallbackData: "/book 2"}

	inlineKeyboard := make([][]types.KeyboardButton, 2)

	inlineKeyboard[0] = service1
	inlineKeyboard[1] = service2

	message := types.SendTelegramMessage{
		ChatId:         update.Message.Chat.Id,
		Text:           markdownText,
		ParseMode:      constants.TelegramMarkdown,
		ProtectContent: true,
		ReplyMarkup:    types.ReplyMarkup{InlineKeyboard: inlineKeyboard},
	}
	err := s.bot.SendMsg(message)

	if err != nil {
		return err
	}

	return nil
}
