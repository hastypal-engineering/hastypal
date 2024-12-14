package telegram

import (
	"github.com/adriein/hastypal/internal/hastypal/shared/constants"
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/service"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"github.com/google/uuid"
	"google.golang.org/api/calendar/v3"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type FinishCommandTelegramService struct {
	bot                    *service.TelegramBot
	googleApi              *service.GoogleApi
	sessionRepository      types.Repository[types.BookingSession]
	notificationRepository types.Repository[types.TelegramNotification]
	bookingRepository      types.Repository[types.Booking]
	googleTokenRepository  types.Repository[types.GoogleToken]
	businessRepository     types.Repository[types.Business]
}

func NewFinishCommandTelegramService(
	bot *service.TelegramBot,
	googleApi *service.GoogleApi,
	sessionRepository types.Repository[types.BookingSession],
	notificationRepository types.Repository[types.TelegramNotification],
	bookingRepository types.Repository[types.Booking],
	googleTokenRepository types.Repository[types.GoogleToken],
	businessRepository types.Repository[types.Business],
) *FinishCommandTelegramService {
	return &FinishCommandTelegramService{
		bot:                    bot,
		googleApi:              googleApi,
		sessionRepository:      sessionRepository,
		notificationRepository: notificationRepository,
		bookingRepository:      bookingRepository,
		googleTokenRepository:  googleTokenRepository,
		businessRepository:     businessRepository,
	}
}

func (s *FinishCommandTelegramService) Execute(update types.TelegramUpdate) error {
	if ackErr := s.ackToTelegramClient(update.CallbackQuery.Id); ackErr != nil {
		return exception.Wrap(
			"s.ackToTelegramClient",
			"finish-command-telegram-service.go",
			ackErr,
		)
	}

	var markdownText strings.Builder

	parsedUrl, parseUrlErr := url.Parse(update.CallbackQuery.Data)

	if parseUrlErr != nil {
		return exception.New(parseUrlErr.Error()).
			Trace("url.Parse", "finish-command-telegram-service.go").
			WithValues([]string{update.CallbackQuery.Data})
	}

	queryParams := parsedUrl.Query()

	sessionId := queryParams.Get("session")

	session, getSessionErr := s.getCurrentSession(sessionId)

	if getSessionErr != nil {
		return exception.Wrap(
			"s.getCurrentSession",
			"finish-command-telegram-service.go",
			getSessionErr,
		)
	}

	if invalidSession := session.EnsureIsValid(); invalidSession != nil {
		return exception.Wrap(
			"session.EnsureIsValid",
			"finish-command-telegram-service.go",
			invalidSession,
		)
	}

	business, getBusinessErr := s.getBusiness(session.BusinessId)

	if getBusinessErr != nil {
		return exception.Wrap(
			"s.getBusiness",
			"finish-command-telegram-service.go",
			getBusinessErr,
		)
	}

	if duplicatedErr := s.ensureNotDuplicatedNotification(session); duplicatedErr != nil {
		return exception.Wrap(
			"s.ensureNotDuplicatedNotification",
			"finish-command-telegram-service.go",
			duplicatedErr,
		)
	}

	booking, createBookingErr := s.createBooking(session, business)

	if createBookingErr != nil {
		return exception.Wrap(
			"s.createBooking",
			"finish-command-telegram-service.go",
			createBookingErr,
		)
	}

	if storeBookingErr := s.storeBooking(booking); storeBookingErr != nil {
		return exception.Wrap(
			"s.storeBooking",
			"finish-command-telegram-service.go",
			storeBookingErr,
		)
	}

	if registerErr := s.registerNotification(booking, update, business); registerErr != nil {
		return exception.Wrap(
			"s.registerNotification",
			"finish-command-telegram-service.go",
			registerErr,
		)
	}

	if registerEventErr := s.registerEventInBusinessCalendar(business, booking); registerEventErr != nil {
		return exception.Wrap(
			"s.registerEventInBusinessCalendar",
			"finish-command-telegram-service.go",
			registerEventErr,
		)
	}

	markdownText.WriteString("![🎉](tg://emoji?id=5368324170671202286) *¡Reserva confirmada\\!*\n\n")
	markdownText.WriteString("![📅](tg://emoji?id=5368324170671202286) ")
	markdownText.WriteString("Te avisaré un día antes para recordarte la cita\n\n")
	markdownText.WriteString("![💙](tg://emoji?id=5368324170671202286) Muchas gracias por la confianza depositada")

	buttons := make([][]types.KeyboardButton, 0)

	message := types.SendTelegramMessage{
		ChatId:         update.CallbackQuery.From.Id,
		Text:           markdownText.String(),
		ParseMode:      constants.TelegramMarkdown,
		ProtectContent: true,
		ReplyMarkup:    types.ReplyMarkup{InlineKeyboard: buttons},
	}

	if botSendMsgErr := s.bot.SendMsg(message); botSendMsgErr != nil {
		return exception.Wrap(
			"s.bot.SendMsg",
			"finish-command-telegram-service.go",
			botSendMsgErr,
		)
	}

	return nil
}

func (s *FinishCommandTelegramService) getCurrentSession(sessionId string) (types.BookingSession, error) {
	filter := types.Filter{
		Name:    "id",
		Operand: constants.Equal,
		Value:   sessionId,
	}

	criteria := types.Criteria{Filters: []types.Filter{filter}}

	session, findOneErr := s.sessionRepository.FindOne(criteria)

	if findOneErr != nil {
		return types.BookingSession{}, exception.Wrap(
			"s.sessionRepository.FindOne",
			"finish-command-telegram-service.go",
			findOneErr,
		)
	}

	return session, nil
}

func (s *FinishCommandTelegramService) ackToTelegramClient(callbackQueryId string) error {
	return s.bot.AnswerCallbackQuery(types.AnswerCallbackQuery{CallbackQueryId: callbackQueryId})
}

func (s *FinishCommandTelegramService) ensureNotDuplicatedNotification(session types.BookingSession) error {
	filter := types.Filter{
		Name:    "session_id",
		Operand: constants.Equal,
		Value:   session.Id,
	}

	criteria := types.Criteria{Filters: []types.Filter{filter}}

	result, err := s.notificationRepository.Find(criteria)

	if err != nil {
		return exception.Wrap("s.notificationRepository.Find", "finish-command-telegram-service.go", err)
	}

	if len(result) > 0 {
		return exception.New("Already saved notification for this booking").
			Trace("ensureNotDuplicatedNotification", "finish-command-telegram-service.go").
			WithValues([]string{session.Id})
	}

	return nil
}

func (s *FinishCommandTelegramService) registerNotification(
	booking types.Booking,
	update types.TelegramUpdate,
	business types.Business,
) error {
	bookingDate, timeParseErr := time.Parse(time.DateTime, booking.When)

	if timeParseErr != nil {
		return exception.New(timeParseErr.Error()).
			Trace("time.Parse", "finish-command-telegram-service.go").
			WithValues([]string{booking.When})
	}

	oneDayBefore := bookingDate.Add(-24 * time.Hour)

	notificationDate := time.Date(
		oneDayBefore.Year(),
		oneDayBefore.Month(),
		oneDayBefore.Day(),
		10, 0, 0, 0, time.UTC,
	)

	notification := types.TelegramNotification{
		Id:           uuid.New().String(),
		SessionId:    booking.SessionId,
		BookingId:    booking.Id,
		BusinessId:   business.Id,
		ScheduledAt:  notificationDate.Format(time.DateTime),
		ChatId:       update.CallbackQuery.From.Id,
		BusinessName: business.Name,
		ServiceName:  "placeholder-test",
		BookingDate:  booking.When,
		Sent:         false,
		CreatedAt:    time.Now().UTC().Format(time.DateTime),
	}

	if err := s.notificationRepository.Save(notification); err != nil {
		return exception.Wrap("s.notificationRepository.Save", "finish-command-telegram-service.go", err)
	}

	return nil
}

func (s *FinishCommandTelegramService) createBooking(
	session types.BookingSession,
	business types.Business,
) (types.Booking, error) {
	bookingDate, timeParseErr := time.Parse(time.DateTime, session.Date)

	if timeParseErr != nil {
		return types.Booking{}, exception.New(timeParseErr.Error()).
			Trace("time.Parse", "finish-command-telegram-service.go").
			WithValues([]string{session.Date})
	}

	stringHours := strings.Split(session.Hour, ":")

	bookingHour, hourConversionErr := strconv.Atoi(stringHours[0])

	if hourConversionErr != nil {
		return types.Booking{}, exception.New(hourConversionErr.Error()).
			Trace("strconv.Atoi(stringHours[0])", "finish-command-telegram-service.go").
			WithValues([]string{stringHours[0]})
	}

	bookingMinutes, minutesConversionErr := strconv.Atoi(stringHours[1])

	if minutesConversionErr != nil {
		return types.Booking{}, exception.New(minutesConversionErr.Error()).
			Trace("strconv.Atoi(stringHours[1])", "finish-command-telegram-service.go").
			WithValues([]string{stringHours[1]})
	}

	bookingDateWithHour := time.Date(
		bookingDate.Year(),
		bookingDate.Month(),
		bookingDate.Day(),
		bookingHour, bookingMinutes, 0, 0, time.UTC,
	)

	booking := types.Booking{
		Id:         uuid.New().String(),
		SessionId:  session.Id,
		BusinessId: business.Id,
		ServiceId:  session.ServiceId,
		When:       bookingDateWithHour.Format(time.DateTime),
		CreatedAt:  time.Now().UTC().Format(time.DateTime),
	}

	return booking, nil
}

func (s *FinishCommandTelegramService) storeBooking(booking types.Booking) error {
	if err := s.bookingRepository.Save(booking); err != nil {
		return exception.Wrap("s.bookingRepository.Save", "finish-command-telegram-service.go", err)
	}

	return nil
}

func (s *FinishCommandTelegramService) getBusinessGoogleToken(business types.Business) (types.GoogleToken, error) {
	filter := types.Filter{
		Name:    "business_id",
		Operand: constants.Equal,
		Value:   business.Id,
	}

	criteria := types.Criteria{Filters: []types.Filter{filter}}

	result, err := s.googleTokenRepository.FindOne(criteria)

	if err != nil {
		return types.GoogleToken{}, exception.Wrap(
			"s.googleTokenRepository.FindOne",
			"finish-command-telegram-service.go",
			err,
		)
	}

	return result, nil
}

func (s *FinishCommandTelegramService) getGoogleCalendarClient(business types.Business) (*calendar.Service, error) {
	token, getTokenErr := s.getBusinessGoogleToken(business)

	if getTokenErr != nil {
		return nil, exception.Wrap(
			"s.getBusinessGoogleToken",
			"finish-command-telegram-service.go",
			getTokenErr,
		)
	}

	client, calendarClientErr := s.googleApi.CalendarClient(token)

	if calendarClientErr != nil {
		return nil, exception.Wrap(
			"s.googleApi.CalendarClient",
			"finish-command-telegram-service.go",
			calendarClientErr,
		)
	}

	return client, nil
}

func (s *FinishCommandTelegramService) registerEventInBusinessCalendar(
	business types.Business,
	booking types.Booking,
) error {
	client, getClientErr := s.getGoogleCalendarClient(business)

	if getClientErr != nil {
		return exception.Wrap(
			"s.getGoogleCalendarClient",
			"finish-command-telegram-service.go",
			getClientErr,
		)
	}

	bookingStartDate, timeParseErr := time.Parse(time.DateTime, booking.When)

	if timeParseErr != nil {
		return exception.New(timeParseErr.Error()).
			Trace("time.Parse", "finish-command-telegram-service.go").
			WithValues([]string{booking.When})
	}

	bookingFinishDate := bookingStartDate.Add(1 * time.Hour)

	event := &calendar.Event{
		Summary:     "Google I/O 2015",
		Location:    "800 Howard St., San Francisco, CA 94103",
		Description: "A chance to hear more about Google's developer products.",
		Start: &calendar.EventDateTime{
			DateTime: bookingStartDate.Format(time.RFC3339),
			TimeZone: "Europe/Madrid",
		},
		End: &calendar.EventDateTime{
			DateTime: bookingFinishDate.Format(time.RFC3339),
			TimeZone: "Europe/Madrid",
		},
		Organizer: &calendar.EventOrganizer{
			Email: "adria.claret@gmail.com",
			Self:  true,
		},
		Status: "confirmed",
	}

	_, insertErr := client.Events.Insert("primary", event).Do()

	if insertErr != nil {
		return exception.New(insertErr.Error()).
			Trace("client.Events.Insert", "finish-command-telegram-service.go")
	}

	return nil
}

func (s *FinishCommandTelegramService) getBusiness(businessId string) (types.Business, error) {
	filters := make([]types.Filter, 1)

	filters[0] = types.Filter{Name: "id", Operand: constants.Equal, Value: businessId}

	criteria := types.Criteria{Filters: filters}

	business, err := s.businessRepository.FindOne(criteria)

	if err != nil {
		return types.Business{}, exception.Wrap(
			"s.businessRepository.FindOne",
			"finish-command-telegram-service.go",
			err,
		)
	}

	return business, nil
}
