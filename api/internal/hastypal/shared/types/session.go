package types

import (
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"time"
)

type BookingSession struct {
	Id         string   `json:"id" db:"id"`
	BusinessId string   `json:"businessId" db:"business_id"`
	ChatId     int      `json:"chatId" db:"chat_id"`
	ServiceId  string   `json:"serviceId" db:"service_id"`
	Date       string   `json:"date" db:"date"`
	Hour       string   `json:"hour" db:"hour"`
	CreatedAt  string   `json:"createdAt" db:"created_at"`
	UpdatedAt  string   `json:"updatedAt" db:"updated_at"`
	Ttl        int64    `json:"ttl" db:"ttl"`
	Booking    *Booking `json:"booking" db:"booking"`
}

// Database

func (s *BookingSession) DatabaseMappings() map[string]string {
	return map[string]string{
		"booking": "booking;session_id",
	}
}

func (s *BookingSession) EnsureIsValid() error {
	updatedAt, err := time.Parse(time.DateTime, s.UpdatedAt)

	if err != nil {
		return exception.
			New(err.Error()).
			Trace("time.Parse", "session.go").
			WithValues([]string{s.CreatedAt})
	}

	maxAllowedDate := updatedAt.Add(time.Duration(300000) * time.Millisecond)

	if maxAllowedDate.Before(time.Now().UTC()) {
		return exception.
			New("The session has expired").
			Trace("maxAllowedDate.Before", "session.go").
			WithValues([]string{s.CreatedAt}).
			Domain()
	}

	return nil
}
