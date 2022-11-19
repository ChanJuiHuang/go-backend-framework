package http

import (
	"fmt"

	"github.com/ChanJuiHuang/go-backend-framework/app/http/route"
	"github.com/gin-gonic/gin"
)

func RouteList() {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	route.AttachApiRoutes(engine)
	route.AttachSchedulerRoutes(engine)

	for _, routeInfo := range engine.Routes() {
		fmt.Printf("method: [%s], path: [%s], handler: [%s]\n", routeInfo.Method, routeInfo.Path, routeInfo.Handler)
	}
}
