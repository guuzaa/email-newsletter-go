package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/guuzaa/email-newsletter/web"
	"net/http"
)

func RegisterWebStaticEmbed(e *gin.Engine) {
	routeWebStatic(e, "/", "/index.html", "/assets/*filepath")
}

func routeWebStatic(e *gin.Engine, paths ...string) {
	staticHandler := http.FileServer(web.NewFileSystem())
	handler := func(c *gin.Context) {
		staticHandler.ServeHTTP(c.Writer, c.Request)
	}
	for _, path := range paths {
		e.GET(path, handler)
		e.HEAD(path, handler)
	}
}
