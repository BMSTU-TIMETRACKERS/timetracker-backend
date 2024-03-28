package delivery

import "time"

type CreateGoalIn struct {
	Name        string    `json:"name" validate:"required"`
	ProjectID   int64     `json:"project_id" validate:"required"`
	TimeSeconds int64     `json:"time_seconds" validate:"required"`
	DateStart   time.Time `json:"date_start" validate:"required"`
	DateEnd     time.Time `json:"date_end" validate:"required"`
}

type CreateGoalOut struct {
	ID int64 `json:"id" validate:"required"`
}

type GoalOut struct {
	ID              int64     `json:"id"`
	ProjectID       int64     `json:"project_id"`
	UserID          int64     `json:"user_id"`
	TimeSeconds     int64     `json:"time_seconds"`
	Name            string    `json:"name"`
	DateStart       time.Time `json:"date_start"`
	DateEnd         time.Time `json:"date_end"`
	DurationSeconds float64   `json:"duration_seconds"`
	Percent         float64   `json:"percent"`
}