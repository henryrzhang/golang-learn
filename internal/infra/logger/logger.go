package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init 根据配置初始化全局 logger
func Init(level, format string) {
	var lvl zerolog.Level
	switch level {
	case "debug":
		lvl = zerolog.DebugLevel
	case "info":
		lvl = zerolog.InfoLevel
	case "warn":
		lvl = zerolog.WarnLevel
	case "error":
		lvl = zerolog.ErrorLevel
	default:
		lvl = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(lvl)

	var w io.Writer = os.Stdout
	if format == "console" {
		w = zerolog.ConsoleWriter{Out: os.Stdout}
	}

	log.Logger = zerolog.New(w).With().Timestamp().Caller().Logger()
}

// L 返回带上下文的 logger
func L() *zerolog.Logger {
	return &log.Logger
}
