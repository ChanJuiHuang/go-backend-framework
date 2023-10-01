package user

import (
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/controller/user"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

func AttachApiRoute(router *gin.Engine) {
	userRouter := router.Group("/api/user")
	{
		userRouter.POST("/register", user.Register)
		userRouter.POST("/login", user.Login)
		userRouter.GET("/me", middleware.Authenticate(), user.Me)
		userRouter.PUT("", middleware.Authenticate(), user.Update)
	}
}
