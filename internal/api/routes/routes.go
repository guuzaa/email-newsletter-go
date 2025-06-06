package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/guuzaa/email-newsletter/internal"
	"github.com/guuzaa/email-newsletter/internal/api/middleware"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, emailClient *internal.EmailClient, baseURL string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.UseLogger())

	r.GET("/", home)

	loginHandler := NewLoginHandler(db)
	r.GET("/login", loginHandler.get)
	r.POST("/login", loginHandler.post)

	r.GET("/health_check", healthCheck)
	confirmSubscriptionHandler := NewConfirmSubscriptionHandler(db)
	r.GET("/subscriptions/confirm", confirmSubscriptionHandler.confirm)

	subscriptionHandler := NewSubscriptionHandler(db, emailClient, baseURL)
	r.POST("/subscriptions", subscriptionHandler.subscribe)

	newslettersHandler := NewNewslettersHandler(db, emailClient)
	r.POST("/newsletters", newslettersHandler.publishNewsletter)

	return r
}
