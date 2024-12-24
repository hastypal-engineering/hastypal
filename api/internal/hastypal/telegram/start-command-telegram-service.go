package telegram

import (
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/shared/constants"
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	"github.com/adriein/hastypal/internal/hastypal/shared/service"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"strings"
	"time"
)

type StartCommandTelegramService struct {
	bot                *service.TelegramBot
	sessionRepository  types.Repository[types.BookingSession]
	businessRepository types.Repository[types.Business]
}

func NewStartCommandTelegramService(
	bot *service.TelegramBot,
	sessionRepository types.Repository[types.BookingSession],
	businessRepository types.Repository[types.Business],
) *StartCommandTelegramService {
	return &StartCommandTelegramService{
		bot:                bot,
		sessionRepository:  sessionRepository,
		businessRepository: businessRepository,
	}
}

func (s *StartCommandTelegramService) Execute(update types.TelegramUpdate) error {
	var markdownText strings.Builder

	businessId := strings.ReplaceAll(update.Message.Text, "/start ", "")

	business, getBusinessErr := s.getBusiness(businessId)

	if getBusinessErr != nil {
		return exception.Wrap(
			"s.getBusiness",
			"start-command-telegram-service",
			getBusinessErr,
		)
	}

	session, createSessionErr := s.createSession(business.Id, update.Message.Chat.Id)

	if createSessionErr != nil {
		return exception.Wrap(
			"s.createSession",
			"start-command-telegram-service",
			createSessionErr,
		)
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

	buttons := make([]types.KeyboardButton, len(services))

	for i, serv := range services {
		markdownText.WriteString(fmt.Sprintf("%s %s\n\n", emoji, serv))

		buttons[i] = types.KeyboardButton{
			Text:         fmt.Sprintf("%s 📅", services[i]),
			CallbackData: fmt.Sprintf("/dates?session=%s&service=%s&page=0", session.Id, "test-short"),
		}
	}

	array := helper.NewArrayHelper[types.KeyboardButton]()

	inlineKeyboard := array.Chunk(buttons, 1)

	message := types.TelegramMessage{
		ChatId:         update.Message.Chat.Id,
		Text:           markdownText.String(),
		ParseMode:      constants.TelegramMarkdown,
		ProtectContent: true,
		ReplyMarkup:    types.ReplyMarkup{InlineKeyboard: inlineKeyboard},
	}

	bookingMessage := types.BookingTelegramMessage{
		BusinessName:     business.Name,
		BookingSessionId: session.Id,
		Message:          message,
	}

	if botSendMsgErr := s.bot.SendMsg(bookingMessage); botSendMsgErr != nil {
		return exception.Wrap(
			"s.bot.SendMsg",
			"start-command-telegram-service",
			botSendMsgErr,
		)
	}

	return nil
}

func (s *StartCommandTelegramService) getBusiness(businessId string) (types.Business, error) {
	filters := make([]types.Filter, 1)

	filters[0] = types.Filter{Name: "id", Operand: constants.Equal, Value: businessId}

	criteria := types.Criteria{Filters: filters}

	business, err := s.businessRepository.FindOne(criteria)

	if err != nil {
		return types.Business{}, exception.Wrap(
			"s.businessRepository.FindOne",
			"start-command-telegram-service",
			err,
		)
	}

	return business, nil
}

func (s *StartCommandTelegramService) createSession(businessId string, chatId int) (types.BookingSession, error) {
	uuidHelper := helper.NewUuidHelper()

	sessionId := uuidHelper.GenerateShort()

	session := types.BookingSession{
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

	if err := s.sessionRepository.Save(session); err != nil {
		return types.BookingSession{}, exception.Wrap(
			"s.sessionRepository.Save",
			"start-command-telegram-service",
			err,
		)
	}

	return session, nil
}
