package utils

import (
	"time"
)

func GetDayInterval(date time.Time) (time.Time, time.Time) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, date.Location())
	return start, end
}
