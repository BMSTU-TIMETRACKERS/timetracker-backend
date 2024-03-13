package delivery

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	usecaseDto "github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/entry/usecase"
	"github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/response"
	"github.com/BMSTU-TIMETRACKERS/timetracker-backend/internal/validator"
)

type usecase interface {
	CreateEntry(ctx context.Context, e usecaseDto.Entry) (int64, error)
	GetUserEntries(ctx context.Context, userID int64) ([]usecaseDto.Entry, error)
	GetUserEntriesForDay(ctx context.Context, userID int64, date time.Time) ([]usecaseDto.Entry, error)
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

	e.POST("/entries/create", handler.CreateEntry)
	e.GET("/me/entries", handler.GetMyEntries)
	// e.POST("/entries/edit", handler.UpdateEntry)
	// e.GET("/entries/:id", handler.GetEntry)
	// e.DELETE("/entries/:id", handler.DeleteEntry)
}

// CreateEntry godoc
// @Summary      Create entry.
// @Description  Create entry.
// @Tags     	 entries
// @Accept	 application/json
// @Produce  application/json
// @Param    entry body CreateEntryIn true "entry info"
// @Success  200 {object} CreateEntryOut "success create entry"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 404 {object} echo.HTTPError "item is not found"
// @Failure 422 {object} echo.HTTPError "unprocessable entity"
// @Router   /entries/create [post]
func (d *Delivery) CreateEntry(c echo.Context) error {
	ctx := context.Background()

	var in CreateEntryIn
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

	entry := usecaseDto.Entry{
		ProjectID: in.ProjectID,
		Name:      in.Name,
		TimeStart: in.TimeStart,
		TimeEnd:   in.TimeEnd,
	}

	entry.UserID = userID
	entryID, err := d.usecase.CreateEntry(ctx, entry)
	if err != nil {
		c.Logger().Errorf("usecase: %v", err)
		return handleUsecaseError(err)
	}

	out := CreateEntryOut{ID: entryID}

	return c.JSON(http.StatusOK, out)
}

// GetMyEntries godoc
// @Summary      Get my entries.
// @Description  Get my entries or get my entries for a day
// @Tags     	 entries
// @Accept	 	application/json
// @Produce  	application/json
// @Param        day    query     string  false  "day for events in YYYY-MM-DD format"
// @Success  200 {object} []EntryOut "success create entry"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Router   /me/entries [get]
func (d *Delivery) GetMyEntries(c echo.Context) error {
	ctx := context.Background()

	userID, ok := c.Get("user_id").(int64)
	if !ok {
		c.Logger().Error("can't parse context user_id")
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			response.ErrorMsgsByCode[http.StatusInternalServerError],
		)
	}

	day := c.QueryParam("day")

	var entries []usecaseDto.Entry
	var err error

	if day == "" {
		entries, err = d.usecase.GetUserEntries(ctx, userID)
	} else {
		date, err := time.Parse("2006-01-02", day)

		if err != nil {
			c.Logger().Errorf("invalid data format, should be YYYY-MM-DD: %v", err)
			return echo.NewHTTPError(
				http.StatusBadRequest,
				response.ErrorMsgsByCode[http.StatusBadRequest],
			)
		}

		entries, err = d.usecase.GetUserEntriesForDay(ctx, userID, date)
	}

	if err != nil {
		c.Logger().Errorf("usecase: %v", err)
		return handleUsecaseError(err)
	}

	out := convertFromUsecaseEntries(entries)

	return c.JSON(http.StatusOK, out)
}

func handleUsecaseError(err error) *echo.HTTPError {
	// Не нашли запись времени.
	if errors.Is(err, usecaseDto.ErrEntryNotFound) {
		return echo.NewHTTPError(
			http.StatusNotFound,
			fmt.Sprintf("%s: %s", response.ErrorMsgsByCode[http.StatusNotFound], "entry"))
	}

	// По дефолту пятисотим.
	return echo.NewHTTPError(
		http.StatusInternalServerError,
		response.ErrorMsgsByCode[http.StatusInternalServerError],
	)
}

func convertFromUsecaseEntries(entries []usecaseDto.Entry) []EntryOut {
	out := make([]EntryOut, 0, len(entries))
	for _, entry := range entries {
		e := EntryOut{
			ID:          entry.ID,
			ProjectID:   entry.ProjectID,
			ProjectName: entry.ProjectName,
			Name:        entry.Name,
			TimeStart:   entry.TimeStart,
			TimeEnd:     entry.TimeEnd,
		}

		out = append(out, e)
	}

	return out
}
