package middleware

import (
	"fmt"
	"path"
	"time"

	"github.com/ChanJuiHuang/go-backend-framework/internal/global"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/provider"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/config"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func AccessLogger() gin.HandlerFunc {
	skipPaths := map[string]bool{
		"/skip-path": true,
	}

	globalConfig := config.Registry.Get("global").(global.Config)
	fileConfig := config.Registry.Get("logger.file").(logger.FileConfig)
	fileConfig.LogPath = path.Join(globalConfig.RootDir, "storage/log/access.log")

	v := config.Registry.GetViper()
	var accessLogger *zap.Logger
	switch logger.Type(v.GetString("logger.type")) {
	case logger.Console:
		accessLogger = provider.Registry.ConsoleLogger()
	case "file":
		var err error
		accessLogger, err = logger.NewFileLogger(fileConfig, logger.JsonEncoder, []zap.Option{}...)
		if err != nil {
			panic(err)
		}
	default:
		accessLogger = provider.Registry.ConsoleLogger()
	}

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

		message := fmt.Sprintf("%s %s", c.Request.Method, path)
		switch {
		case status < 400:
			accessLogger.Info(message, fields...)
		case status >= 400 && status < 500:
			accessLogger.Warn(message, fields...)
		case status >= 500:
			accessLogger.Error(message, fields...)
		}
	}
}
