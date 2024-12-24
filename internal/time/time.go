package time

import "time"

const DateOnly = time.DateOnly
const DateTime = time.DateTime
const Kitchen = time.Kitchen
const RFC3339 = time.RFC3339
const TimeOnly = time.TimeOnly

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
