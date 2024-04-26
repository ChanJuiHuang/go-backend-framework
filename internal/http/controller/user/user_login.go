package user

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/database"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/user"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/argon2"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/authentication"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// @tags user
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true "csrf token"
// @param request body user.UserLoginRequest true "login user"
// @success 200 {object} response.Response{data=user.TokenData}
// @failure 400 {object} response.ErrorResponse "code: 400-001(issue access token failed), 400-002(request validation failed), 400-003(email is wrong), 400-004(password is wrong)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(csrf token mismatch)"
// @failure 500 {object} response.ErrorResponse "code: 500-001"
// @router /api/user/login [post]
func Login(c *gin.Context) {
	logger := service.Registry.Get("logger").(*zap.Logger)
	reqBody := new(UserLoginRequest)
	if err := c.ShouldBindJSON(reqBody); err != nil {
		errResp := response.NewErrorResponse(response.RequestValidationFailed, errors.WithStack(err), nil)
		logger.Warn(response.RequestValidationFailed, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	u, err := user.Get(database.NewTx("users"), "email = ?", reqBody.Email)
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

	encoder := schema.NewEncoder()
	values := url.Values{}
	userQuery := user.Query{
		UserId: u.Id,
	}
	if err := encoder.Encode(userQuery, values); err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	authenticator := service.Registry.Get("authentication.authenticator").(*authentication.Authenticator)
	accessToken, err := authenticator.IssueAccessToken(values.Encode())
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	data := TokenData{}
	data.Fill(fmt.Sprintf("Bearer %s", accessToken))
	respBody := response.NewResponse(data)
	c.JSON(http.StatusOK, respBody)
}
