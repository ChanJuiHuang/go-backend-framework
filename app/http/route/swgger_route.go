package route

import (
	"github.com/ChanJuiHuang/go-backend-framework/app/config"
	_ "github.com/ChanJuiHuang/go-backend-framework/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// type [http://localhost:8080/swagger/index.html] in browser to watch the swagger api doc
func AddSwaggerRoute(router *gin.Engine) {
	if !config.App().Debug {
		return
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
