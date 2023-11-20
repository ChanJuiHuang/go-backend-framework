package admin

import (
	"net/http"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/provider"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AdminDeletePolicySubjectRequest struct {
	Subjects []string `json:"subjects" binding:"required"`
}

type AdminDeletePolicySubjectResponse struct {
	Subjects []string `json:"subjects" validate:"required"`
}

// @tags admin
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true "csrf token"
// @param Authorization header string true "bearer token"
// @param request body admin.AdminDeletePolicySubjectRequest true "delete policy subject"
// @success 200 {object} admin.AdminDeletePolicySubjectResponse
// @failure 401 {object} response.ErrorResponse "code: 401-001(access token is wrong)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(casbin authorization failed)"
// @failure 500 {object} response.ErrorResponse "code: 500-001"
// @router /api/admin/policy/subject [delete]
func DeletePolicySubject(c *gin.Context) {
	reqBody := new(AdminDeletePolicySubjectRequest)
	logger := provider.Registry.Get("logger").(*zap.Logger)
	if err := c.ShouldBindJSON(reqBody); err != nil {
		errResp := response.NewErrorResponse(response.RequestValidationFailed, errors.WithStack(err), nil)
		logger.Warn(response.RequestValidationFailed, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	database := provider.Registry.Get("database").(*gorm.DB)
	err := database.Transaction(func(tx *gorm.DB) error {
		tx1 := tx.Table("casbin_rules").
			Where("ptype = ?", "p").
			Where(map[string]any{"v0": reqBody.Subjects}).
			Delete(&struct{}{})
		if err := tx1.Error; err != nil {
			logger.Warn(err.Error())
			return err
		}

		tx2 := tx.Table("casbin_rules").
			Where("ptype = ?", "g").
			Where(map[string]any{"v1": reqBody.Subjects}).
			Delete(&struct{}{})
		if err := tx2.Error; err != nil {
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

	enforcer := provider.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
	if err := enforcer.LoadPolicy(); err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}
	c.JSON(http.StatusOK, &AdminDeletePolicySubjectResponse{
		Subjects: enforcer.GetAllSubjects(),
	})
}
