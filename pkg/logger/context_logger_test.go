package logger

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestContextLogger_BasicLogging(t *testing.T) {
	var buf bytes.Buffer
	logger := NewContextLogger()
	logger.SetOutput(&buf)

	// Add some log entries
	logger.Infof("First message")
	logger.Warnf("Warning message")
	logger.Errorf("Error message")

	// Flush should output all entries
	logger.Flush()

	output := buf.String()
	if !strings.Contains(output, "First message") {
		t.Error("Expected 'First message' in output")
	}
	if !strings.Contains(output, "Warning message") {
		t.Error("Expected 'Warning message' in output")
	}
	if !strings.Contains(output, "Error message") {
		t.Error("Expected 'Error message' in output")
	}
}

func TestContextLogger_AddFields(t *testing.T) {
	var buf bytes.Buffer
	logger := NewContextLogger()
	logger.SetOutput(&buf)

	// Add context fields
	logger.AddContextValue("request_id", "12345")
	logger.AddContextValue("user_id", "user123")

	// Add multiple fields
	logger.AddContextValues(map[string]interface{}{
		"session_id": "session456",
		"ip_address": "192.168.1.1",
	})

	logger.Infof("Test message")
	logger.Flush()

	output := buf.String()
	// Context fields should appear in the context section, not in individual log entries
	if !strings.Contains(output, "request_id") {
		t.Error("Expected 'request_id' in output")
	}
	if !strings.Contains(output, "12345") {
		t.Error("Expected '12345' in output")
	}
	if !strings.Contains(output, "session_id") {
		t.Error("Expected 'session_id' in output")
	}
}

func TestContextLogger_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := NewContextLogger()
	logger.SetOutput(&buf)
	logger.SetLevel(WarnLevel) // Only WARN and above

	logger.Debugf("Debug message")
	logger.Infof("Info message")
	logger.Warnf("Warning message")
	logger.Errorf("Error message")

	logger.Flush()

	output := buf.String()
	if strings.Contains(output, "Debug message") {
		t.Error("Debug message should be filtered out")
	}
	if strings.Contains(output, "Info message") {
		t.Error("Info message should be filtered out")
	}
	if !strings.Contains(output, "Warning message") {
		t.Error("Warning message should be included")
	}
	if !strings.Contains(output, "Error message") {
		t.Error("Error message should be included")
	}
}

func TestContextLogger_EmptyFlush(t *testing.T) {
	var buf bytes.Buffer
	logger := NewContextLogger()
	logger.SetOutput(&buf)

	// Flush without any log entries
	logger.Flush()

	if buf.Len() > 0 {
		t.Error("Expected no output when flushing empty logger")
	}
}

func TestContextLogger_NilOutput(t *testing.T) {
	logger := NewContextLogger()
	logger.SetOutput(nil)

	// Should not panic
	logger.Infof("Test message")
	logger.Flush()
}

func TestContextLogger_SeverityCalculation(t *testing.T) {
	var buf bytes.Buffer
	logger := NewContextLogger()
	logger.SetOutput(&buf)

	logger.Debugf("Debug message")
	logger.Infof("Info message")
	logger.Criticalf("Critical message")
	logger.Warnf("Warning message")

	logger.Flush()

	output := buf.String()
	if !strings.Contains(output, "CRITICAL") {
		t.Error("Expected highest severity to be CRITICAL")
	}
}

func TestWithLogger_AndFromContext(t *testing.T) {
	ctx := context.Background()
	logger := NewContextLogger()

	// Add logger to context
	ctx = WithLogger(ctx, logger)

	// Retrieve logger from context
	retrievedLogger := FromContext(ctx)

	if retrievedLogger != logger {
		t.Error("Expected to retrieve the same logger instance")
	}
}

func TestFromContext_WithoutLogger(t *testing.T) {
	ctx := context.Background()

	// Should return a new logger when none exists in context
	logger := FromContext(ctx)

	if logger == nil {
		t.Error("Expected a new logger to be created")
	}
}

func TestContextAPI_Functions(t *testing.T) {
	var buf bytes.Buffer
	ctx := context.Background()
	logger := NewContextLogger()
	logger.SetOutput(&buf)

	ctx = WithLogger(ctx, logger)

	// Test context API functions
	AddContextValue(ctx, "test_key", "test_value")
	Infof(ctx, "Test message")
	Errorf(ctx, "Error message")

	FlushContext(ctx)

	output := buf.String()
	if !strings.Contains(output, "Test message") {
		t.Error("Expected 'Test message' in output")
	}
	if !strings.Contains(output, "Error message") {
		t.Error("Expected 'Error message' in output")
	}
	if !strings.Contains(output, "test_key") {
		t.Error("Expected 'test_key' in output")
	}
}

func TestIsHigherSeverity(t *testing.T) {
	tests := []struct {
		level1   string
		level2   string
		expected bool
	}{
		{"ERROR", "INFO", true},
		{"CRITICAL", "ERROR", true},
		{"DEBUG", "INFO", false},
		{"INFO", "INFO", false},
		{"WARN", "DEBUG", true},
	}

	for _, test := range tests {
		result := isHigherSeverity(test.level1, test.level2)
		if result != test.expected {
			t.Errorf("isHigherSeverity(%s, %s) = %v, expected %v",
				test.level1, test.level2, result, test.expected)
		}
	}
}
