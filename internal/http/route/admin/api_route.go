package admin

import (
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/controller/admin"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

type AdminRouter struct {
	router *gin.RouterGroup
}

func NewAdminRouter(router *gin.Engine) *AdminRouter {
	return &AdminRouter{
		router: router.Group("api/admin"),
	}
}

func (ar *AdminRouter) AttachRoutes() {
	ar.router.GET(
		"user/:userId/grouping-policy",
		middleware.Authenticate(),
		middleware.Authorize(),
		admin.GetUserGroupingPolicy,
	)
	ar.AttachPolicyRoutes()
	ar.AttachPolicySubjectRoutes()
	ar.AttachGroupingPolicyRoutes()
}

func (ar *AdminRouter) AttachPolicyRoutes() {
	policyRouter := ar.router.Group("policy")
	policyRouter.POST(
		"",
		middleware.Authenticate(),
		middleware.Authorize(),
		admin.CreatePolicy,
	)
	policyRouter.DELETE(
		"",
		middleware.Authenticate(),
		middleware.Authorize(),
		admin.DeletePolicy,
	)
	policyRouter.POST(
		"reload",
		middleware.Authenticate(),
		middleware.Authorize(),
		admin.ReloadPolicy,
	)
}

func (ar *AdminRouter) AttachPolicySubjectRoutes() {
	policySubjectRouter := ar.router.Group("policy/subject")
	policySubjectRouter.GET(
		"",
		middleware.Authenticate(),
		middleware.Authorize(),
		admin.SearchPolicySubject,
	)
	policySubjectRouter.GET(
		":subject",
		middleware.Authenticate(),
		middleware.Authorize(),
		admin.GetPolicySubject,
	)
	policySubjectRouter.GET(
		":subject/user",
		middleware.Authenticate(),
		middleware.Authorize(),
		admin.GetPolicySubjectUser,
	)
	policySubjectRouter.DELETE(
		"",
		middleware.Authenticate(),
		middleware.Authorize(),
		admin.DeletePolicySubject,
	)
}

func (ar *AdminRouter) AttachGroupingPolicyRoutes() {
	groupingPolicyRouter := ar.router.Group("grouping-policy")
	groupingPolicyRouter.POST(
		"",
		middleware.Authenticate(),
		middleware.Authorize(),
		admin.CreateGroupingPolicy,
	)
	groupingPolicyRouter.DELETE(
		"",
		middleware.Authenticate(),
		middleware.Authorize(),
		admin.DeleteGroupingPolicy,
	)
}
