package logger

import (
	"context"
)

// contextKey is a private type for context keys to avoid collisions
type contextKey string

const (
	// loggerContextKey is the key used to store the logger in context
	loggerContextKey contextKey = "logger"
)

// WithLogger returns a new context with the logger attached
func WithLogger(ctx context.Context, logger *ContextLogger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

// FromContext retrieves the logger from the context
// If no logger is found, it returns a new ContextLogger
func FromContext(ctx context.Context) *ContextLogger {
	if logger, ok := ctx.Value(loggerContextKey).(*ContextLogger); ok {
		return logger
	}
	// Return a new logger if none is found in context
	return NewContextLogger()
}

// AddContextValue adds a field to the logger in the context
func AddContextValue(ctx context.Context, key string, value interface{}) {
	logger := FromContext(ctx)
	logger.AddContextValue(key, value)
}

// AddContextValues adds multiple fields to the logger in the context
func AddContextValues(ctx context.Context, fields map[string]interface{}) {
	logger := FromContext(ctx)
	logger.AddContextValues(fields)
}

// Infof logs an info message using the logger from context
func Infof(ctx context.Context, format string, args ...interface{}) {
	logger := FromContext(ctx)
	logger.Infof(format, args...)
}

// Debugf logs a debug message using the logger from context
func Debugf(ctx context.Context, format string, args ...interface{}) {
	logger := FromContext(ctx)
	logger.Debugf(format, args...)
}

// Warnf logs a warning message using the logger from context
func Warnf(ctx context.Context, format string, args ...interface{}) {
	logger := FromContext(ctx)
	logger.Warnf(format, args...)
}

// Errorf logs an error message using the logger from context
func Errorf(ctx context.Context, format string, args ...interface{}) {
	logger := FromContext(ctx)
	logger.Errorf(format, args...)
}

// Criticalf logs a critical message using the logger from context
func Criticalf(ctx context.Context, format string, args ...interface{}) {
	logger := FromContext(ctx)
	logger.Criticalf(format, args...)
}

// FlushContext flushes the logger from the context
func FlushContext(ctx context.Context) {
	logger := FromContext(ctx)
	logger.Flush()
}
