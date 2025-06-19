package main

import (
	"os"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/registrar"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/scheduler"
	"github.com/chan-jui-huang/go-backend-framework/v2/pkg/app"
	"github.com/chan-jui-huang/go-backend-framework/v2/pkg/booter"
	"github.com/chan-jui-huang/go-backend-framework/v2/pkg/booter/config"
	"github.com/chan-jui-huang/go-backend-framework/v2/pkg/booter/service"
)

func init() {
	booter.Boot(
		func() {},
		booter.NewConfigWithCommand,
		&registrar.RegisterExecutor,
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
		scheduler.Start,
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
			logger.Info("app is terminating")
		},
		scheduler.Stop,
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
