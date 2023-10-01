package admin

import (
	"net/http"

	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/provider"
	"github.com/gin-gonic/gin"
)

type AdminSearchPolicySubjectResponse struct {
	Subjects []string `json:"subjects" validate:"required"`
}

// @tags admin
// @accept json
// @produce json
// @param Authorization header string true "bearer token"
// @success 200 {object} admin.AdminSearchPolicySubjectResponse
// @failure 401 {object} response.ErrorResponse "code: 401-001(access token is wrong)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(casbin authorization failed)"
// @failure 500 {object} response.ErrorResponse "code: 500-001"
// @router /api/admin/policy/subject [get]
func SearchPolicySubject(c *gin.Context) {
	respBody := new(AdminSearchPolicySubjectResponse)
	respBody.Subjects = provider.Registry.Casbin().GetAllSubjects()

	c.JSON(http.StatusOK, respBody)
}
