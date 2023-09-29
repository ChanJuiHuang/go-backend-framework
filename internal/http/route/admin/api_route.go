package admin

import (
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/controller/admin"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

func AttachApiRoute(router *gin.Engine) {
	adminRouter := router.Group("/api/admin")
	{
		adminRouter.POST(
			"/policy",
			middleware.Authenticate(),
			middleware.Authorize(),
			admin.CreatePolicy,
		)
		adminRouter.GET(
			"/policy/subject",
			middleware.Authenticate(),
			middleware.Authorize(),
			admin.SearchPolicySubject,
		)
		adminRouter.GET(
			"/policy/subject/:subject",
			middleware.Authenticate(),
			middleware.Authorize(),
			admin.GetPolicySubject,
		)
	}
}
