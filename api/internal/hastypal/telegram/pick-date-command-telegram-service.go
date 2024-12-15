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

type PickDateCommandTelegramService struct {
	bot               *service.TelegramBot
	sessionRepository types.Repository[types.BookingSession]
	bookingRepository types.Repository[types.Booking]
	translation       *translation.Translation
}

func NewPickDateCommandTelegramService(
	bot *service.TelegramBot,
	sessionRepository types.Repository[types.BookingSession],
	bookingRepository types.Repository[types.Booking],
	translation *translation.Translation,
) *PickDateCommandTelegramService {
	return &PickDateCommandTelegramService{
		bot:               bot,
		sessionRepository: sessionRepository,
		bookingRepository: bookingRepository,
		translation:       translation,
	}
}

func (s *PickDateCommandTelegramService) Execute(update types.TelegramUpdate) error {
	if ackErr := s.ackToTelegramClient(update.CallbackQuery.Id); ackErr != nil {
		return exception.Wrap(
			"s.ackToTelegramClient",
			"pick-date-command-telegram-service.go",
			ackErr,
		)
	}

	var markdownText strings.Builder

	parsedUrl, parseUrlErr := url.Parse(update.CallbackQuery.Data)

	if parseUrlErr != nil {
		return exception.New(parseUrlErr.Error()).
			Trace("url.Parse(update.CallbackQuery.Data)", "pick-date-command-telegram-service.go").
			WithValues([]string{update.CallbackQuery.Data})
	}

	queryParams := parsedUrl.Query()

	sessionId := queryParams.Get("session")
	serviceId := queryParams.Get("service")

	session, getSessionErr := s.getCurrentSession(sessionId)

	if getSessionErr != nil {
		return exception.Wrap(
			"s.getCurrentSession",
			"pick-date-command-telegram-service.go",
			getSessionErr,
		)
	}

	if invalidSession := session.EnsureIsValid(); invalidSession != nil {
		message := types.SendTelegramMessage{ChatId: update.CallbackQuery.From.Id}

		expiredSessionMessage := message.SessionExpired()

		if botSendMsgErr := s.bot.SendMsg(expiredSessionMessage); botSendMsgErr != nil {
			return exception.Wrap(
				"s.bot.SendMsg",
				"pick-date-command-telegram-service.go",
				botSendMsgErr,
			)
		}

		return nil
	}

	if updateSessionErr := s.updateSession(session, serviceId); updateSessionErr != nil {
		return exception.Wrap(
			"s.updateSession",
			"pick-date-command-telegram-service.go",
			updateSessionErr,
		)
	}

	commandInformation := fmt.Sprintf(
		"%s tiene disponibles para:\n\n![🔸](tg://emoji?id=5368324170671202286) %s\n\n",
		"Hastypal Business Test",
		"Corte de pelo y barba express 18€",
	)

	processInstructions := "*Selecciona un día para ver las horas disponibles:*\n\n"

	markdownText.WriteString("![📅](tg://emoji?id=5368324170671202286) A continuación puedes ver las fechas que ")
	markdownText.WriteString(commandInformation)
	markdownText.WriteString(processInstructions)

	location, loadLocationErr := time.LoadLocation("Europe/Madrid")

	if loadLocationErr != nil {
		return exception.New(loadLocationErr.Error()).
			Trace("time.LoadLocation", "pick-date-command-telegram-service.go").
			WithValues([]string{"Europe/Madrid"})
	}

	time.Local = location

	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 07, 0, 0, 0, location)

	buttons := make([]types.KeyboardButton, 15)

	for i := 0; i < 15; i++ {
		newDate := today.AddDate(0, 0, i)

		hasAvailableSlots, err := s.dateHasAvailableSlots(newDate)

		if err != nil {
			return exception.Wrap("s.dateHasAvailableSlots", "pick-date-command-telegram-service.go", err)
		}

		if !hasAvailableSlots {
			continue
		}

		dateParts := strings.Split(newDate.Format(time.RFC822), " ")

		day := dateParts[0]
		month := s.translation.GetSpanishMonthShortForm(newDate.Month())

		buttons[i] = types.KeyboardButton{
			Text:         fmt.Sprintf("%s %s", day, month),
			CallbackData: fmt.Sprintf("/hours?session=%s&date=%s", sessionId, newDate.Format(time.DateOnly)),
		}
	}

	backButton := types.KeyboardButton{
		Text:         "Atrás",
		CallbackData: fmt.Sprintf("/service?sessionId=%s", session.Id),
	}

	buttons = append(buttons, backButton)

	array := helper.NewArrayHelper[types.KeyboardButton]()

	inlineKeyboard := array.Chunk(buttons, 5)

	message := types.SendTelegramMessage{
		ChatId:         update.CallbackQuery.From.Id,
		Text:           markdownText.String(),
		ParseMode:      constants.TelegramMarkdown,
		ProtectContent: true,
		ReplyMarkup:    types.ReplyMarkup{InlineKeyboard: inlineKeyboard},
	}

	if botSendMsgErr := s.bot.SendMsg(message); botSendMsgErr != nil {
		return exception.Wrap(
			"s.bot.SendMsg",
			"pick-date-command-telegram-service.go",
			botSendMsgErr,
		)
	}

	return nil
}

