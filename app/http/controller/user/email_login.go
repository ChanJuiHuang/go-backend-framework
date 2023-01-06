package user

import (
	"net/http"

	"github.com/ChanJuiHuang/go-backend-framework/app/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/app/module/user/model"
	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/ChanJuiHuang/go-backend-framework/util"
	"github.com/gin-gonic/gin"
)

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type EmailLoginResponse struct {
	AccessToken  string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// @tags user
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true "csrf token"
// @param request body user.UserLoginRequest true "email login"
// @success 200 {object} user.EmailLoginResponse
// @failure 400 {object} response.SwaggerErrorResponse
// @failure 404 {object} response.SwaggerErrorResponse
// @failure 503 {object} response.SwaggerErrorResponse
// @router /api/token [post]
func EmailLogin(c *gin.Context) {
	requestBody := new(UserLoginRequest)
	if err := c.ShouldBindJSON(requestBody); err != nil {
		response.NewHttpErrorResponse(response.ErrRequestValidationFailed).
			MakePreviousMessage(util.WrapError(err)).
			ToJson(c)
		return
	}

	emailUser := &model.EmailUser{}
	db := provider.App.DB.Where("email = ?", requestBody.Email).First(emailUser)
	if util.IsDatabaseError(db) {
		response.NewHttpErrorResponse(response.ErrLoginEmailIsWrong).
			MakePreviousMessage(util.WrapError(db.Error)).
			ToJson(c)
		return
	}
	if db.RowsAffected != 1 {
		response.NewHttpErrorResponse(response.ErrRecordNotFound).
			ToJson(c)
		return
	}

	if !util.VerifyArgon2IdHash(requestBody.Password, emailUser.Password) {
		response.NewHttpErrorResponse(response.ErrLoginPasswordIsWrong).
			ToJson(c)
		return
	}

	c.JSON(http.StatusOK, &EmailLoginResponse{
		AccessToken:  util.IssueAccessToken(emailUser.UserId),
		RefreshToken: util.IssueRefreshToken(emailUser.UserId),
	})
}
