package user

import (
	"fmt"
	"net/http"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/database"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/model"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/user"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/argon2"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/authentication"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type UserRegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=8,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=1234567890"`
}

// @tags user
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true "csrf token"
// @param request body user.UserRegisterRequest true "register user"
// @success 200 {object} response.Response{data=user.TokenData}
// @failure 400 {object} response.ErrorResponse "code: 400-001(issue access token failed), 400-002(request validation failed)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(csrf token mismatch)"
// @failure 500 {object} response.ErrorResponse "code: 500-001"
// @router /api/user/register [post]
func Register(c *gin.Context) {
	logger := service.Registry.Get("logger").(*zap.Logger)
	reqBody := new(UserRegisterRequest)
	if err := c.ShouldBindJSON(reqBody); err != nil {
		errResp := response.NewErrorResponse(response.RequestValidationFailed, errors.WithStack(err), nil)
		logger.Warn(response.RequestValidationFailed, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	u := &model.User{
		Name:     reqBody.Name,
		Email:    reqBody.Email,
		Password: argon2.MakeArgon2IdHash(reqBody.Password),
	}
	if err := user.Create(database.NewTx("users"), u); err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, err, nil)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	authenticator := service.Registry.Get("authentication.authenticator").(*authentication.Authenticator)
	accessToken, err := authenticator.IssueAccessToken(fmt.Sprintf("%v", u.Id))
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
