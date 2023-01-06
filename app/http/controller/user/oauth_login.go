package user

import (
	"net/http"

	"github.com/ChanJuiHuang/go-backend-framework/app/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/app/module/user"
	"github.com/ChanJuiHuang/go-backend-framework/util"
	"github.com/gin-gonic/gin"
)

type OauthLoginRequest struct {
	Code string `json:"code" binding:"required"`
}

type OauthLoginResponse struct {
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
// @failure 401 {object} response.SwaggerErrorResponse
// @failure 503 {object} response.SwaggerErrorResponse
// @router /api/oauth/{provider}/token [post]
// @param provider path string true "oauth provider" enums(google)
func OauthLogin(c *gin.Context) {
	requestBody := new(OauthLoginRequest)
	if err := c.ShouldBindJSON(requestBody); err != nil {
		response.NewHttpErrorResponse(response.ErrRequestValidationFailed).
			MakePreviousMessage(util.WrapError(err)).
			ToJson(c)
		return
	}

	oauthProvider := util.NewOauthProvider(util.OauthProvider(c.Param("provider")))
	if err := oauthProvider.GetAccessToken(requestBody.Code); err != nil {
		response.NewHttpErrorResponse(response.ErrOauthLoginFailed).
			MakePreviousMessage(util.WrapError(err)).
			ToJson(c)
		return
	}
	if err := oauthProvider.SetIdAndEmail(); err != nil {
		response.NewHttpErrorResponse(response.ErrOauthLoginFailed).
			MakePreviousMessage(util.WrapError(err)).
			ToJson(c)
		return
	}

	userId, err := user.GetUserId(oauthProvider)
	if err != nil {
		response.NewHttpErrorResponse(response.ErrOauthLoginFailed).
			MakePreviousMessage(util.WrapError(err)).
			ToJson(c)
		return
	}
	if userId != 0 {
		sendOauthLoginResponse(c, userId)
		return
	}
	userId, err = user.CreateOauthUser(oauthProvider)
	if err != nil {
		response.NewHttpErrorResponse(response.ErrOauthLoginFailed).
			MakePreviousMessage(util.WrapError(err)).
			ToJson(c)
		return
	}

	sendOauthLoginResponse(c, userId)
}

func sendOauthLoginResponse(c *gin.Context, userId uint) {
	c.JSON(http.StatusOK, &EmailLoginResponse{
		AccessToken:  util.IssueAccessToken(userId),
		RefreshToken: util.IssueRefreshToken(userId),
	})
}
