package provider

import (
	"github.com/ChanJuiHuang/go-backend-framework/pkg/provider"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/mold/v4/modifiers"
)

func RegisterService() {
	logger, consoleLogger, fileLogger := ProvideLogger()
	db := ProvideDB()

	provider.Registry.Register(map[string]any{
		"logger":         logger,
		"consoleLogger":  consoleLogger,
		"fileLogger":     fileLogger,
		"database":       db,
		"redis":          ProvideRedis(),
		"authenticator":  ProvideAuthenticator(),
		"casbinEnforcer": ProvideCasbinEnforcer(db),
		"formDecoder":    form.NewDecoder(),
		"modifier":       modifiers.New(),
	})
}
