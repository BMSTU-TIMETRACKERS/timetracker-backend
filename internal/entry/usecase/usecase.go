package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	repo "github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/entry/repository"
	"github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/utils"
)

var ErrEntryNotFound = errors.New("entry not found")

type repository interface {
	CreateEntry(ctx context.Context, entry repo.Entry) (int64, error)
	GetUserEntries(ctx context.Context, userID int64) ([]repo.Entry, error)
	GetUserEntriesForInterval(ctx context.Context, userID int64, start time.Time, end time.Time) ([]repo.Entry, error)

	GetProjectsInfo(ctx context.Context, projectIDs []int64) ([]repo.ProjectInfo, error)
}

type Usecase struct {
	repository repository
}

func NewUsecase(repository repository) *Usecase {
	return &Usecase{
		repository: repository,
	}
}

func (u *Usecase) CreateEntry(ctx context.Context, entry Entry) (int64, error) {
	id, err := u.repository.CreateEntry(ctx, convertToRepoEntry(entry))

	if err != nil {
		return 0, fmt.Errorf("repo create entry: %v", err)
	}

	return id, err
}

func (u *Usecase) GetUserEntries(ctx context.Context, userID int64) ([]Entry, error) {
	repoEntries, err := u.repository.GetUserEntries(ctx, userID)
	if err != nil {
		if errors.Is(err, repo.ErrEntryNotFound) {
			return nil, fmt.Errorf("repo get user entries: %w", ErrEntryNotFound)
		}
		return nil, fmt.Errorf("repo get user entries: %v", err)
	}

	entries := convertToEntries(repoEntries)

	err = u.enrichEntries(ctx, entries)
	if err != nil {
		return nil, fmt.Errorf("enrich parcels: %v", err)
	}

	return entries, nil
}

func (u *Usecase) GetUserEntriesForDay(ctx context.Context, userID int64, date time.Time) ([]Entry, error) {
	startDay, endDay := utils.GetDayInterval(date)
	repoEntries, err := u.repository.GetUserEntriesForInterval(ctx, userID, startDay, endDay)
	if err != nil {
		if errors.Is(err, repo.ErrEntryNotFound) {
			return nil, fmt.Errorf("repo get user entries for day: %w", ErrEntryNotFound)
		}
		return nil, fmt.Errorf("repo get user entries for day: %v", err)
	}

	entries := convertToEntries(repoEntries)

	err = u.enrichEntries(ctx, entries)
	if err != nil {
		return nil, fmt.Errorf("enrich parcels: %v", err)
	}

	return entries, nil
}

func (u *Usecase) enrichEntries(ctx context.Context, entries []Entry) error {
	var projectIDs []int64
	for _, e := range entries {
		projectIDs = append(projectIDs, e.ProjectID)
	}

	projectInfo, err := u.repository.GetProjectsInfo(ctx, projectIDs)
	if err != nil {
		return fmt.Errorf("")
	}

	projectIdName := make(map[int64]string)
	for _, info := range projectInfo {
		projectIdName[info.ID] = info.Name
	}

	for id, _ := range entries {
		entries[id].ProjectName = projectIdName[entries[id].ProjectID]
	}

	fmt.Println(entries)

	return nil
}

func convertToRepoEntries(entries []Entry) []repo.Entry {
	repoEntries := make([]repo.Entry, 0, len(entries))
	for _, e := range entries {
		repoEntries = append(repoEntries, convertToRepoEntry(e))
	}

	return repoEntries
}

func convertToRepoEntry(entry Entry) repo.Entry {
	return repo.Entry{
		ID:        entry.ID,
		UserID:    entry.UserID,
		ProjectID: entry.ProjectID,
		Name:      entry.Name,
		TimeStart: entry.TimeStart,
		TimeEnd:   entry.TimeEnd,
	}
}

func convertToEntry(e repo.Entry) Entry {
	return Entry{
		ID:          e.ID,
		UserID:      e.UserID,
		ProjectID:   e.ProjectID,
		Name:        e.Name,
		TimeStart:   e.TimeStart,
		TimeEnd:     e.TimeEnd,
		ProjectName: "",
	}
}

func convertToEntries(entries []repo.Entry) []Entry {
	repoEntries := make([]Entry, 0, len(entries))
	for _, e := range entries {
		repoEntries = append(repoEntries, convertToEntry(e))
	}

	return repoEntries
}
