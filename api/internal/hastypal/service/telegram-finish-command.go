package service

import (
	"github.com/adriein/hastypal/internal/hastypal/constants"
	"github.com/adriein/hastypal/internal/hastypal/types"
	"github.com/google/uuid"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type TelegramFinishCommandService struct {
	bot                    *TelegramBot
	sessionRepository      types.Repository[types.BookingSession]
	notificationRepository types.Repository[types.TelegramNotification]
	bookingRepository      types.Repository[types.Booking]
}

func NewTelegramFinishCommandService(
	bot *TelegramBot,
	sessionRepository types.Repository[types.BookingSession],
	notificationRepository types.Repository[types.TelegramNotification],
	bookingRepository types.Repository[types.Booking],
) *TelegramFinishCommandService {
	return &TelegramFinishCommandService{
		bot:                    bot,
		sessionRepository:      sessionRepository,
		notificationRepository: notificationRepository,
		bookingRepository:      bookingRepository,
	}
}

func (s *TelegramFinishCommandService) Execute(business types.Business, update types.TelegramUpdate) error {
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

	if registerErr := s.registerNotification(session, update, business); registerErr != nil {
		return registerErr
	}

	if createBookingErr := s.createBooking(session, business); createBookingErr != nil {
		return createBookingErr
	}

	markdownText.WriteString("![🎉](tg://emoji?id=5368324170671202286) *¡Reserva confirmada\\!*\n\n")
	markdownText.WriteString("Te avisaremos un día antes para recordarte la cita ")
	markdownText.WriteString("![📅](tg://emoji?id=5368324170671202286)")

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

func (s *TelegramFinishCommandService) getCurrentSession(sessionId string) (types.BookingSession, error) {
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

func (s *TelegramFinishCommandService) ackToTelegramClient(callbackQueryId string) error {
	return s.bot.AnswerCallbackQuery(types.AnswerCallbackQuery{CallbackQueryId: callbackQueryId})
}

func (s *TelegramFinishCommandService) ensureNotDuplicatedNotification(session types.BookingSession) error {
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

func (s *TelegramFinishCommandService) registerNotification(
	session types.BookingSession,
	update types.TelegramUpdate,
	business types.Business,
) error {
	bookingDate, timeParseErr := time.Parse(time.DateTime, session.Date)

	if timeParseErr != nil {
		return types.ApiError{
			Msg:      timeParseErr.Error(),
			Function: "registerNotification -> time.Parse()",
			File:     "service/telegram-finish-command.go",
			Values:   []string{session.Date},
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
		SessionId:   session.Id,
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

func (s *TelegramFinishCommandService) createBooking(session types.BookingSession, business types.Business) error {
	bookingDate, timeParseErr := time.Parse(time.DateTime, session.Date)

	if timeParseErr != nil {
		return types.ApiError{
			Msg:      timeParseErr.Error(),
			Function: "createBooking -> time.Parse()",
			File:     "service/telegram-finish-command.go",
			Values:   []string{session.Date},
		}
	}

	stringHours := strings.Split(session.Hour, ":")

	bookingHour, hourConversionErr := strconv.Atoi(stringHours[0])

	if hourConversionErr != nil {
		return types.ApiError{
			Msg:      hourConversionErr.Error(),
			Function: "createBooking -> strconv.Atoi()",
			File:     "service/telegram-finish-command.go",
			Values:   []string{stringHours[0]},
		}
	}

	bookingMinutes, minutesConversionErr := strconv.Atoi(stringHours[1])

	if minutesConversionErr != nil {
		return types.ApiError{
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

	if saveErr := s.bookingRepository.Save(booking); saveErr != nil {
		return saveErr
	}

	return nil
}