func (s *PickDateCommandTelegramService) getCurrentSession(sessionId string) (types.BookingSession, error) {
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
			"pick-date-command-telegram-service.go",
			findOneErr,
		)
	}

	return session, nil
}

func (s *PickDateCommandTelegramService) ackToTelegramClient(callbackQueryId string) error {
	return s.bot.AnswerCallbackQuery(types.AnswerCallbackQuery{CallbackQueryId: callbackQueryId})
}

func (s *PickDateCommandTelegramService) updateSession(actualSession types.BookingSession, serviceId string) error {
	updatedSession := types.BookingSession{
		Id:         actualSession.Id,
		BusinessId: actualSession.BusinessId,
		ChatId:     actualSession.ChatId,
		ServiceId:  serviceId,
		Date:       "",
		Hour:       "",
		CreatedAt:  actualSession.CreatedAt,
		UpdatedAt:  time.Now().UTC().Format(time.DateTime), //We refresh the created at on purpose
		Ttl:        actualSession.Ttl,
	}

	if err := s.sessionRepository.Update(updatedSession); err != nil {
		return exception.Wrap(
			"s.repository.Update",
			"pick-date-command-telegram-service.go",
			err,
		)
	}

	return nil
}

func (s *PickDateCommandTelegramService) similarSessionCriteria(date, openTime, closeTime time.Time) types.Criteria {
	return types.NewCriteria().
		Equal("date", date.Format(time.DateTime)).
		GreaterThanOrEqual("hour", openTime.Format(time.TimeOnly)).
		LessThanOrEqual("hour", closeTime.Format(time.TimeOnly))
}

func (s *PickDateCommandTelegramService) bookingCriteria(sessionId string) types.Criteria {
	return types.NewCriteria().
		Equal("session_id", sessionId)
}

func (s *PickDateCommandTelegramService) dateHasAvailableSlots(date time.Time) (bool, error) {
	openTime := time.Date(0, 0, 0, 8, 0, 0, 0, time.UTC)
	closeTime := time.Date(0, 0, 0, 19, 0, 0, 0, time.UTC)

	totalHoursOpened := closeTime.Sub(openTime)

	openSessionsCriteria := s.similarSessionCriteria(date, openTime, closeTime)

	sessions, findSessionErr := s.sessionRepository.Find(openSessionsCriteria)

	if findSessionErr != nil {
		return false, exception.Wrap(
			"s.repository.Find",
			"pick-date-command-telegram-service.go",
			findSessionErr,
		)
	}

	bookingsCounter := 0

	for _, session := range sessions {
		if expiredSession := session.EnsureIsValid(); expiredSession != nil {
			bookingCriteria := s.bookingCriteria(session.Id)

			_, findOneBookingErr := s.bookingRepository.FindOne(bookingCriteria)

			var bookingNotFoundErr exception.HastypalError

			if findOneBookingErr != nil && errors.As(findOneBookingErr, &bookingNotFoundErr) {
				if bookingNotFoundErr.IsDomain() {
					continue
				}

				return false, exception.Wrap(
					"s.bookingRepository.FindOne",
					"pick-date-command-telegram-service.go",
					findOneBookingErr,
				)
			}

			bookingsCounter++

			continue
		}

		bookingsCounter++
	}

	hasAvailableSlots := bookingsCounter != int(totalHoursOpened.Hours())

	return hasAvailableSlots, nil
}
