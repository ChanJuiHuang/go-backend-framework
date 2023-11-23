package middleware

import (
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/config"
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
