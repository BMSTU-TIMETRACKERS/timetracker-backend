package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	ErrEntryNotFound       = errors.New("entry not found")
	ErrProjectInfoNotFound = errors.New("error project info not found")
)

type repository struct {
	db    *sqlx.DB
	close func() error
}

func (r *repository) CreateEntry(ctx context.Context, entry Entry) (int64, error) {
	query := `INSERT INTO entries
				(
					id,
					user_id,
					project_id,
					name,
					time_start,
					time_end
				) VALUES ($1, $2, $3, $4, $5, $6);`

	res, err := r.db.Exec(
		query,
		entry.ID,
		entry.UserID,
		entry.ProjectID,
		entry.Name,
		entry.TimeStart,
		entry.TimeEnd,
	)

	if err != nil {
		return 0, fmt.Errorf("exec: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("last insert id: %v", err)
	}

	return id, nil
}
