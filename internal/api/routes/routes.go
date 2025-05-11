package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/guuzaa/email-newsletter/internal"
	"github.com/guuzaa/email-newsletter/internal/middleware"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, emailClient *internal.EmailClient) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.UseLogger())
	r.Use(middleware.RequestID())

	r.GET("/health_check", healthCheck)

	subscriptionHandler := NewSubscriptionHandler(db)
	r.POST("/subscriptions", subscriptionHandler.subscribe)

	return r
}
