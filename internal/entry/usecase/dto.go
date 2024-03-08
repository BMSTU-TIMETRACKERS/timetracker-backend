package usecase

import (
	"time"
)

type Entry struct {
	ID        int64
	UserID    int64
	ProjectID int64
	Name      string
	TimeStart time.Time
	TimeEnd   time.Time

	// Поля только для чтения.
	ProjectName string
}
