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
	notifications, getNotificationsErr := s.getNotifications()

	if getNotificationsErr != nil {
		return exception.Wrap("s.getNotifications", "send-notification-service.go", getNotificationsErr)
	}

	location, loadLocationErr := time.LoadLocation("Europe/Madrid")

	if loadLocationErr != nil {
		return exception.New(loadLocationErr.Error()).
			Trace("time.LoadLocation", "send-notification-service.go").
			WithValues([]string{"Europe/Madrid"})
	}

	time.Local = location

	spanishMonths := map[time.Month]string{
		time.January:   "Enero",
		time.February:  "Febrero",
		time.March:     "Marzo",
		time.April:     "Abril",
		time.May:       "Mayo",
		time.June:      "Junio",
		time.July:      "Julio",
		time.August:    "Agosto",
		time.September: "Septiembre",
		time.October:   "Octubre",
		time.November:  "Noviembre",
		time.December:  "Diciembre",
	}

	for _, notification := range notifications {
		var markdownText strings.Builder

		welcome := "![👋](tg://emoji?id=5368324170671202286) Hola y "

		apologyForBother := fmt.Sprintf(
			"perdona las molestias, recuerda que tienes una cita en %s ",
			notification.BusinessName,
		)

		bookingService := fmt.Sprintf(
			"para %s ",
			"Corte de pelo y barba express 18€",
		)

		parsedBookingDate, bookingDateParseErr := time.Parse(time.DateTime, notification.BookingDate)

		if bookingDateParseErr != nil {
			return exception.New(bookingDateParseErr.Error()).
				Trace("time.Parse", "send-notification-service.go").
				WithValues([]string{notification.BookingDate})
		}

		bookingDate := fmt.Sprintf("el día %d de %s, %d a las %02d:%02d %s",
			parsedBookingDate.Day(),
			spanishMonths[parsedBookingDate.Month()],
			parsedBookingDate.Year(),
			parsedBookingDate.Hour(),
			parsedBookingDate.Minute(),
			parsedBookingDate.Format("PM"),
		)

		markdownText.WriteString(welcome)
		markdownText.WriteString(apologyForBother)
		markdownText.WriteString(bookingService)
		markdownText.WriteString(bookingDate)

		message := types.SendTelegramMessage{
			ChatId:         notification.ChatId,
			Text:           markdownText.String(),
			ParseMode:      constants.TelegramMarkdown,
			ProtectContent: true,
			ReplyMarkup:    types.ReplyMarkup{InlineKeyboard: make([][]types.KeyboardButton, 0)},
		}

		if botSendMsgErr := s.bot.SendMsg(message); botSendMsgErr != nil {
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
