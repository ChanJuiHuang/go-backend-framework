package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type attachApiRouteFunc func(router *gin.Engine)

var apiRouteGroups = []attachApiRouteFunc{}

// @produce json
// @success 200 {string} string "{"message": "pong"}"
// @router /api/ping [get]
func AttachApiRoutes(router *gin.Engine) {
	router.GET("/api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	for _, apiRouteGroup := range apiRouteGroups {
		apiRouteGroup(router)
	}
}
