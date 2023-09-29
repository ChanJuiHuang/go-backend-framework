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
		adminRouter.DELETE(
			"/policy",
			middleware.Authenticate(),
			middleware.Authorize(),
			admin.DeletePolicy,
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
		adminRouter.DELETE(
			"/policy/subject",
			middleware.Authenticate(),
			middleware.Authorize(),
			admin.DeletePolicySubject,
		)
	}
	{
		adminRouter.POST(
			"/grouping-policy",
			middleware.Authenticate(),
			middleware.Authorize(),
			admin.CreateGroupingPolicy,
		)
		adminRouter.GET(
			"/grouping-policy/:userId",
			middleware.Authenticate(),
			middleware.Authorize(),
			admin.GetGroupingPolicy,
		)
		adminRouter.DELETE(
			"/grouping-policy",
			middleware.Authenticate(),
			middleware.Authorize(),
			admin.DeleteGroupingPolicy,
		)
	}
}
