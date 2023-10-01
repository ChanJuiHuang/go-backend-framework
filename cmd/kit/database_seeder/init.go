package main

import (
	"os"
	"path"
	"strings"

	"github.com/ChanJuiHuang/go-backend-framework/internal/global"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/provider"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/authentication"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/database"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/logger"
	"github.com/spf13/viper"
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

func registerConfig(globalConfig global.Config) {
	byteYaml, err := os.ReadFile(path.Join(globalConfig.RootDir, "config.yml"))
	if err != nil {
		panic(err)
	}
	stringYaml := os.ExpandEnv(string(byteYaml))

	v := viper.New()
	v.SetConfigType("yaml")
	v.ReadConfig(strings.NewReader(stringYaml))

	config.Registry.SetViper(v)
	config.Registry.Register(map[string]any{
		"logger.console":               &logger.ConsoleConfig{},
		"logger.file":                  &logger.FileConfig{},
		"database":                     &database.Config{},
		"authentication.authenticator": &authentication.Config{},
	})
}

func registerProvider() {
	logger, consoleLogger, fileLogger := provider.ProvideLogger()
	db := provider.ProvideDB()

	provider.Registry.Register(
		logger,
		consoleLogger,
		fileLogger,
		db,
		nil,
		provider.ProvideAuthenticator(),
		nil,
	)
}