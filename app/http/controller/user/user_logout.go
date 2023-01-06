package user

import (
	"net/http"
	"strings"
	"time"

	"github.com/ChanJuiHuang/go-backend-framework/app/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/app/module/user"
	"github.com/ChanJuiHuang/go-backend-framework/app/module/user/model"
	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/ChanJuiHuang/go-backend-framework/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type UserLogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required,jwt"`
	Device       string `json:"device" binding:"required,oneof=web ios android"`
}

// @tags user
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true "csrf token"
// @param request body user.UserLogoutRequest true "logout"
// @success 204
// @failure 400 {object} response.SwaggerErrorResponse
// @failure 401 {object} response.SwaggerErrorResponse
// @failure 503 {object} response.SwaggerErrorResponse
// @router /api/token [delete]
func Logout(c *gin.Context) {
	requestBody := new(UserLogoutRequest)
	if err := c.ShouldBindJSON(requestBody); err != nil {
		response.NewHttpErrorResponse(response.ErrRequestValidationFailed).
			MakePreviousMessage(util.WrapError(err)).
			ToJson(c)
		return
	}

	refreshTokenString := requestBody.RefreshToken
	refreshToken, _, err := jwt.NewParser().ParseUnverified(refreshTokenString, jwt.MapClaims{})
	if err != nil {
		response.NewHttpErrorResponse(response.ErrJwtAuthenticationFailed).
			MakePreviousMessage(util.WrapError(err)).
			ToJson(c)
		return
	}

	userId := c.MustGet("user_id").(uint)
	db := provider.App.DB.Create(&model.RefreshTokenRecord{
		RefreshToken: requestBody.RefreshToken,
		UserId:       userId,
		Device:       requestBody.Device,
		ExpireAt:     time.Unix(int64(refreshToken.Claims.(jwt.MapClaims)["exp"].(float64)), 0),
	})
	if db.Error != nil {
		response.NewHttpErrorResponse(response.ErrJwtAuthenticationFailed).
			MakePreviousMessage(util.WrapError(db.Error)).
			ToJson(c)
		return
	}

	accessTokenString := strings.Replace(c.GetHeader("Authorization"), "Bearer ", "", 1)
	if err := user.RecordAccessToken(accessTokenString, userId); err != nil {
		response.NewHttpErrorResponse(response.ErrJwtAuthenticationFailed).
			MakePreviousMessage(util.WrapError(err)).
			ToJson(c)
		return
	}

	c.Status(http.StatusNoContent)
}
