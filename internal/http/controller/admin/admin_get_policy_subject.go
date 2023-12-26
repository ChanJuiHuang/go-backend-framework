package admin

import (
	"net/http"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

type AdminGetPolicySubjectData struct {
	Subject string `json:"subject" validate:"required"`
	Rules   []Rule `json:"rules" validate:"required"`
}

// @tags admin
// @summary get permissions in the role
// @description get permissions in the role
// @accept json
// @produce json
// @param Authorization header string true "bearer token"
// @param subject path string true "subject"
// @success 200 {object} response.Response{data=admin.AdminGetPolicySubjectData}
// @failure 401 {object} response.ErrorResponse "code: 401-001(access token is wrong)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(casbin authorization failed)"
// @failure 500 {object} response.ErrorResponse "code: 500-001"
// @router /api/admin/policy/subject/{subject} [get]
func GetPolicySubject(c *gin.Context) {
	subject := c.Param("subject")
	enforcer := service.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
	policies := enforcer.GetFilteredPolicy(0, subject)
	rules := make([]Rule, len(policies))

	for i := 0; i < len(rules); i++ {
		rules[i].Object = policies[i][1]
		rules[i].Action = policies[i][2]
	}

	respBody := response.NewResponse(AdminGetPolicySubjectData{
		Subject: subject,
		Rules:   rules,
	})
	c.JSON(http.StatusOK, respBody)
}
