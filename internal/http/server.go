package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/middleware"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/route"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	server *http.Server
	config ServerConfig
}

type ServerConfig struct {
	Address             string
	GracefulShutdownTtl time.Duration
}

func NewEngine() (*gin.Engine, error) {
	engine := gin.New()
	engine.RemoteIPHeaders = []string{
		"X-Forwarded-For",
		"X-Real-IP",
	}
	err := engine.SetTrustedProxies([]string{
		"0.0.0.0/0",
		"::/0",
	})

	return engine, err
}

func NewServer(config ServerConfig) *Server {
	srv := &Server{
		server: &http.Server{
			Addr:              config.Address,
			ReadHeaderTimeout: 30 * time.Minute,
		},
		config: config,
	}

	return srv
}

func (srv *Server) Run() {
	engine, err := NewEngine()
	if err != nil {
		panic(err)
	}
	middleware.AttachGlobalMiddleware(engine)

	routers := []route.Router{
		route.NewApiRouter(engine),
		route.NewSwaggerRouter(engine),
	}
	for _, router := range routers {
		router.AttachRoutes()
	}

	srv.server.Handler = engine.Handler()
	logger := service.Registry.Get("logger").(*zap.Logger)

	if err := srv.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(err.Error())
	}
}

func (srv *Server) GracefulShutdown(ctx context.Context) {
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), srv.config.GracefulShutdownTtl)
	defer cancel()

	logger := service.Registry.Get("logger").(*zap.Logger)
	logger.Info("server start to shutdown")
	if err := srv.server.Shutdown(ctx); err != nil {
		logger.Error(err.Error())
	}
	logger.Info("server end to shutdown")
}
