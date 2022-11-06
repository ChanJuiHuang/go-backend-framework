package middleware

import (
	"net/http"

	"github.com/ChanJuiHuang/go-backend-framework/app/http/response"
	"github.com/ChanJuiHuang/go-backend-framework/util"
	"github.com/gin-gonic/gin"
)

func VerifyCsrfToken() gin.HandlerFunc {
	skipPaths := map[string]bool{
		"/skip-path":                      true,
		"/scheduler/welcome":              true,
		"/scheduler/refresh-token-record": true,
	}

	return func(c *gin.Context) {
		setCsrfToken(c)
		if isReadingHttpMethod(c) ||
			skipPaths[c.Request.URL.Path] ||
			verifyCsrfToken(c) {
			c.Next()
			return
		}

		response.NewHttpErrorResponse(response.ErrCsrfTokenMismatch).
			AbortWithJson(c)
	}
}

func setCsrfToken(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "XSRF-TOKEN",
		Value:    util.RandomString(20),
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
	})
}

func isReadingHttpMethod(c *gin.Context) bool {
	methods := map[string]bool{
		http.MethodGet:     true,
		http.MethodHead:    true,
		http.MethodOptions: true,
	}
	return methods[c.Request.Method]
}

func verifyCsrfToken(c *gin.Context) bool {
	csrfCookie, _ := c.Cookie("XSRF-TOKEN")
	csrfHeader := c.GetHeader("X-XSRF-TOKEN")
	if csrfCookie == csrfHeader &&
		csrfCookie != "" &&
		csrfHeader != "" {

		return true
	}

	return false
}
