package internal

import (
	"io"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"runtime/debug"
)

var (
	once sync.Once
	log  zerolog.Logger
)

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

// GetContextLogger returns a logger with request ID from context
func GetContextLogger(c *gin.Context) zerolog.Logger {
	logger := Logger()
	requestID, exists := c.Get("requestID")
	if exists {
		logger = logger.With().Interface("ID", requestID).Logger()
	}
	return logger
}
