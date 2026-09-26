// Package business
package business

import (
	"time"

	"github.com/rotisserie/eris"
)

type Business struct {
	ID           int
	PublicID     string
	Name         string
	ContactPhone string
	Email        string
	Address      string
	Country      string
	Password     string
	Lang         string
	DateAdd      time.Time
	DateUpd      time.Time
}

type ServiceCatalog struct {
	ID          int
	Name        string
	Description string
	Price       float64
	Currency    string
	Duration    string
	BusinessID  int
	DateAdd     time.Time
	DateUpd     time.Time
}

type BusinessSchedule struct {
	WeeklySchedule []*OperatingDay
	Holidays       []*Holiday
	Overrides      []*ScheduleOverride
}

type OperatingDay struct {
	DayOfWeek time.Weekday
	IsClosed  bool
	TimeSlots []TimeSlot
}

type TimeSlot struct {
	OpenTime  string
	CloseTime string
}

func parseTimeSlot(slot TimeSlot) (time.Time, time.Time, error) {
	openTime, err := time.Parse(timeOnlyFormat, slot.OpenTime)
	if err != nil {
		return time.Time{}, time.Time{}, eris.Wrapf(err, "Invalid time slot open time %s", slot.OpenTime)
	}

	closeTime, err := time.Parse(timeOnlyFormat, slot.CloseTime)
	if err != nil {
		return time.Time{}, time.Time{}, eris.Wrapf(err, "Invalid time slot close time %s", slot.CloseTime)
	}

	return openTime, closeTime, nil
}

type Holiday struct {
	ID          int
	Name        string
	StartDate   time.Time
	EndDate     time.Time
	IsRecurring bool
	IsClosed    bool
	Type        string
}

type ScheduleOverride struct {
	Date      time.Time
	IsClosed  bool
	TimeSlots []TimeSlot
	Reason    string
}
