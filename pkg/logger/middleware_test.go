package logger

import (
	"encoding/json"
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
		Fields:    make(map[string]interface{}),
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

	// Add middleware that adds a field
	middleware := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Fields["middleware"] = "test-value"
		next(entry)
	}
	chain.Add(middleware)

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "test message",
		Fields:    make(map[string]interface{}),
	}

	var processedEntry *LogEntry
	chain.Process(entry, func(e *LogEntry) {
		processedEntry = e
	})

	// Verify the field was added
	if processedEntry.Fields["middleware"] != "test-value" {
		t.Errorf("Expected field 'middleware' to be 'test-value', got: %v", processedEntry.Fields["middleware"])
	}
}

func TestMiddlewareChain_Process_MultipleMiddleware(t *testing.T) {
	chain := NewMiddlewareChain()

	// Add middleware that adds fields
	middleware1 := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Fields["middleware1"] = "value1"
		next(entry)
	}

	middleware2 := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Fields["middleware2"] = "value2"
		next(entry)
	}

	middleware3 := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Fields["middleware3"] = "value3"
		next(entry)
	}

	chain.Add(middleware1)
	chain.Add(middleware2)
	chain.Add(middleware3)

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "test message",
		Fields:    make(map[string]interface{}),
	}

	var processedEntry *LogEntry
	chain.Process(entry, func(e *LogEntry) {
		processedEntry = e
	})

	// Verify all fields were added
	expectedFields := map[string]string{
		"middleware1": "value1",
		"middleware2": "value2",
		"middleware3": "value3",
	}

	for key, expectedValue := range expectedFields {
		if processedEntry.Fields[key] != expectedValue {
			t.Errorf("Expected field '%s' to be '%s', got: %v", key, expectedValue, processedEntry.Fields[key])
		}
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
		entry.Fields["middleware1"] = "executed"
		// Don't call next() - this should stop the chain
	}

	middleware2 := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Fields["middleware2"] = "executed"
		next(entry)
	}

	chain.Add(middleware1)
	chain.Add(middleware2)

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "test message",
		Fields:    make(map[string]interface{}),
	}

	finalCalled := false
	chain.Process(entry, func(e *LogEntry) {
		finalCalled = true
	})

	// Verify first middleware executed
	if entry.Fields["middleware1"] != "executed" {
		t.Error("Expected middleware1 to execute")
	}

	// Verify second middleware did not execute
	if _, exists := entry.Fields["middleware2"]; exists {
		t.Error("Expected middleware2 to not execute when previous middleware skips next()")
	}

	// Verify final function was not called
	if finalCalled {
		t.Error("Expected final function to not be called when middleware skips next()")
	}
}

// TestMiddlewareInterface tests the basic middleware interface
func TestMiddlewareInterface(t *testing.T) {
	// Create a simple middleware that adds a field
	middleware := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Fields["test-field"] = "test-value"
		next(entry)
	}

	// Create a test entry
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "test message",
		Fields:    make(map[string]interface{}),
	}

	// Test the middleware
	var processedEntry *LogEntry
	middleware(entry, func(e *LogEntry) {
		processedEntry = e
	})

	// Verify the field was added
	if processedEntry.Fields["test-field"] != "test-value" {
		t.Errorf("Expected field 'test-field' to be 'test-value', got: %v", processedEntry.Fields["test-field"])
	}
}

// TestMiddlewareChain tests the middleware chain functionality
func TestMiddlewareChain(t *testing.T) {
	chain := NewMiddlewareChain()

	// Add middleware that adds fields
	middleware1 := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Fields["middleware1"] = "value1"
		next(entry)
	}

	middleware2 := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Fields["middleware2"] = "value2"
		next(entry)
	}

	chain.Add(middleware1)
	chain.Add(middleware2)

	// Test entry
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "test message",
		Fields:    make(map[string]interface{}),
	}

	// Process through chain
	var processedEntry *LogEntry
	chain.Process(entry, func(e *LogEntry) {
		processedEntry = e
	})

	// Verify both fields were added
	if processedEntry.Fields["middleware1"] != "value1" {
		t.Errorf("Expected field 'middleware1' to be 'value1', got: %v", processedEntry.Fields["middleware1"])
	}
	if processedEntry.Fields["middleware2"] != "value2" {
		t.Errorf("Expected field 'middleware2' to be 'value2', got: %v", processedEntry.Fields["middleware2"])
	}
}

// TestMiddlewareChainClear tests the Clear functionality
func TestMiddlewareChainClear(t *testing.T) {
	chain := NewMiddlewareChain()

	// Add a middleware
	middleware := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Fields["test"] = "value"
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
		Fields:    make(map[string]interface{}),
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
		entry.Fields["global1"] = "value1"
		next(entry)
	}

	middleware2 := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Fields["global2"] = "value2"
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
		Fields:    make(map[string]interface{}),
	}

	var processedEntry *LogEntry
	processWithGlobalMiddleware(entry, func(e *LogEntry) {
		processedEntry = e
	})

	// Verify both middleware were applied
	if processedEntry.Fields["global1"] != "value1" {
		t.Errorf("Expected field 'global1' to be 'value1', got: %v", processedEntry.Fields["global1"])
	}
	if processedEntry.Fields["global2"] != "value2" {
		t.Errorf("Expected field 'global2' to be 'value2', got: %v", processedEntry.Fields["global2"])
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
		entry.Fields["middleware-field"] = "middleware-value"
		next(entry)
	})

	// Log a message
	logger.Infof("test message")

	// Parse the output to verify middleware was applied
	output := buf.String()
	if !strings.Contains(output, "middleware-field") {
		t.Error("Expected middleware field to be present in log output")
	}
	if !strings.Contains(output, "middleware-value") {
		t.Error("Expected middleware value to be present in log output")
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

	// Add middleware that adds a field
	AddMiddleware(func(entry *LogEntry, next func(*LogEntry)) {
		entry.Fields["middleware_field"] = "middleware_value"
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

	// Parse JSON to verify middleware was applied
	var logData map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logData); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	// Check that the middleware field was added
	runtime, ok := logData["runtime"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected runtime section in log output")
	}

	lines, ok := runtime["lines"].([]interface{})
	if !ok || len(lines) == 0 {
		t.Fatal("Expected lines array in runtime section")
	}

	firstLine, ok := lines[0].(map[string]interface{})
	if !ok {
		t.Fatal("Expected first line to be an object")
	}

	fields, ok := firstLine["fields"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected fields object in log entry")
	}

	if fields["middleware_field"] != "middleware_value" {
		t.Errorf("Expected middleware field to be added, got: %v", fields)
	}

	// Clean up
	ClearMiddleware()
}
