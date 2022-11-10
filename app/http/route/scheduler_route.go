package route

import (
	"net/http"

	"github.com/ChanJuiHuang/go-backend-framework/app/http/middleware"
	"github.com/gin-gonic/gin"
)

// @produce json
// @success 200 {string} string
// @failure 401 {object} response.SwaggerErrorResponse
// @router /scheduler/refresh-token-record [get]
func AttachSchedulerRoutes(router *gin.Engine) {
	schedulerRoutes := router.Group("/scheduler").
		Use(
			middleware.Authenticate(),
			middleware.Authorize(),
		)
	{
		schedulerRoutes.GET("/welcome", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "welcome",
			})
		})
	}
}
