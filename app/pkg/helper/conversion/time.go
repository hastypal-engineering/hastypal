// Package conversion
package conversion

import (
	"time"

	"github.com/rotisserie/eris"
)

func DateToDB(date time.Time) *time.Time {
	if date.IsZero() {
		return nil
	}

	return &date
}

func HourToDB(hour time.Duration) any {
	if hour == 0 {
		return nil
	}

	return time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC).Add(hour)
}

func HourFromDB(hour time.Time) time.Duration {
	return time.Duration(hour.Hour())*time.Hour +
		time.Duration(hour.Minute())*time.Minute +
		time.Duration(hour.Second())*time.Second
}

// CombineDayAndTime builds a single timestamp using the calendar day of day and
// the hour, minutes and seconds of clock, which is expected to only carry a time.
func CombineDayAndTime(day time.Time, clock time.Time) time.Time {
	return time.Date(
		day.Year(),
		day.Month(),
		day.Day(),
		clock.Hour(),
		clock.Minute(),
		clock.Second(),
		0,
		day.Location(),
	)
}

func StringToTime(s string, format string) (time.Time, error) {
	date, err := time.Parse(format, s)
	if err != nil {
		return time.Now(), eris.Wrapf(err, "Error converting string %s to time with format %s", s, format)
	}

	return date, nil
}

func TimeToMinFromMidnight(t time.Time) int {
	return t.Hour()*60 + t.Minute()
}
