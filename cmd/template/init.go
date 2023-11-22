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
		RootDir: wd,
		Debug:   false,
		Testing: false,
	}

	return globalConfig
}

func registerGlobalConfig(globalConfig *global.Config) {
	config.Registry.Set("global", globalConfig)
}
