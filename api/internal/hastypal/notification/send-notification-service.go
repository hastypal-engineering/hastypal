package notification

import (
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/shared/constants"
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/service"
	"github.com/adriein/hastypal/internal/hastypal/shared/translation"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"strings"
	"time"
)

type SendNotificationService struct {
	repository  types.Repository[types.TelegramNotification]
	bot         *service.TelegramBot
	translation *translation.Translation
}

func NewSendNotificationService(
	repository types.Repository[types.TelegramNotification],
	bot *service.TelegramBot,
	translation *translation.Translation,
) *SendNotificationService {
	return &SendNotificationService{
		repository:  repository,
		bot:         bot,
		translation: translation,
	}
}

func (s *SendNotificationService) Execute() error {
	notifications, getNotificationsErr := s.getNotifications()

	if getNotificationsErr != nil {
		return exception.Wrap("s.getNotifications", "send-notification-service.go", getNotificationsErr)
	}

	for _, notification := range notifications {
		var markdownText strings.Builder

		parsedBookingDate, bookingDateParseErr := time.Parse(time.DateTime, notification.BookingDate)

		if bookingDateParseErr != nil {
			return exception.New(bookingDateParseErr.Error()).
				Trace("time.Parse", "send-notification-service.go").
				WithValues([]string{notification.BookingDate})
		}

		bookedService := fmt.Sprintf(
			"![🟢](tg://emoji?id=5368324170671202286) %s\n\n",
			"Corte de pelo y barba express 18€",
		)

		date := fmt.Sprintf(
			"![📅](tg://emoji?id=5368324170671202286) %d %s\n\n",
			parsedBookingDate.Day(),
			s.translation.GetSpanishMonth(parsedBookingDate.Month()),
		)

		hour := fmt.Sprintf(
			"![⌚️](tg://emoji?id=5368324170671202286) %02d:%02dH",
			parsedBookingDate.Hour(),
			parsedBookingDate.Minute(),
		)

		markdownText.WriteString("![⏰](tg://emoji?id=5368324170671202286) *¡Recordatorio de tu próxima cita\\!*\n\n")
		markdownText.WriteString(bookedService)
		markdownText.WriteString(date)
		markdownText.WriteString(hour)

		message := types.TelegramMessage{
			ChatId:         notification.ChatId,
			Text:           markdownText.String(),
			ParseMode:      constants.TelegramMarkdown,
			ProtectContent: true,
			ReplyMarkup:    types.ReplyMarkup{InlineKeyboard: make([][]types.KeyboardButton, 0)},
		}

		bookingMessage := types.BookingTelegramMessage{
			BusinessName:     notification.BusinessName,
			BookingSessionId: notification.SessionId,
			Message:          message,
		}

		if botSendMsgErr := s.bot.SendMsg(bookingMessage); botSendMsgErr != nil {
			return exception.Wrap(
				"s.bot.SendMsg",
				"send-notification-service.go",
				botSendMsgErr,
			)
		}

		notification.MarkAsSent()

		if updateErr := s.repository.Update(notification); updateErr != nil {
			return exception.Wrap("s.repository.Update", "send-notification-service.go", updateErr)
		}
	}

	return nil
}

func (s *SendNotificationService) getNotifications() ([]types.TelegramNotification, error) {
	scheduledAtFilter := types.Filter{
		Name:    "scheduled_at",
		Operand: constants.LessThanOrEqual,
		Value:   time.Now().UTC().Format(time.DateTime),
	}

	notSentFilter := types.Filter{
		Name:    "sent",
		Operand: constants.Equal,
		Value:   false,
	}

	criteria := types.Criteria{Filters: []types.Filter{scheduledAtFilter, notSentFilter}}

	notifications, findErr := s.repository.Find(criteria)

	if findErr != nil {
		return nil, exception.Wrap("s.repository.Find", "send-notification-service.go", findErr)
	}

	return notifications, nil
}
