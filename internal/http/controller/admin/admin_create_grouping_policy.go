package admin

import (
	"net/http"
	"strconv"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/provider"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type AdminCreateGroupingPolicyRequest struct {
	UserId   uint     `json:"user_id" binding:"required"`
	Subjects []string `json:"subjects" binding:"required"`
}

type AdminCreateGroupingPolicyResponse struct {
	UserId   uint     `json:"user_id" validate:"required"`
	Subjects []string `json:"subjects" validate:"required"`
}

// @tags admin
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true "csrf token"
// @param Authorization header string true "bearer token"
// @param request body admin.AdminCreateGroupingPolicyRequest true "create grouping policy"
// @success 200 {object} admin.AdminCreateGroupingPolicyResponse
// @failure 400 {object} response.ErrorResponse "code: 400-001(add grouping policy is failed), 400-002(request validation failed)"
// @failure 401 {object} response.ErrorResponse "code: 401-001(access token is wrong)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(csrf token mismatch, casbin authorization failed)"
// @failure 500 {object} response.ErrorResponse "code: 500-001"
// @router /api/admin/grouping-policy [post]
func CreateGroupingPolicy(c *gin.Context) {
	reqBody := new(AdminCreateGroupingPolicyRequest)
	logger := provider.Registry.Get("logger").(*zap.Logger)
	if err := c.ShouldBindJSON(reqBody); err != nil {
		errResp := response.NewErrorResponse(response.RequestValidationFailed, errors.WithStack(err), nil)
		logger.Warn(response.RequestValidationFailed, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	userId := strconv.Itoa(int(reqBody.UserId))
	groupingPolicies := make([][]string, 0, len(reqBody.Subjects))
	for _, subject := range reqBody.Subjects {
		groupingPolicies = append(groupingPolicies, []string{userId, subject})
	}

	enforcer := provider.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
	result, err := enforcer.AddGroupingPolicies(groupingPolicies)
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}
	if !result {
		errResp := response.NewErrorResponse(response.OneOfGroupingPolicyIsRepeat, errors.New(response.OneOfGroupingPolicyIsRepeat), nil)
		logger.Warn(response.OneOfGroupingPolicyIsRepeat, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	groupingPolicies = enforcer.GetFilteredGroupingPolicy(0, userId)
	subjects := make([]string, len(groupingPolicies))
	for i := 0; i < len(subjects); i++ {
		subjects[i] = groupingPolicies[i][1]
	}

	if err := enforcer.LoadPolicy(); err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}
	c.JSON(http.StatusOK, &AdminCreateGroupingPolicyResponse{
		UserId:   reqBody.UserId,
		Subjects: subjects,
	})
}
