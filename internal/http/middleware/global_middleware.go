package middleware

import (
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/config"
	"github.com/gin-gonic/gin"
)

func AttachGlobalMiddleware(router *gin.Engine) {
	csrfConfig := config.Registry.Get("middleware.csrf").(CsrfConfig)

	handlerFunctions := []gin.HandlerFunc{
		AccessLogger(),
		Recover(),
		VerifyCsrfToken(csrfConfig),
	}

	router.Use(
		handlerFunctions...,
	)
}
