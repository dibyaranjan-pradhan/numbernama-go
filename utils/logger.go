package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

// Logger is the single app logger.
// Use Log(ctx) for request-scoped logging (adds traceID + caller).
// Use GetLogger() only when no context is available.
type Logger struct {
	zl *zerolog.Logger
}

var globalLogger *Logger

// InitLogger initializes the global logger. Called by wire; do not call from main.
func InitLogger(level string, pretty bool) {
	var logLevel zerolog.Level
	switch level {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	default:
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	var zl zerolog.Logger
	if pretty {
		// NOTE: LOG_PRETTY=true outputs plain text — do NOT use on EC2.
		// Promtail expects JSON. Set LOG_PRETTY=false in all EC2 environments.
		zl = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().CallerWithSkipFrameCount(4).Logger()
	} else {
		// Production default: structured JSON with caller (filename:line) on every log line.
		zl = zerolog.New(os.Stdout).With().Timestamp().CallerWithSkipFrameCount(4).Logger()
	}
	globalLogger = &Logger{zl: &zl}
}

// GetLogger returns the global logger. Safe after wire has run.
// Prefer Log(ctx) over GetLogger() wherever a context is available.
func GetLogger() *Logger {
	if globalLogger == nil {
		InitLogger("info", false)
	}
	return globalLogger
}

// Log is a package-level shortcut for request-scoped logging.
// Replaces the verbose utils.Log(ctx) pattern.
//
// Usage:
//
//	utils.Log(ctx).Infof("Service foo: enter")
//	utils.Log(ctx).WithError(err).Error("something failed")
func Log(ctx context.Context) *Logger {
	return GetLogger().WithContext(ctx)
}

// WithContext returns a Logger that adds traceID from ctx to every log line.
func (l *Logger) WithContext(ctx context.Context) *Logger {
	if l == nil || l.zl == nil {
		return l
	}
	if ctx == nil {
		return l
	}
	traceID, _ := GetTraceID(ctx)
	if traceID == "" {
		return l
	}
	sub := l.zl.With().Str(string("traceID"), traceID).Logger()
	return &Logger{zl: &sub}
}

// WithField returns a new Logger with an additional structured key-value field.
// Use this to attach domain-specific context to a log line.
//
// Usage:
//
//	utils.Log(ctx).WithField("userID", userID).Infof("user fetched")
func (l *Logger) WithField(key string, val any) *Logger {
	if l == nil || l.zl == nil {
		return l
	}
	sub := l.zl.With().Interface(key, val).Logger()
	return &Logger{zl: &sub}
}

// WithError returns a new Logger with a structured "error" field.
// Use this instead of embedding error strings in the message.
//
// Usage:
//
//	utils.Log(ctx).WithError(err).Error("failed to fetch user")
func (l *Logger) WithError(err error) *Logger {
	if l == nil || l.zl == nil || err == nil {
		return l
	}
	sub := l.zl.With().Str("error", err.Error()).Logger()
	return &Logger{zl: &sub}
}

func (l *Logger) logIfReady(level string, msg string) {
	if l == nil || l.zl == nil {
		return
	}
	switch level {
	case "debug":
		l.zl.Debug().Msg(msg)
	case "info":
		l.zl.Info().Msg(msg)
	case "warn":
		l.zl.Warn().Msg(msg)
	case "error":
		l.zl.Error().Msg(msg)
	case "fatal":
		l.zl.Fatal().Msg(msg)
	default:
		l.zl.Info().Msg(msg)
	}
}

// Debug logs at debug level.
// NOTE: debug logs are NOT dropped by Promtail currently.
// Add a Promtail drop rule for level=debug when log volume becomes a concern.
func (l *Logger) Debug(msg string) {
	l.logIfReady("debug", msg)
}

// Debugf logs at debug level with formatting.
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logIfReady("debug", fmt.Sprintf(format, v...))
}

// Info logs at info level.
func (l *Logger) Info(msg string) {
	l.logIfReady("info", msg)
}

// Infof logs at info level with formatting.
func (l *Logger) Infof(format string, v ...interface{}) {
	l.logIfReady("info", fmt.Sprintf(format, v...))
}

// Warn logs at warn level.
func (l *Logger) Warn(msg string) {
	l.logIfReady("warn", msg)
}

// Warnf logs at warn level with formatting.
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logIfReady("warn", fmt.Sprintf(format, v...))
}

// Error logs at error level.
func (l *Logger) Error(msg string) {
	l.logIfReady("error", msg)
}

// Errorf logs at error level with formatting.
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logIfReady("error", fmt.Sprintf(format, v...))
}

// Fatal logs at fatal level and exits. Do not use in request path.
func (l *Logger) Fatal(msg string) {
	if l != nil && l.zl != nil {
		l.zl.Fatal().Msg(msg)
	}
	os.Exit(1)
}

// Fatalf logs at fatal level with formatting and exits. Do not use in request path.
func (l *Logger) Fatalf(format string, v ...interface{}) {
	if l != nil && l.zl != nil {
		l.zl.Fatal().Msg(fmt.Sprintf(format, v...))
	}
	os.Exit(1)
}
