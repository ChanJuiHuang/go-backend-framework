package test

import (
	"os"
	"path"
	"runtime"
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
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime caller cannot get file information")
	}
	wd := path.Join(path.Dir(file), "../..")

	globalConfig := newGlobalConfig(wd)
	registerGlobalConfig(globalConfig)
	setEnv(*globalConfig)
	registerConfig(*globalConfig)

	registerProvider()

	HttpHandler = NewHttpHandler()
	Migration = NewMigration()
}

func newGlobalConfig(rootDir string) *global.Config {
	return &global.Config{
		RootDir:  rootDir,
		Timezone: "UTC",
		Debug:    false,
		Testing:  true,
	}
}

func registerGlobalConfig(globalConfig *global.Config) {
	config.Registry.Set("global", globalConfig)
}

func setEnv(globalConfig global.Config) {
	err := godotenv.Load(path.Join(globalConfig.RootDir, ".env.testing"))
	if err != nil {
		panic(err)
	}

	err = os.Setenv("TZ", globalConfig.Timezone)
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)
}

func registerConfig(globalConfig global.Config) {
	byteYaml, err := os.ReadFile(path.Join(globalConfig.RootDir, "config.testing.yml"))
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
		"middleware.rateLimit":         &middleware.RateLimitConfig{},
	})
}

func registerProvider() {
	logger, consoleLogger, _ := provider.ProvideLogger()
	db := provider.ProvideDB()

	provider.Registry.Register(
		logger,
		consoleLogger,
		nil,
		db,
		provider.ProvideRedis(),
		provider.ProvideAuthenticator(),
		provider.ProvideCasbin(db),
	)
}
