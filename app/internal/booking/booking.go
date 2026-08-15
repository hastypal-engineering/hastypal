package booking

import (
	"time"

	"github.com/rotisserie/eris"
)

var BookingSessionExpired = eris.New("Booking session expired")

type Booking struct {
	ID         string
	SessionID  string
	BusinessID int
	ServiceID  string
	Date       time.Time
	DateAdd    time.Time
	DateUpd    time.Time
}

type Session struct {
	Id         string
	BusinessId int
	ChatId     int
	ServiceId  string
	Date       string
	Hour       string
	Ttl        int64
	SlotIndex  int
	DateAdd    time.Time
	DateUpd    time.Time
}

func (s *Session) EnsureIsValid() error {
	maxAllowedDate := s.DateUpd.Add(time.Duration(300000) * time.Millisecond)

	if maxAllowedDate.Before(time.Now().UTC()) {
		return BookingSessionExpired
	}

	return nil
}

func (s *Session) Refresh() {
	s.DateUpd = time.Now().UTC()
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
