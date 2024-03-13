package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var (
	ErrEntryNotFound       = errors.New("entry not found")
	ErrProjectInfoNotFound = errors.New("error project info not found")
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

func (r *Repository) CreateEntry(_ context.Context, entry Entry) (int64, error) {
	query := `INSERT INTO entries
				(
					user_id,
					project_id,
					name,
					time_start,
					time_end
				) VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	var id int64
	err := r.db.QueryRow(
		query,
		entry.UserID,
		entry.ProjectID,
		entry.Name,
		entry.TimeStart,
		entry.TimeEnd,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("exec: %v", err)
	}

	return id, nil
}

func (r *Repository) GetUserEntries(_ context.Context, userID int64) ([]Entry, error) {
	rows, err := r.db.Query(
		`SELECT 
			id,
			user_id,
			project_id,
			name,
			time_start,
			time_end
		FROM entries
		WHERE user_id = $1`, userID)

	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	var entries []Entry
	for rows.Next() {
		var entry Entry
		if err = rows.Scan(
			&entry.ID,
			&entry.UserID,
			&entry.ProjectID,
			&entry.Name,
			&entry.TimeStart,
			&entry.TimeEnd,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", rows.Err())
		}

		entries = append(entries, entry)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows err: %w", rows.Err())
	}

	if len(entries) == 0 {
		return nil, ErrEntryNotFound
	}

	return entries, nil
}

func (r *Repository) GetUserEntriesForInterval(
	_ context.Context,
	userID int64,
	start time.Time,
	end time.Time) ([]Entry, error) {
	rows, err := r.db.Query(
		`SELECT 
			id,
			user_id,
			project_id,
			name,
			time_start,
			time_end
		FROM entries
		WHERE user_id = $1 AND time_start BETWEEN $2 AND $3`, userID, start, end)

	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	var entries []Entry
	for rows.Next() {
		var entry Entry
		if err = rows.Scan(
			&entry.ID,
			&entry.UserID,
			&entry.ProjectID,
			&entry.Name,
			&entry.TimeStart,
			&entry.TimeEnd,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", rows.Err())
		}

		entries = append(entries, entry)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows err: %w", rows.Err())
	}

	if len(entries) == 0 {
		return nil, ErrEntryNotFound
	}

	return entries, nil
}

func (r *Repository) GetProjectsInfo(_ context.Context, projectIDs []int64) ([]ProjectInfo, error) {
	rows, err := r.db.Query(
		`SELECT 
			id,
			name
		FROM projects
		WHERE id = ANY($1)`, pq.Array(projectIDs))

	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	var projectInfos []ProjectInfo
	for rows.Next() {
		var info ProjectInfo
		if err = rows.Scan(
			&info.ID,
			&info.Name,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", rows.Err())
		}

		projectInfos = append(projectInfos, info)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows err: %w", rows.Err())
	}

	if len(projectInfos) == 0 {
		return nil, ErrProjectInfoNotFound
	}

	return projectInfos, nil
}
