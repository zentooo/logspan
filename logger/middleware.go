package logger

// Middleware defines the interface for log entry processing middleware
// Middleware receives a LogEntry and a next function, processes the entry,
// and calls next to continue the chain
type Middleware func(entry *LogEntry, next func(*LogEntry))

// MiddlewareChain manages a chain of middleware functions
type MiddlewareChain struct {
	middlewares []Middleware
}

// NewMiddlewareChain creates a new middleware chain
func NewMiddlewareChain() *MiddlewareChain {
	return &MiddlewareChain{
		middlewares: make([]Middleware, 0),
	}
}

// Add appends a middleware to the chain
func (mc *MiddlewareChain) Add(middleware Middleware) {
	mc.middlewares = append(mc.middlewares, middleware)
}

// Process executes the middleware chain on a log entry
// The final function in the chain is called when all middleware have been processed
func (mc *MiddlewareChain) Process(entry *LogEntry, final func(*LogEntry)) {
	if len(mc.middlewares) == 0 {
		final(entry)
		return
	}

	// Build the chain from the end backwards
	next := final
	for i := len(mc.middlewares) - 1; i >= 0; i-- {
		middleware := mc.middlewares[i]
		currentNext := next
		next = func(e *LogEntry) {
			middleware(e, currentNext)
		}
	}

	// Execute the chain
	next(entry)
}

// Clear removes all middleware from the chain
func (mc *MiddlewareChain) Clear() {
	mc.middlewares = mc.middlewares[:0]
}

// Count returns the number of middleware in the chain
func (mc *MiddlewareChain) Count() int {
	return len(mc.middlewares)
}
