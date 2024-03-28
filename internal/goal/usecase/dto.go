package delivery

import "time"

type Goal struct {
	ID          int64
	ProjectID   int64
	UserID      int64
	TimeSeconds int64
	Name        string
	DateStart   time.Time
	DateEnd     time.Time

	DurationSeconds float64
	Percent         float64
}
