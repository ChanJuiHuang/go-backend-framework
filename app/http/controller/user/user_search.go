package user

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/ChanJuiHuang/go-backend-framework/app/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/app/module/user/model"
	"github.com/ChanJuiHuang/go-backend-framework/app/provider"
	"github.com/ChanJuiHuang/go-backend-framework/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserSearchRequest struct {
	util.PaginationRequest
	Filter *struct {
		Name  *string `form:"name" key:"name"`
		Email *string `form:"email" key:"email" validate:"omitempty,email"`
	} `form:"filter"`
	OrderBy string `form:"order_by" validate:"omitempty,oneof=email"`
}

type UserSearchData struct {
	Id        uint      `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
	EmailUser struct {
		Email     string    `json:"email" validate:"required"`
		UserId    uint      `json:"user_id" validate:"required"`
		CreatedAt time.Time `json:"created_at" validate:"required"`
		UpdatedAt time.Time `json:"updated_at" validate:"required"`
	} `json:"email_user" validate:"required"`
}

type UserSearchResponse struct {
	Data []UserSearchData `json:"data" validate:"required"`
	util.PaginationResponse
}

func MakeUserSearchResponse(paginator util.Paginator, users []model.User) *UserSearchResponse {
	res := &UserSearchResponse{
		Data: make([]UserSearchData, len(users)),
		PaginationResponse: util.PaginationResponse{
			Page:     paginator.Page,
			PerPage:  paginator.PerPage,
			Total:    paginator.Total,
			LastPage: paginator.LastPage,
		},
	}
	for index, user := range users {
		res.Data[index].Id = user.Id
		res.Data[index].Name = user.Name
		res.Data[index].CreatedAt = user.CreatedAt
		res.Data[index].UpdatedAt = user.UpdatedAt
		res.Data[index].EmailUser.Email = user.EmailUser.Email
		res.Data[index].EmailUser.UserId = user.EmailUser.UserId
		res.Data[index].EmailUser.CreatedAt = user.EmailUser.CreatedAt
		res.Data[index].EmailUser.UpdatedAt = user.EmailUser.UpdatedAt
	}

	return res
}

// @tags user
// @accept json
// @produce json
// @param request query user.UserSearchRequest false "search users"
// @param filter.name query string false "name"
// @param filter.email query string false "email"
// @success 200  {object}  user.UserSearchResponse
// @failure 400 {object} response.SwaggerErrorResponse
// @failure 503 {object} response.SwaggerErrorResponse
// @router /api/user [get]
func Search(c *gin.Context) {
	queryStrings := new(UserSearchRequest)
	if err := provider.App.FormDecoder.Decode(queryStrings, c.Request.URL.Query()); err != nil {
		response.NewHttpErrorResponse(response.ErrBadRequest).
			MakePreviousMessage(util.WrapError(err)).
			ToJson(c)
		return
	}

	if err := provider.App.Modifier.Struct(context.Background(), queryStrings); err != nil {
		response.NewHttpErrorResponse(response.ErrBadRequest).
			MakePreviousMessage(util.WrapError(err)).
			ToJson(c)
		return
	}
	if err := provider.App.Validator.Struct(queryStrings); err != nil {
		response.NewHttpErrorResponse(response.ErrRequestValidationFailed).
			MakePreviousMessage(util.WrapError(err)).
			ToJson(c)
		return
	}
	db := provider.App.DB.Table("users").
		Joins("join email_users on users.id = email_users.user_id").
		Preload("EmailUser")
	whereConditions := util.WhereCondition{
		"name": func(db *gorm.DB, value any) {
			db.Where("users.name LIKE ?", fmt.Sprintf("%%%s%%", *value.(*string)))
		},
		"email": func(db *gorm.DB, value any) {
			db.Where("email_users.email = ?", value)
		},
	}
	paginator := util.NewPaginator(
		db,
		whereConditions,
		map[string]string{"email": "email_users.email"},
	)
	users := []model.User{}
	paginator.AddWhereConditions(reflect.ValueOf(queryStrings.Filter)).
		SetPagination(queryStrings.Page, queryStrings.PerPage).
		SetOrderBy(queryStrings.OrderBy).
		CalculateTotalAndLastPage().
		GetData(&users, "users.*")

	c.JSON(http.StatusOK, MakeUserSearchResponse(paginator, users))
}
