package repository

type Project struct {
	ID     int64  `db:"id"`
	Name   string `db:"name"`
	UserID int64  `db:"user_id"`
}
