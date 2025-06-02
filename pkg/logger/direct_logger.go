package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/zentooo/logspan/pkg/formatter"
)

// DirectLogger implements the Logger interface for direct logging without context
type DirectLogger struct {
	output    io.Writer
	minLevel  LogLevel
	formatter formatter.Formatter
	mu        sync.Mutex
}

// NewDirectLogger creates a new DirectLogger instance
func NewDirectLogger() *DirectLogger {
	// Get global config to determine formatter settings
	config := GetConfig()
	var jsonFormatter formatter.Formatter
	if config.PrettifyJSON {
		jsonFormatter = formatter.NewJSONFormatterWithIndent("  ")
	} else {
		jsonFormatter = formatter.NewJSONFormatter()
	}

	return &DirectLogger{
		output:    os.Stdout,
		minLevel:  InfoLevel,
		formatter: jsonFormatter,
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

// SetFormatter sets the formatter for the logger
func (l *DirectLogger) SetFormatter(f formatter.Formatter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.formatter = f
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
	}

	// Process through global middleware chain
	processWithGlobalMiddleware(entry, func(processedEntry *LogEntry) {
		// Use the formatter (default or explicitly set)
		jsonData, err := formatLogOutput([]*LogEntry{processedEntry}, make(map[string]interface{}), processedEntry.Timestamp, processedEntry.Timestamp, l.formatter)
		if err != nil {
			fmt.Fprintf(l.output, "Error formatting log: %v\n", err)
			return
		}

		fmt.Fprintf(l.output, "%s\n", jsonData)
	})
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
