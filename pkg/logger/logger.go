package logger

import "sync"

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
)

func init() {
	globalMiddlewareChain = NewMiddlewareChain()
}

// AddMiddleware adds a middleware to the global middleware chain
// This middleware will be applied to all log entries processed by the logger
func AddMiddleware(middleware Middleware) {
	middlewareMutex.Lock()
	defer middlewareMutex.Unlock()
	globalMiddlewareChain.Add(middleware)
}

// ClearMiddleware removes all middleware from the global chain
func ClearMiddleware() {
	middlewareMutex.Lock()
	defer middlewareMutex.Unlock()
	globalMiddlewareChain.Clear()
}

// GetMiddlewareCount returns the number of middleware in the global chain
func GetMiddlewareCount() int {
	middlewareMutex.RLock()
	defer middlewareMutex.RUnlock()
	return globalMiddlewareChain.Count()
}

// processWithGlobalMiddleware processes a log entry through the global middleware chain
func processWithGlobalMiddleware(entry *LogEntry, final func(*LogEntry)) {
	middlewareMutex.RLock()
	defer middlewareMutex.RUnlock()
	globalMiddlewareChain.Process(entry, final)
}

// D is the global direct logger instance
// Usage: logger.D.Infof("message", args...)
var D Logger = NewDirectLogger()
