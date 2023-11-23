package main

import (
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http"
	"github.com/ChanJuiHuang/go-backend-framework/internal/registrar"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/app"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
)

func init() {
	booter.Boot(
		func() {},
		booter.NewConfigWithCommand,
		&registrar.ConfigRegistrar,
		&registrar.ServiceRegistrar,
	)
}

// @title Example API
// @version 1.0
// @schemes http https
// @host localhost:8080
func main() {
	httpServer := http.NewServer(config.Registry.Get("httpServer").(http.ServerConfig))
	logger := service.Registry.Get("logger").(*zap.Logger)
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
