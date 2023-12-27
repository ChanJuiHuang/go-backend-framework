package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserGetPolicyData struct {
	Rules []string `json:"rules" validate:"required"`
}

// @tags user
// @accept json
// @produce json
// @param Authorization header string true "bearer token"
// @success 200 {object} response.Response{data=UserGetPolicyData}
// @failure 400 {object} response.ErrorResponse "code: 400-001(get user policy failed)"
// @failure 401 {object} response.ErrorResponse "code: 401-001(access token is wrong)"
// @failure 500 {object} response.ErrorResponse "code: 500-001"
// @router /api/user/policy [get]
func GetPolicy(c *gin.Context) {
	enforcer := service.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
	userId := c.GetUint("user_id")
	policies, err := enforcer.GetImplicitPermissionsForUser(strconv.FormatInt(int64(userId), 10))
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, err, nil)
		logger := service.Registry.Get("logger").(*zap.Logger)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	data := &UserGetPolicyData{
		Rules: make([]string, 0, len(policies)),
	}
	for i := 0; i < len(policies); i++ {
		data.Rules = append(data.Rules, fmt.Sprintf("%s %s", policies[i][2], policies[i][1]))
	}
	c.JSON(http.StatusOK, response.NewResponse(data))
}
