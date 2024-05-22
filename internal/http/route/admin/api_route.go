package admin

import (
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/controller/admin"
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
	r.AttachHttpApiRoutes()
	r.AttachPermissionRoutes()
	r.AttachPolicyRoutes()
}

func (r *Router) AttachHttpApiRoutes() {
	httpApiRouter := r.router.Group("http-api")
	httpApiRouter.GET("", httpapi.Search)
}

func (r *Router) AttachPermissionRoutes() {
	permissionRouter := r.router.Group("permission")
	permissionRouter.POST("", permission.Create)
	permissionRouter.GET("", permission.Search)
	permissionRouter.GET(":id", permission.Get)
	permissionRouter.PUT(":id", permission.Update)
	permissionRouter.DELETE("", permission.Delete)

	roleRouter := r.router.Group("role")
	roleRouter.POST("", permission.CreateRole)
	roleRouter.GET("", permission.SearchRoles)
	roleRouter.PUT(":id", permission.UpdateRole)
	roleRouter.DELETE("", permission.DeleteRoles)
}

func (r *Router) AttachPolicyRoutes() {
	policyRouter := r.router.Group("policy")
	policyRouter.POST("reload", admin.ReloadPolicy)
}
