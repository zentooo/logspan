package logger

import (
	"encoding/json"
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

// log writes a log entry with the given level and message in structured format
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

	now := time.Now()
	entry := &LogEntry{
		Timestamp: now,
		Level:     level.String(),
		Message:   fmt.Sprintf(format, args...),
		Fields:    make(map[string]interface{}),
		Tags:      make([]string, 0),
	}

	// Create structured log output similar to context logger
	output := map[string]interface{}{
		"type":    "request",
		"context": map[string]interface{}{},
		"runtime": map[string]interface{}{
			"severity":  level.String(),
			"startTime": now.Format(time.RFC3339Nano),
			"endTime":   now.Format(time.RFC3339Nano),
			"elapsed":   0, // Direct log has no elapsed time
			"lines":     []*LogEntry{entry},
		},
		"config": map[string]interface{}{
			"elapsedUnit": "ms",
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(output)
	if err != nil {
		// Fallback to simple output if JSON marshaling fails
		fmt.Fprintf(l.output, "Error marshaling log: %v\n", err)
		return
	}

	fmt.Fprintf(l.output, "%s\n", jsonData)
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
