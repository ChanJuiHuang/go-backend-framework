package main

import (
	_ "github.com/joho/godotenv/autoload"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/provider"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/app"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/config"
)

func init() {
	globalConfig := newGlobalConfig()
	registerGlobalConfig(globalConfig)
	setEnv(*globalConfig)
	registerConfig(*globalConfig)

	registerProvider()
}

// @title Example API
// @version 1.0
// @schemes http https
// @host localhost:8080
func main() {
	httpServer := http.NewServer(config.Registry.Get("httpServer").(http.ServerConfig))
	app := app.New(
		[]app.StartingCallback{
			httpServer.GracefulShutdown,
		},
		[]app.StartedCallback{
			func() {
				provider.Registry.Logger().Info("app is started")
			},
		},
		[]app.TerminatedCallback{
			func() {
				provider.Registry.Logger().Info("app is terminated")
			},
		},
		httpServer.Run,
	)

	app.Run()
}
