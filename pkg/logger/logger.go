package logger

import (
	"sync"
	"time"

	"github.com/zentooo/logspan/pkg/formatter"
)

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

// Global middleware management
var (
	globalMiddlewareChain *MiddlewareChain
	middlewareMutex       sync.RWMutex
	middlewareOnce        sync.Once
)

// ensureMiddlewareChain ensures the global middleware chain is initialized
func ensureMiddlewareChain() {
	middlewareOnce.Do(func() {
		globalMiddlewareChain = NewMiddlewareChain()
	})
}

// AddMiddleware adds a middleware to the global middleware chain
// This middleware will be applied to all log entries processed by the logger
func AddMiddleware(middleware Middleware) {
	ensureMiddlewareChain()
	middlewareMutex.Lock()
	defer middlewareMutex.Unlock()
	globalMiddlewareChain.Add(middleware)
}

// ClearMiddleware removes all middleware from the global chain
func ClearMiddleware() {
	ensureMiddlewareChain()
	middlewareMutex.Lock()
	defer middlewareMutex.Unlock()
	globalMiddlewareChain.Clear()
}

// GetMiddlewareCount returns the number of middleware in the global chain
func GetMiddlewareCount() int {
	ensureMiddlewareChain()
	middlewareMutex.RLock()
	defer middlewareMutex.RUnlock()
	return globalMiddlewareChain.Count()
}

// processWithGlobalMiddleware processes a log entry through the global middleware chain
func processWithGlobalMiddleware(entry *LogEntry, final func(*LogEntry)) {
	ensureMiddlewareChain()
	middlewareMutex.RLock()
	defer middlewareMutex.RUnlock()
	globalMiddlewareChain.Process(entry, final)
}

// D is the global direct logger instance
// Usage: logger.D.Infof("message", args...)
var D Logger = NewDirectLogger()

// createDefaultFormatter creates a default formatter based on global configuration
func createDefaultFormatter() formatter.Formatter {
	config := GetConfig()
	if config.PrettifyJSON {
		return formatter.NewJSONFormatterWithIndent("  ")
	}
	return formatter.NewJSONFormatter()
}

// formatLogOutput creates a LogOutput structure and formats it using the given formatter
// If formatter is nil, uses default JSONFormatter
func formatLogOutput(entries []*LogEntry, contextFields map[string]interface{}, startTime, endTime time.Time, f formatter.Formatter) ([]byte, error) {
	elapsed := endTime.Sub(startTime).Milliseconds()

	// Find the highest severity level
	maxSeverity := DebugLevel
	for _, entry := range entries {
		entryLevel := ParseLogLevel(entry.Level)
		maxSeverity = GetHigherLevel(entryLevel, maxSeverity)
	}

	// Convert logger.LogEntry to formatter.LogEntry
	formatterEntries := make([]*formatter.LogEntry, len(entries))
	for i, entry := range entries {
		formatterEntries[i] = &formatter.LogEntry{
			Timestamp: entry.Timestamp,
			Level:     entry.Level,
			Message:   entry.Message,
			Funcname:  entry.Funcname,
			Filename:  entry.Filename,
			Fileline:  entry.Fileline,
		}
	}

	// Create LogOutput structure
	logOutput := &formatter.LogOutput{
		Type:    "request",
		Context: contextFields,
		Runtime: formatter.RuntimeInfo{
			Severity:  maxSeverity.String(),
			StartTime: startTime.Format(time.RFC3339Nano),
			EndTime:   endTime.Format(time.RFC3339Nano),
			Elapsed:   elapsed,
			Lines:     formatterEntries,
		},
	}

	// Use provided formatter or default JSONFormatter
	if f == nil {
		// Use default JSONFormatter with prettify setting from global config
		f = createDefaultFormatter()
	}

	return f.Format(logOutput)
}
