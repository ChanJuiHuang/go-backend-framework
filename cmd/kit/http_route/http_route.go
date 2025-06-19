package main

import (
	"fmt"

	"github.com/ChanJuiHuang/go-backend-framework/v2/internal/http/route"
	"github.com/ChanJuiHuang/go-backend-framework/v2/internal/registrar"
	"github.com/ChanJuiHuang/go-backend-framework/v2/pkg/booter"
	"github.com/gin-gonic/gin"
)

func init() {
	booter.Boot(
		func() {},
		booter.NewDefaultConfig,
		&registrar.SimpleRegisterExecutor,
	)
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	routers := []route.Router{
		route.NewApiRouter(engine),
		route.NewSwaggerRouter(engine),
	}
	for _, router := range routers {
		router.AttachRoutes()
	}

	for _, routeInfo := range engine.Routes() {
		fmt.Printf("method: [%s], path: [%s], handler: [%s]\n", routeInfo.Method, routeInfo.Path, routeInfo.Handler)
	}
}
