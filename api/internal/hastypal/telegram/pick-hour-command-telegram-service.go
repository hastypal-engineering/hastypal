package telegram

import (
	"errors"
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/shared/constants"
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	"github.com/adriein/hastypal/internal/hastypal/shared/service"
	"github.com/adriein/hastypal/internal/hastypal/shared/translation"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"net/url"
	"strings"
	"time"
)

type PickHourCommandTelegramService struct {
	bot                *service.TelegramBot
	sessionRepository  types.Repository[types.BookingSession]
	bookingRepository  types.Repository[types.Booking]
	businessRepository types.Repository[types.Business]
	translation        *translation.Translation
}

func NewPickHourCommandTelegramService(
	bot *service.TelegramBot,
	sessionRepository types.Repository[types.BookingSession],
	bookingRepository types.Repository[types.Booking],
	businessRepository types.Repository[types.Business],
	translation *translation.Translation,
) *PickHourCommandTelegramService {
	return &PickHourCommandTelegramService{
		bot:                bot,
		sessionRepository:  sessionRepository,
		bookingRepository:  bookingRepository,
		businessRepository: businessRepository,
		translation:        translation,
	}
}

func (s *PickHourCommandTelegramService) Execute(update types.TelegramUpdate) error {
	if ackErr := s.ackToTelegramClient(update.CallbackQuery.Id); ackErr != nil {
		return exception.Wrap(
			"s.ackToTelegramClient",
			"pick-hour-command-telegram-service.go",
			ackErr,
		)
	}

	var markdownText strings.Builder

	location, loadLocationErr := time.LoadLocation("Europe/Madrid")

	if loadLocationErr != nil {
		return exception.New(loadLocationErr.Error()).
			Trace("time.LoadLocation", "pick-hour-command-telegram-service.go").
			WithValues([]string{"Europe/Madrid"})
	}

	time.Local = location

	parsedUrl, parseUrlErr := url.Parse(update.CallbackQuery.Data)

	if parseUrlErr != nil {
		return exception.New(parseUrlErr.Error()).
			Trace("url.Parse(update.CallbackQuery.Data)", "pick-hour-command-telegram-service.go").
			WithValues([]string{update.CallbackQuery.Data})
	}

	queryParams := parsedUrl.Query()

	stringSelectedDate := fmt.Sprintf("%s %s", queryParams.Get("date"), "07:00:00")
	sessionId := queryParams.Get("session")

	selectedDate, timeParseErr := time.Parse(time.DateTime, stringSelectedDate)

	if timeParseErr != nil {
		return exception.New(timeParseErr.Error()).
			Trace(
				"time.Parse(time.DateTime, stringSelectedDate)",
				"pick-hour-command-telegram-service.go",
			).
			WithValues([]string{stringSelectedDate})
	}

	session, getSessionErr := s.getCurrentSession(sessionId)

	if getSessionErr != nil {
		return exception.Wrap(
			"s.getCurrentSession",
			"pick-hour-command-telegram-service.go",
			getSessionErr,
		)
	}

	business, getBusinessErr := s.getBusiness(session.BusinessId)

	if getBusinessErr != nil {
		return exception.Wrap(
			"s.getBusiness",
			"pick-hour-command-telegram-service.go",
			getBusinessErr,
		)
	}

	if invalidSession := session.EnsureIsValid(); invalidSession != nil {
		message := types.TelegramMessage{ChatId: update.CallbackQuery.From.Id}

		expiredSessionMessage := message.SessionExpired()

		bookingExpiredSessionMessage := types.BookingTelegramMessage{
			BusinessName:     business.Name,
			BookingSessionId: session.Id,
			Message:          expiredSessionMessage,
		}

		if botSendMsgErr := s.bot.SendMsg(bookingExpiredSessionMessage); botSendMsgErr != nil {
			return exception.Wrap(
				"s.bot.SendMsg",
				"pick-hour-command-telegram-service.go",
				botSendMsgErr,
			)
		}

		return nil
	}

	if updateSessionErr := s.updateSession(session, stringSelectedDate); updateSessionErr != nil {
		return exception.Wrap(
			"s.updateSession",
			"pick-hour-command-telegram-service.go",
			updateSessionErr,
		)
	}

	dateParts := strings.Split(selectedDate.Format(time.RFC822), " ")

	day := dateParts[0]
	month := s.translation.GetSpanishMonthShortForm(selectedDate.Month())

	welcome := fmt.Sprintf(
		"![⌚️](tg://emoji?id=5368324170671202286) Las horas disponibles para:\n\n",
	)

	selectedService := fmt.Sprintf(
		"![🔸](tg://emoji?id=5368324170671202286) %s\n\n",
		"Corte de pelo y barba express 18€",
	)

	date := fmt.Sprintf(
		"![📅](tg://emoji?id=5368324170671202286) %s\n\n",
		fmt.Sprintf("%s %s", day, month),
	)

	processInstructions := "*Selecciona una hora y te escribiré un resumen para que puedas confirmar la reserva*\n\n"

	markdownText.WriteString(welcome)
	markdownText.WriteString(selectedService)
	markdownText.WriteString(date)
	markdownText.WriteString(processInstructions)

	buttons := make([]types.KeyboardButton, 12)

	for i := 8; i <= len(buttons)+7; i++ {
		hour := fmt.Sprintf("%02d:00", i)

		criteria := s.similarSessionCriteria(selectedDate, hour)

		isAlreadyPicked, err := s.isHourAlreadyPicked(criteria)

		if err != nil {
			return exception.Wrap("s.isHourAlreadyPicked", "pick-hour-command-telegram-service.go", err)
		}

		if isAlreadyPicked {
			continue
		}

		buttons[i-8] = types.KeyboardButton{
			Text: fmt.Sprintf("%s", hour),
			CallbackData: fmt.Sprintf(
				"/confirmation?session=%s&hour=%s",
				sessionId,
				hour,
			),
		}
	}

	backButton := types.KeyboardButton{
		Text:         "Atrás",
		CallbackData: fmt.Sprintf("/dates?session=%s&service=%s", session.Id, "test-short"),
	}

	buttons = append(buttons, backButton)

	array := helper.NewArrayHelper[types.KeyboardButton]()

	inlineKeyboard := array.Chunk(buttons, 3)

	message := types.TelegramMessage{
		ChatId:         update.CallbackQuery.From.Id,
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
			"pick-hour-command-telegram-service.go",
			botSendMsgErr,
		)
	}

	return nil
}

