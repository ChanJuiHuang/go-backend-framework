package middleware

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/provider"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Recover() gin.HandlerFunc {
	logger := provider.Registry.Get("logger").(*zap.Logger)

	return func(c *gin.Context) {
		defer func() {
			err := recover()
			if err == nil {
				return
			}
			// condition that warrants a panic stack trace.
			var isBrokenPipe bool
			if ne, ok := err.(*net.OpError); ok {
				var se *os.SyscallError
				if errors.As(ne, &se) {
					if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
						isBrokenPipe = true
					}
				}
			}

			errResp := response.NewErrorResponse(response.InternalServerError, errors.New(fmt.Sprintf("%v", err)), nil)
			logger.Error(response.InternalServerError, errResp.MakeLogFields(c.Request)...)
			if isBrokenPipe {
				c.Abort()
			} else {
				c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
			}
		}()
		c.Next()
	}
}
