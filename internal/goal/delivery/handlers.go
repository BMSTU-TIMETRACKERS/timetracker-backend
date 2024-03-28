package delivery

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	usecaseDto "github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/goal/usecase"
	"github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/response"
	"github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/validator"
)

type usecase interface {
	CreateGoal(ctx context.Context, e usecaseDto.Goal) (int64, error)
	GetGoals(ctx context.Context, userID, projectID int64) ([]usecaseDto.Goal, error)
}

type Delivery struct {
	usecase usecase

	logger echo.Logger
}

func RegisterHandlers(
	e *echo.Echo,
	usecase usecase,
	logger echo.Logger,
) {
	handler := &Delivery{
		usecase: usecase,

		logger: logger,
	}

	e.POST("/goals/create", handler.CreateGoal)
	e.GET("/me/projects/:project_id/goals", handler.GetMyGoals)
}

// CreateGoal godoc
// @Summary      Create goal.
// @Description  Create goal.
// @Tags     	 goals
// @Accept	 application/json
// @Produce  application/json
// @Param    goal body CreateGoalIn true "goal info"
// @Success  200 {object} CreateGoalOut "success create goal"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 404 {object} echo.HTTPError "item is not found"
// @Failure 422 {object} echo.HTTPError "unprocessable entity"
// @Router   /goals/create [post]
func (d *Delivery) CreateGoal(c echo.Context) error {
	ctx := context.Background()

	var in CreateGoalIn
	err := c.Bind(&in)

	if err != nil {
		c.Logger().Errorf("bind request: %v", err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, response.ErrorMsgsByCode[http.StatusUnprocessableEntity])
	}

	if ok, err := validator.IsRequestValid(&in); !ok {
		c.Logger().Errorf("validation: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, response.ErrorMsgsByCode[http.StatusBadRequest])
	}

	// Получаем userID (проставляется в auth мидлваре).
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		c.Logger().Error("can't parse context user_id")
		return echo.NewHTTPError(http.StatusInternalServerError, response.ErrorMsgsByCode[http.StatusInternalServerError])
	}

	goal := usecaseDto.Goal{
		ProjectID:   in.ProjectID,
		TimeSeconds: in.TimeSeconds,
		Name:        in.Name,
		DateStart:   in.DateStart,
		DateEnd:     in.DateEnd,
	}

	goal.UserID = userID
	goalID, err := d.usecase.CreateGoal(ctx, goal)
	if err != nil {
		c.Logger().Errorf("usecase: %v", err)
		return handleUsecaseError(err)
	}

	out := CreateGoalOut{ID: goalID}

	return c.JSON(http.StatusOK, out)
}

// GetMyGoals godoc
// @Summary      Get my goals for project.
// @Description  Get my goals for project
// @Tags     	 goals
// @Accept	 	application/json
// @Produce  	application/json
// @Param    project_id path int true "project id"
// @Success  200 {object} []GoalOut "success get goals"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Router   /me/projects/{project_id}/goals [get]
func (d *Delivery) GetMyGoals(c echo.Context) error {
	ctx := context.Background()

	projectIDStr := c.Param("project_id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		c.Logger().Errorf("parse int: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, response.ErrorMsgsByCode[http.StatusBadRequest])
	}

	userID, ok := c.Get("user_id").(int64)
	if !ok {
		c.Logger().Error("can't parse context user_id")
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			response.ErrorMsgsByCode[http.StatusInternalServerError],
		)
	}

	var goals []usecaseDto.Goal

	goals, err = d.usecase.GetGoals(ctx, userID, projectID)

	if err != nil {
		c.Logger().Errorf("usecase: %v", err)
		return handleUsecaseError(err)
	}

	out := convertFromUsecaseEntries(goals)

	return c.JSON(http.StatusOK, out)
}

func handleUsecaseError(err error) *echo.HTTPError {
	// Не нашли цель.
	if errors.Is(err, usecaseDto.ErrGoalNotFound) {
		return echo.NewHTTPError(
			http.StatusNotFound,
			fmt.Sprintf("%s: %s", response.ErrorMsgsByCode[http.StatusNotFound], "goal"))
	}

	// По дефолту пятисотим.
	return echo.NewHTTPError(
		http.StatusInternalServerError,
		response.ErrorMsgsByCode[http.StatusInternalServerError],
	)
}

func convertFromUsecaseEntries(goals []usecaseDto.Goal) []GoalOut {
	out := make([]GoalOut, 0, len(goals))
	for _, goal := range goals {
		g := GoalOut{
			ID:              goal.ID,
			ProjectID:       goal.ProjectID,
			UserID:          goal.UserID,
			TimeSeconds:     goal.TimeSeconds,
			Name:            goal.Name,
			DateStart:       goal.DateStart,
			DateEnd:         goal.DateEnd,
			DurationSeconds: goal.DurationSeconds,
			Percent:         goal.Percent,
		}

		out = append(out, g)
	}

	return out
}
