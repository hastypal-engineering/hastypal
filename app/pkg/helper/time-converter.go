package helper

import "time"

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
