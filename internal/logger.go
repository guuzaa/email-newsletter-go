package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var (
	once sync.Once
	log  zerolog.Logger
)

func Logger() zerolog.Logger {
	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano
		logLevel, err := zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
		if err != nil {
			logLevel = zerolog.TraceLevel // default to TRACE
			if gin.Mode() == gin.ReleaseMode {
				logLevel = zerolog.WarnLevel
			}
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
	requestID := GetRequestID(c.Request.Context())
	return logger.With().Str("ID", requestID).Logger()
}

func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value("requestID").(string); ok {
		return requestID
	}
	return "unknown"
}

// Adopted from https://github.com/go-mods/zerolog-gorm, licensed under MIT License
type GormLogger struct {
	// SlowThreshold is the delay which define the query as slow
	SlowThreshold time.Duration

	// IgnoreRecordNotFoundError is to ignore when the record is not found
	IgnoreRecordNotFoundError bool

	// FieldsOrder defines the order of fields in output.
	FieldsOrder []string

	// FieldsExclude defines contextual fields to not display in output.
	FieldsExclude []string
}

var (
	// TimestampFieldName is the field name used for the timestamp field.
	TimestampFieldName = zerolog.TimestampFieldName

	// DurationFieldName is the field name used for the duration field.
	DurationFieldName = "elapsed"

	// FileFieldName is the field name used for the file field.
	FileFieldName = "file"

	// SqlFieldName is the field name used for the sql field.
	SqlFieldName = "sql"

	// RowsFieldName is the field name used for the rows field.
	RowsFieldName = "rows"
)

// GormLogger implements the logger.Interface
var _ logger.Interface = &GormLogger{}

// NewGormLogger creates and initializes a new GormLogger.
func NewGormLogger() *GormLogger {
	l := &GormLogger{
		FieldsOrder: gormDefaultFieldsOrder(),
	}

	return l
}

// gormDefaultFieldsOrder defines the default order of fields
func gormDefaultFieldsOrder() []string {
	return []string{
		TimestampFieldName,
		DurationFieldName,
		FileFieldName,
		SqlFieldName,
		RowsFieldName,
	}
}

// isExcluded check if a field is excluded from the output
func (l GormLogger) isExcluded(field string) bool {
	if l.FieldsExclude == nil {
		return false
	}
	for _, f := range l.FieldsExclude {
		if f == field {
			return true
		}
	}

	return false
}

// LogMode log mode
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

// Info print info
func (l GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	requestID := GetRequestID(ctx)
	zerolog.Ctx(ctx).Info().Str("ID", requestID).Msg(fmt.Sprintf(msg, data...))
}

// Warn print warn messages
func (l GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	requestID := GetRequestID(ctx)
	zerolog.Ctx(ctx).Warn().Str("ID", requestID).Msg(fmt.Sprintf(msg, data...))
}

// Error print error messages
func (l GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	requestID := GetRequestID(ctx)
	zerolog.Ctx(ctx).Error().Str("ID", requestID).Msg(fmt.Sprintf(msg, data...))
}

// Trace print sql message
func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {

	// get zerolog from context
	z := zerolog.Ctx(ctx)

	// return if zerolog is disabled
	if z.GetLevel() == zerolog.Disabled {
		return
	}

	if l.FieldsOrder == nil {
		l.FieldsOrder = gormDefaultFieldsOrder()
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	var event *zerolog.Event
	var eventError bool
	var eventWarn bool

	// set message level
	if err != nil && !(l.IgnoreRecordNotFoundError && errors.Is(err, gorm.ErrRecordNotFound)) {
		eventError = true
		event = z.Error()
	} else if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		eventWarn = true
		event = z.Warn()
	} else {
		event = z.Trace()
	}

	// add fields
	for _, f := range l.FieldsOrder {
		// add time field
		if f == TimestampFieldName && !l.isExcluded(f) {
			event.Time(TimestampFieldName, begin)
		}

		// add duration field
		if f == DurationFieldName && !l.isExcluded(f) {
			var durationFieldName string
			switch zerolog.DurationFieldUnit {
			case time.Nanosecond:
				durationFieldName = DurationFieldName + "_ns"
			case time.Microsecond:
				durationFieldName = DurationFieldName + "_us"
			case time.Millisecond:
				durationFieldName = DurationFieldName + "_ms"
			case time.Second:
				durationFieldName = DurationFieldName
			case time.Minute:
				durationFieldName = DurationFieldName + "_min"
			case time.Hour:
				durationFieldName = DurationFieldName + "_hr"
			default:
				z.Error().Interface("zerolog.DurationFieldUnit", zerolog.DurationFieldUnit).Msg("unknown value for DurationFieldUnit")
				durationFieldName = DurationFieldName
			}
			event.Dur(durationFieldName, elapsed)
		}

		// add file field
		if f == FileFieldName && !l.isExcluded(f) {
			event.Str("file", utils.FileWithLineNum())
		}

		// add sql field
		if f == SqlFieldName && !l.isExcluded(f) {
			if sql != "" {
				event.Str("sql", sql)
			}
		}

		// add rows field
		if f == RowsFieldName && !l.isExcluded(f) {
			if rows > -1 {
				event.Int64("rows", rows)
			}
		}
	}

	// post the message
	if eventError {
		event.Msg("SQL error")
	} else if eventWarn {
		event.Msg("SQL slow query")
	} else {
		event.Msg("SQL")
	}
}
