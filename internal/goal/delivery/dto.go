package delivery

import "time"

type CreateGoalIn struct {
	Name        string    `json:"name" validate:"required" example:"Потратить 100часов на разработку"` // Название цели.
	ProjectID   int64     `json:"project_id" validate:"required" example:"1"`                          // Идентификатор проекта.
	TimeSeconds int64     `json:"time_seconds" validate:"required" example:"360000"`                   // Требуемое(целевое) время в секундах.
	DateStart   time.Time `json:"date_start" validate:"required" example:"2024-03-23T00:00:00Z"`       // Дата начала цели.
	DateEnd     time.Time `json:"date_end" validate:"required" example:"2024-04-23T00:00:00Z"`         // Дата окончания цели.
}

type CreateGoalOut struct {
	ID int64 `json:"id" validate:"required" example:"1"` // Идентификатор цели.
}

type GoalOut struct {
	ID              int64     `json:"id" example:"1"`                                  // Идентификатор цели.
	ProjectID       int64     `json:"project_id" example:"1"`                          // Идентификатор проекта.
	UserID          int64     `json:"user_id" example:"1"`                             // Идентификатор пользователя.
	TimeSeconds     int64     `json:"time_seconds" example:"360000"`                   // Требуемое(целевое) время в секундах.
	Name            string    `json:"name" example:"Потратить 100часов на разработку"` // Название цели.
	DateStart       time.Time `json:"date_start" example:"2024-03-23T00:00:00Z"`       // Дата начала цели.
	DateEnd         time.Time `json:"date_end" example:"2024-04-23T00:00:00Z"`         // Дата окончания цели.
	DurationSeconds float64   `json:"duration_seconds" example:"36000"`                // Прогресс: количество секунд потреченных на эту цель
	Percent         float64   `json:"percent" example:"10"`                            // Прогресс: процент выполнения цели.
}
