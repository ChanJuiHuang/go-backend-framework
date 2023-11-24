package logger

import (
	"time"

	"go.uber.org/zap/zapcore"
)

type Level string

const (
	Debug  Level = "debug"
	Info         = "info"
	Warn         = "warn"
	Error        = "error"
	Dpanic       = "dpanic"
	Panic        = "panic"
	Fatal        = "fatal"
)

type Type string

const (
	Console Type = "console"
	File         = "file"
)

type Config struct {
	Type       Type
	LogPath    string
	MaxSize    int64
	MaxBackups int
	MaxAge     time.Duration
	Compress   bool
	Level      Level
}

func GetLevel(level Level) zapcore.Level {
	switch level {
	case Debug:
		return zapcore.DebugLevel
	case Info:
		return zapcore.InfoLevel
	case Warn:
		return zapcore.WarnLevel
	case Error:
		return zapcore.ErrorLevel
	case Dpanic:
		return zapcore.DPanicLevel
	case Panic:
		return zapcore.PanicLevel
	case Fatal:
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}
