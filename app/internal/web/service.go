package web

import (
	"github.com/adriein/hastypal/internal/booking"
	"github.com/adriein/hastypal/internal/business"
	"github.com/adriein/hastypal/internal/google"
	"github.com/adriein/hastypal/internal/reminder"
	"github.com/adriein/hastypal/internal/translation"
)

type WebService interface{}

type Service struct {
	business business.BusinessService
	booking  booking.BookingService
	lang     translation.TranslationService
	reminder reminder.ReminderService
	google   google.GoogleService
}

func NewService(
	business business.BusinessService,
	booking booking.BookingService,
	lang translation.TranslationService,
	reminder reminder.ReminderService,
	google google.GoogleService,
) *Service {
	return &Service{
		business: business,
		booking:  booking,
		lang:     lang,
		reminder: reminder,
		google:   google,
	}
}
