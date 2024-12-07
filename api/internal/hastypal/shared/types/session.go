package types

import (
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"time"
)

type BookingSession struct {
	Id         string   `json:"id"`
	BusinessId string   `json:"business_id"`
	ChatId     int      `json:"chat_id"`
	ServiceId  string   `json:"service"`
	Date       string   `json:"date"`
	Hour       string   `json:"hour"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
	Ttl        int64    `json:"ttl"`
	Booking    *Booking `json:"booking"`
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
