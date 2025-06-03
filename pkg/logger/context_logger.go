package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/zentooo/logspan/pkg/formatter"
)

// ContextLogger implements context-based logging with log aggregation
type ContextLogger struct {
	entries    []*LogEntry
	fields     map[string]interface{}
	startTime  time.Time
	output     io.Writer
	minLevel   LogLevel
	formatter  formatter.Formatter
	maxEntries int // Maximum number of entries before auto-flush
	mu         sync.Mutex
}

// NewContextLogger creates a new ContextLogger instance
func NewContextLogger() *ContextLogger {
	// Get global config to determine formatter settings
	config := GetConfig()
	var jsonFormatter formatter.Formatter
	if config.PrettifyJSON {
		jsonFormatter = formatter.NewJSONFormatterWithIndent("  ")
	} else {
		jsonFormatter = formatter.NewJSONFormatter()
	}

	return &ContextLogger{
		entries:    make([]*LogEntry, 0),
		fields:     make(map[string]interface{}),
		startTime:  time.Now(),
		output:     os.Stdout,
		minLevel:   config.MinLevel,
		formatter:  jsonFormatter,
		maxEntries: config.MaxLogEntries,
	}
}

// SetOutput sets the output writer for the logger
func (l *ContextLogger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = w
}

// SetLevel sets the minimum log level for filtering
func (l *ContextLogger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.minLevel = level
}

// isLevelEnabled checks if the given level should be logged
func (l *ContextLogger) isLevelEnabled(level LogLevel) bool {
	return level >= l.minLevel
}

// addEntry adds a log entry to the context logger
func (l *ContextLogger) addEntry(level LogLevel, message string) {
	if !l.isLevelEnabled(level) {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     level.String(),
		Message:   message,
	}

	// Process through global middleware chain
	processWithGlobalMiddleware(entry, func(processedEntry *LogEntry) {
		l.entries = append(l.entries, processedEntry)

		// Check if we need to auto-flush due to entry limit
		if l.maxEntries > 0 && len(l.entries) >= l.maxEntries {
			l.flushInternal()
		}
	})
}

// flushInternal performs the flush operation without acquiring the mutex
// This method assumes the mutex is already held by the caller
func (l *ContextLogger) flushInternal() {
	if l.output == nil || len(l.entries) == 0 {
		return
	}

	endTime := time.Now()

	// Use the formatter (default or explicitly set)
	jsonData, err := formatLogOutput(l.entries, l.fields, l.startTime, endTime, l.formatter)
	if err != nil {
		// Fallback to simple output if formatting fails
		_, _ = fmt.Fprintf(l.output, "Error formatting log: %v\n", err)
		return
	}

	if _, err := fmt.Fprintf(l.output, "%s\n", jsonData); err != nil {
		// If writing fails, try to write an error message
		// This is a best-effort attempt since the output might be broken
		// We intentionally ignore any error from this fallback write
		_, _ = fmt.Fprintf(l.output, "Error writing log output: %v\n", err)
	}

	// Clear entries after flushing and reset start time
	l.entries = l.entries[:0] // Clear slice but keep capacity
	l.startTime = time.Now()  // Reset start time for next batch
}

// Flush outputs all accumulated log entries as a single JSON
func (l *ContextLogger) Flush() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.flushInternal()
}

// AddContextValue adds a field to the context
func (l *ContextLogger) AddContextValue(key string, value interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.fields[key] = value
}

// AddContextValues adds multiple fields to the context
func (l *ContextLogger) AddContextValues(fields map[string]interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for k, v := range fields {
		l.fields[k] = v
	}
}

// SetFormatter sets the formatter for the logger
func (l *ContextLogger) SetFormatter(f formatter.Formatter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.formatter = f
}

// Debugf logs a debug message
func (l *ContextLogger) Debugf(format string, args ...interface{}) {
	l.addEntry(DebugLevel, fmt.Sprintf(format, args...))
}

// Infof logs an info message
func (l *ContextLogger) Infof(format string, args ...interface{}) {
	l.addEntry(InfoLevel, fmt.Sprintf(format, args...))
}

// Warnf logs a warning message
func (l *ContextLogger) Warnf(format string, args ...interface{}) {
	l.addEntry(WarnLevel, fmt.Sprintf(format, args...))
}

// Errorf logs an error message
func (l *ContextLogger) Errorf(format string, args ...interface{}) {
	l.addEntry(ErrorLevel, fmt.Sprintf(format, args...))
}

// Criticalf logs a critical message
func (l *ContextLogger) Criticalf(format string, args ...interface{}) {
	l.addEntry(CriticalLevel, fmt.Sprintf(format, args...))
}
