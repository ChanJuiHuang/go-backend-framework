package database

import (
	"fmt"
	"time"

	"gorm.io/gorm/logger"
)

type Driver string

const (
	MySql  Driver = "mysql"
	PgSql         = "pgsql"
	Sqlite        = "sqlite"
)

func GetDriver(driver Driver) Driver {
	switch driver {
	case MySql:
		return MySql
	case PgSql:
		return PgSql
	case Sqlite:
		return Sqlite
	default:
		panic(fmt.Sprintf("[%s] is invalid driver", driver))
	}
}

type LogLevel string

const (
	Info   LogLevel = "info"
	Warn            = "warn"
	Error           = "error"
	Silent          = "silent"
)

func GetGormLogLevel(level LogLevel) logger.LogLevel {
	switch level {
	case Info:
		return logger.Info
	case Warn:
		return logger.Warn
	case Error:
		return logger.Error
	case Silent:
		return logger.Silent
	default:
		return logger.Info
	}
}

type Config struct {
	Driver          Driver
	Username        string
	Password        string
	Host            string
	Port            string
	Database        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	LogLevel        LogLevel
}
