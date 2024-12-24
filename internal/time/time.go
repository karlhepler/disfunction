package time

import "time"

var Now = time.Now

func StartOfDay(t time.Time) time.Time {
	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		0, 0, 0, 0, // hour, minute, second, nanosecond
		t.Location(),
	)
}
