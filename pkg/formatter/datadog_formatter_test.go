package formatter

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDataDogFormatter_Format(t *testing.T) {
	formatter := NewDataDogFormatter()

	entry := &LogEntry{
		Timestamp: time.Date(2023, 10, 27, 10, 0, 0, 0, time.UTC),
		Level:     "INFO",
		Message:   "test message",
	}

	output := &LogOutput{
		Type:    "request",
		Context: map[string]interface{}{"request_id": "123"},
		Runtime: RuntimeInfo{
			Severity:  "INFO",
			StartTime: "2023-10-27T10:00:00Z",
			EndTime:   "2023-10-27T10:00:01Z",
			Elapsed:   1000,
			Lines:     []*LogEntry{entry},
		},
		Config: ConfigInfo{
			ElapsedUnit: "ms",
		},
	}

	result, err := formatter.Format(output)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Verify it's valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal(result, &parsed); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	// Verify DataDog format structure
	if parsed["status"] != "info" {
		t.Errorf("Expected status 'info', got %v", parsed["status"])
	}
	if parsed["logger"] != "logspan" {
		t.Errorf("Expected logger 'logspan', got %v", parsed["logger"])
	}
}

func TestDataDogFormatter_ConvertSeverityToStatus(t *testing.T) {
	formatter := NewDataDogFormatter()

	tests := []struct {
		severity string
		expected string
	}{
		{"DEBUG", "debug"},
		{"INFO", "info"},
		{"WARN", "warn"},
		{"ERROR", "error"},
		{"CRITICAL", "critical"},
		{"UNKNOWN", "info"}, // default case
	}

	for _, test := range tests {
		result := formatter.convertSeverityToStatus(test.severity)
		if result != test.expected {
			t.Errorf("convertSeverityToStatus(%s) = %s, expected %s", test.severity, result, test.expected)
		}
	}
}

func TestDataDogFormatter_CreateSummaryMessage(t *testing.T) {
	formatter := NewDataDogFormatter()

	// Test empty lines
	output := &LogOutput{
		Runtime: RuntimeInfo{
			Lines: []*LogEntry{},
		},
	}
	result := formatter.createSummaryMessage(output)
	if result != "Empty log context" {
		t.Errorf("Expected 'Empty log context', got %s", result)
	}

	// Test single line
	output.Runtime.Lines = []*LogEntry{
		{Message: "single message"},
	}
	result = formatter.createSummaryMessage(output)
	if result != "single message" {
		t.Errorf("Expected 'single message', got %s", result)
	}

	// Test multiple lines
	output.Runtime.Lines = []*LogEntry{
		{Message: "first message"},
		{Message: "second message"},
	}
	result = formatter.createSummaryMessage(output)
	if result != "Log context with multiple entries" {
		t.Errorf("Expected 'Log context with multiple entries', got %s", result)
	}
}

func TestDataDogFormatter_FormatWithIndent(t *testing.T) {
	formatter := NewDataDogFormatterWithIndent("  ")

	// Create simple test output
	output := &LogOutput{
		Type:    "request",
		Context: map[string]interface{}{},
		Runtime: RuntimeInfo{
			Severity:  "INFO",
			StartTime: "2023-10-27T09:59:58.123456Z",
			EndTime:   "2023-10-27T10:00:00.223456Z",
			Elapsed:   0,
			Lines:     []*LogEntry{},
		},
		Config: ConfigInfo{
			ElapsedUnit: "ms",
		},
	}

	// Format the output
	result, err := formatter.Format(output)
	if err != nil {
		t.Fatalf("Format() error = %v", err)
	}

	// Verify that the result contains indentation
	resultStr := string(result)
	if !containsIndentation(resultStr) {
		t.Errorf("Expected indented JSON, got: %s", resultStr)
	}

	// Verify that the result is still valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal(result, &parsed); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}
}
