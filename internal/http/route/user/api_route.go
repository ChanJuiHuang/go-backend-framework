package user

import (
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/controller/user"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	router *gin.RouterGroup
}

func NewUserRouter(router *gin.Engine) *UserRouter {
	return &UserRouter{
		router: router.Group("api/user"),
	}
}

func (ur *UserRouter) AttachRoutes() {
	ur.router.POST("register", user.Register)
	ur.router.POST("login", user.Login)
	ur.router.GET("me", middleware.Authenticate(), user.Me)
	ur.router.PUT("", middleware.Authenticate(), user.Update)
	ur.router.PUT("password", middleware.Authenticate(), user.UpdatePassword)
}
