package user

import (
	"net/http"

	"github.com/ChanJuiHuang/go-backend-framework/v2/internal/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/v2/internal/pkg/database"
	"github.com/ChanJuiHuang/go-backend-framework/v2/internal/pkg/user"
	"github.com/ChanJuiHuang/go-backend-framework/v2/pkg/booter/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @tags user
// @accept json
// @produce json
// @param Authorization header string true "bearer token"
// @success 200 {object} response.Response{data=UserData}
// @failure 400 {object} response.ErrorResponse "code: 400-001(Bad Request)"
// @failure 401 {object} response.ErrorResponse "code: 401-001(Unauthorized)"
// @failure 500 {object} response.ErrorResponse "code: 500-001(Internal Server Error)"
// @router /api/user/me [get]
func Me(c *gin.Context) {
	u, err := user.Get(database.NewTx("Roles.Permissions"), "id = ?", c.GetUint("user_id"))
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, err, nil)
		logger := service.Registry.Get("logger").(*zap.Logger)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	data := UserData{}
	data.Fill(u)
	respBody := response.NewResponse(data)
	c.JSON(http.StatusOK, respBody)
}
