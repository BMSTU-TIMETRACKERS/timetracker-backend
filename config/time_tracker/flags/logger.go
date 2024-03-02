package flags

import (
	"github.com/labstack/echo/v4"
	elog "github.com/labstack/gommon/log"
)

type LoggerFlags struct {
	LogLevel      uint8  `toml:"level"`
	LogHeader     string `toml:"header"`
	LogHttpFormat string `toml:"log-http-format"`
}

func (f LoggerFlags) Init(e *echo.Echo) echo.Logger {
	e.Logger.SetLevel(elog.Lvl(f.LogLevel))
	e.Logger.SetHeader(f.LogHeader)
	return e.Logger
}
