package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @produce json
// @success 200 {string} string
// @router /api/ping [get]
func AttachApiRoutes(router *gin.Engine) {
	apiRoutes := router.Group("/api")
	{
		apiRoutes.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}
}
