package telegram

import (
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/shared/constants"
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
}

func NewFinishCommandTelegramService(
	bot *service.TelegramBot,
	googleApi *service.GoogleApi,
	sessionRepository types.Repository[types.BookingSession],
	notificationRepository types.Repository[types.TelegramNotification],
	bookingRepository types.Repository[types.Booking],
	googleTokenRepository types.Repository[types.GoogleToken],
) *FinishCommandTelegramService {
	return &FinishCommandTelegramService{
		bot:                    bot,
		googleApi:              googleApi,
		sessionRepository:      sessionRepository,
		notificationRepository: notificationRepository,
		bookingRepository:      bookingRepository,
		googleTokenRepository:  googleTokenRepository,
	}
}

func (s *FinishCommandTelegramService) Execute(business types.Business, update types.TelegramUpdate) error {
	if ackErr := s.ackToTelegramClient(update.CallbackQuery.Id); ackErr != nil {
		return ackErr
	}

	var markdownText strings.Builder

	parsedUrl, parseUrlErr := url.Parse(update.CallbackQuery.Data)

	if parseUrlErr != nil {
		return types.ApiError{
			Msg:      parseUrlErr.Error(),
			Function: "Execute -> url.Parse()",
			File:     "telegram-finish-command.go",
			Values:   []string{update.CallbackQuery.Data},
		}
	}

	queryParams := parsedUrl.Query()

	sessionId := queryParams.Get("session")

	session, getSessionErr := s.getCurrentSession(sessionId)

	if getSessionErr != nil {
		return getSessionErr
	}

	if invalidSession := session.EnsureIsValid(); invalidSession != nil {
		return invalidSession
	}

	if duplicatedErr := s.ensureNotDuplicatedNotification(session); duplicatedErr != nil {
		return duplicatedErr
	}

	booking, createBookingErr := s.createBooking(session, business)

	if createBookingErr != nil {
		return createBookingErr
	}

	if storeBookingErr := s.storeBooking(booking); storeBookingErr != nil {
		return storeBookingErr
	}

	if registerErr := s.registerNotification(booking, update, business); registerErr != nil {
		return registerErr
	}

	event := &calendar.Event{
		Summary:     "Google I/O 2015",
		Location:    "800 Howard St., San Francisco, CA 94103",
		Description: "A chance to hear more about Google's developer products.",
		Start: &calendar.EventDateTime{
			DateTime: "2024-11-20T09:00:00-07:00",
			TimeZone: "America/Los_Angeles",
		},
		End: &calendar.EventDateTime{
			DateTime: "2024-11-20T17:00:00-07:00",
			TimeZone: "America/Los_Angeles",
		},
		Organizer: &calendar.EventOrganizer{
			Email: "adria.claret@gmail.com",
			Self:  true,
		},
		Status: "confirmed",
	}

	if registerEventErr := s.registerEventInBusinessCalendar(business, event); registerEventErr != nil {
		return registerEventErr
	}

	markdownText.WriteString("![🎉](tg://emoji?id=5368324170671202286) *¡Reserva confirmada\\!*\n\n")
	markdownText.WriteString("Te avisaremos un día antes para recordarte la cita ")
	markdownText.WriteString("![📅](tg://emoji?id=5368324170671202286)\n\n")

	buttons := make([][]types.KeyboardButton, 0)

	message := types.SendTelegramMessage{
		ChatId:         update.CallbackQuery.From.Id,
		Text:           markdownText.String(),
		ParseMode:      constants.TelegramMarkdown,
		ProtectContent: true,
		ReplyMarkup:    types.ReplyMarkup{InlineKeyboard: buttons},
	}

	if botSendMsgErr := s.bot.SendMsg(message); botSendMsgErr != nil {
		return botSendMsgErr
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
		return types.BookingSession{}, findOneErr
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
		return err
	}

	if len(result) > 0 {
		return types.ApiError{
			Msg:      "Already saved notification for this booking",
			Function: "ensureNotDuplicatedNotification",
			File:     "service/telegram-finish-command.go",
			Values:   []string{session.Id},
		}
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
		return types.ApiError{
			Msg:      timeParseErr.Error(),
			Function: "registerNotification -> time.Parse()",
			File:     "service/telegram-finish-command.go",
			Values:   []string{booking.When},
		}
	}

	oneDayBefore := bookingDate.Add(-24 * time.Hour)

	notificationDate := time.Date(
		oneDayBefore.Year(),
		oneDayBefore.Month(),
		oneDayBefore.Day(),
		10, 0, 0, 0, time.UTC,
	)

	notification := types.TelegramNotification{
		Id:          uuid.New().String(),
		SessionId:   booking.SessionId,
		BookingId:   booking.Id,
		BusinessId:  business.Id,
		ScheduledAt: notificationDate.Format(time.DateTime),
		ChatId:      update.CallbackQuery.From.Id,
		From:        business.Name,
		CreatedAt:   time.Now().Format(time.DateTime),
	}

	if err := s.notificationRepository.Save(notification); err != nil {
		return err
	}

	return nil
}

