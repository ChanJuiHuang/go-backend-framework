package admin

import (
	"net/http"

	"github.com/ChanJuiHuang/go-backend-framework/pkg/provider"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

type AdminGetPolicySubjectResponse struct {
	Subject string `json:"subject" validate:"required"`
	Rules   []Rule `json:"rules" validate:"required"`
}

// @tags admin
// @accept json
// @produce json
// @param Authorization header string true "bearer token"
// @param subject path string true "subject"
// @success 200 {object} admin.AdminGetPolicySubjectResponse
// @failure 401 {object} response.ErrorResponse "code: 401-001(access token is wrong)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(casbin authorization failed)"
// @failure 500 {object} response.ErrorResponse "code: 500-001"
// @router /api/admin/policy/subject/{subject} [get]
func GetPolicySubject(c *gin.Context) {
	subject := c.Param("subject")
	enforcer := provider.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
	policies := enforcer.GetFilteredPolicy(0, subject)
	rules := make([]Rule, len(policies))

	for i := 0; i < len(rules); i++ {
		rules[i].Object = policies[i][1]
		rules[i].Action = policies[i][2]
	}

	c.JSON(http.StatusOK, &AdminGetPolicySubjectResponse{
		Subject: subject,
		Rules:   rules,
	})
}
