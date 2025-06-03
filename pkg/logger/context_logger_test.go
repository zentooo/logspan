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

func TestContextLogger_AutoFlushOnMaxEntries(t *testing.T) {
	// Save original config
	originalConfig := GetConfig()
	defer func() {
		Init(originalConfig)
	}()

	// Set up config with small max entries for testing
	testConfig := DefaultConfig()
	testConfig.MaxLogEntries = 3
	Init(testConfig)

	var buf bytes.Buffer
	logger := NewContextLogger()
	logger.SetOutput(&buf)

	// Add context field
	logger.AddContextValue("test_id", "auto_flush_test")

	// Add entries up to the limit - should not flush yet
	logger.Infof("Message 1")
	logger.Infof("Message 2")

	// Buffer should be empty (no auto-flush yet)
	if buf.Len() > 0 {
		t.Error("Expected no output before reaching max entries")
	}

	// Add third entry - should trigger auto-flush
	logger.Infof("Message 3")

	// Buffer should now contain the flushed output
	output := buf.String()
	if len(output) == 0 {
		t.Error("Expected output after reaching max entries")
	}

	// Verify all three messages are in the output
	if !strings.Contains(output, "Message 1") {
		t.Error("Expected 'Message 1' in auto-flushed output")
	}
	if !strings.Contains(output, "Message 2") {
		t.Error("Expected 'Message 2' in auto-flushed output")
	}
	if !strings.Contains(output, "Message 3") {
		t.Error("Expected 'Message 3' in auto-flushed output")
	}

	// Verify context field is included
	if !strings.Contains(output, "test_id") {
		t.Error("Expected context field in auto-flushed output")
	}

	// Reset buffer and add more entries to test that logger continues working
	buf.Reset()
	logger.Infof("Message 4")
	logger.Infof("Message 5")

	// Should not flush yet
	if buf.Len() > 0 {
		t.Error("Expected no output before reaching max entries again")
	}

	// Add third entry to trigger another auto-flush
	logger.Infof("Message 6")

	// Should have new output
	output2 := buf.String()
	if len(output2) == 0 {
		t.Error("Expected output after reaching max entries again")
	}

	// Verify new messages are in the second output
	if !strings.Contains(output2, "Message 4") {
		t.Error("Expected 'Message 4' in second auto-flushed output")
	}
	if !strings.Contains(output2, "Message 6") {
		t.Error("Expected 'Message 6' in second auto-flushed output")
	}
}

func TestContextLogger_NoAutoFlushWhenMaxEntriesZero(t *testing.T) {
	// Save original config
	originalConfig := GetConfig()
	defer func() {
		Init(originalConfig)
	}()

	// Set up config with MaxLogEntries = 0 (no limit)
	testConfig := DefaultConfig()
	testConfig.MaxLogEntries = 0
	Init(testConfig)

	var buf bytes.Buffer
	logger := NewContextLogger()
	logger.SetOutput(&buf)

	// Add many entries - should not auto-flush
	for i := 0; i < 10; i++ {
		logger.Infof("Message %d", i+1)
	}

	// Buffer should be empty (no auto-flush)
	if buf.Len() > 0 {
		t.Error("Expected no auto-flush when MaxLogEntries is 0")
	}

	// Manual flush should work
	logger.Flush()
	output := buf.String()
	if len(output) == 0 {
		t.Error("Expected output after manual flush")
	}

	// Verify all messages are in the output
	if !strings.Contains(output, "Message 1") {
		t.Error("Expected 'Message 1' in manually flushed output")
	}
	if !strings.Contains(output, "Message 10") {
		t.Error("Expected 'Message 10' in manually flushed output")
	}
}

func TestContextLogger_AutoFlushWithLevelFiltering(t *testing.T) {
	// Save original config
	originalConfig := GetConfig()
	defer func() {
		Init(originalConfig)
	}()

	// Set up config with small max entries for testing
	testConfig := DefaultConfig()
	testConfig.MaxLogEntries = 2
	Init(testConfig)

	var buf bytes.Buffer
	logger := NewContextLogger()
	logger.SetOutput(&buf)
	logger.SetLevel(WarnLevel) // Only WARN and above

	// Add entries that will be filtered out - should not count toward limit
	logger.Debugf("Debug message 1")
	logger.Infof("Info message 1")
	logger.Debugf("Debug message 2")

	// Buffer should be empty (no entries added due to filtering)
	if buf.Len() > 0 {
		t.Error("Expected no output when all entries are filtered")
	}

	// Add entries that pass the filter
	logger.Warnf("Warning message 1")
	logger.Errorf("Error message 1")

	// Should trigger auto-flush (2 entries that passed filter)
	output := buf.String()
	if len(output) == 0 {
		t.Error("Expected output after reaching max entries with filtering")
	}

	// Verify only the non-filtered messages are in the output
	if strings.Contains(output, "Debug message") {
		t.Error("Debug messages should be filtered out")
	}
	if strings.Contains(output, "Info message") {
		t.Error("Info messages should be filtered out")
	}
	if !strings.Contains(output, "Warning message 1") {
		t.Error("Expected 'Warning message 1' in output")
	}
	if !strings.Contains(output, "Error message 1") {
		t.Error("Expected 'Error message 1' in output")
	}
}
