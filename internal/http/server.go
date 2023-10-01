package http

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/middleware"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/route"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/provider"
	"github.com/gin-gonic/gin"
)

type Server struct {
	server *http.Server
	config ServerConfig
}

type ServerConfig struct {
	Address             string
	GracefulShutdownTtl time.Duration
}

func NewEngine() *gin.Engine {
	engine := gin.New()
	engine.RemoteIPHeaders = []string{
		"X-Forwarded-For",
		"X-Real-IP",
	}
	engine.SetTrustedProxies([]string{
		"0.0.0.0/0",
		"::/0",
	})

	return engine
}

func NewServer(config ServerConfig) *Server {
	srv := &Server{
		server: &http.Server{
			Addr: config.Address,
		},
		config: config,
	}

	return srv
}

func (srv *Server) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	engine := NewEngine()
	middleware.AttachGlobalMiddleware(engine)
	route.AttachApiRoutes(engine)
	route.AttachSwaggerRoute(engine)
	srv.server.Handler = engine.Handler()

	if err := srv.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		provider.Registry.Logger().Error(err.Error())
	}
}

func (srv *Server) GracefulShutdown(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), srv.config.GracefulShutdownTtl)
	defer cancel()

	logger := provider.Registry.Logger()
	logger.Info("server start to shutdown")
	if err := srv.server.Shutdown(ctx); err != nil {
		logger.Error(err.Error())
	}
	logger.Info("server end to shutdown")
}
