package types

import "time"

type BookingSession struct {
	Id         string `json:"id"`
	BusinessId string `json:"business_id"`
	ChatId     int    `json:"chat_id"`
	ServiceId  string `json:"service"`
	Date       string `json:"date"`
	Hour       string `json:"hour"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	Ttl        int64  `json:"ttl"`
}

func (s *BookingSession) EnsureIsValid() error {
	updatedAt, err := time.Parse(time.DateTime, s.UpdatedAt)

	if err != nil {
		return ApiError{
			Msg:      err.Error(),
			Function: "time.Parse",
			File:     "types/session.go",
			Values:   []string{s.CreatedAt},
		}
	}

	maxAllowedDate := updatedAt.Add(time.Duration(300000) * time.Millisecond)

	if maxAllowedDate.Before(time.Now().UTC()) {
		return ApiError{
			Msg:      "The session has expired",
			Function: "maxAllowedDate.Before",
			File:     "types/session.go",
			Domain:   true,
		}
	}

	return nil
}
