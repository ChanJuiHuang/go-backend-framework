package admin

import (
	httpapi "github.com/ChanJuiHuang/go-backend-framework/internal/http/controller/admin/http_api"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/controller/admin/permission"
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
	r.AttachHttpApiRoutes()
	r.AttachPermissionRoutes()
	r.AttachPolicyRoutes()
	r.AttachPolicySubjectRoutes()
	r.AttachGroupingPolicyRoutes()
}

func (r *Router) AttachHttpApiRoutes() {
	httpApiRouter := r.router.Group("http-api")
	httpApiRouter.GET("", httpapi.Search)
}

func (r *Router) AttachPermissionRoutes() {
	permissionRouter := r.router.Group("permission")
	permissionRouter.POST("", permission.Create)
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
