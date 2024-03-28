package delivery

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	usecaseDto "github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/project/usecase"
	"github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/response"
	"github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/validator"
)

type usecase interface {
	CreateProject(ctx context.Context, project usecaseDto.Project) (int64, error)
	GetUserProjects(ctx context.Context, userID int64) ([]usecaseDto.Project, error)
	ProjectsStats(ctx context.Context, userID int64, timeStart, timeEnd time.Time) (usecaseDto.AllProjectsStat, error)
	ProjectStat(ctx context.Context, projectID int64, userID int64, timeStart, timeEnd time.Time) (usecaseDto.AllProjectEntriesStat, error)
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

	e.POST("/projects/create", handler.CreateProject)
	e.GET("/me/projects", handler.GetMyProjects)
	e.GET("/me/projects/stat", handler.GetProjectsStat)
	e.GET("/me/projects/:id/stat", handler.GetProjectStat)
}

// CreateProject godoc
// @Summary      Create project.
// @Description  Create project.
// @Tags     	 projects
// @Accept	 application/json
// @Produce  application/json
// @Param    project body CreateProjectIn true "project info"
// @Success  200 {object} CreateProjectOut "success create project"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 422 {object} echo.HTTPError "unprocessable entity"
// @Router   /projects/create [post]
func (d *Delivery) CreateProject(c echo.Context) error {
	ctx := context.Background()

	var in CreateProjectIn
	err := c.Bind(&in)

	if err != nil {
		c.Logger().Errorf("bind request: %v", err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, response.ErrorMsgsByCode[http.StatusUnprocessableEntity])
	}

	if ok, err := validator.IsRequestValid(&in); !ok {
		c.Logger().Errorf("validation: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, response.ErrorMsgsByCode[http.StatusBadRequest])
	}

	userID, ok := c.Get("user_id").(int64)
	if !ok {
		c.Logger().Error("can't parse context user_id")
		return echo.NewHTTPError(http.StatusInternalServerError, response.ErrorMsgsByCode[http.StatusInternalServerError])
	}

	project := usecaseDto.Project{
		Name: in.Name,
	}

	project.UserID = userID
	projectID, err := d.usecase.CreateProject(ctx, project)
	if err != nil {
		c.Logger().Errorf("usecase: %v", err)
		return handleUsecaseError(err)
	}

	out := CreateProjectOut{ID: projectID}

	return c.JSON(http.StatusOK, out)
}

// GetMyProjects godoc
// @Summary      Get my projects.
// @Description  Get my projects or get my projects for a day
// @Tags     	 projects
// @Accept	 	application/json
// @Produce  	application/json
// @Success  200 {object} []ProjectOut "success get projects"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Router   /me/projects [get]
func (d *Delivery) GetMyProjects(c echo.Context) error {
	ctx := context.Background()

	userID, ok := c.Get("user_id").(int64)
	if !ok {
		c.Logger().Error("can't parse context user_id")
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			response.ErrorMsgsByCode[http.StatusInternalServerError],
		)
	}

	var projects []usecaseDto.Project
	var err error

	projects, err = d.usecase.GetUserProjects(ctx, userID)

	if err != nil {
		c.Logger().Errorf("usecase: %v", err)
		return handleUsecaseError(err)
	}

	out := convertFromUsecaseProjects(projects)

	return c.JSON(http.StatusOK, out)
}

// GetProjectsStat godoc
// @Summary      Get project stats.
// @Description  Get project stats
// @Tags     	 projects
// @Accept	 	application/json
// @Produce  	application/json
// @Param        time_start    query     string  false  "RFC3339 format"
// @Param        time_end    query     string  false  "RFC3339 format"
// @Success  200 {object} ProjectsStatOut "success"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Router   /me/projects/stat [get]
func (d *Delivery) GetProjectsStat(c echo.Context) error {
	ctx := context.Background()

	userID, ok := c.Get("user_id").(int64)
	if !ok {
		c.Logger().Error("can't parse context user_id")
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			response.ErrorMsgsByCode[http.StatusInternalServerError],
		)
	}

	timeStartStr := c.QueryParam("time_start")
	timeEndStr := c.QueryParam("time_end")

	timeStart := time.Time{}
	timeEnd := time.Now()

	if timeStartStr != "" {
		// Намеренный скип ошибки.
		timeStart, _ = time.Parse(time.RFC3339, timeStartStr)
	}

	if timeEndStr != "" {
		// Намеренный скип ошибки.
		timeEnd, _ = time.Parse(time.RFC3339, timeEndStr)
	}

	projectsStat, err := d.usecase.ProjectsStats(ctx, userID, timeStart, timeEnd)

	if err != nil {
		c.Logger().Errorf("usecase: %v", err)
		return handleUsecaseError(err)
	}

	out := convertFromUsecaseProjectsStat(projectsStat)

	return c.JSON(http.StatusOK, out)
}

