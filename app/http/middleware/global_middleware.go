package middleware

import (
	"github.com/ChanJuiHuang/go-backend-framework/app/config"
	"github.com/gin-gonic/gin"
)

func AttachGlobalMiddleware(router *gin.Engine) {
	handlerFunctions := []gin.HandlerFunc{
		RateLimit(),
		VerifyCsrfToken(),
	}

	switch config.Log().Type {
	case config.Console:
		handlerFunctions = append(
			[]gin.HandlerFunc{
				RecordAccessLog(),
				gin.Recovery(),
			},
			handlerFunctions...,
		)
	default:
		handlerFunctions = append(
			[]gin.HandlerFunc{
				RecordAccessLog(),
				RecoverServer(),
			},
			handlerFunctions...,
		)
	}

	router.Use(
		handlerFunctions...,
	)
}
