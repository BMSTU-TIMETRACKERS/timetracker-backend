package repository

import (
	"context"
	"errors"

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
	return 0, nil
}

func (r *Repository) GetUserProjects(ctx context.Context, userID int64) ([]Project, error) {
	return []Project{}, nil
}
