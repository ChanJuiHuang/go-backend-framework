package main

import (
	"os"
	"syscall"

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

	startingCallbacks := []func(){
		func() {
			logger.Info("app is starting")
		},
	}
	startedCallbacks := []func(){
		func() {
			logger.Info("app is started")
		},
	}
	signalCallbacks := []app.SignalCallback{
		{
			Signals: []os.Signal{
				syscall.SIGINT,
				syscall.SIGTERM,
			},
			CallbackFunc: httpServer.GracefulShutdown,
		},
	}
	asyncCallbacks := []func(){}
	terminatedCallbacks := []func(){
		func() {
			logger.Info("app is terminated")
		},
	}

	app := app.New(
		startingCallbacks,
		startedCallbacks,
		signalCallbacks,
		asyncCallbacks,
		terminatedCallbacks,
	)

	app.Run(httpServer.Run)
}
