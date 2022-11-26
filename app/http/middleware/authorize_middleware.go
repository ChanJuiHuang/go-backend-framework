package middleware

import (
	"strconv"

	"github.com/ChanJuiHuang/go-backend-framework/app/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/ChanJuiHuang/go-backend-framework/util"
	"github.com/gin-gonic/gin"
)

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := uint64(c.MustGet("user_id").(uint))
		ok, err := provider.App.Casbin.Enforce(strconv.FormatUint(userId, 10), c.Request.URL.Path, c.Request.Method)

		if err != nil {
			response.NewHttpErrorResponse(response.ErrAuthorizationFailed).
				MakePreviousMessage(util.WrapError(err)).
				AbortWithJson(c)
			return
		}

		if !ok {
			response.NewHttpErrorResponse(response.ErrAuthorizationFailed).
				AbortWithJson(c)
			return
		}

		c.Next()
	}
}
