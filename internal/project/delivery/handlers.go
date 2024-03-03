package delivery

import (
	"net/http"

	"github.com/labstack/echo/v4"

	usecaseDto "github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/project/usecase"
	"github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/response"
	"github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/validator"
)

type usecase interface {
	CreateProject(e usecaseDto.Project) (int64, error)
	GetUserProjects(userID int64) ([]usecaseDto.Project, error)
}

type Delivery struct {
	usecase usecase

	logger echo.Logger
}

func NewDelivery(
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
	// e.POST("/projects/edit", handler.UpdateEntry)
	// e.GET("/projects/:id", handler.GetEntry)
	// e.DELETE("/projects/:id", handler.DeleteEntry)
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
	projectID, err := d.usecase.CreateProject(project)
	if err != nil {
		c.Logger().Errorf("usecase: %v", err)
		return handleUsecaseError(err)
	}

	out := CreateProjectOut{ID: projectID}

	return c.JSON(http.StatusOK, response.Response{Body: out})
}

// GetMyProjects godoc
// @Summary      Get my projects.
// @Description  Get my projects or get my projects for a day
// @Tags     	 projects
// @Accept	 	application/json
// @Produce  	application/json
// @Success  200 {object} []ProjectOut "success create project"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Router   /me/projects [get]
func (d *Delivery) GetMyProjects(c echo.Context) error {
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

	projects, err = d.usecase.GetUserProjects(userID)

	if err != nil {
		c.Logger().Errorf("usecase: %v", err)
		return handleUsecaseError(err)
	}

	out := convertFromUsecaseProjects(projects)

	return c.JSON(http.StatusOK, response.Response{Body: out})
}

func handleUsecaseError(err error) *echo.HTTPError {
	// TODO во время юзкейсов
	// causeErr := errors.Cause(err)
	// switch {
	// case errors.Is(causeErr, models.ErrNotFound):
	// 	return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
	// case errors.Is(causeErr, models.ErrBadRequest):
	// 	return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	// case errors.Is(causeErr, models.ErrPermissionDenied):
	// 	return echo.NewHTTPError(http.StatusForbidden, models.ErrPermissionDenied.Error())
	// default:
	// 	return echo.NewHTTPError(http.StatusInternalServerError, causeErr.Error())
	// }
	return nil
}

func convertFromUsecaseProjects(projects []usecaseDto.Project) []ProjectOut {
	out := make([]ProjectOut, 0, len(projects))
	for _, project := range projects {
		out = append(out, convertFromUsecaseProject(project))
	}

	return out
}

func convertFromUsecaseProject(project usecaseDto.Project) ProjectOut {
	return ProjectOut{
		ID:   project.ID,
		Name: project.Name,
	}
}
