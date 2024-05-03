package admin

import (
	"net/http"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	casbinrule "github.com/ChanJuiHuang/go-backend-framework/internal/pkg/casbin_rule"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AdminDeletePolicySubjectRequest struct {
	Subjects []string `json:"subjects" binding:"required"`
}

type AdminDeletePolicySubjectData struct {
	Subjects []string `json:"subjects" mapstructure:"subjects" validate:"required"`
}

// @tags admin
// @summary delete roles
// @description delete roles
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true "csrf token"
// @param Authorization header string true "bearer token"
// @param request body admin.AdminDeletePolicySubjectRequest true "delete policy subject"
// @success 200 {object} response.Response{data=admin.AdminDeletePolicySubjectData}
// @failure 401 {object} response.ErrorResponse "code: 401-001(access token is wrong)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(casbin authorization failed)"
// @failure 500 {object} response.ErrorResponse "code: 500-001"
// @router /api/admin/policy/subject [delete]
func DeletePolicySubject(c *gin.Context) {
	reqBody := new(AdminDeletePolicySubjectRequest)
	logger := service.Registry.Get("logger").(*zap.Logger)
	if err := c.ShouldBindJSON(reqBody); err != nil {
		errResp := response.NewErrorResponse(response.RequestValidationFailed, errors.WithStack(err), nil)
		logger.Warn(response.RequestValidationFailed, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	database := service.Registry.Get("database").(*gorm.DB)
	err := database.Transaction(func(tx *gorm.DB) error {
		if err := casbinrule.Delete(tx, "ptype = ? AND v0 IN ?", "p", reqBody.Subjects); err != nil {
			logger.Warn(err.Error())
			return err
		}

		if err := casbinrule.Delete(tx, "ptype = ? AND v1 IN ?", "g", reqBody.Subjects); err != nil {
			logger.Warn(err.Error())
			return err
		}

		return nil
	})

	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	enforcer := service.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
	if err := enforcer.LoadPolicy(); err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	respBody := response.NewResponse(AdminDeletePolicySubjectData{
		Subjects: enforcer.GetAllSubjects(),
	})
	c.JSON(http.StatusOK, respBody)
}
