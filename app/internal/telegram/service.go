package telegram

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/adriein/hastypal/internal/booking"
	"github.com/adriein/hastypal/internal/business"
	"github.com/adriein/hastypal/internal/translation"
	"github.com/adriein/hastypal/pkg/constants"
	"github.com/adriein/hastypal/pkg/helper/array"
	"github.com/adriein/hastypal/pkg/helper/reflection"
	"github.com/rotisserie/eris"
)

type TelegramService interface {
	HandleMessage(ctx context.Context, update TelegramUpdate) error
}

type Service struct {
	business business.BusinessService
	booking  booking.BookingService
	lang     translation.TranslationService
	bot      *TelegramBot
}

func NewService(
	business business.BusinessService,
	booking booking.BookingService,
	lang translation.TranslationService,
	bot *TelegramBot,
) *Service {
	return &Service{
		business: business,
		booking:  booking,
		lang:     lang,
		bot:      bot,
	}
}

/*
================================================================================
TELEGRAM UPDATE HANLDER
================================================================================
*/

func (s *Service) HandleMessage(ctx context.Context, update TelegramUpdate) error {
	if reflection.HasField(update, constants.TelegramMessageField) {
		if err := s.resolveBotCommand(ctx, update); err != nil {
			return eris.Wrap(err, "Error resolving bot command")
		}

		return nil
	}

	if err := s.resolveCallbackQueryCommand(ctx, update); err != nil {
		return eris.Wrap(err, "Error resolving callback query command")
	}

	return nil
}

func (s *Service) resolveBotCommand(ctx context.Context, update TelegramUpdate) error {
	if !reflection.HasField(update, constants.TelegramMessageField) {
		return nil
	}

	txtArr := strings.Split(update.Message.Text, " ")
	command := txtArr[0]

	switch command {
	case constants.StartCommand:
		return s.startConversation(ctx, update)
	}

	return nil
}

func (s *Service) resolveCallbackQueryCommand(ctx context.Context, update TelegramUpdate) error {
	if !reflection.HasField(update, constants.TelegramCallbackQueryField) {
		return nil
	}

	url, err := url.Parse(update.CallbackQuery.Data)

	if err != nil {
		return eris.Wrap(err, "Error parsing the callback query url")
	}

	switch url.Path {
	case constants.ServiceCommand:
		return s.showServices(ctx, update)
	case constants.DatesCommand:
		return s.showDates(ctx, update)
	case constants.HoursCommand:
		return nil
	case constants.ConfirmationCommand:
		return nil
	case constants.FinishCommand:
		return nil
	}

	return nil
}

/*
================================================================================
TELEGRAM START CONVERSATION COMMAND
================================================================================
*/

