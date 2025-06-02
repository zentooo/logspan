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
		Tags:      []string{},
	}

	finalCalled := false
	final := func(e *LogEntry) {
		finalCalled = true
		if e != entry {
			t.Error("Final function should receive the same entry")
		}
	}

	chain.Process(entry, final)

	if !finalCalled {
		t.Error("Final function should be called for empty chain")
	}
}

func TestMiddlewareChain_Process_SingleMiddleware(t *testing.T) {
	chain := NewMiddlewareChain()
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "test message",
		Fields:    make(map[string]interface{}),
		Tags:      []string{},
	}

	middlewareCalled := false
	finalCalled := false

	middleware := func(e *LogEntry, next func(*LogEntry)) {
		middlewareCalled = true
		// Modify the entry
		e.Message = "modified by middleware"
		next(e)
	}

	final := func(e *LogEntry) {
		finalCalled = true
		if e.Message != "modified by middleware" {
			t.Errorf("Expected message to be modified, got: %s", e.Message)
		}
	}

	chain.Add(middleware)
	chain.Process(entry, final)

	if !middlewareCalled {
		t.Error("Middleware should be called")
	}
	if !finalCalled {
		t.Error("Final function should be called")
	}
}

func TestMiddlewareChain_Process_MultipleMiddleware(t *testing.T) {
	chain := NewMiddlewareChain()
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "original",
		Fields:    make(map[string]interface{}),
		Tags:      []string{},
	}

	var executionOrder []string

	middleware1 := func(e *LogEntry, next func(*LogEntry)) {
		executionOrder = append(executionOrder, "middleware1")
		e.Message += " -> middleware1"
		next(e)
	}

	middleware2 := func(e *LogEntry, next func(*LogEntry)) {
		executionOrder = append(executionOrder, "middleware2")
		e.Message += " -> middleware2"
		next(e)
	}

	final := func(e *LogEntry) {
		executionOrder = append(executionOrder, "final")
		expectedMessage := "original -> middleware1 -> middleware2"
		if e.Message != expectedMessage {
			t.Errorf("Expected message: %s, got: %s", expectedMessage, e.Message)
		}
	}

	chain.Add(middleware1)
	chain.Add(middleware2)
	chain.Process(entry, final)

	expectedOrder := []string{"middleware1", "middleware2", "final"}
	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("Expected %d executions, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, expected := range expectedOrder {
		if executionOrder[i] != expected {
			t.Errorf("Expected execution order[%d]: %s, got: %s", i, expected, executionOrder[i])
		}
	}
}

func TestMiddlewareChain_Clear(t *testing.T) {
	chain := NewMiddlewareChain()

	middleware := func(entry *LogEntry, next func(*LogEntry)) {
		next(entry)
	}

	chain.Add(middleware)
	chain.Add(middleware)

	if chain.Count() != 2 {
		t.Errorf("Expected 2 middleware before clear, got %d", chain.Count())
	}

	chain.Clear()

	if chain.Count() != 0 {
		t.Errorf("Expected 0 middleware after clear, got %d", chain.Count())
	}
}

func TestMiddlewareChain_Process_MiddlewareCanSkipNext(t *testing.T) {
	chain := NewMiddlewareChain()
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "test message",
		Fields:    make(map[string]interface{}),
		Tags:      []string{},
	}

	finalCalled := false

	// Middleware that doesn't call next
	middleware := func(e *LogEntry, next func(*LogEntry)) {
		// Intentionally not calling next to test chain interruption
		e.Message = "blocked by middleware"
	}

	final := func(e *LogEntry) {
		finalCalled = true
	}

	chain.Add(middleware)
	chain.Process(entry, final)

	if finalCalled {
		t.Error("Final function should not be called when middleware doesn't call next")
	}

	if entry.Message != "blocked by middleware" {
		t.Errorf("Expected message to be modified by middleware, got: %s", entry.Message)
	}
}

// TestMiddlewareInterface tests the basic middleware interface
func TestMiddlewareInterface(t *testing.T) {
	// Create a simple middleware that adds a tag
	middleware := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Tags = append(entry.Tags, "test-tag")
		next(entry)
	}

	// Create a test entry
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "test message",
		Fields:    make(map[string]interface{}),
		Tags:      make([]string, 0),
	}

	// Test the middleware
	var processedEntry *LogEntry
	middleware(entry, func(e *LogEntry) {
		processedEntry = e
	})

	// Verify the tag was added
	if len(processedEntry.Tags) != 1 || processedEntry.Tags[0] != "test-tag" {
		t.Errorf("Expected tag 'test-tag' to be added, got: %v", processedEntry.Tags)
	}
}

// TestMiddlewareChain tests the middleware chain functionality
func TestMiddlewareChain(t *testing.T) {
	chain := NewMiddlewareChain()

	// Add middleware that adds tags
	middleware1 := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Tags = append(entry.Tags, "middleware1")
		next(entry)
	}

	middleware2 := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Tags = append(entry.Tags, "middleware2")
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
		Tags:      make([]string, 0),
	}

	// Process through chain
	var processedEntry *LogEntry
	chain.Process(entry, func(e *LogEntry) {
		processedEntry = e
	})

	// Verify both tags were added in order
	expectedTags := []string{"middleware1", "middleware2"}
	if len(processedEntry.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(processedEntry.Tags))
	}

	for i, expectedTag := range expectedTags {
		if processedEntry.Tags[i] != expectedTag {
			t.Errorf("Expected tag[%d] to be '%s', got '%s'", i, expectedTag, processedEntry.Tags[i])
		}
	}
}

// TestMiddlewareChainClear tests the Clear functionality
func TestMiddlewareChainClear(t *testing.T) {
	chain := NewMiddlewareChain()

	// Add a middleware
	middleware := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Tags = append(entry.Tags, "test")
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
		Tags:      make([]string, 0),
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
		entry.Tags = append(entry.Tags, "global1")
		next(entry)
	}

	middleware2 := func(entry *LogEntry, next func(*LogEntry)) {
		entry.Tags = append(entry.Tags, "global2")
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
		Tags:      make([]string, 0),
	}

	var processedEntry *LogEntry
	processWithGlobalMiddleware(entry, func(e *LogEntry) {
		processedEntry = e
	})

	// Verify both middleware were applied
	expectedTags := []string{"global1", "global2"}
	if len(processedEntry.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(processedEntry.Tags))
	}

	for i, expectedTag := range expectedTags {
		if processedEntry.Tags[i] != expectedTag {
			t.Errorf("Expected tag[%d] to be '%s', got '%s'", i, expectedTag, processedEntry.Tags[i])
		}
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

	// Add middleware that adds a tag
	AddMiddleware(func(entry *LogEntry, next func(*LogEntry)) {
		entry.Tags = append(entry.Tags, "middleware-tag")
		next(entry)
	})

	// Log a message
	logger.Infof("test message")

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

	// Check that the middleware tag was added
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

	tags, ok := firstLine["tags"].([]interface{})
	if !ok {
		t.Fatal("Expected tags array in log entry")
	}

	if len(tags) != 1 || tags[0] != "middleware-tag" {
		t.Errorf("Expected middleware tag to be added, got: %v", tags)
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
