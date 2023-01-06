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

type UserCreateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=8,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=1234567890"`
}

func (r *UserCreateRequest) transform() (*model.User, *model.EmailUser) {
	user := &model.User{
		Name: r.Name,
	}
	emailUser := &model.EmailUser{
		Email:    r.Email,
		Password: util.MakeArgon2IdHash(r.Password),
	}

	return user, emailUser
}

type UserRegisterResponse struct {
	AccessToken  string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// @tags user
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true "csrf token"
// @param request body user.UserCreateRequest true "register user"
// @success 200  {object}  user.UserRegisterResponse
// @failure 400 {object} response.SwaggerErrorResponse
// @failure 503 {object} response.SwaggerErrorResponse
// @router /api/user [post]
func Register(c *gin.Context) {
	requestBody := new(UserCreateRequest)
	if err := c.ShouldBindJSON(requestBody); err != nil {
		response.NewHttpErrorResponse(response.ErrRequestValidationFailed).
			MakePreviousMessage(util.WrapError(err)).
			ToJson(c)
		return
	}
	user, emailUser := requestBody.transform()
	err := provider.App.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return util.WrapError(err)
		}

		emailUser.UserId = user.Id
		if err := tx.Create(emailUser).Error; err != nil {
			return util.WrapError(err)
		}

		return nil
	})

	if err != nil {
		response.NewHttpErrorResponse(response.ErrBadRequest).
			MakePreviousMessage(util.WrapError(err)).
			ToJson(c)
		return
	}

	// Todo: send register email

	c.JSON(http.StatusOK, &UserRegisterResponse{
		AccessToken:  util.IssueAccessToken(user.Id),
		RefreshToken: util.IssueRefreshToken(user.Id),
	})
}
