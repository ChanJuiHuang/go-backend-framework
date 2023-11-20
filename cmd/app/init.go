package main

import (
	"flag"
	"os"

	"github.com/ChanJuiHuang/go-backend-framework/internal/global"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/config"
)

func newGlobalConfig() *global.Config {
	globalConfig := &global.Config{}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	flag.StringVar(&globalConfig.RootDir, "rootDir", wd, "root directory which the executable file in")
	flag.StringVar(&globalConfig.Timezone, "timezone", "UTC", "timezone")
	flag.BoolVar(&globalConfig.Debug, "debug", false, "debug mode")
	flag.BoolVar(&globalConfig.Testing, "testing", false, "does run in testing mode")
	flag.Parse()

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
