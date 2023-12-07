package admin

import (
	"net/http"
	"strconv"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type AdminGetGroupingPolicyData struct {
	UserId   uint     `json:"user_id" mapstructure:"user_id" validate:"required"`
	Subjects []string `json:"subjects" mapstructure:"subjects" validate:"required"`
}

// @tags admin
// @accept json
// @produce json
// @param Authorization header string true "bearer token"
// @success 200 {object} response.Response{data=admin.AdminGetGroupingPolicyData}
// @failure 401 {object} response.ErrorResponse "code: 401-001(access token is wrong)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(casbin authorization failed)"
// @failure 500 {object} response.ErrorResponse "code: 500-001"
// @router /api/admin/grouping-policy [get]
func GetGroupingPolicy(c *gin.Context) {
	userId := c.Param("userId")
	enforcer := service.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
	groupingPolicies := enforcer.GetFilteredGroupingPolicy(0, userId)

	id, err := strconv.Atoi(userId)
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil)
		logger := service.Registry.Get("logger").(*zap.Logger)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}
	data := AdminGetGroupingPolicyData{
		UserId:   uint(id),
		Subjects: make([]string, len(groupingPolicies)),
	}

	for i := 0; i < len(groupingPolicies); i++ {
		data.Subjects[i] = groupingPolicies[i][1]
	}

	c.JSON(http.StatusOK, response.NewResponse(data))
}
