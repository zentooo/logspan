package logger

import (
	"fmt"
	"os"
	"time"
)

// DirectLogger implements the Logger interface for direct logging without context
type DirectLogger struct {
	*BaseLogger
}

// NewDirectLogger creates a new DirectLogger instance
func NewDirectLogger() *DirectLogger {
	base := newBaseLogger()
	base.output = os.Stdout // Set default output for DirectLogger
	return &DirectLogger{
		BaseLogger: &base,
	}
}

// logf writes a log entry with the given level and message in structured format
func (l *DirectLogger) logf(level LogLevel, format string, args ...interface{}) {
	if !l.isLevelEnabled(level) {
		return
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	// Do nothing if output is nil
	if l.output == nil {
		return
	}

	now := time.Now()
	entry := &LogEntry{
		Timestamp: now,
		Level:     level.String(),
		Message:   fmt.Sprintf(format, args...),
	}

	// Add source information if enabled
	config := GetConfig()
	if config.EnableSourceInfo {
		// Skip levels: getSourceInfo(0) -> logf(1) -> Infof/Debugf/etc(2) -> actual caller(3)
		sourceInfo := getSourceInfo(3)
		entry.Funcname = sourceInfo.Funcname
		entry.Filename = sourceInfo.Filename
		entry.Fileline = sourceInfo.Fileline
	}

	// Process through global middleware chain
	processWithGlobalMiddleware(entry, func(processedEntry *LogEntry) {
		// Use the formatter (default or explicitly set)
		jsonData, err := formatLogOutput([]*LogEntry{processedEntry}, make(map[string]interface{}), processedEntry.Timestamp, processedEntry.Timestamp, l.formatter)
		if err != nil {
			_, _ = fmt.Fprintf(l.output, "Error formatting log: %v\n", err)
			return
		}

		if _, err := fmt.Fprintf(l.output, "%s\n", jsonData); err != nil {
			// If writing fails, try to write an error message
			// This is a best-effort attempt since the output might be broken
			// We intentionally ignore any error from this fallback write
			_, _ = fmt.Fprintf(l.output, "Error writing log output: %v\n", err)
		}
	})
}

// Debugf logs a debug message
func (l *DirectLogger) Debugf(format string, args ...interface{}) {
	l.logf(DebugLevel, format, args...)
}

// Infof logs an info message
func (l *DirectLogger) Infof(format string, args ...interface{}) {
	l.logf(InfoLevel, format, args...)
}

// Warnf logs a warning message
func (l *DirectLogger) Warnf(format string, args ...interface{}) {
	l.logf(WarnLevel, format, args...)
}

// Errorf logs an error message
func (l *DirectLogger) Errorf(format string, args ...interface{}) {
	l.logf(ErrorLevel, format, args...)
}

// Criticalf logs a critical message
func (l *DirectLogger) Criticalf(format string, args ...interface{}) {
	l.logf(CriticalLevel, format, args...)
}
