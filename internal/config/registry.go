package config

import (
	"os"
	"path"
	"strings"

	"github.com/ChanJuiHuang/go-backend-framework/internal/global"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/middleware"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/authentication"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/database"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/logger"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/redis"
	"github.com/spf13/viper"
)

func RegisterConfigWithFile(globalConfig global.Config, filename string) {
	byteYaml, err := os.ReadFile(path.Join(globalConfig.RootDir, filename))
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

func RegisterConfig(globalConfig global.Config) {
	RegisterConfigWithFile(globalConfig, "config.yml")
}
