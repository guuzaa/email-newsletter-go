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
		traceID := getRequestID(c)
		logger.Trace().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("ip", c.ClientIP()).
			Str("userAgent", c.Request.UserAgent()).
			Str("latency", time.Since(t).String()).
			Str("statusCode", strconv.Itoa(c.Writer.Status())).
			Str("ID", traceID)
		c.Next()
	}
}

// GetContextLogger returns a logger with request ID from context
func GetContextLogger(c *gin.Context) zerolog.Logger {
	logger := internal.Logger()
	requestID := getRequestID(c)
	return logger.With().Interface("ID", requestID).Logger()
}