// GetProjectStat godoc
// @Summary      Get project entries stat.
// @Description  Get project entries stat
// @Tags     	 projects
// @Accept	 	application/json
// @Produce  	application/json
// @Param id  path int  true  "project ID"
// @Param        time_start    query     string  false  "RFC3339 format"
// @Param        time_end    query     string  false  "RFC3339 format"
// @Success  200 {object}  ProjectEntriesStatOut "success"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Router   /me/projects/{id}/stat [get]
func (d *Delivery) GetProjectStat(c echo.Context) error {
	ctx := context.Background()

	userID, ok := c.Get("user_id").(int64)
	if !ok {
		c.Logger().Error("can't parse context user_id")
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			response.ErrorMsgsByCode[http.StatusInternalServerError],
		)
	}

	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.Logger().Errorf("parse int: %v", err)
		return echo.NewHTTPError(
			http.StatusBadRequest,
			response.ErrorMsgsByCode[http.StatusBadRequest],
		)
	}

	timeStartStr := c.QueryParam("time_start")
	timeEndStr := c.QueryParam("time_end")

	timeStart := time.Time{}
	timeEnd := time.Now()

	if timeStartStr != "" {
		// Намеренный скип ошибки.
		timeStart, _ = time.Parse(time.RFC3339, timeStartStr)
	}

	if timeEndStr != "" {
		// Намеренный скип ошибки.
		timeEnd, _ = time.Parse(time.RFC3339, timeEndStr)
	}

	projectEntriesStat, err := d.usecase.ProjectStat(ctx, projectID, userID, timeStart, timeEnd)

	if err != nil {
		c.Logger().Errorf("usecase: %v", err)
		return handleUsecaseError(err)
	}

	out := convertFromUsecaseProjectEntriesStat(projectEntriesStat)

	return c.JSON(http.StatusOK, out)
}

func handleUsecaseError(err error) *echo.HTTPError {
	// Не нашли проект.
	if errors.Is(err, usecaseDto.ErrProjectNotFound) {
		return echo.NewHTTPError(
			http.StatusNotFound,
			fmt.Sprintf("%s: %s", response.ErrorMsgsByCode[http.StatusNotFound], "project"))
	}
	if errors.Is(err, usecaseDto.ErrProjectExists) {
		return echo.NewHTTPError(http.StatusBadRequest, response.ErrorMsgsByCode[http.StatusBadRequest])
	}

	// По дефолту пятисотим.
	return echo.NewHTTPError(
		http.StatusInternalServerError,
		response.ErrorMsgsByCode[http.StatusInternalServerError],
	)
}

func convertFromUsecaseProjects(projects []usecaseDto.Project) []ProjectOut {
	out := make([]ProjectOut, 0, len(projects))
	for _, project := range projects {
		out = append(out, convertFromUsecaseProject(project))
	}

	return out
}

func convertFromUsecaseProjectsStat(stat usecaseDto.AllProjectsStat) ProjectsStatOut {
	projectsOut := make([]ProjectStat, 0, len(stat.ProjectsStat))
	for _, s := range stat.ProjectsStat {
		projectsOut = append(projectsOut, ProjectStat{
			ID:              s.ProjectID,
			Name:            s.ProjectName,
			DurationInSec:   s.ProjectDurationInSec,
			PercentDuration: s.ProjectDurationPercent,
		})
	}

	return ProjectsStatOut{
		TotalDurationInSec: stat.TotalDurationInSec,
		Projects:           projectsOut,
	}
}

func convertFromUsecaseProjectEntriesStat(stat usecaseDto.AllProjectEntriesStat) ProjectEntriesStatOut {
	entriesOut := make([]ProjectEntriesStat, 0, len(stat.EntriesStat))
	for _, s := range stat.EntriesStat {
		entriesOut = append(entriesOut, ProjectEntriesStat{
			Name:            s.EntryName,
			DurationInSec:   s.EntryDurationInSec,
			PercentDuration: s.EntryDurationPercent,
		})
	}

	return ProjectEntriesStatOut{
		TotalDurationInSec: stat.TotalDurationInSec,
		Entries:            entriesOut,
	}
}

func convertFromUsecaseProject(project usecaseDto.Project) ProjectOut {
	return ProjectOut{
		ID:   project.ID,
		Name: project.Name,
	}
}
