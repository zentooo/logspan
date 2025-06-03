package logger

import (
	"testing"
)

func TestMiddlewareManager_EnsureMiddlewareChain(t *testing.T) {
	// Clear any existing middleware
	ClearMiddleware()

	// Test that ensureMiddlewareChain initializes the chain
	ensureMiddlewareChain()
	if globalMiddlewareChain == nil {
		t.Error("Expected globalMiddlewareChain to be initialized")
	}

	// Test that multiple calls don't reinitialize
	oldChain := globalMiddlewareChain
	ensureMiddlewareChain()
	if globalMiddlewareChain != oldChain {
		t.Error("Expected globalMiddlewareChain to not be reinitialized")
	}
}

func TestMiddlewareManager_AddMiddleware(t *testing.T) {
	// Clear any existing middleware
	ClearMiddleware()

	// Test adding middleware
	middleware1 := func(entry *LogEntry, next func(*LogEntry)) {
		next(entry)
	}

	AddMiddleware(middleware1)
	if GetMiddlewareCount() != 1 {
		t.Errorf("Expected middleware count to be 1, got %d", GetMiddlewareCount())
	}

	// Test adding another middleware
	middleware2 := func(entry *LogEntry, next func(*LogEntry)) {
		next(entry)
	}

	AddMiddleware(middleware2)
	if GetMiddlewareCount() != 2 {
		t.Errorf("Expected middleware count to be 2, got %d", GetMiddlewareCount())
	}
}

func TestMiddlewareManager_ClearMiddleware(t *testing.T) {
	// Clear any existing middleware first
	ClearMiddleware()

	// Verify we start with 0 middleware
	if GetMiddlewareCount() != 0 {
		t.Errorf("Expected initial middleware count to be 0, got %d", GetMiddlewareCount())
	}

	// Add some middleware
	middleware := func(entry *LogEntry, next func(*LogEntry)) {
		next(entry)
	}
	AddMiddleware(middleware)
	AddMiddleware(middleware)

	if GetMiddlewareCount() != 2 {
		t.Errorf("Expected middleware count to be 2, got %d", GetMiddlewareCount())
	}

	// Clear middleware
	ClearMiddleware()
	if GetMiddlewareCount() != 0 {
		t.Errorf("Expected middleware count to be 0 after clear, got %d", GetMiddlewareCount())
	}
}

func TestMiddlewareManager_GetMiddlewareCount(t *testing.T) {
	// Clear any existing middleware
	ClearMiddleware()

	// Test initial count
	if GetMiddlewareCount() != 0 {
		t.Errorf("Expected initial middleware count to be 0, got %d", GetMiddlewareCount())
	}

	// Add middleware and test count
	middleware := func(entry *LogEntry, next func(*LogEntry)) {
		next(entry)
	}

	for i := 1; i <= 5; i++ {
		AddMiddleware(middleware)
		if GetMiddlewareCount() != i {
			t.Errorf("Expected middleware count to be %d, got %d", i, GetMiddlewareCount())
		}
	}
}

func TestMiddlewareManager_ProcessWithGlobalMiddleware(t *testing.T) {
	// Clear any existing middleware
	ClearMiddleware()

	// Test processing without middleware
	entry := &LogEntry{
		Level:   "INFO",
		Message: "test message",
	}

	processed := false
	processWithGlobalMiddleware(entry, func(e *LogEntry) {
		processed = true
		if e != entry {
			t.Error("Expected same entry to be passed to final function")
		}
	})

	if !processed {
		t.Error("Expected final function to be called")
	}

	// Test processing with middleware
	ClearMiddleware()
	middlewareCalled := false
	AddMiddleware(func(e *LogEntry, next func(*LogEntry)) {
		middlewareCalled = true
		e.Message = "modified by middleware"
		next(e)
	})

	processed = false
	processWithGlobalMiddleware(entry, func(e *LogEntry) {
		processed = true
		if e.Message != "modified by middleware" {
			t.Errorf("Expected message to be modified by middleware, got: %s", e.Message)
		}
	})

	if !middlewareCalled {
		t.Error("Expected middleware to be called")
	}
	if !processed {
		t.Error("Expected final function to be called")
	}
}

func TestMiddlewareManager_ConcurrentAccess(t *testing.T) {
	// Clear any existing middleware
	ClearMiddleware()

	// Test concurrent access to middleware management
	done := make(chan bool, 10)

	// Start multiple goroutines that add middleware
	for i := 0; i < 5; i++ {
		go func() {
			middleware := func(entry *LogEntry, next func(*LogEntry)) {
				next(entry)
			}
			AddMiddleware(middleware)
			done <- true
		}()
	}

	// Start multiple goroutines that read middleware count
	for i := 0; i < 5; i++ {
		go func() {
			GetMiddlewareCount()
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify that all middleware was added
	if GetMiddlewareCount() != 5 {
		t.Errorf("Expected middleware count to be 5, got %d", GetMiddlewareCount())
	}
}
