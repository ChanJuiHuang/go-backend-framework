package middleware

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/ChanJuiHuang/go-backend-framework/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RecoverServer() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			logger := provider.App.Logger

			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne, &se) {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				if logger != nil {
					stack := util.Stacktrace(1)
					httpRequest, _ := httputil.DumpRequest(c.Request, false)
					headers := strings.Split(string(httpRequest), "\r\n")
					for idx, header := range headers {
						current := strings.Split(header, ":")
						if current[0] == "Authorization" {
							headers[idx] = current[0] + ": *"
						}
					}
					headersToStr := strings.Join(headers, "\r\n")
					if brokenPipe {
						logger.Error(
							fmt.Sprintf("%s", err),
							zap.String("headers", headersToStr),
						)
					} else if gin.IsDebugging() {
						logger.Error(
							fmt.Sprintf("%s", err),
							zap.String("headers", headersToStr),
							zap.String("stacktrace", string(stack)),
						)
					} else {
						logger.Error(
							fmt.Sprintf("%s", err),
							zap.String("stacktrace", string(stack)),
						)
					}
				}
				if brokenPipe {
					// If the connection is dead, we can't write a status to it.
					logger.Error(fmt.Sprintf("%s", err))
					c.Abort()
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()
		c.Next()
	}
}
