package notification

import (
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/shared/constants"
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/service"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"strings"
	"time"
)

type SendNotificationService struct {
	repository types.Repository[types.TelegramNotification]
	bot        *service.TelegramBot
}

func NewSendNotificationService(
	repository types.Repository[types.TelegramNotification],
	bot *service.TelegramBot,
) *SendNotificationService {
	return &SendNotificationService{
		repository: repository,
		bot:        bot,
	}
}

func (s *SendNotificationService) Execute() error {
	filter := types.Filter{
		Name:    "scheduled_at",
		Operand: constants.LessThanOrEqual,
		Value:   time.Now().Format(time.DateTime),
	}

	criteria := types.Criteria{Filters: []types.Filter{filter}, Join: []string{"booking_id"}}

	notifications, findErr := s.repository.Find(criteria)

	if findErr != nil {
		return exception.Wrap("s.repository.Find", "send-notification-service.go", findErr)
	}

	for _, notification := range notifications {
		var markdownText strings.Builder

		welcome := fmt.Sprintf(
			"Hola ![👋](tg://emoji?id=5368324170671202286), recuerda que tienes una cita en %s\n\n",
			notification.BusinessName,
		)

		bookingService := fmt.Sprintf(
			"![🟢](tg://emoji?id=5368324170671202286) Para %s\n\n",
			notification.ServiceName,
		)

		bookingDate := fmt.Sprintf(
			"![⌚️](tg://emoji?id=5368324170671202286) En la siguiente fecha %s\n\n",
			notification.BookingDate,
		)

		markdownText.WriteString(welcome)
		markdownText.WriteString(bookingService)
		markdownText.WriteString(bookingDate)

		message := types.SendTelegramMessage{
			ChatId:         notification.ChatId,
			Text:           markdownText.String(),
			ParseMode:      constants.TelegramMarkdown,
			ProtectContent: true,
		}

		if botSendMsgErr := s.bot.SendMsg(message); botSendMsgErr != nil {
			return exception.Wrap(
				"s.bot.SendMsg",
				"send-notification-service.go",
				botSendMsgErr,
			)
		}
	}

	return nil
}
