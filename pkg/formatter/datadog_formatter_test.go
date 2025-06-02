package formatter

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDataDogFormatter_Format(t *testing.T) {
	formatter := NewDataDogFormatter()

	// Create test log entry
	entry := &LogEntry{
		Timestamp: time.Date(2023, 10, 27, 9, 59, 59, 123456000, time.UTC),
		Level:     "INFO",
		Message:   "test message",
		Fields:    map[string]interface{}{"key": "value"},
		Tags:      []string{"tag1", "tag2"},
	}

	// Create test log output
	output := &LogOutput{
		Type: "request",
		Context: map[string]interface{}{
			"REQUEST_ID": "test-request-id",
		},
		Runtime: RuntimeInfo{
			Severity:  "INFO",
			StartTime: "2023-10-27T09:59:58.123456Z",
			EndTime:   "2023-10-27T10:00:00.223456Z",
			Elapsed:   2100,
			Lines:     []*LogEntry{entry},
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

	// Verify that the result is valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal(result, &parsed); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	// Verify DataDog Standard Attributes
	if parsed["timestamp"] != "2023-10-27T09:59:58.123456Z" {
		t.Errorf("Expected timestamp '2023-10-27T09:59:58.123456Z', got %v", parsed["timestamp"])
	}

	if parsed["status"] != "info" {
		t.Errorf("Expected status 'info', got %v", parsed["status"])
	}

	if parsed["message"] != "test message" {
		t.Errorf("Expected message 'test message', got %v", parsed["message"])
	}

	if parsed["logger"] != "logspan" {
		t.Errorf("Expected logger 'logspan', got %v", parsed["logger"])
	}

	if parsed["duration"] != float64(2100) {
		t.Errorf("Expected duration 2100, got %v", parsed["duration"])
	}

	// Verify context fields are included as custom attributes
	if parsed["REQUEST_ID"] != "test-request-id" {
		t.Errorf("Expected REQUEST_ID 'test-request-id', got %v", parsed["REQUEST_ID"])
	}

	// Verify lines array
	lines, ok := parsed["lines"].([]interface{})
	if !ok {
		t.Fatalf("Lines is not an array")
	}
	if len(lines) != 1 {
		t.Fatalf("Expected 1 line, got %d", len(lines))
	}

	line, ok := lines[0].(map[string]interface{})
	if !ok {
		t.Fatalf("Line is not a map")
	}
	if line["status"] != "info" {
		t.Errorf("Expected line status 'info', got %v", line["status"])
	}
	if line["message"] != "test message" {
		t.Errorf("Expected line message 'test message', got %v", line["message"])
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
