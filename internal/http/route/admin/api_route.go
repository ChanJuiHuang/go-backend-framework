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
		router: router.Group(
			"api/admin",
			middleware.Authenticate(),
			middleware.Authorize(),
		),
	}
}

func (r *Router) AttachRoutes() {
	r.router.GET("user/:userId/grouping-policy", admin.GetUserGroupingPolicy)
	r.AttachPolicyRoutes()
	r.AttachPolicySubjectRoutes()
	r.AttachGroupingPolicyRoutes()
}

func (r *Router) AttachPolicyRoutes() {
	policyRouter := r.router.Group("policy")
	policyRouter.POST("", admin.CreatePolicy)
	policyRouter.DELETE("", admin.DeletePolicy)
	policyRouter.POST("reload", admin.ReloadPolicy)
}

func (r *Router) AttachPolicySubjectRoutes() {
	policySubjectRouter := r.router.Group("policy/subject")
	policySubjectRouter.GET("", admin.SearchPolicySubject)
	policySubjectRouter.GET(":subject", admin.GetPolicySubject)
	policySubjectRouter.GET(":subject/user", admin.GetPolicySubjectUser)
	policySubjectRouter.DELETE("", admin.DeletePolicySubject)
}

func (r *Router) AttachGroupingPolicyRoutes() {
	groupingPolicyRouter := r.router.Group("grouping-policy")
	groupingPolicyRouter.POST("", admin.CreateGroupingPolicy)
	groupingPolicyRouter.DELETE("", admin.DeleteGroupingPolicy)
}
