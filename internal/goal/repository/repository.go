package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	ErrEntryNotFound = errors.New("entry not found")
)

type Repository struct {
	db    *sqlx.DB
	close func() error
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
		close: func() error {
			return db.Close()
		},
	}
}

func (r *Repository) CreateGoal(_ context.Context, goal Goal) (int64, error) {
	query := `INSERT INTO goals
				(
				 	project_id,
					user_id,
				 	time_seconds,
					name,
					date_start,
					date_end
				) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`

	var id int64
	err := r.db.QueryRow(
		query,
		goal.ProjectID,
		goal.UserID,
		goal.TimeSeconds,
		goal.Name,
		goal.DateStart,
		goal.DateEnd,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("exec: %v", err)
	}

	return id, nil
}

func (r *Repository) GetGoals(_ context.Context, userID, projectID int64) ([]Goal, error) {
	rows, err := r.db.Query(`
SELECT g.id,
       g.project_id,
       g.user_id,
       g.name,
       g.time_seconds,
       g.date_start,
       g.date_end,
       COALESCE((SELECT JSON_AGG(
                       JSON_BUILD_OBJECT(
                               'entry_start', e.time_start::timestamptz,
                               'entry_end', e.time_end::timestamptz
                       )
               )
        FROM entries e
        WHERE e.project_id = g.project_id
          AND (e.time_end::date <= g.date_end AND e.time_end::date >= g.date_start OR
               e.time_start::date >= g.date_start AND e.time_start::date <= g.date_end OR
               e.time_start::date < g.date_start AND e.time_end::date > g.date_end)), JSON_ARRAY()) AS entries
FROM goals g
WHERE g.user_id = $1 AND g.project_id = $2
ORDER BY g.date_start`, userID, projectID)

	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	var goals []Goal
	for rows.Next() {
		var goal Goal
		var entriesJSON string

		if err = rows.Scan(
			&goal.ID,
			&goal.ProjectID,
			&goal.UserID,
			&goal.Name,
			&goal.TimeSeconds,
			&goal.DateStart,
			&goal.DateEnd,
			&entriesJSON,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", rows.Err())
		}

		err = json.Unmarshal([]byte(entriesJSON), &goal.Entries)
		if err != nil {
			return nil, fmt.Errorf("unmarshal: %w", err)
		}

		goals = append(goals, goal)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows err: %w", rows.Err())
	}

	if len(goals) == 0 {
		return nil, ErrEntryNotFound
	}

	return goals, nil
}
