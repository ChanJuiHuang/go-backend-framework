package main

import (
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http"
	internalProvider "github.com/ChanJuiHuang/go-backend-framework/internal/provider"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/app"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/provider"
)

func init() {
	globalConfig := newGlobalConfig()
	registerGlobalConfig(globalConfig)
	setEnv(*globalConfig)
	registerConfig(*globalConfig)
	internalProvider.RegisterService()
}

// @title Example API
// @version 1.0
// @schemes http https
// @host localhost:8080
func main() {
	httpServer := http.NewServer(config.Registry.Get("httpServer").(http.ServerConfig))
	logger := provider.Registry.Get("logger").(*zap.Logger)
	app := app.New(
		[]app.StartingCallback{
			httpServer.GracefulShutdown,
		},
		[]app.StartedCallback{
			func() {
				logger.Info("app is started")
			},
		},
		[]app.TerminatedCallback{
			func() {
				logger.Info("app is terminated")
			},
		},
		httpServer.Run,
	)

	app.Run()
}
