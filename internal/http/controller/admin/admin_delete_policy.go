package admin

import (
	"net/http"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/provider"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type AdminDeletePolicyRequest struct {
	Subject string `json:"subject" binding:"required"`
	Rules   []Rule `json:"rules" binding:"required,dive"`
}

type AdminDeletePolicyResponse struct {
	Subject string `json:"subject" validate:"required"`
	Rules   []Rule `json:"rules" validate:"required"`
}

// @tags admin
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true "csrf token"
// @param Authorization header string true "bearer token"
// @param request body admin.AdminDeletePolicyRequest true "delete policy"
// @success 200 {object} admin.AdminDeletePolicyResponse
// @failure 400 {object} response.ErrorResponse "code: 400-001(delete policy failed), 400-002(request validation failed)"
// @failure 401 {object} response.ErrorResponse "code: 401-001(access token is wrong)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(csrf token mismatch, casbin authorization failed)"
// @failure 500 {object} response.ErrorResponse "code: 500-001"
// @router /api/admin/policy [delete]
func DeletePolicy(c *gin.Context) {
	reqBody := new(AdminDeletePolicyRequest)
	logger := provider.Registry.Logger()
	if err := c.ShouldBindJSON(reqBody); err != nil {
		errResp := response.NewErrorResponse(response.RequestValidationFailed, errors.WithStack(err), nil)
		logger.Warn(response.RequestValidationFailed, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	policies := make([][]string, 0, len(reqBody.Rules))
	for _, rule := range reqBody.Rules {
		policies = append(policies, []string{reqBody.Subject, rule.Object, rule.Action})
	}

	enforcer := provider.Registry.Casbin()
	result, err := enforcer.RemovePolicies(policies)
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}
	if !result {
		provider.Registry.Logger().Warn("the policies has been deleted PROBABLY")
	}

	policies = enforcer.GetFilteredPolicy(0, reqBody.Subject)
	rules := make([]Rule, len(policies))
	for i := 0; i < len(rules); i++ {
		rules[i].Object = policies[i][1]
		rules[i].Action = policies[i][2]
	}

	if err := enforcer.LoadPolicy(); err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}
	c.JSON(http.StatusOK, &AdminDeletePolicyResponse{
		Subject: reqBody.Subject,
		Rules:   rules,
	})
}
