package logger

// Logger defines the interface for logging operations
type Logger interface {
	// Debugf logs a debug message
	Debugf(format string, args ...interface{})

	// Infof logs an info message
	Infof(format string, args ...interface{})

	// Warnf logs a warning message
	Warnf(format string, args ...interface{})

	// Errorf logs an error message
	Errorf(format string, args ...interface{})

	// Criticalf logs a critical message
	Criticalf(format string, args ...interface{})
}

// D is the global direct logger instance
// Usage: logger.D.Infof("message", args...)
var D Logger = NewDirectLogger()
