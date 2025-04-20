package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.GET("/health_check", healthCheck)

	subscriptionHandler := NewSubscriptionHandler(db)
	r.POST("/subscriptions", subscriptionHandler.subscribe)

	return r
}
