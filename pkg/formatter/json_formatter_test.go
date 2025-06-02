package formatter

import (
	"encoding/json"
	"testing"
	"time"
)

func TestJSONFormatter_Format(t *testing.T) {
	formatter := NewJSONFormatter()

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
			Tags:      map[string]int{"tag1": 1, "tag2": 1},
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

	// Verify basic structure
	if parsed["type"] != "request" {
		t.Errorf("Expected type 'request', got %v", parsed["type"])
	}

	context, ok := parsed["context"].(map[string]interface{})
	if !ok {
		t.Fatalf("Context is not a map")
	}
	if context["REQUEST_ID"] != "test-request-id" {
		t.Errorf("Expected REQUEST_ID 'test-request-id', got %v", context["REQUEST_ID"])
	}

	runtime, ok := parsed["runtime"].(map[string]interface{})
	if !ok {
		t.Fatalf("Runtime is not a map")
	}
	if runtime["severity"] != "INFO" {
		t.Errorf("Expected severity 'INFO', got %v", runtime["severity"])
	}
}

func TestJSONFormatter_FormatWithIndent(t *testing.T) {
	formatter := NewJSONFormatterWithIndent("  ")

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

func TestJSONFormatter_FormatCompact(t *testing.T) {
	formatter := NewJSONFormatter()

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

	// Verify that the result is compact (no unnecessary whitespace)
	resultStr := string(result)
	if containsIndentation(resultStr) {
		t.Errorf("Expected compact JSON, got: %s", resultStr)
	}

	// Verify that the result is valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal(result, &parsed); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}
}

func TestJSONFormatter_FormatEmptyOutput(t *testing.T) {
	formatter := NewJSONFormatter()

	// Create empty output
	output := &LogOutput{}

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
}

// containsIndentation checks if the string contains indentation (multiple spaces or tabs)
func containsIndentation(s string) bool {
	// Check for newlines followed by spaces (indented JSON)
	for i := 0; i < len(s)-1; i++ {
		if s[i] == '\n' && (s[i+1] == ' ' || s[i+1] == '\t') {
			return true
		}
	}
	return false
}