func (s *Service) startConversation(ctx context.Context, update TelegramUpdate) error {
	var markdownText strings.Builder

	businessIdRaw := strings.ReplaceAll(update.Message.Text, "/start ", "")

	businessID, err := strconv.Atoi(businessIdRaw)

	if err != nil {
		return eris.Wrap(err, "Error converting business ID to int")
	}

	business, err := s.business.GetBusinessByID(ctx, businessID)

	if err != nil {
		return eris.Wrap(err, "Error fetching business")
	}

	sessionID, err := s.booking.InitSession(ctx, businessID, update.Message.Chat.Id)

	if err != nil {
		return eris.Wrap(err, "Error creating a session for this conversation")
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

	buttons := make([]KeyboardButton, len(services))

	for i, serv := range services {
		markdownText.WriteString(fmt.Sprintf("%s %s\n\n", emoji, serv))

		buttons[i] = KeyboardButton{
			Text:         fmt.Sprintf("%s 📅", services[i]),
			CallbackData: fmt.Sprintf("/dates?session=%s&service=%s&page=0", sessionID, "test-short"),
		}
	}

	inlineKeyboard := array.Chunk(buttons, 1)

	message := TelegramMessage{
		ChatId:         update.Message.Chat.Id,
		Text:           markdownText.String(),
		ParseMode:      constants.TelegramMarkdown,
		ProtectContent: true,
		ReplyMarkup:    ReplyMarkup{InlineKeyboard: inlineKeyboard},
	}

	bookingMessage := BookingTelegramMessage{
		BusinessName:     business.Name,
		BookingSessionId: sessionID,
		Message:          message,
	}

	if err := s.bot.SendMsg(bookingMessage); err != nil {
		return eris.Wrap(err, "Error sending message to telegram")
	}

	return nil
}

/*
================================================================================
TELEGRAM SHOW SERVICES COMMAND
================================================================================
*/

func (s *Service) showServices(ctx context.Context, update TelegramUpdate) error {
	ack := AnswerCallbackQuery{CallbackQueryId: update.CallbackQuery.Id}

	if err := s.bot.AnswerCallbackQuery(ack); err != nil {
		return eris.Wrap(err, "Error acking telegram conversation")
	}

	var markdownText strings.Builder

	parsedUrl, err := url.Parse(update.CallbackQuery.Data)

	if err != nil {
		return eris.Wrap(err, "Error parsing url")
	}

	queryParams := parsedUrl.Query()

	sessionID := queryParams.Get("sessionId")

	session, err := s.booking.GetCurrentSession(ctx, sessionID)

	if err != nil {
		return eris.Wrap(err, "Error fetching current booking session")
	}

	business, err := s.business.GetBusinessByID(ctx, session.BusinessId)

	if err != nil {
		return eris.Wrap(err, "Error fetching business")
	}

	if err := session.EnsureIsValid(); err != nil {
		message := TelegramMessage{ChatId: update.CallbackQuery.From.Id}

		expiredSessionMessage := message.SessionExpired()

		bookingExpiredSessionMessage := BookingTelegramMessage{
			BusinessName:     business.Name,
			BookingSessionId: session.Id,
			Message:          expiredSessionMessage,
		}

		if err := s.bot.SendMsg(bookingExpiredSessionMessage); err != nil {
			return eris.Wrap(err, "Error sending message to telegram")
		}

		return nil
	}

	if err := s.booking.RefreshSession(ctx, session); err != nil {
		return eris.Wrap(err, "Error refreshing the current session")
	}

	services := []string{
		"Corte de pelo y barba express 18€",
		"Corte de pelo y barba premium 22€",
	}

	emoji := "![🔸](tg://emoji?id=5368324170671202286)"

	markdownText.WriteString("*Te muestro a continuación los servicios que ofrecemos:*\n\n")

	buttons := make([]KeyboardButton, len(services))

	for i, serv := range services {
		markdownText.WriteString(fmt.Sprintf("%s %s\n\n", emoji, serv))

		buttons[i] = KeyboardButton{
			Text:         fmt.Sprintf("%s 📅", services[i]),
			CallbackData: fmt.Sprintf("/dates?session=%s&service=%s&page=0", session.Id, "test-short"),
		}
	}

	inlineKeyboard := array.Chunk(buttons, 1)

	message := TelegramMessage{
		ChatId:         update.CallbackQuery.From.Id,
		Text:           markdownText.String(),
		ParseMode:      constants.TelegramMarkdown,
		ProtectContent: true,
		ReplyMarkup:    ReplyMarkup{InlineKeyboard: inlineKeyboard},
	}

	bookingMessage := BookingTelegramMessage{
		BusinessName:     business.Name,
		BookingSessionId: session.Id,
		Message:          message,
	}

	if err := s.bot.SendMsg(bookingMessage); err != nil {
		return eris.Wrap(err, "Error sending telegram msg")
	}

	return nil
}

/*
================================================================================
TELEGRAM SHOW DATES COMMAND
================================================================================
*/

func (s *Service) showDates(ctx context.Context, update TelegramUpdate) error {
	ack := AnswerCallbackQuery{CallbackQueryId: update.CallbackQuery.Id}

	if err := s.bot.AnswerCallbackQuery(ack); err != nil {
		return eris.Wrap(err, "Error acking telegram conversation")
	}

	var markdownText strings.Builder

	parsedUrl, err := url.Parse(update.CallbackQuery.Data)

	if err != nil {
		return eris.Wrap(err, "Error parsing url")
	}

	queryParams := parsedUrl.Query()

	sessionID := queryParams.Get("session")
	serviceId := queryParams.Get("service")
	page := queryParams.Get("page")

	currentPage, err := strconv.Atoi(page)

	if err != nil {
		return eris.Wrap(err, "Error converting string to int")
	}

	session, err := s.booking.GetCurrentSession(ctx, sessionID)

	if err != nil {
		return eris.Wrap(err, "Error fetching current booking session")
	}

	business, err := s.business.GetBusinessByID(ctx, session.BusinessId)

	if err != nil {
		return eris.Wrap(err, "Error fetching business")
	}

	if err := session.EnsureIsValid(); err != nil {
		message := TelegramMessage{ChatId: update.CallbackQuery.From.Id}

		expiredSessionMessage := message.SessionExpired()

		bookingExpiredSessionMessage := BookingTelegramMessage{
			BusinessName:     business.Name,
			BookingSessionId: session.Id,
			Message:          expiredSessionMessage,
		}

		if err := s.bot.SendMsg(bookingExpiredSessionMessage); err != nil {
			return eris.Wrap(err, "Error sending message to telegram")
		}

		return nil
	}

	if err := s.booking.RefreshSession(ctx, session); err != nil {
		return eris.Wrap(err, "Error refreshing the current session")
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

	location, err := time.LoadLocation("Europe/Madrid")

	if err != nil {
		return eris.Wrap(err, "Error loading time location")
	}

	startDate := time.Now().In(location)
	startDateWithHour := time.Date(
		startDate.Year(),
		startDate.Month(),
		startDate.Day(),
		07,
		0,
		0,
		0,
		location,
	)

	startDateWithHour = startDateWithHour.AddDate(0, 0, constants.DaysPerPage*currentPage)

	buttons := make([]KeyboardButton, 15)

	for i := 0; i < 15; i++ {
		newDate := startDateWithHour.AddDate(0, 0, i)

		sessions, err := s.booking.GetSessionsOnDate(ctx, newDate)

		if err != nil {
			return eris.Wrap(err, "Error getting all sessions for a specific date")
		}

		schedule := booking.NewDaySchedule(newDate)
		schedule.ApplyActiveSessions(sessions)

		if !schedule.HasAnyAvailableSlot() {
			continue
		}

		dateParts := strings.Split(newDate.Format(time.RFC822), " ")

		day := dateParts[0]
		month := s.lang.GetSpanishMonthShortForm(newDate.Month())

		buttons[i] = KeyboardButton{
			Text:         fmt.Sprintf("%s %s", day, month),
			CallbackData: fmt.Sprintf("/hours?session=%s&date=%s", sessionID, newDate.Format(time.DateOnly)),
		}
	}

	inlineKeyboard := array.Chunk(buttons, 3)

	inlineKeyboard = s.addNavigationButtons(session.Id, serviceId, currentPage, inlineKeyboard)

	message := TelegramMessage{
		ChatId:         update.CallbackQuery.From.Id,
		Text:           markdownText.String(),
		ParseMode:      constants.TelegramMarkdown,
		ProtectContent: true,
		ReplyMarkup:    ReplyMarkup{InlineKeyboard: inlineKeyboard},
	}

	bookingMessage := BookingTelegramMessage{
		BusinessName:     business.Name,
		BookingSessionId: session.Id,
		Message:          message,
	}

	if err := s.bot.SendMsg(bookingMessage); err != nil {
		return eris.Wrap(err, "Error sending telegram msg")
	}

	return nil
}

func (s *Service) addNavigationButtons(
	sessionID string,
	serviceID string,
	currentPage int,
	inlineKeyboard [][]KeyboardButton,
) [][]KeyboardButton {
	navigationButtons := make([]KeyboardButton, 0, 3)

	if currentPage == constants.MinAllowedDatePage {
		moreDaysButton := KeyboardButton{
			Text:         "Más fechas",
			CallbackData: fmt.Sprintf("/dates?session=%s&service=%s&page=%d", sessionID, serviceID, currentPage+1),
		}

		backButton := KeyboardButton{
			Text:         "Atrás",
			CallbackData: fmt.Sprintf("/service?session=%s", sessionID),
		}

		navigationButtons = append(navigationButtons, moreDaysButton, backButton)

		navigationKeyboard := array.Chunk(navigationButtons, 1)

		return append(inlineKeyboard, navigationKeyboard...)
	}

	if currentPage == constants.MaxAllowedDatePage {
		lessDaysButton := KeyboardButton{
			Text:         "Menos fechas",
			CallbackData: fmt.Sprintf("/dates?session=%s&service=%s&page=%d", sessionID, serviceID, currentPage-1),
		}

		backButton := KeyboardButton{
			Text:         "Atrás",
			CallbackData: fmt.Sprintf("/service?session=%s", sessionID),
		}

		navigationButtons = append(navigationButtons, lessDaysButton, backButton)

		navigationKeyboard := array.Chunk(navigationButtons, 1)

		return append(inlineKeyboard, navigationKeyboard...)
	}

	lessDaysButton := KeyboardButton{
		Text:         "Menos fechas",
		CallbackData: fmt.Sprintf("/dates?session=%s&service=%s&page=%d", sessionID, serviceID, currentPage-1),
	}

	moreDaysButton := KeyboardButton{
		Text:         "Más fechas",
		CallbackData: fmt.Sprintf("/dates?session=%s&service=%s&page=%d", sessionID, serviceID, currentPage+1),
	}

	backButton := KeyboardButton{
		Text:         "Atrás",
		CallbackData: fmt.Sprintf("/service?session=%s", sessionID),
	}

	navigationButtons = append(navigationButtons, lessDaysButton, moreDaysButton, backButton)

	navigationKeyboard := array.Chunk(navigationButtons, 1)

	return append(inlineKeyboard, navigationKeyboard...)
}

/*
================================================================================
TELEGRAM SHOW HOURS COMMAND
================================================================================
*/

func (s *Service) showHours(ctx context.Context, update TelegramUpdate) error {
	ack := AnswerCallbackQuery{CallbackQueryId: update.CallbackQuery.Id}

	if err := s.bot.AnswerCallbackQuery(ack); err != nil {
		return eris.Wrap(err, "Error acking telegram conversation")
	}

	var markdownText strings.Builder

	parsedUrl, err := url.Parse(update.CallbackQuery.Data)

	if err != nil {
		return eris.Wrap(err, "Error parsing URL")
	}

	queryParams := parsedUrl.Query()

	stringSelectedDate := fmt.Sprintf("%s %s", queryParams.Get("date"), "07:00:00")
	sessionID := queryParams.Get("session")

	selectedDate, err := time.Parse(time.DateTime, stringSelectedDate)

	if err != nil {
		return eris.Wrap(err, "Error parsing time")
	}

	session, err := s.booking.GetCurrentSession(ctx, sessionID)

	if err != nil {
		return eris.Wrap(err, "Error fetching current booking session")
	}

	business, err := s.business.GetBusinessByID(ctx, session.BusinessId)

	if err != nil {
		return eris.Wrap(err, "Error fetching business")
	}

	if err := session.EnsureIsValid(); err != nil {
		message := TelegramMessage{ChatId: update.CallbackQuery.From.Id}

		expiredSessionMessage := message.SessionExpired()

		bookingExpiredSessionMessage := BookingTelegramMessage{
			BusinessName:     business.Name,
			BookingSessionId: session.Id,
			Message:          expiredSessionMessage,
		}

		if err := s.bot.SendMsg(bookingExpiredSessionMessage); err != nil {
			return eris.Wrap(err, "Error sending message to telegram")
		}

		return nil
	}

	if err := s.booking.RefreshSession(ctx, session); err != nil {
		return eris.Wrap(err, "Error refreshing the current session")
	}

	dateParts := strings.Split(selectedDate.Format(time.RFC822), " ")

	day := dateParts[0]
	month := s.lang.GetSpanishMonthShortForm(selectedDate.Month())

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

	buttons := make([]KeyboardButton, 12)

	for i := 8; i <= len(buttons)+7; i++ {
		hour := fmt.Sprintf("%02d:00", i)

		session, err := s.booking.GetSessionOnHour(ctx, selectedDate)

		if err != nil {
			return eris.Wrap(err, "Error getting all sessions for a specific date")
		}

		if session != nil {
			continue
		}

		buttons[i-8] = KeyboardButton{
			Text: fmt.Sprintf("%s", hour),
			CallbackData: fmt.Sprintf(
				"/confirmation?session=%s&hour=%s",
				sessionID,
				hour,
			),
		}
	}

	backButton := KeyboardButton{
		Text:         "Atrás",
		CallbackData: fmt.Sprintf("/dates?session=%s&service=%s", session.Id, "test-short"),
	}

	buttons = append(buttons, backButton)

	inlineKeyboard := array.Chunk(buttons, 3)

	message := TelegramMessage{
		ChatId:         update.CallbackQuery.From.Id,
		Text:           markdownText.String(),
		ParseMode:      constants.TelegramMarkdown,
		ProtectContent: true,
		ReplyMarkup:    ReplyMarkup{InlineKeyboard: inlineKeyboard},
	}

	bookingMessage := BookingTelegramMessage{
		BusinessName:     business.Name,
		BookingSessionId: session.Id,
		Message:          message,
	}

	if err := s.bot.SendMsg(bookingMessage); err != nil {
		return eris.Wrap(err, "Error sending telegram msg")
	}

	return nil
}
