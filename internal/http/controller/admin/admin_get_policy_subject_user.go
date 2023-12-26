package admin

import (
	"net/http"
	"strconv"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AdminGetPolicySubjectUserData struct {
	UserIds []uint `json:"user_ids" mapstructure:"user_ids" validate:"required"`
}

// @tags admin
// @summary get user ids in the role
// @description get user ids in the role
// @accept json
// @produce json
// @param Authorization header string true "bearer token"
// @param subject path string true "subject"
// @success 200 {object} response.Response{data=admin.AdminGetPolicySubjectUserData}
// @failure 401 {object} response.ErrorResponse "code: 401-001(access token is wrong)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(casbin authorization failed)"
// @failure 500 {object} response.ErrorResponse "code: 500-001"
// @router /api/admin/policy/subject/{subject}/user [get]
func GetPolicySubjectUser(c *gin.Context) {
	enforcer := service.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
	logger := service.Registry.Get("logger").(*zap.Logger)
	userIds, err := enforcer.GetUsersForRole(c.Param("subject"))
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, err, nil)
		logger.Error(errResp.Message, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	data := &AdminGetPolicySubjectUserData{
		UserIds: make([]uint, 0, len(userIds)),
	}

	for i := 0; i < len(userIds); i++ {
		userId, err := strconv.ParseUint(userIds[i], 10, 64)
		if err != nil {
			errResp := response.NewErrorResponse(response.BadRequest, err, nil)
			logger.Error(errResp.Message, errResp.MakeLogFields(c.Request)...)
			c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
			return
		}
		data.UserIds = append(data.UserIds, uint(userId))
	}

	c.JSON(http.StatusOK, response.NewResponse(data))
}
