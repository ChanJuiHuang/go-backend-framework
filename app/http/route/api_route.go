package route

import (
	"net/http"

	"github.com/ChanJuiHuang/go-backend-framework/app/http/controller/user"
	"github.com/ChanJuiHuang/go-backend-framework/app/http/middleware"
	"github.com/gin-gonic/gin"
)

// @produce json
// @success 200 {string} string
// @router /api/ping [get]
func AttachApiRoutes(router *gin.Engine) {
	apiRoutes := router.Group("/api")
	{
		apiRoutes.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		userRoutes := apiRoutes.Group("/user")
		{
			userRoutes.POST("", user.Register)
			userRoutes.GET("", user.Search)
			userRoutes.PUT("",
				middleware.Authenticate(),
				user.Update,
			)
		}

		tokenRoutes := apiRoutes.Group("/token")
		{
			tokenRoutes.POST("", user.EmailLogin)
			tokenRoutes.PUT("", user.RefreshToken)
			tokenRoutes.DELETE("",
				middleware.Authenticate(),
				user.Logout,
			)
		}

		oauthRoutes := apiRoutes.Group("/oauth/:provider")
		{
			oauthRoutes.POST("/token", user.OauthLogin)
		}
	}
}
