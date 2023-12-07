package user

import (
	"net/http"
	"time"

	"github.com/ChanJuiHuang/go-backend-framework/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/internal/pkg/user"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserMeData struct {
	Id        uint      `json:"id" mapstructure:"id" validate:"required"`
	Name      string    `json:"name" mapstructure:"name" validate:"required"`
	Email     string    `json:"email" mapstructure:"email" validate:"required"`
	CreatedAt time.Time `json:"created_at" mapstructure:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" mapstructure:"updated_at" validate:"required"`
}

// @tags user
// @accept json
// @produce json
// @param Authorization header string true "bearer token"
// @success 200 {object} response.Response{data=UserMeData}
// @failure 400 {object} response.ErrorResponse "code: 400-001(get user failed)"
// @failure 401 {object} response.ErrorResponse "code: 401-001(access token is wrong)"
// @failure 500 {object} response.ErrorResponse "code: 500-001"
// @router /api/user/me [get]
func Me(c *gin.Context) {
	u, err := user.Get("id = ?", c.GetUint("user_id"))
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, err, nil)
		logger := service.Registry.Get("logger").(*zap.Logger)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	respBody := response.NewResponse(UserMeData{
		Id:        u.Id,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	})
	c.JSON(http.StatusOK, respBody)
}
