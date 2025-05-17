package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/guuzaa/email-newsletter/internal"
	"github.com/guuzaa/email-newsletter/internal/middleware"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, emailClient *internal.EmailClient, baseURL string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.UseLogger())
	r.Use(middleware.RequestID())

	r.GET("/health_check", healthCheck)
	confirmSubscriptionHandler := NewConfirmSubscriptionHandler(db)
	r.GET("/subscriptions/confirm", confirmSubscriptionHandler.confirm)

	subscriptionHandler := NewSubscriptionHandler(db, emailClient, baseURL)
	r.POST("/subscriptions", subscriptionHandler.subscribe)

	return r
}
