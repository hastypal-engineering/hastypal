package service

import (
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/constants"
	"github.com/adriein/hastypal/internal/hastypal/helper"
	"github.com/adriein/hastypal/internal/hastypal/types"
	"strings"
	"time"
)

type TelegramStartCommandService struct {
	bot        *TelegramBot
	repository types.Repository[types.BookingSession]
}

func NewTelegramStartCommandService(
	bot *TelegramBot,
	repository types.Repository[types.BookingSession],
) *TelegramStartCommandService {
	return &TelegramStartCommandService{
		bot:        bot,
		repository: repository,
	}
}

func (s *TelegramStartCommandService) Execute(business types.Business, update types.TelegramUpdate) error {
	var markdownText strings.Builder

	uuidHelper := helper.NewUuidHelper()

	sessionId := uuidHelper.GenerateShort()

	session := types.BookingSession{
		Id:         sessionId,
		BusinessId: business.Id,
		ChatId:     update.Message.Chat.Id,
		ServiceId:  "",
		Date:       "",
		Hour:       "",
		CreatedAt:  time.Now().Format(time.DateTime),
		Ttl:        time.Minute.Milliseconds() * 5,
	}

	if err := s.repository.Save(session); err != nil {
		return err
	}

	welcome := fmt.Sprintf(
		"*Hola %s ![👋](tg://emoji?id=5368324170671202286), soy HastypalBot el ayudante de %s\\.*\n\n",
		update.Message.From.FirstName,
		"Hastypal Business Test",
	)

	services := []string{
		"Corte de pelo y barba express 18€",
		"Corte de pelo y barba premium 22€",
	}

	emoji := "![🔸](tg://emoji?id=5368324170671202286)"

	markdownText.WriteString(welcome)
	markdownText.WriteString("*Te muestro a continuación los servicios que ofrecemos:*\n\n")

	buttons := make([]types.KeyboardButton, len(services))

	for i, service := range services {
		markdownText.WriteString(fmt.Sprintf("%s %s\n\n", emoji, service))

		buttons[i] = types.KeyboardButton{
			Text:         fmt.Sprintf("%s 📅", services[i]),
			CallbackData: fmt.Sprintf("/dates?session=%s&service=%s", sessionId, "test-short"),
		}
	}

	array := helper.NewArrayHelper[types.KeyboardButton]()

	inlineKeyboard := array.Chunk(buttons, 1)

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
