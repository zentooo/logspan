package logger

import (
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
