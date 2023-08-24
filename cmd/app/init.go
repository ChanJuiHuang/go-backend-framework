package main

import (
	"flag"
	"os"
	"path"
	"strings"

	"github.com/ChanJuiHuang/go-backend-framework/internal/global"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/middleware"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/provider"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/authentication"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/database"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/logger"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/redis"
	"github.com/spf13/viper"
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
		"redis":                        &redis.Config{},
		"authentication.authenticator": &authentication.Config{},
		"httpServer":                   &http.ServerConfig{},
		"middleware.csrf":              &middleware.CsrfConfig{},
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
		provider.ProvideRedis(),
		provider.ProvideAuthenticator(),
		provider.ProvideCasbin(db),
	)
}
