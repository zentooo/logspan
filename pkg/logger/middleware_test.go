package logger

import (
	"strings"
	"testing"
	"time"
)

func TestNewMiddlewareChain(t *testing.T) {
	chain := NewMiddlewareChain()
	if chain == nil {
		t.Fatal("NewMiddlewareChain should return a non-nil chain")
	}
	if chain.Count() != 0 {
		t.Errorf("New chain should have 0 middleware, got %d", chain.Count())
	}
}

func TestMiddlewareChain_Add(t *testing.T) {
	chain := NewMiddlewareChain()

	// Add first middleware
	middleware1 := func(entry *LogEntry, next func(*LogEntry)) {
		next(entry)
	}
	chain.Add(middleware1)

	if chain.Count() != 1 {
		t.Errorf("Expected 1 middleware, got %d", chain.Count())
	}

	// Add second middleware
	middleware2 := func(entry *LogEntry, next func(*LogEntry)) {
		next(entry)
	}
	chain.Add(middleware2)

	if chain.Count() != 2 {
		t.Errorf("Expected 2 middleware, got %d", chain.Count())
	}
}

func TestMiddlewareChain_Process_EmptyChain(t *testing.T) {
	chain := NewMiddlewareChain()

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "test message",
	}

	var processedEntry *LogEntry
	chain.Process(entry, func(e *LogEntry) {
		processedEntry = e
	})

	if processedEntry != entry {
		t.Error("Expected entry to be unchanged when processing through empty chain")
	}
}

func TestMiddlewareChain_Process_SingleMiddleware(t *testing.T) {
	chain := NewMiddlewareChain()

	// Add middleware that modifies the message
	middleware := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Message = entry.Message + " [middleware]"
		next(entry)
	}
	chain.Add(middleware)

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "test message",
	}

	var processedEntry *LogEntry
	chain.Process(entry, func(e *LogEntry) {
		processedEntry = e
	})

	// Verify the message was modified
	expected := "test message [middleware]"
	if processedEntry.Message != expected {
		t.Errorf("Expected message '%s', got '%s'", expected, processedEntry.Message)
	}
}

func TestMiddlewareChain_Process_MultipleMiddleware(t *testing.T) {
	chain := NewMiddlewareChain()

	// Add middleware that modify the message
	middleware1 := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Message = entry.Message + " [middleware1]"
		next(entry)
	}

	middleware2 := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Message = entry.Message + " [middleware2]"
		next(entry)
	}

	middleware3 := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Message = entry.Message + " [middleware3]"
		next(entry)
	}

	chain.Add(middleware1)
	chain.Add(middleware2)
	chain.Add(middleware3)

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "test message",
	}

	var processedEntry *LogEntry
	chain.Process(entry, func(e *LogEntry) {
		processedEntry = e
	})

	// Verify all middleware modified the message
	expected := "test message [middleware1] [middleware2] [middleware3]"
	if processedEntry.Message != expected {
		t.Errorf("Expected message '%s', got '%s'", expected, processedEntry.Message)
	}
}

func TestMiddlewareChain_Clear(t *testing.T) {
	chain := NewMiddlewareChain()

	// Add some middleware
	middleware := func(entry *LogEntry, next func(*LogEntry)) {
		next(entry)
	}
	chain.Add(middleware)
	chain.Add(middleware)

	if chain.Count() != 2 {
		t.Errorf("Expected count to be 2, got %d", chain.Count())
	}

	chain.Clear()

	if chain.Count() != 0 {
		t.Errorf("Expected count to be 0 after clear, got %d", chain.Count())
	}
}

func TestMiddlewareChain_Process_MiddlewareCanSkipNext(t *testing.T) {
	chain := NewMiddlewareChain()

	// Add middleware that skips calling next
	middleware1 := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Message = entry.Message + " [middleware1]"
		// Don't call next() - this should stop the chain
	}

	middleware2 := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Message = entry.Message + " [middleware2]"
		next(entry)
	}

	chain.Add(middleware1)
	chain.Add(middleware2)

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "test message",
	}

	finalCalled := false
	chain.Process(entry, func(e *LogEntry) {
		finalCalled = true
	})

	// Verify first middleware executed
	if !strings.Contains(entry.Message, "[middleware1]") {
		t.Error("Expected middleware1 to execute")
	}

	// Verify second middleware did not execute
	if strings.Contains(entry.Message, "[middleware2]") {
		t.Error("Expected middleware2 to not execute when previous middleware skips next()")
	}

	// Verify final function was not called
	if finalCalled {
		t.Error("Expected final function to not be called when middleware skips next()")
	}
}

// TestMiddlewareInterface tests the basic middleware interface
func TestMiddlewareInterface(t *testing.T) {
	// Create a simple middleware that modifies the message
	middleware := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Message = entry.Message + " [test-middleware]"
		next(entry)
	}

	// Create a test entry
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "test message",
	}

	// Test the middleware
	var processedEntry *LogEntry
	middleware(entry, func(e *LogEntry) {
		processedEntry = e
	})

	// Verify the message was modified
	expected := "test message [test-middleware]"
	if processedEntry.Message != expected {
		t.Errorf("Expected message '%s', got '%s'", expected, processedEntry.Message)
	}
}

