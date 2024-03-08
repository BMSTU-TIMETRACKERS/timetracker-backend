package repository

import (
	"time"
)

type ProjectInfo struct {
	ID   int64  `gorm:"column:id;default:null"`
	Name string `gorm:"column:name;default:null"`
}

type Entry struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	ProjectID int64     `db:"project_id;default:null"`
	Name      string    `db:"name"`
	TimeStart time.Time `db:"time_start"`
	TimeEnd   time.Time `db:"time_end"`
}

func (Entry) TableName() string {
	return "entry"
}
