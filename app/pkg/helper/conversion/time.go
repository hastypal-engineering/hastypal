// Package conversion
package conversion

import (
	"time"

	"github.com/rotisserie/eris"
)

func DateToDB(date time.Time) any {
	if date.IsZero() {
		return nil
	}

	return date
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

// DurationToDB returns the amount of seconds of a duration, the value expected
// by make_interval(secs => $n) to persist an INTERVAL column
func DurationToDB(duration time.Duration) float64 {
	return duration.Seconds()
}

// DurationFromDB rebuilds a duration out of the epoch seconds returned by
// EXTRACT(EPOCH FROM $column) when reading an INTERVAL column
func DurationFromDB(seconds float64) time.Duration {
	return time.Duration(seconds * float64(time.Second))
}

func StringToTime(s string, format string) (time.Time, error) {
	date, err := time.Parse(format, s)
	if err != nil {
		return time.Now(), eris.Wrapf(err, "Error converting string %s to time with format %s", s, format)
	}

	return date, nil
}
