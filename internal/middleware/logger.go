package middleware

import (
	"io"
	"os"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var (
	once sync.Once
	log  zerolog.Logger
)

func UseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := Logger()
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
			Msgf("ID: %s, request", requestID)
	}
}

func Logger() zerolog.Logger {
	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano
		logLevel, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
		if err != nil {
			logLevel = int(zerolog.TraceLevel) // default to TRACE
		}

		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}

		var gitRevision string
		buildInfo, ok := debug.ReadBuildInfo()
		if ok {
			for _, v := range buildInfo.Settings {
				if v.Key == "vcs.revision" {
					gitRevision = v.Value
					break
				}
			}
		}

		log = zerolog.New(output).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			Str("gitRevision", gitRevision).
			Str("goVersion", buildInfo.GoVersion).
			Logger()
	})
	return log
}
