package time_tracker

import (
	"github.com/BMSTU-TIMETRACKERS/timetracker-backend/config/time_tracker/flags"
	"github.com/labstack/echo/v4"
)

type Base struct {
	Logger   flags.LoggerFlags `toml:"logger"`
	services *baseServices
}

type baseServices struct {
	Logger echo.Logger
	// Tracer          *otel.Tracer
	// MetricsRegistry *metrics.Registry
}

func (b *Base) Init(e *echo.Echo) (*baseServices, error) {
	services := &baseServices{}
	logger := b.Logger.Init(e)
	services.Logger = logger
	b.services = services

	return services, nil
}
