package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type attachApiRouteFunc func(router *gin.Engine)

var apiRouteGroups = []attachApiRouteFunc{}

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
