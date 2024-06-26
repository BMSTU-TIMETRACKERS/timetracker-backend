package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	entryRepoDto "github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/entry/repository"
	repo "github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/project/repository"
)

var (
	ErrProjectNotFound = errors.New("project not found")
	ErrProjectExists   = errors.New("project with that name already exists")
)

type repository interface {
	CreateProject(ctx context.Context, project repo.Project) (int64, error)
	GetUserProjects(ctx context.Context, userID int64) ([]repo.Project, error)
	ClearUserData(ctx context.Context, userID int64) error
	GetProjectByName(ctx context.Context, userID int64, projectName string) (repo.Project, error)
}

type entryRepository interface {
	GetProjectEntriesForInterval(
		_ context.Context,
		userID int64,
		projectID int64,
		start time.Time,
		end time.Time) ([]entryRepoDto.Entry, error)
	GetProjectEntries(
		_ context.Context,
		userID int64,
		projectID int64,
	) ([]entryRepoDto.Entry, error)
}

type Usecase struct {
	repository      repository
	entryRepository entryRepository
}

func NewUsecase(repository repository, entryRepository entryRepository) *Usecase {
	return &Usecase{
		repository:      repository,
		entryRepository: entryRepository,
	}
}

func (u *Usecase) CreateProject(ctx context.Context, project Project) (int64, error) {
	oldProject, err := u.repository.GetProjectByName(ctx, project.UserID, project.Name)
	if err != nil && !errors.Is(err, repo.ErrProjectNotFound) {
		return 0, fmt.Errorf("repo get project by name: %v", err)
	}
	if oldProject.ID != 0 {
		return 0, ErrProjectExists
	}

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

func (u *Usecase) ProjectStat(ctx context.Context, projectID int64, userID int64, timeStart, timeEnd time.Time) (AllProjectEntriesStat, error) {
	projectEntries, err := u.entryRepository.GetProjectEntriesForInterval(ctx, userID, projectID, timeStart, timeEnd)
	if err != nil {
		return AllProjectEntriesStat{}, fmt.Errorf("get project entries error: %w", err)
	}

	// Собираем уникальные имена записей, чтобы сложить по ним стату.
	entriesStatMap := make(map[string]ProjectEntrieInfo)
	totalDurationSec := float64(0)
	for _, e := range projectEntries {
		durationSec := e.TimeEnd.Sub(e.TimeStart).Seconds()

		// Структура при получении из мапы копируется, поэтому надо ее переприсваивать.
		tmpInfo := entriesStatMap[e.Name]
		tmpInfo.EntryName = e.Name
		tmpInfo.EntryDurationInSec += durationSec

		entriesStatMap[e.Name] = tmpInfo

		totalDurationSec += durationSec
	}

	// Собираем стату в массив.
	entriesStat := make([]ProjectEntrieInfo, 0, len(entriesStatMap))
	for _, v := range entriesStatMap {
		entriesStat = append(entriesStat, ProjectEntrieInfo{
			EntryName:            v.EntryName,
			EntryDurationInSec:   v.EntryDurationInSec,
			EntryDurationPercent: calculatePercentDuration(v.EntryDurationInSec, totalDurationSec),
		})
	}

	return AllProjectEntriesStat{
		TotalDurationInSec: totalDurationSec,
		EntriesStat:        entriesStat,
	}, nil
}

func (u *Usecase) ProjectsStats(ctx context.Context, userID int64, timeStart, timeEnd time.Time) (AllProjectsStat, error) {
	repoProjects, err := u.repository.GetUserProjects(ctx, userID)
	if err != nil {
		if errors.Is(err, repo.ErrProjectNotFound) {
			return AllProjectsStat{}, nil
		}
		return AllProjectsStat{}, fmt.Errorf("repo get user projects: %v", err)
	}

	generalStat := AllProjectsStat{
		TotalDurationInSec: 0,
		ProjectsStat:       nil,
	}

	projectStats := make([]ProjectStatInfo, 0, len(repoProjects))
	for _, project := range repoProjects {
		stat, err := u.getProjectStat(ctx, userID, project, timeStart, timeEnd)
		if err != nil {
			if errors.Is(err, entryRepoDto.ErrEntryNotFound) {
				continue
			}
			return AllProjectsStat{}, fmt.Errorf("get project stat: %v", err)
		}
		projectStats = append(projectStats, stat)
		generalStat.TotalDurationInSec += stat.ProjectDurationInSec
	}

	generalStat.ProjectsStat = projectStats

	for idx := range generalStat.ProjectsStat {
		generalStat.ProjectsStat[idx].ProjectDurationPercent = calculatePercentDuration(
			generalStat.ProjectsStat[idx].ProjectDurationInSec,
			generalStat.TotalDurationInSec,
		)
	}

	return generalStat, nil
}

func (u *Usecase) ClearUserData(ctx context.Context, userID int64) error {
	if err := u.repository.ClearUserData(ctx, userID); err != nil {
		return fmt.Errorf("repo clear user data: %v", err)
	}

	return nil
}

func (u *Usecase) getProjectStat(ctx context.Context, userID int64, project repo.Project, timeStart, timeEnd time.Time) (ProjectStatInfo, error) {
	projectEntries, err := u.entryRepository.GetProjectEntriesForInterval(ctx, userID, project.ID, timeStart, timeEnd)
	if err != nil {
		return ProjectStatInfo{}, fmt.Errorf("get project entries error: %w", err)
	}

	projectDuration := calculateProjectDuration(projectEntries)

	return ProjectStatInfo{
		ProjectID:              project.ID,
		ProjectName:            project.Name,
		ProjectDurationInSec:   projectDuration.Seconds(),
		ProjectDurationPercent: 0,
	}, err
}

func calculateProjectDuration(entries []entryRepoDto.Entry) time.Duration {
	totalDuration := time.Duration(0)
	for _, e := range entries {
		totalDuration += e.TimeEnd.Sub(e.TimeStart)
	}

	return totalDuration
}

func calculatePercentDuration(duration float64, totalDuration float64) float64 {
	return float64(duration) / float64(totalDuration) * 100
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
