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
	markdownText := "*Los servicios ofrecidos por Hastypal Business test son:*\n\nCorte de pelo y barba express 18€\n" +
		"Corte de pelo y barba premium 22€\n"

	message := types.SendTelegramMessage{
		ChatId:         update.Message.Chat.Id,
		Text:           markdownText,
		ParseMode:      constants.TelegramMarkdown,
		ProtectContent: true,
	}
	err := s.bot.SendMsg(message)

	if err != nil {
		return err
	}

	return nil
}