func (s *PickHourCommandTelegramService) ackToTelegramClient(callbackQueryId string) error {
	return s.bot.AnswerCallbackQuery(types.AnswerCallbackQuery{CallbackQueryId: callbackQueryId})
}

func (s *PickHourCommandTelegramService) getCurrentSession(sessionId string) (types.BookingSession, error) {
	filter := types.Filter{
		Name:    "id",
		Operand: constants.Equal,
		Value:   sessionId,
	}

	criteria := types.Criteria{Filters: []types.Filter{filter}}

	session, findOneErr := s.sessionRepository.FindOne(criteria)

	if findOneErr != nil {
		return types.BookingSession{}, exception.Wrap(
			"s.repository.FindOne",
			"pick-hour-command-telegram-service.go",
			findOneErr,
		)
	}

	return session, nil
}

func (s *PickHourCommandTelegramService) updateSession(actualSession types.BookingSession, date string) error {
	updatedSession := types.BookingSession{
		Id:         actualSession.Id,
		BusinessId: actualSession.BusinessId,
		ChatId:     actualSession.ChatId,
		ServiceId:  actualSession.ServiceId,
		Date:       date,
		Hour:       "",
		CreatedAt:  actualSession.CreatedAt,
		UpdatedAt:  time.Now().UTC().Format(time.DateTime), //We refresh the created at on purpose
		Ttl:        actualSession.Ttl,
	}

	if err := s.sessionRepository.Update(updatedSession); err != nil {
		return exception.Wrap(
			"s.repository.Update",
			"pick-hour-command-telegram-service.go",
			err,
		)
	}

	return nil
}

func (s *PickHourCommandTelegramService) isHourAlreadyPicked(criteria types.Criteria) (bool, error) {
	session, findOneErr := s.sessionRepository.FindOne(criteria)

	var sessionNotFoundErr exception.HastypalError

	if findOneErr != nil && errors.As(findOneErr, &sessionNotFoundErr) {
		if sessionNotFoundErr.IsDomain() {
			return false, nil
		}

		return false, exception.Wrap(
			"s.sessionRepository.FindOne",
			"pick-hour-command-telegram-service.go",
			findOneErr,
		)
	}

	if invalidSession := session.EnsureIsValid(); invalidSession != nil {
		bookingCriteria := s.bookingCriteria(session.Id)

		_, findOneBookingErr := s.bookingRepository.FindOne(bookingCriteria)

		var bookingNotFoundErr exception.HastypalError

		if findOneBookingErr != nil && errors.As(findOneErr, &bookingNotFoundErr) {
			if sessionNotFoundErr.IsDomain() {
				return false, nil
			}

			return false, exception.Wrap(
				"s.bookingRepository.FindOne",
				"pick-hour-command-telegram-service.go",
				findOneErr,
			)
		}

		return true, nil
	}

	return true, nil
}

func (s *PickHourCommandTelegramService) similarSessionCriteria(date time.Time, hour string) types.Criteria {
	return types.NewCriteria().
		Equal("date", date.Format(time.DateTime)).
		Equal("hour", hour)
}

func (s *PickHourCommandTelegramService) bookingCriteria(sessionId string) types.Criteria {
	return types.NewCriteria().
		Equal("session_id", sessionId)
}

func (s *PickHourCommandTelegramService) getBusiness(businessId string) (types.Business, error) {
	filters := make([]types.Filter, 1)

	filters[0] = types.Filter{Name: "id", Operand: constants.Equal, Value: businessId}

	criteria := types.Criteria{Filters: filters}

	business, err := s.businessRepository.FindOne(criteria)

	if err != nil {
		return types.Business{}, exception.Wrap(
			"s.businessRepository.FindOne",
			"pick-hour-command-telegram-service.go",
			err,
		)
	}

	return business, nil
}
