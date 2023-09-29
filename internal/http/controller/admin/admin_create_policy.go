package admin

import (
	"net/http"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/provider"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Rule struct {
	Object string `json:"object" binding:"required"`
	Action string `json:"action" binding:"required,oneof=GET POST PUT PATCH DELETE"`
}

type AdminCreatePolicyRequest struct {
	Subject string `json:"subject" binding:"required"`
	Rules   []Rule `json:"rules" binding:"required,dive"`
}

type AdminCreatePolicyResponse struct {
	Subject string `json:"subject" validate:"required"`
	Rules   []Rule `json:"rules" validate:"required"`
}

// @tags admin
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true "csrf token"
// @param Authorization header string true "bearer token"
// @param request body admin.AdminCreatePolicyRequest true "create policy"
// @success 200 {object} admin.AdminCreatePolicyResponse
// @failure 400 {object} response.ErrorResponse "code: 400-001(add policy is failed), 400-002(request validation failed), 400-005(one of policy is repeat)"
// @failure 401 {object} response.ErrorResponse "code: 401-001(access token is wrong)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(csrf token mismatch, casbin authorization failed)"
// @failure 500 {object} response.ErrorResponse "code: 500-001"
// @router /api/admin/policy [post]
func CreatePolicy(c *gin.Context) {
	reqBody := new(AdminCreatePolicyRequest)
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
	result, err := enforcer.AddPolicies(policies)
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}
	if !result {
		errResp := response.NewErrorResponse(response.OneOfPolicyIsRepeat, errors.New(response.OneOfPolicyIsRepeat), nil)
		logger.Warn(response.OneOfPolicyIsRepeat, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
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
	c.JSON(http.StatusOK, &AdminCreatePolicyResponse{
		Subject: reqBody.Subject,
		Rules:   rules,
	})
}