package logging

import (
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

const (
	LogFormatLogfmt = "logfmt"
	LogFormatJSON   = "json"
)

type LevelLogger struct {
	log.Logger
	LogLevel string
}

func NewLogger(logLevel, logFormat, debugName string) log.Logger {
	var (
		logger log.Logger
		lvl    level.Option
	)

	switch logLevel {
	case "error":
		lvl = level.AllowError()
	case "warn":
		lvl = level.AllowWarn()
	case "info":
		lvl = level.AllowInfo()
	case "debug":
		lvl = level.AllowDebug()
	default:
		panic("invalid log level")
	}

	switch logFormat {
	case LogFormatLogfmt:
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	case LogFormatJSON:
		logger = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	}

	logger = level.NewFilter(logger, lvl)

	if debugName != "" {
		logger = log.With(logger, "name", debugName)
	}

	return LevelLogger{
		Logger:   logger,
		LogLevel: logLevel,
	}

}
