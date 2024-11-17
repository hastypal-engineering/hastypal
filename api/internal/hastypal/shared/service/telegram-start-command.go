package service

import (
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/shared/constants"
	helper2 "github.com/adriein/hastypal/internal/hastypal/shared/helper"
	types2 "github.com/adriein/hastypal/internal/hastypal/shared/types"
	"strings"
	"time"
)

type TelegramStartCommandService struct {
	bot        *TelegramBot
	repository types2.Repository[types2.BookingSession]
}

func NewTelegramStartCommandService(
	bot *TelegramBot,
	repository types2.Repository[types2.BookingSession],
) *TelegramStartCommandService {
	return &TelegramStartCommandService{
		bot:        bot,
		repository: repository,
	}
}

func (s *TelegramStartCommandService) Execute(business types2.Business, update types2.TelegramUpdate) error {
	var markdownText strings.Builder

	session, createSessionErr := s.createSession(business.Id, update.Message.Chat.Id)

	if createSessionErr != nil {
		return createSessionErr
	}

	welcome := fmt.Sprintf(
		"Hola %s ![👋](tg://emoji?id=5368324170671202286), soy HastypalBot el ayudante de %s\\.\n\n",
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

	buttons := make([]types2.KeyboardButton, len(services))

	for i, service := range services {
		markdownText.WriteString(fmt.Sprintf("%s %s\n\n", emoji, service))

		buttons[i] = types2.KeyboardButton{
			Text:         fmt.Sprintf("%s 📅", services[i]),
			CallbackData: fmt.Sprintf("/dates?session=%s&service=%s", session.Id, "test-short"),
		}
	}

	array := helper2.NewArrayHelper[types2.KeyboardButton]()

	inlineKeyboard := array.Chunk(buttons, 1)

	message := types2.SendTelegramMessage{
		ChatId:         update.Message.Chat.Id,
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

func (s *TelegramStartCommandService) createSession(businessId string, chatId int) (types2.BookingSession, error) {
	uuidHelper := helper2.NewUuidHelper()

	sessionId := uuidHelper.GenerateShort()

	session := types2.BookingSession{
		Id:         sessionId,
		BusinessId: businessId,
		ChatId:     chatId,
		ServiceId:  "",
		Date:       "",
		Hour:       "",
		CreatedAt:  time.Now().UTC().Format(time.DateTime),
		UpdatedAt:  time.Now().UTC().Format(time.DateTime),
		Ttl:        time.Minute.Milliseconds() * 5,
	}

	if err := s.repository.Save(session); err != nil {
		return types2.BookingSession{}, err
	}

	return session, nil
}