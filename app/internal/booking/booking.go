package booking

import (
	"time"

	"github.com/rotisserie/eris"
)

var BookingSessionExpired = eris.New("Booking session expired")

type Booking struct {
	Id         string `json:"id"`
	SessionId  string `json:"sessionId"`
	BusinessId string `json:"businessId"`
	ServiceId  string `json:"serviceId"`
	When       string `json:"when"`
	CreatedAt  string `json:"createdAt"`
}

type Session struct {
	Id         string `json:"id"`
	BusinessId int    `json:"businessId"`
	ChatId     int    `json:"chatId"`
	ServiceId  string `json:"serviceId"`
	Date       string `json:"date"`
	Hour       string `json:"hour"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
	Ttl        int64  `json:"ttl"`
}

func (s *Session) EnsureIsValid() error {
	updatedAt, err := time.Parse(time.DateTime, s.UpdatedAt)

	if err != nil {
		return eris.Wrap(err, "Error parsing datetime")
	}

	maxAllowedDate := updatedAt.Add(time.Duration(300000) * time.Millisecond)

	if maxAllowedDate.Before(time.Now().UTC()) {
		return BookingSessionExpired
	}

	return nil
}
