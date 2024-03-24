package delivery

import "time"

type Entry struct {
	TimeStart time.Time `json:"entry_start"`
	TimeEnd   time.Time `json:"entry_end"`
}

type Goal struct {
	ID          int64     `db:"id"`
	ProjectID   int64     `db:"project_id"`
	UserID      int64     `db:"user_id"`
	TimeSeconds int64     `db:"time_seconds"`
	Name        string    `db:"name"`
	DateStart   time.Time `db:"date_start"`
	DateEnd     time.Time `db:"date_end"`

	Entries []Entry `db:"entries"`
}
