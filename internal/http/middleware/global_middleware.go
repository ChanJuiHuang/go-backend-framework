package middleware

import (
	"github.com/ChanJuiHuang/go-backend-framework/pkg/config"
	"github.com/gin-gonic/gin"
)

func AttachGlobalMiddleware(router *gin.Engine) {
	csrfConfig := config.Registry.Get("middleware.csrf").(CsrfConfig)

	handlerFunctions := []gin.HandlerFunc{
		VerifyCsrfToken(csrfConfig),
	}

	router.Use(
		handlerFunctions...,
	)
}
