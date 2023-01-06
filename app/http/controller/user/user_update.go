package user

import (
	"net/http"

	"github.com/ChanJuiHuang/go-backend-framework/app/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/app/module/user/model"
	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/ChanJuiHuang/go-backend-framework/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserUpdateRequest struct {
	Name     string `json:"name"`
	Password string `json:"password" binding:"omitempty,gte=8,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=1234567890"`
}

func (r *UserUpdateRequest) transform() (*model.User, string) {
	user := &model.User{Name: r.Name}
	password := ""
	if r.Password != "" {
		password = util.MakeArgon2IdHash(r.Password)
	}

	return user, password
}

// @tags user
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true "csrf token"
// @param request body user.UserUpdateRequest true "update user"
// @success 204
// @failure 400 {object} response.SwaggerErrorResponse
// @failure 401 {object} response.SwaggerErrorResponse
// @failure 503 {object} response.SwaggerErrorResponse
// @router /api/user [put]
func Update(c *gin.Context) {
	requestBody := new(UserUpdateRequest)
	if err := c.ShouldBindJSON(requestBody); err != nil {
		response.NewHttpErrorResponse(response.ErrRequestValidationFailed).
			MakePreviousMessage(util.WrapError(err)).
			ToJson(c)
		return
	}
	user, password := requestBody.transform()

	err := provider.App.DB.Transaction(func(tx *gorm.DB) error {
		userId := c.MustGet("user_id").(uint)
		err := tx.Where("id = ?", userId).
			Updates(user).
			Error
		if err != nil {
			return util.WrapError(err)
		}

		if password != "" {
			err = tx.Table("email_users").
				Where("user_id = ?", userId).
				Update("password", password).
				Error
			if err != nil {
				return util.WrapError(err)
			}
		}

		return nil
	})
	if err != nil {
		response.NewHttpErrorResponse(response.ErrBadRequest).
			MakePreviousMessage(util.WrapError(err)).
			ToJson(c)
		return
	}

	c.Status(http.StatusNoContent)
}
