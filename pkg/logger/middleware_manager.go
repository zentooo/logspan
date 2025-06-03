package logger

import (
	"sync"
)

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
