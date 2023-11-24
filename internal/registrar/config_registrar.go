package registrar

import (
	"github.com/ChanJuiHuang/go-backend-framework/internal/http"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/middleware"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/authentication"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/database"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/logger"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/redis"
)

var ConfigRegistrar configRegistrar

type configRegistrar struct{}

func (*configRegistrar) Register() {
	config.Registry.RegisterMany(map[string]any{
		"logger.console":               &logger.Config{},
		"logger.file":                  &logger.Config{},
		"logger.access":                &logger.Config{},
		"database":                     &database.Config{},
		"redis":                        &redis.Config{},
		"authentication.authenticator": &authentication.Config{},
		"httpServer":                   &http.ServerConfig{},
		"middleware.csrf":              &middleware.CsrfConfig{},
		"middleware.rateLimit":         &middleware.RateLimitConfig{},
	})
}
