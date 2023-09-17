package middleware

import (
	"github.com/gin-gonic/gin"
)

func AttachGlobalMiddleware(router *gin.Engine) {
	handlerFunctions := []gin.HandlerFunc{}

	router.Use(
		handlerFunctions...,
	)
}