func (s *FinishCommandTelegramService) createBooking(
	session types.BookingSession,
	business types.Business,
) (types.Booking, error) {
	bookingDate, timeParseErr := time.Parse(time.DateTime, session.Date)

	if timeParseErr != nil {
		return types.Booking{}, types.ApiError{
			Msg:      timeParseErr.Error(),
			Function: "createBooking -> time.Parse()",
			File:     "service/telegram-finish-command.go",
			Values:   []string{session.Date},
		}
	}

	stringHours := strings.Split(session.Hour, ":")

	bookingHour, hourConversionErr := strconv.Atoi(stringHours[0])

	if hourConversionErr != nil {
		return types.Booking{}, types.ApiError{
			Msg:      hourConversionErr.Error(),
			Function: "createBooking -> strconv.Atoi()",
			File:     "service/telegram-finish-command.go",
			Values:   []string{stringHours[0]},
		}
	}

	bookingMinutes, minutesConversionErr := strconv.Atoi(stringHours[1])

	if minutesConversionErr != nil {
		return types.Booking{}, types.ApiError{
			Msg:      minutesConversionErr.Error(),
			Function: "createBooking -> strconv.Atoi()",
			File:     "service/telegram-finish-command.go",
			Values:   []string{stringHours[1]},
		}
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
		CreatedAt:  time.Now().Format(time.DateTime),
	}

	return booking, nil
}

func (s *FinishCommandTelegramService) storeBooking(booking types.Booking) error {
	if saveErr := s.bookingRepository.Save(booking); saveErr != nil {
		return saveErr
	}

	return nil
}

func (s *FinishCommandTelegramService) getBusinessGoogleToken(business types.Business) (types.GoogleToken, error) {
	filter := types.Filter{
		Name:    "business_id",
		Operand: constants.Equal,
		Value:   "testId",
	}

	criteria := types.Criteria{Filters: []types.Filter{filter}}

	result, err := s.googleTokenRepository.FindOne(criteria)

	if err != nil {
		return types.GoogleToken{}, err
	}

	return result, nil
}

func (s *FinishCommandTelegramService) getGoogleCalendarClient(business types.Business) (*calendar.Service, error) {
	token, getTokenErr := s.getBusinessGoogleToken(business)

	if getTokenErr != nil {
		return nil, getTokenErr
	}

	client, calendarClientErr := s.googleApi.CalendarClient(token)

	if calendarClientErr != nil {
		return nil, calendarClientErr
	}

	return client, nil
}

func (s *FinishCommandTelegramService) registerEventInBusinessCalendar(
	business types.Business,
	event *calendar.Event,
) error {
	client, getClientErr := s.getGoogleCalendarClient(business)

	if getClientErr != nil {
		return getClientErr
	}

	_, insertErr := client.Events.Insert("primary", event).Do()

	if insertErr != nil {
		return types.ApiError{
			Msg:      insertErr.Error(),
			Function: "registerEventInBusinessCalendar -> client.Events.Insert()",
			File:     "service/telegram-finish-command.go",
		}
	}

	return nil
}

func (s *FinishCommandTelegramService) createInviteLink(event *calendar.Event) (strings.Builder, error) {
	var builder strings.Builder

	startParsedTime, startParseErr := time.Parse(time.RFC3339, event.Start.DateTime)

	if startParseErr != nil {
		return builder, types.ApiError{
			Msg:      startParseErr.Error(),
			Function: "createInviteLink -> startParsedTime -> time.Parse()",
			File:     "service/telegram-finish-command.go",
		}
	}

	endParsedTime, endParseErr := time.Parse(time.RFC3339, event.End.DateTime)

	if endParseErr != nil {
		return builder, types.ApiError{
			Msg:      endParseErr.Error(),
			Function: "createInviteLink -> endParsedTime -> time.Parse()",
			File:     "service/telegram-finish-command.go",
		}
	}

	builder.WriteString("https://calendar.google.com/calendar/r/eventedit?action=TEMPLATE")
	builder.WriteString(fmt.Sprintf("&dates=%s", startParsedTime.UTC().Format("20060102T150405Z")))
	builder.WriteString("%2F")
	builder.WriteString(fmt.Sprintf("%s", endParsedTime.UTC().Format("20060102T150405Z")))
	builder.WriteString(fmt.Sprintf("&stz=%s", event.Start.TimeZone))
	builder.WriteString(fmt.Sprintf("&etz=%s", event.End.TimeZone))
	builder.WriteString(fmt.Sprintf(
		"&details=%s",
		url.QueryEscape("A chance to hear more about Google's developer products\\."),
	))
	builder.WriteString(fmt.Sprintf(
		"&location=%s",
		url.QueryEscape("800 Howard St., San Francisco, CA 94103"),
	))
	builder.WriteString(fmt.Sprintf("&text=%s", url.QueryEscape("Google I/O 2015")))

	return builder, nil
}
