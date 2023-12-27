package admin

import (
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/controller/admin"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	router *gin.RouterGroup
}

func NewRouter(router *gin.Engine) *Router {
	return &Router{
		router: router.Group("api/admin"),
	}
}

func (r *Router) AttachRoutes() {
	r.router.GET(
		"user/:userId/grouping-policy",
		middleware.Authenticate(),
		middleware.Authorize(),
		admin.GetUserGroupingPolicy,
	)
	r.AttachPolicyRoutes()
	r.AttachPolicySubjectRoutes()
	r.AttachGroupingPolicyRoutes()
}

func (r *Router) AttachPolicyRoutes() {
	policyRouter := r.router.Group("policy")
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

func (r *Router) AttachPolicySubjectRoutes() {
	policySubjectRouter := r.router.Group("policy/subject")
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

func (r *Router) AttachGroupingPolicyRoutes() {
	groupingPolicyRouter := r.router.Group("grouping-policy")
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
