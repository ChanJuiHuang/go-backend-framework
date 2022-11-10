package middleware

import (
	"time"

	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func RecordAccessLog() gin.HandlerFunc {
	skipPaths := map[string]bool{
		"/skip-path": true,
	}
	logger := provider.NewAccessLogger()

	return func(c *gin.Context) {
		now := time.Now()
		path := c.Request.URL.Path
		c.Next()

		if skipPaths[path] {
			return
		}
		latency := time.Since(now)
		status := c.Writer.Status()
		fields := []zapcore.Field{
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("referer", c.Request.Referer()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
		}

		switch {
		case status < 400:
			logger.Info(path, fields...)
		case status >= 400 && status < 500:
			logger.Warn(path, fields...)
		case status >= 500:
			logger.Error(path, fields...)
		}
	}
}
