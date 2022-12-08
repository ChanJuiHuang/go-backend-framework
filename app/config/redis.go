package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type redisConfig struct {
	Enabled         bool          `required:"true"`
	Address         string        `required:"true"`
	Password        string        `required:"true"`
	DB              int           `required:"true"`
	MinIdleConns    int           `required:"true" split_words:"true"`
	ConnMaxLifetime time.Duration `required:"true" split_words:"true"`
}

var redisCfg *redisConfig

func Redis() *redisConfig {
	if redisCfg == nil {
		redisCfg = new(redisConfig)
		err := envconfig.Process("redis", redisCfg)

		if err != nil {
			panic(err)
		}
	}

	return redisCfg
}
