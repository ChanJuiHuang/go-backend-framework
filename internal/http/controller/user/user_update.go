package user

import (
	"net/http"
	"time"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/user"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type UserUpdateRequest struct {
	Name  string `json:"name" structs:"name,omitempty"`
	Email string `json:"email" structs:"email,omitempty"`
}

type UserUpdateData struct {
	Id        uint      `json:"id" mapstructure:"id" validate:"required"`
	Name      string    `json:"name" mapstructure:"name" validate:"required"`
	Email     string    `json:"email" mapstructure:"email" validate:"required"`
	CreatedAt time.Time `json:"created_at" mapstructure:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" mapstructure:"updated_at" validate:"required"`
}

// @tags user
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true "csrf token"
// @param Authorization header string true "bearer token"
// @param request body user.UserUpdateRequest true "update user"
// @success 200 {object} response.Response{data=user.UserUpdateData}
// @failure 400 {object} response.ErrorResponse "code: 400-001(update user failed, get user failed), 400-002(request validation failed)"
// @failure 401 {object} response.ErrorResponse "code: 401-001(access token is wrong)"
// @failure 500 {object} response.ErrorResponse "code: 500-001"
// @router /api/user [put]
func Update(c *gin.Context) {
	logger := service.Registry.Get("logger").(*zap.Logger)
	reqBody := new(UserUpdateRequest)
	if err := c.ShouldBindJSON(reqBody); err != nil {
		errResp := response.NewErrorResponse(response.RequestValidationFailed, errors.WithStack(err), nil)
		logger.Warn(response.RequestValidationFailed, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	values := structs.Map(reqBody)
	_, err := user.Update(c.GetUint("user_id"), values)
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, err, nil)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	u, err := user.Get("id = ?", c.GetUint("user_id"))
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, err, nil)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	respBody := response.NewResponse(UserUpdateData{
		Id:        u.Id,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	})
	c.JSON(http.StatusOK, respBody)
}
