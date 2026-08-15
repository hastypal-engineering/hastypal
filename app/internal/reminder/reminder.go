package reminder

import "time"

type Reminder struct {
	ID          string
	BookingID   string
	ScheduledAt time.Time
	Sent        bool
	SentAt      string
	DateAdd     time.Time
	DateUpd     time.Time
}
