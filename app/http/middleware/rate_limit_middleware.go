package middleware

import (
	"errors"

	"github.com/ChanJuiHuang/go-backend-framework/app/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/util"
	"github.com/gin-gonic/gin"
)

var errTokenBucketIsEmpty error = errors.New("token bucket is empty")

func RateLimit() gin.HandlerFunc {
	limiter := util.NewRateLimiter()
	skipPaths := map[string]bool{
		"/skip-path": true,
	}

	return func(c *gin.Context) {
		if skipPaths[c.Request.URL.Path] || limiter.Allow() {
			c.Next()
			return
		}
		response.NewHttpErrorResponse(response.ErrServiceUnavailable).
			MakePreviousMessage(errTokenBucketIsEmpty).
			AbortWithStatus(c)
	}
}
