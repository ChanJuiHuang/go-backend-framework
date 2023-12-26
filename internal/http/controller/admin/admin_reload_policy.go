package admin

import (
	"net/http"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// @tags admin
// @summary reload the policies
// @description reload the policies
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true "csrf token"
// @param Authorization header string true "bearer token"
// @success 204 "no content"
// @failure 400 {object} response.ErrorResponse "code: 400-001(load policy is failed)"
// @failure 401 {object} response.ErrorResponse "code: 401-001(access token is wrong)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(csrf token mismatch, casbin authorization failed)"
// @failure 500 {object} response.ErrorResponse "code: 500-001"
// @router /api/admin/policy/reload [post]
func ReloadPolicy(c *gin.Context) {
	enforcer := service.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
	if err := enforcer.LoadPolicy(); err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil)
		logger := service.Registry.Get("logger").(*zap.Logger)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	c.Status(http.StatusNoContent)
}
