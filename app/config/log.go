package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type LogType string

const (
	Console LogType = "console"
	File            = "file"
)

type LogLevel string

const (
	Debug  LogLevel = "debug"
	Info            = "info"
	Warn            = "warn"
	Error           = "error"
	Dpanic          = "dpanic"
	Panic           = "panic"
	Fatal           = "fatal"
	Silent          = "silent"
)

type logConfig struct {
	AccessLogPath string        `required:"true" split_words:"true"`
	AppLogPath    string        `required:"true" split_words:"true"`
	MaxSize       int64         `required:"true" split_words:"true"`
	MaxBackups    int           `required:"true" split_words:"true"`
	MaxAge        time.Duration `required:"true" split_words:"true"`
	Compress      bool          `required:"true"`
	Level         LogLevel      `required:"true"`
	Type          LogType       `required:"true"`
}

var logCfg *logConfig

func Log() *logConfig {
	if logCfg == nil {
		logCfg = new(logConfig)
		err := envconfig.Process("log", logCfg)

		if err != nil {
			panic(err)
		}
		logCfg.MaxSize = logCfg.MaxSize << 20
	}
	switch logCfg.Type {
	case Console:
	case File:
	default:
		panic(fmt.Sprintf("Log type [%s] does not exist", logCfg.Type))
	}

	switch logCfg.Level {
	case Debug:
	case Info:
	case Warn:
	case Error:
	case Dpanic:
	case Panic:
	case Fatal:
	case Silent:
	default:
		panic(fmt.Sprintf("Log level [%s] does not exist", logCfg.Level))
	}

	return logCfg
}
