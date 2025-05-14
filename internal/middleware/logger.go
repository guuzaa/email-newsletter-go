package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/guuzaa/email-newsletter/internal"
	"github.com/rs/zerolog"
)

func UseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := internal.Logger()
		t := time.Now()
		c.Next()
		requestID := c.Value("requestID")
		if requestID == nil {
			requestID = "unknown"
		}
		logger.Trace().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("ip", c.ClientIP()).
			Str("userAgent", c.Request.UserAgent()).
			Str("latency", time.Since(t).String()).
			Str("statusCode", strconv.Itoa(c.Writer.Status())).
			Msgf("ID= %s", requestID)
	}
}

// GetContextLogger returns a logger with request ID from context
func GetContextLogger(c *gin.Context) zerolog.Logger {
	logger := internal.Logger()
	if requestID, exists := c.Get("requestID"); exists {
		logger = logger.With().Interface("ID", requestID).Logger()
	}
	return logger
}
