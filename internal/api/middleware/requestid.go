package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/guuzaa/email-newsletter/internal"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.NewString()
		c.Set("requestID", requestID)
		ctx := context.WithValue(c.Request.Context(), "requestID", requestID)
		ctx = internal.Logger().With().Str("ID", requestID).Logger().WithContext(ctx)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
