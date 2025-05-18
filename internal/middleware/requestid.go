package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.NewString()
		ctx := context.WithValue(c.Request.Context(), "requestID", requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func getRequestID(c *gin.Context) string {
	requestID := c.Request.Context().Value("requestID")
	if requestID, ok := requestID.(string); ok {
		return requestID
	}
	return uuid.NewString()
}
