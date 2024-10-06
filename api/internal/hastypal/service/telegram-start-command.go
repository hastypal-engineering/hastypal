package service

import (
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
	return nil
}
