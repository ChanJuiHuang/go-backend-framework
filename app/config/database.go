package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type DatabaseDriver string

const (
	MySql  DatabaseDriver = "mysql"
	PgSql                 = "pgsql"
	Sqlite                = "sqlite"
)

type databaseConfig struct {
	Enabled         bool           `required:"true"`
	Driver          DatabaseDriver `required:"true"`
	Username        string
	Password        string
	Host            string
	Port            string
	Database        string        `required:"true"`
	MaxOpenConns    int           `required:"true" split_words:"true"`
	MaxIdleConns    int           `required:"true" split_words:"true"`
	ConnMaxLifetime time.Duration `required:"true" split_words:"true"`
}

var databaseCfg *databaseConfig

func Database() *databaseConfig {
	if databaseCfg == nil {
		databaseCfg = new(databaseConfig)
		err := envconfig.Process("db", databaseCfg)

		if err != nil {
			panic(err)
		}
	}
	if !databaseCfg.Enabled {
		return databaseCfg
	}

	switch databaseCfg.Driver {
	case MySql:
	case PgSql:
	case Sqlite:
	default:
		panic(fmt.Sprintf("Database driver [%s] does not exist", databaseCfg.Driver))
	}

	return databaseCfg
}
