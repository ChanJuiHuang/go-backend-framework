package user

import (
	"net/http"
	"time"

	"github.com/ChanJuiHuang/go-backend-framework/app/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/app/module/user/model"
	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/ChanJuiHuang/go-backend-framework/util"
	"github.com/gin-gonic/gin"
)

// @tags user
// @produce json
// @success 204
// @failure 400 {object} response.SwaggerErrorResponse
// @failure 401 {object} response.SwaggerErrorResponse
// @router /scheduler/refresh-token-record [delete]
func DeleteRefreshTokenRecord(c *gin.Context) {
	db := provider.App.DB.Where("expire_at < ?", time.Now()).
		Delete(&model.RefreshTokenRecord{})
	if util.IsDatabaseError(db) {
		response.NewHttpErrorResponse(response.ErrBadRequest).
			MakePreviousMessage(util.WrapError(db.Error)).
			ToJson(c)
		return
	}

	c.Status(http.StatusNoContent)
}
