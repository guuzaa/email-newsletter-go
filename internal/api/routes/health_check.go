package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guuzaa/email-newsletter/internal/api/middleware"
)

func healthCheck(c *gin.Context) {
	log := middleware.GetContextLogger(c)
	log.Debug().Msg("health check")

	c.String(http.StatusOK, "")
}
