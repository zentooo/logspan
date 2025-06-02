package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// DirectLogger implements the Logger interface for direct logging without context
type DirectLogger struct {
	output   io.Writer
	minLevel LogLevel
	mu       sync.Mutex
}

// NewDirectLogger creates a new DirectLogger instance
func NewDirectLogger() *DirectLogger {
	return &DirectLogger{
		output:   os.Stdout,
		minLevel: InfoLevel, // Default to INFO level
	}
}

// SetOutput sets the output writer for the logger
func (l *DirectLogger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = w
}

// SetLevel sets the minimum log level for filtering
func (l *DirectLogger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.minLevel = level
}

// SetLevelFromString sets the minimum log level from a string
func (l *DirectLogger) SetLevelFromString(level string) {
	l.SetLevel(ParseLogLevel(level))
}

// isLevelEnabled checks if the given level should be logged
func (l *DirectLogger) isLevelEnabled(level LogLevel) bool {
	return level >= l.minLevel
}

// log writes a log entry with the given level and message
func (l *DirectLogger) log(level LogLevel, format string, args ...interface{}) {
	if !l.isLevelEnabled(level) {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// nilの出力先の場合は何もしない
	if l.output == nil {
		return
	}

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     level.String(),
		Message:   fmt.Sprintf(format, args...),
	}

	// TODO: Add formatter support
	fmt.Fprintf(l.output, "[%s] %s: %s\n", entry.Timestamp.Format(time.RFC3339), entry.Level, entry.Message)
}

// Debugf logs a debug message
func (l *DirectLogger) Debugf(format string, args ...interface{}) {
	l.log(DebugLevel, format, args...)
}

// Infof logs an info message
func (l *DirectLogger) Infof(format string, args ...interface{}) {
	l.log(InfoLevel, format, args...)
}

// Warnf logs a warning message
func (l *DirectLogger) Warnf(format string, args ...interface{}) {
	l.log(WarnLevel, format, args...)
}

// Errorf logs an error message
func (l *DirectLogger) Errorf(format string, args ...interface{}) {
	l.log(ErrorLevel, format, args...)
}

// Criticalf logs a critical message
func (l *DirectLogger) Criticalf(format string, args ...interface{}) {
	l.log(CriticalLevel, format, args...)
}
