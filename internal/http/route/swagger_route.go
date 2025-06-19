package route

import (
	_ "github.com/chan-jui-huang/go-backend-framework/v2/docs"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/config"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type SwaggerRouter struct {
	router *gin.RouterGroup
}

func NewSwaggerRouter(router *gin.Engine) *SwaggerRouter {
	return &SwaggerRouter{
		router: router.Group(""),
	}
}

// type [http://localhost:8080/swagger/index.html] in browser to watch the swagger api doc
func (sr *SwaggerRouter) AttachRoutes() {
	booterConfig := config.Registry.Get("booter").(booter.Config)
	if !booterConfig.Debug {
		return
	}
	sr.router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
