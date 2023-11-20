package user

import (
	"fmt"
	"net/http"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/user"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/argon2"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/authentication"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/provider"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	AccessToken  string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// @tags user
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true "csrf token"
// @param request body user.UserLoginRequest true "login user"
// @success 200 {object} user.UserLoginResponse
// @failure 400 {object} response.ErrorResponse "code: 400-001(issue access token failed, issue refresh token failed), 400-002(request validation failed), 400-003(email is wrong), 400-004(password is wrong)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(csrf token mismatch)"
// @failure 500 {object} response.ErrorResponse "code: 500-001"
// @router /api/user/login [post]
func Login(c *gin.Context) {
	logger := provider.Registry.Get("logger").(*zap.Logger)
	reqBody := new(UserLoginRequest)
	if err := c.ShouldBindJSON(reqBody); err != nil {
		errResp := response.NewErrorResponse(response.RequestValidationFailed, errors.WithStack(err), nil)
		logger.Warn(response.RequestValidationFailed, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	u, err := user.Get("email = ?", reqBody.Email)
	if err != nil {
		errResp := response.NewErrorResponse(response.EmailIsWrong, err, nil)
		logger.Warn(response.EmailIsWrong, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}
	if !argon2.VerifyArgon2IdHash(reqBody.Password, u.Password) {
		errResp := response.NewErrorResponse(response.PasswordIsWrong, errors.New(response.PasswordIsWrong), nil)
		logger.Warn(response.PasswordIsWrong, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	authenticator := provider.Registry.Get("authenticator").(*authentication.Authenticator)
	accessToken, err := authenticator.IssueAccessToken(fmt.Sprintf("%v", u.Id))
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}
	refreshToken, err := authenticator.IssueRefreshToken(fmt.Sprintf("%v", u.Id))
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	c.JSON(http.StatusOK, &UserLoginResponse{
		AccessToken:  fmt.Sprintf("Bearer %s", accessToken),
		RefreshToken: refreshToken,
	})
}
