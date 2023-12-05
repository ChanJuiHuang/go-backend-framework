package main

import (
	"fmt"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/route"
	"github.com/ChanJuiHuang/go-backend-framework/internal/registrar"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter"
	"github.com/gin-gonic/gin"
)

func init() {
	booter.Boot(
		func() {},
		booter.NewDefaultConfig,
		&registrar.RegisterExecutor,
	)
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	route.AttachApiRoutes(engine)

	for _, routeInfo := range engine.Routes() {
		fmt.Printf("method: [%s], path: [%s], handler: [%s]\n", routeInfo.Method, routeInfo.Path, routeInfo.Handler)
	}
}
