package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guuzaa/email-newsletter/internal/api/middleware"
	"github.com/guuzaa/email-newsletter/web"
)

// Deprecated: use vue3 frontend instead
func home(c *gin.Context) {
	log := middleware.GetContextLogger(c)
	log.Trace().Msg("home page")
	c.Data(http.StatusOK, "text/html; charset=utf-8", web.HomeHTML)
}
