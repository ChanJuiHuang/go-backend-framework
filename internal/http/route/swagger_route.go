package route

import (
	_ "github.com/ChanJuiHuang/go-backend-framework/docs"
	"github.com/ChanJuiHuang/go-backend-framework/internal/global"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/config"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// type [http://localhost:8080/swagger/index.html] in browser to watch the swagger api doc
func AttachSwaggerRoute(router *gin.Engine) {
	globalConfig := config.Registry.Get("global").(global.Config)
	if !globalConfig.Debug {
		return
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
