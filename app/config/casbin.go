package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type casbinConfig struct {
	Enabled              bool          `required:"true"`
	ModelPath            string        `required:"true" split_words:"true"`
	AutoLoadTimeInterval time.Duration `required:"true" split_words:"true"`
}

var casbinCfg *casbinConfig

func Casbin() *casbinConfig {
	if casbinCfg == nil {
		casbinCfg = new(casbinConfig)
		err := envconfig.Process("casbin", casbinCfg)

		if err != nil {
			panic(err)
		}
	}

	return casbinCfg
}
