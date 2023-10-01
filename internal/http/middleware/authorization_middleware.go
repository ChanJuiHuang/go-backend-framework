package middleware

import (
	"strconv"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/provider"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func Authorize() gin.HandlerFunc {
	logger := provider.Registry.Logger()
	return func(c *gin.Context) {
		userId := c.GetUint("user_id")
		ok, err := provider.Registry.Casbin().Enforce(strconv.FormatUint(uint64(userId), 10), c.Request.URL.Path, c.Request.Method)
		if err != nil {
			errResp := response.NewErrorResponse(response.Forbidden, errors.WithStack(err), nil)
			logger.Warn(response.Forbidden, errResp.MakeLogFields(c.Request)...)
			c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
			return
		}

		if !ok {
			errResp := response.NewErrorResponse(response.Forbidden, errors.New("casbin authorization failed"), nil)
			logger.Warn(response.Forbidden, errResp.MakeLogFields(c.Request)...)
			c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
			return
		}

		c.Next()
	}
}
