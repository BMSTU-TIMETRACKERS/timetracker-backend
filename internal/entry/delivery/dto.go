package delivery

import "time"

type CreateEntryIn struct {
	ProjectID int64     `json:"project_id" validate:"required"`
	Name      string    `json:"name"`
	TimeStart time.Time `json:"time_start" validate:"required"`
	TimeEnd   time.Time `json:"time_end" validate:"required"`
}

type CreateEntryOut struct {
	ID int64 `json:"id" validate:"required"`
}

type EntryOut struct {
	ID        int64     `json:"id"`
	ProjectID int64     `json:"project_id"`
	Name      string    `json:"name"`
	TimeStart time.Time `json:"time_start"`
	TimeEnd   time.Time `json:"time_end"`
}
