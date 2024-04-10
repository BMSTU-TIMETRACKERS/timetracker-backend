package delivery

import "time"

type CreateEntryIn struct {
	ProjectID int64     `json:"project_id" validate:"required" example:"1"`                    // Идентификатор проекта.
	Name      string    `json:"name" example:"task1"`                                          // Название записи.
	TimeStart time.Time `json:"time_start" validate:"required" example:"2024-03-23T15:04:05Z"` // Время начала записи.
	TimeEnd   time.Time `json:"time_end" validate:"required" example:"2024-03-23T19:04:05Z"`   // Время окончания записи.
}

type CreateEntryOut struct {
	ID int64 `json:"id" validate:"required" example:"1"` // Идентификатор записи.
}

type EntryOut struct {
	ID          int64     `json:"id" example:"1"`                            // Идентификатор записи.
	ProjectID   int64     `json:"project_id" example:"1"`                    // Идентификатор проекта.
	ProjectName string    `json:"project_name" example:"work"`               // Название проекта.
	Name        string    `json:"name" example:"task1"`                      // Название записи.
	TimeStart   time.Time `json:"time_start" example:"2024-03-23T15:04:05Z"` // Время начала записи.
	TimeEnd     time.Time `json:"time_end" example:"2024-03-23T19:04:05Z"`   // Время окончания записи.
}