// TestMiddlewareChain tests the middleware chain functionality
func TestMiddlewareChain(t *testing.T) {
	chain := NewMiddlewareChain()

	// Add middleware that modify the message
	middleware1 := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Message = entry.Message + " [chain1]"
		next(entry)
	}

	middleware2 := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Message = entry.Message + " [chain2]"
		next(entry)
	}

	chain.Add(middleware1)
	chain.Add(middleware2)

	// Test entry
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "test message",
	}

	// Process through chain
	var processedEntry *LogEntry
	chain.Process(entry, func(e *LogEntry) {
		processedEntry = e
	})

	// Verify both middleware modified the message
	expected := "test message [chain1] [chain2]"
	if processedEntry.Message != expected {
		t.Errorf("Expected message '%s', got '%s'", expected, processedEntry.Message)
	}
}

// TestMiddlewareChainClear tests the Clear functionality
func TestMiddlewareChainClear(t *testing.T) {
	chain := NewMiddlewareChain()

	// Add a middleware
	middleware := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Message = entry.Message + " [test]"
		next(entry)
	}
	chain.Add(middleware)

	if chain.Count() != 1 {
		t.Errorf("Expected count to be 1, got %d", chain.Count())
	}

	// Clear the chain
	chain.Clear()

	if chain.Count() != 0 {
		t.Errorf("Expected count to be 0 after clear, got %d", chain.Count())
	}
}

// TestMiddlewareChainEmpty tests processing with empty chain
func TestMiddlewareChainEmpty(t *testing.T) {
	chain := NewMiddlewareChain()

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "test message",
	}

	// Process through empty chain
	var processedEntry *LogEntry
	chain.Process(entry, func(e *LogEntry) {
		processedEntry = e
	})

	// Entry should be unchanged
	if processedEntry != entry {
		t.Error("Expected entry to be unchanged when processing through empty chain")
	}
}

// TestGlobalMiddlewareManagement tests the global middleware management functions
func TestGlobalMiddlewareManagement(t *testing.T) {
	// Clear any existing middleware
	ClearMiddleware()

	// Test initial state
	if GetMiddlewareCount() != 0 {
		t.Errorf("Expected initial middleware count to be 0, got %d", GetMiddlewareCount())
	}

	// Add middleware
	middleware1 := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Message = entry.Message + " [global1]"
		next(entry)
	}

	middleware2 := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Message = entry.Message + " [global2]"
		next(entry)
	}

	AddMiddleware(middleware1)
	if GetMiddlewareCount() != 1 {
		t.Errorf("Expected middleware count to be 1, got %d", GetMiddlewareCount())
	}

	AddMiddleware(middleware2)
	if GetMiddlewareCount() != 2 {
		t.Errorf("Expected middleware count to be 2, got %d", GetMiddlewareCount())
	}

	// Test processing through global middleware
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "test message",
	}

	var processedEntry *LogEntry
	processWithGlobalMiddleware(entry, func(e *LogEntry) {
		processedEntry = e
	})

	// Verify both middleware were applied
	if processedEntry.Message != "test message [global1] [global2]" {
		t.Errorf("Expected message '%s', got '%s'", "test message [global1] [global2]", processedEntry.Message)
	}

	// Test clear
	ClearMiddleware()
	if GetMiddlewareCount() != 0 {
		t.Errorf("Expected middleware count to be 0 after clear, got %d", GetMiddlewareCount())
	}
}

// TestDirectLoggerWithMiddleware tests that DirectLogger integrates with global middleware
func TestDirectLoggerWithMiddleware(t *testing.T) {
	// Clear any existing middleware
	ClearMiddleware()

	// Create a buffer to capture output
	var buf strings.Builder
	logger := NewDirectLogger()
	logger.SetOutput(&buf)

	// Add middleware that adds a field
	AddMiddleware(func(entry *LogEntry, next func(*LogEntry)) {
		entry.Message = entry.Message + " [middleware-field]"
		next(entry)
	})

	// Log a message
	logger.Infof("test message")

	// Parse the output to verify middleware was applied
	output := buf.String()
	if !strings.Contains(output, "middleware-field") {
		t.Error("Expected middleware field to be present in log output")
	}

	// Clean up
	ClearMiddleware()
}

// TestContextLoggerWithMiddleware tests that ContextLogger integrates with global middleware
func TestContextLoggerWithMiddleware(t *testing.T) {
	// Clear any existing middleware
	ClearMiddleware()

	// Create a buffer to capture output
	var buf strings.Builder
	logger := NewContextLogger()
	logger.SetOutput(&buf)

	// Add middleware that modifies the message
	AddMiddleware(func(entry *LogEntry, next func(*LogEntry)) {
		entry.Message = entry.Message + " [middleware_field]"
		next(entry)
	})

	// Log a message and flush
	logger.Infof("test message")
	logger.Flush()

	// Parse the output
	output := buf.String()
	if output == "" {
		t.Fatal("Expected output, got empty string")
	}

	// Check that the middleware modified the message
	if !strings.Contains(output, "[middleware_field]") {
		t.Error("Expected middleware to modify the message")
	}

	// Clean up
	ClearMiddleware()
}
