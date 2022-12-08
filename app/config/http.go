package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type httpConfig struct {
	Address                     string        `required:"true"`
	GracefulShutdownWaitingTime time.Duration `required:"true" split_words:"true"`
}

var httpCfg *httpConfig

func Http() *httpConfig {
	if httpCfg == nil {
		httpCfg = new(httpConfig)
		err := envconfig.Process("http", httpCfg)

		if err != nil {
			panic(err)
		}
	}

	return httpCfg
}
