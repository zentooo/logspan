package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestAddContextValues(t *testing.T) {
	ctx := context.Background()
	contextLogger := NewContextLogger()
	ctx = WithLogger(ctx, contextLogger)

	// Test adding multiple context values
	fields := map[string]interface{}{
		"user_id":    "user123",
		"request_id": "req456",
		"session":    "sess789",
	}

	AddContextValues(ctx, fields)

	// Verify the values were added
	if contextLogger.fields["user_id"] != "user123" {
		t.Errorf("Expected user_id to be 'user123', got %v", contextLogger.fields["user_id"])
	}
	if contextLogger.fields["request_id"] != "req456" {
		t.Errorf("Expected request_id to be 'req456', got %v", contextLogger.fields["request_id"])
	}
	if contextLogger.fields["session"] != "sess789" {
		t.Errorf("Expected session to be 'sess789', got %v", contextLogger.fields["session"])
	}
}

func TestDebugf(t *testing.T) {
	var buf bytes.Buffer
	ctx := context.Background()
	contextLogger := NewContextLogger()
	contextLogger.SetOutput(&buf)
	contextLogger.SetLevel(DebugLevel) // Enable debug level
	ctx = WithLogger(ctx, contextLogger)

	// Test debug logging
	Debugf(ctx, "Debug message: %s", "test")

	// Flush to get output
	FlushContext(ctx)

	output := buf.String()
	if !strings.Contains(output, "Debug message: test") {
		t.Errorf("Expected output to contain 'Debug message: test', got: %s", output)
	}
	if !strings.Contains(output, "DEBUG") {
		t.Errorf("Expected output to contain 'DEBUG', got: %s", output)
	}
}

func TestWarnf(t *testing.T) {
	var buf bytes.Buffer
	ctx := context.Background()
	contextLogger := NewContextLogger()
	contextLogger.SetOutput(&buf)
	contextLogger.SetLevel(DebugLevel) // Enable all levels
	ctx = WithLogger(ctx, contextLogger)

	// Test warning logging
	Warnf(ctx, "Warning message: %s", "test")

	// Flush to get output
	FlushContext(ctx)

	output := buf.String()
	if !strings.Contains(output, "Warning message: test") {
		t.Errorf("Expected output to contain 'Warning message: test', got: %s", output)
	}
	if !strings.Contains(output, "WARN") {
		t.Errorf("Expected output to contain 'WARN', got: %s", output)
	}
}

func TestCriticalf(t *testing.T) {
	var buf bytes.Buffer
	ctx := context.Background()
	contextLogger := NewContextLogger()
	contextLogger.SetOutput(&buf)
	contextLogger.SetLevel(DebugLevel) // Enable all levels
	ctx = WithLogger(ctx, contextLogger)

	// Test critical logging
	Criticalf(ctx, "Critical message: %s", "test")

	// Flush to get output
	FlushContext(ctx)

	output := buf.String()
	if !strings.Contains(output, "Critical message: test") {
		t.Errorf("Expected output to contain 'Critical message: test', got: %s", output)
	}
	if !strings.Contains(output, "CRITICAL") {
		t.Errorf("Expected output to contain 'CRITICAL', got: %s", output)
	}
}

func TestContextAPI_AllLevels(t *testing.T) {
	var buf bytes.Buffer
	ctx := context.Background()
	contextLogger := NewContextLogger()
	contextLogger.SetOutput(&buf)
	contextLogger.SetLevel(DebugLevel) // Enable all levels
	ctx = WithLogger(ctx, contextLogger)

	// Add some context
	AddContextValue(ctx, "test_id", "123")

	// Test all logging levels
	Debugf(ctx, "Debug: %d", 1)
	Infof(ctx, "Info: %d", 2)
	Warnf(ctx, "Warn: %d", 3)
	Errorf(ctx, "Error: %d", 4)
	Criticalf(ctx, "Critical: %d", 5)

	// Flush to get output
	FlushContext(ctx)

	output := buf.String()

	// Parse JSON to verify structure
	var logData map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logData); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	// Check context
	context, ok := logData["context"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected context field to be present")
	}
	if context["test_id"] != "123" {
		t.Errorf("Expected test_id to be '123', got %v", context["test_id"])
	}

	// Check runtime
	runtime, ok := logData["runtime"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected runtime field to be present")
	}

	// Check severity (should be CRITICAL as it's the highest)
	if runtime["severity"] != "CRITICAL" {
		t.Errorf("Expected severity to be 'CRITICAL', got %v", runtime["severity"])
	}

	// Check lines
	lines, ok := runtime["lines"].([]interface{})
	if !ok {
		t.Fatal("Expected lines field to be present")
	}
	if len(lines) != 5 {
		t.Errorf("Expected 5 log lines, got %d", len(lines))
	}

	// Verify each line
	expectedMessages := []string{"Debug: 1", "Info: 2", "Warn: 3", "Error: 4", "Critical: 5"}
	expectedLevels := []string{"DEBUG", "INFO", "WARN", "ERROR", "CRITICAL"}

	for i, line := range lines {
		lineMap, ok := line.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected line %d to be a map", i)
		}

		if lineMap["message"] != expectedMessages[i] {
			t.Errorf("Expected line %d message to be '%s', got %v", i, expectedMessages[i], lineMap["message"])
		}
		if lineMap["level"] != expectedLevels[i] {
			t.Errorf("Expected line %d level to be '%s', got %v", i, expectedLevels[i], lineMap["level"])
		}
	}
}

func TestContextAPI_WithoutLoggerInContext(t *testing.T) {
	ctx := context.Background() // No logger attached

	// This should create a new logger automatically
	Infof(ctx, "Test message")

	// Since we don't have access to the auto-created logger, we can't easily test the output
	// But we can verify that the function doesn't panic
	// The function should work without error even when no logger is in context
}

func TestAddContextValues_EmptyMap(t *testing.T) {
	ctx := context.Background()
	contextLogger := NewContextLogger()
	ctx = WithLogger(ctx, contextLogger)

	// Test adding empty map
	fields := map[string]interface{}{}
	AddContextValues(ctx, fields)

	// Should not cause any issues
	if len(contextLogger.fields) != 0 {
		t.Errorf("Expected no context fields, got %d", len(contextLogger.fields))
	}
}

func TestAddContextValues_NilMap(t *testing.T) {
	ctx := context.Background()
	contextLogger := NewContextLogger()
	ctx = WithLogger(ctx, contextLogger)

	// Test adding nil map
	AddContextValues(ctx, nil)

	// Should not cause any issues
	if len(contextLogger.fields) != 0 {
		t.Errorf("Expected no context fields, got %d", len(contextLogger.fields))
	}
}
