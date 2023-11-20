package main

import (
	"os"

	"github.com/ChanJuiHuang/go-backend-framework/internal/global"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/config"
)

func newGlobalConfig() *global.Config {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	globalConfig := &global.Config{
		RootDir:  wd,
		Timezone: "UTC",
		Debug:    false,
		Testing:  false,
	}

	return globalConfig
}

func registerGlobalConfig(globalConfig *global.Config) {
	config.Registry.Set("global", globalConfig)
}

func setEnv(globalConfig global.Config) {
	err := os.Setenv("TZ", globalConfig.Timezone)
	if err != nil {
		panic(err)
	}
}
