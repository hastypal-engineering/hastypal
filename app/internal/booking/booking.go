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
	Ttl        int64  `json:"ttl"`
	SlotIndex  int    `json:"slotIndex"`
	DateAdd    string `json:"createdAt"`
	DateUpd    string `json:"updatedAt"`
}

func (s *Session) EnsureIsValid() error {
	dateUpd, err := time.Parse(time.DateTime, s.DateUpd)

	if err != nil {
		return eris.Wrap(err, "Error parsing datetime")
	}

	maxAllowedDate := dateUpd.Add(time.Duration(300000) * time.Millisecond)

	if maxAllowedDate.Before(time.Now().UTC()) {
		return BookingSessionExpired
	}

	return nil
}

func (s *Session) Refresh() {
	s.DateUpd = time.Now().UTC().Format(time.DateTime)
}

type Slot struct {
	Index     int
	StartTime time.Time
	IsBooked  bool
	IsLocked  bool
	Available bool
}

type DaySchedule struct {
	WorkDayStart time.Time
	Slots        []Slot
}

func NewDaySchedule(dayStart time.Time) *DaySchedule {
	schedule := &DaySchedule{
		WorkDayStart: dayStart,
	}

	//TODO: 16 is hardcoded and needs to be updated with real business config
	for i := 0; i < 16; i++ {
		slotTime := dayStart.Add(time.Duration(i) * 30)
		schedule.Slots[i] = Slot{
			Index:     i,
			StartTime: slotTime,
			IsBooked:  false,
			IsLocked:  false,
			Available: true,
		}
	}

	return schedule
}

/*func (ds *DaySchedule) ApplyBookings(bookings []Booking) {
	for _, booking := range bookings {
		if booking.SlotIndex >= 0 && booking.SlotIndex < TotalSlots {
			ds.Slots[booking.SlotIndex].IsBooked = true
			ds.Slots[booking.SlotIndex].Available = false
		}
	}
}*/

func (ds *DaySchedule) ApplyActiveSessions(sessions []*Session) {
	for _, session := range sessions {
		if err := session.EnsureIsValid(); err != nil {
			continue
		}

		ds.Slots[session.SlotIndex].IsLocked = true
		ds.Slots[session.SlotIndex].Available = false
	}
}

func (ds *DaySchedule) HasAnyAvailableSlot() bool {
	for _, slot := range ds.Slots {
		if slot.Available {
			return true
		}
	}

	return false
}
