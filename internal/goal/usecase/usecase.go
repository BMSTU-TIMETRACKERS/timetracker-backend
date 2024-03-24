package delivery

import (
	"context"
	"errors"
	"fmt"
	"time"

	repo "github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/goal/repository"
)

var ErrGoalNotFound = errors.New("goal not found")

type repository interface {
	CreateGoal(ctx context.Context, goal repo.Goal) (int64, error)
	GetGoals(ctx context.Context, userID, projectID int64) ([]repo.Goal, error)
}

type Usecase struct {
	repository repository
}

func NewUsecase(repository repository) *Usecase {
	return &Usecase{
		repository: repository,
	}
}

func (u *Usecase) CreateGoal(ctx context.Context, goal Goal) (int64, error) {
	id, err := u.repository.CreateGoal(ctx, convertToRepoGoal(goal))

	if err != nil {
		return 0, fmt.Errorf("repo create goal: %v", err)
	}

	return id, err
}

func (u *Usecase) GetGoals(ctx context.Context, userID, projectID int64) ([]Goal, error) {
	goals, err := u.repository.GetGoals(ctx, userID, projectID)
	if err != nil {
		return nil, fmt.Errorf("repo get goals: %v", err)
	}

	var res []Goal

	for _, goal := range goals {
		var duration time.Duration

		for _, entry := range goal.Entries {
			var timeStart, timeEnd time.Time
			if entry.TimeStart.Before(goal.DateStart) {
				timeStart = goal.DateStart
			} else {
				timeStart = entry.TimeStart
			}

			if goal.DateEnd.AddDate(0, 0, 1).Before(entry.TimeEnd) {
				timeEnd = goal.DateEnd
			} else {
				timeEnd = entry.TimeEnd
			}

			duration += timeEnd.Sub(timeStart)
		}

		durationSeconds := duration.Seconds()
		percent := durationSeconds / float64(goal.TimeSeconds) * 100
		res = append(res, Goal{
			ID:              goal.ID,
			ProjectID:       goal.ProjectID,
			UserID:          goal.UserID,
			TimeSeconds:     goal.TimeSeconds,
			Name:            goal.Name,
			DateStart:       goal.DateStart,
			DateEnd:         goal.DateEnd,
			DurationSeconds: duration.Seconds(),
			Percent:         percent,
		})
	}

	return res, nil
}

func convertToRepoGoal(goal Goal) repo.Goal {
	return repo.Goal{
		ID:          goal.ID,
		ProjectID:   goal.ProjectID,
		UserID:      goal.UserID,
		TimeSeconds: goal.TimeSeconds,
		Name:        goal.Name,
		DateStart:   goal.DateStart,
		DateEnd:     goal.DateEnd,
	}
}
