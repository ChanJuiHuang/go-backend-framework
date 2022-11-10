package http

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ChanJuiHuang/go-backend-framework/app/config"
	"github.com/ChanJuiHuang/go-backend-framework/app/http/middleware"
	"github.com/ChanJuiHuang/go-backend-framework/app/http/route"
	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/gin-gonic/gin"
)

func GetEngine() *gin.Engine {
	engine := gin.New()
	engine.RemoteIPHeaders = []string{
		"X-Forwarded-For",
		"X-Real-IP",
	}
	engine.SetTrustedProxies([]string{
		"0.0.0.0/0",
		"::/0",
	})
	middleware.AttachGlobalMiddleware(engine)
	route.AttachApiRoutes(engine)
	route.AttachSchedulerRoutes(engine)
	route.AddSwaggerRoute(engine)

	return engine
}

func gracefulShutdown(ctx context.Context, srv *http.Server, wg *sync.WaitGroup) {
	defer wg.Done()
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), config.Http().GracefulShutdownWaitingTime)
	defer cancel()

	provider.App.Logger.Info("start shutdown server")
	if err := srv.Shutdown(ctx); err != nil {
		provider.App.Logger.Error(err.Error())
	}
	provider.App.Logger.Info("end shutdown server")
}

func autoLoadCasbinPolicy(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	casbinCfg := config.Casbin()
	ticker := time.NewTicker(casbinCfg.AutoLoadTimeInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !casbinCfg.Enabled {
				continue
			}
			provider.App.Logger.Info("load casbin policy")
			if err := provider.App.Casbin.LoadPolicy(); err != nil {
				provider.App.Logger.Error(err.Error())
			}
		case <-ctx.Done():
			provider.App.Logger.Info("stop auto load casbin policy")
			return
		}
	}
}

func beforeRunServerCallback(srv *http.Server, wg *sync.WaitGroup) {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	go func() {
		<-ctx.Done()
		stop()
	}()

	wg.Add(2)
	go gracefulShutdown(ctx, srv, wg)
	go autoLoadCasbinPolicy(ctx, wg)
}

func RunServer() {
	srv := http.Server{
		Addr:    config.Http().Address,
		Handler: GetEngine(),
	}
	var wg sync.WaitGroup
	beforeRunServerCallback(&srv, &wg)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		provider.App.Logger.Error(err.Error())
	}
	wg.Wait()
}
