package test

import (
	"net/http"

	pkgHttp "github.com/ChanJuiHuang/go-backend-framework/internal/http"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/middleware"
	"github.com/ChanJuiHuang/go-backend-framework/internal/http/route"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/config"
	"github.com/gin-gonic/gin"
)

type httpHandler struct {
	engine *gin.Engine
}

var HttpHandler *httpHandler

func NewHttpHandler() *httpHandler {
	handler := &httpHandler{
		engine: pkgHttp.NewEngine(),
	}
	handler.AttachGlobalMiddleware()
	route.AttachApiRoutes(handler.engine)
	route.AttachSwaggerRoute(handler.engine)

	return handler
}

func (handler *httpHandler) AttachGlobalMiddleware() {
	csrfConfig := config.Registry.Get("middleware.csrf").(middleware.CsrfConfig)

	handlerFunctions := []gin.HandlerFunc{
		middleware.VerifyCsrfToken(csrfConfig),
	}

	handler.engine.Use(
		handlerFunctions...,
	)
}

func (handler *httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler.engine.ServeHTTP(w, req)
}

func AddCsrfToken(req *http.Request) {
	config := config.Registry.Get("middleware.csrf").(middleware.CsrfConfig)
	cookie := &http.Cookie{
		Name:     config.Cookie.Name,
		Value:    "1234567890",
		Path:     config.Cookie.Path,
		Domain:   config.Cookie.Domain,
		MaxAge:   config.Cookie.MaxAge,
		Secure:   config.Cookie.Secure,
		HttpOnly: config.Cookie.HttpOnly,
		SameSite: config.Cookie.SameSite,
	}
	req.AddCookie(cookie)
	req.Header.Set(config.Header, "1234567890")
}

func AddBearerToken(req *http.Request, token string) {
	req.Header.Set("Authorization", token)
}
