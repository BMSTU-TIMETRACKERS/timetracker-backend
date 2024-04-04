package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	ErrProjectNotFound = errors.New("project not found")
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

func (r *Repository) CreateProject(ctx context.Context, project Project) (int64, error) {
	query := `INSERT INTO projects
				(
					user_id,
					name
				) VALUES ($1, $2) RETURNING id;`

	var id int64
	err := r.db.QueryRow(
		query,
		project.UserID,
		project.Name,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("query row: %v", err)
	}

	return id, nil
}

func (r *Repository) GetUserProjects(ctx context.Context, userID int64) ([]Project, error) {
	rows, err := r.db.Query(
		`SELECT 
			id,
			user_id,
			name
		FROM projects
		WHERE user_id = $1`, userID)

	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	var projects []Project
	for rows.Next() {
		var project Project
		if err = rows.Scan(
			&project.ID,
			&project.UserID,
			&project.Name,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", rows.Err())
		}

		projects = append(projects, project)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows err: %w", rows.Err())
	}

	if len(projects) == 0 {
		return nil, ErrProjectNotFound
	}

	return projects, nil
}

func (r *Repository) ClearUserData(ctx context.Context, userID int64) error {
	query := `
				DELETE FROM entries;
				DELETE FROM goals;
				DELETE FROM projects;`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("exec context: %v", err)
	}

	return nil
}
