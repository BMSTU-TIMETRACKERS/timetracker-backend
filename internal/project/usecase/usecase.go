package usecase

import (
	"context"
	"errors"
	"fmt"

	repo "github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/project/repository"
)

var ErrProjectNotFound = errors.New("project not found")

type repository interface {
	CreateProject(ctx context.Context, project repo.Project) (int64, error)
	GetUserProjects(ctx context.Context, userID int64) ([]repo.Project, error)
}

type Usecase struct {
	repository repository
}

func NewUsecase(repository repository) *Usecase {
	return &Usecase{
		repository: repository,
	}
}

func (u *Usecase) CreateProject(ctx context.Context, project Project) (int64, error) {
	id, err := u.repository.CreateProject(ctx, convertToRepoProject(project))

	if err != nil {
		return 0, fmt.Errorf("repo create project: %v", err)
	}

	return id, err
}

func convertToRepoProject(project Project) repo.Project {
	return repo.Project{
		ID:     project.ID,
		UserID: project.UserID,
		Name:   project.Name,
	}
}

func (u *Usecase) GetUserProjects(ctx context.Context, userID int64) ([]Project, error) {
	repoProjects, err := u.repository.GetUserProjects(ctx, userID)
	if err != nil {
		if errors.Is(err, repo.ErrProjectNotFound) {
			return []Project{}, nil
		}
		return nil, fmt.Errorf("repo get user projects: %v", err)
	}

	projects := convertToProjects(repoProjects)

	return projects, nil
}

func convertToProjects(projects []repo.Project) []Project {
	repoProjects := make([]Project, 0, len(projects))
	for _, project := range projects {
		repoProjects = append(repoProjects, convertToProject(project))
	}

	return repoProjects
}

func convertToProject(e repo.Project) Project {
	return Project{
		ID:     e.ID,
		UserID: e.UserID,
		Name:   e.Name,
	}
}
