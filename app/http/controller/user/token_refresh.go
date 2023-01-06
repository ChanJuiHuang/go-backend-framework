package user

import (
	"net/http"
	"time"

	"github.com/ChanJuiHuang/go-backend-framework/app/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/app/module/user"
	"github.com/ChanJuiHuang/go-backend-framework/app/module/user/model"
	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/ChanJuiHuang/go-backend-framework/util"
	"github.com/gin-gonic/gin"
)

type TokenRefreshRequest struct {
	AccessToken  string `json:"access_token" binding:"required,jwt"`
	RefreshToken string `json:"refresh_token" binding:"required,jwt"`
	Device       string `json:"device" binding:"required,oneof=web ios android"`
}

type TokenRefreshResponse struct {
	AccessToken  string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// @tags user
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true "csrf token"
// @param request body user.TokenRefreshRequest true "refresh token"
// @success 200 {object} user.TokenRefreshResponse
// @failure 400 {object} response.SwaggerErrorResponse
// @failure 401 {object} response.SwaggerErrorResponse
// @failure 503 {object} response.SwaggerErrorResponse
// @router /api/token [put]
func RefreshToken(c *gin.Context) {
	requestBody := new(TokenRefreshRequest)
	if err := c.ShouldBindJSON(requestBody); err != nil {
		response.NewHttpErrorResponse(response.ErrRequestValidationFailed).
			MakePreviousMessage(util.WrapError(err)).
			ToJson(c)
		return
	}

	userId, expireAt, err := user.VerifyRefreshToken(requestBody.RefreshToken)
	if err != nil {
		response.NewHttpErrorResponse(response.ErrJwtAuthenticationFailed).
			MakePreviousMessage(util.WrapError(err)).
			ToJson(c)
		return
	}

	db := provider.App.DB.Create(&model.RefreshTokenRecord{
		RefreshToken: requestBody.RefreshToken,
		UserId:       uint(userId),
		Device:       requestBody.Device,
		ExpireAt:     time.Unix(expireAt, 0),
	})
	if db.Error != nil {
		response.NewHttpErrorResponse(response.ErrJwtAuthenticationFailed).
			MakePreviousMessage(util.WrapError(db.Error)).
			ToJson(c)
		return
	}

	if err := user.RecordAccessToken(requestBody.AccessToken, userId); err != nil {
		response.NewHttpErrorResponse(response.ErrJwtAuthenticationFailed).
			MakePreviousMessage(util.WrapError(err)).
			ToJson(c)
		return
	}

	c.JSON(http.StatusOK, &TokenRefreshResponse{
		AccessToken:  util.IssueAccessToken(userId),
		RefreshToken: util.IssueRefreshToken(userId),
	})
}
