package formatter

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestNewContextFlattenFormatter(t *testing.T) {
	formatter := NewContextFlattenFormatter()
	if formatter == nil {
		t.Fatal("NewContextFlattenFormatter() returned nil")
	}
	if formatter.Indent != "" {
		t.Errorf("Expected empty indent, got %q", formatter.Indent)
	}
}

func TestNewContextFlattenFormatterWithIndent(t *testing.T) {
	indent := "  "
	formatter := NewContextFlattenFormatterWithIndent(indent)
	if formatter == nil {
		t.Fatal("NewContextFlattenFormatterWithIndent() returned nil")
	}
	if formatter.Indent != indent {
		t.Errorf("Expected indent %q, got %q", indent, formatter.Indent)
	}
}

func TestContextFlattenFormatter_Format_Basic(t *testing.T) {
	formatter := NewContextFlattenFormatter()

	output := &LogOutput{
		Type: "request",
		Context: map[string]interface{}{
			"request_id": "req-123",
			"user_id":    "user-456",
		},
		Runtime: RuntimeInfo{
			Severity:  "INFO",
			StartTime: "2023-10-27T09:59:58.123456+09:00",
			EndTime:   "2023-10-27T10:00:00.223456+09:00",
			Elapsed:   2100,
			Lines: []*LogEntry{
				{
					Timestamp: time.Date(2023, 10, 27, 9, 59, 59, 123456000, time.FixedZone("JST", 9*60*60)),
					Level:     "INFO",
					Message:   "Request processing started",
				},
			},
		},
	}

	result, err := formatter.Format(output)
	if err != nil {
		t.Fatalf("Format() error = %v", err)
	}

	// Parse the result to verify structure
	var parsed map[string]interface{}
	if err := json.Unmarshal(result, &parsed); err != nil {
		t.Fatalf("Failed to parse result JSON: %v", err)
	}

	// Check that context fields are at top level
	if parsed["user_id"] != "user-456" {
		t.Errorf("Expected user_id=user-456, got %v", parsed["user_id"])
	}
	if parsed["request_id"] != "req-123" {
		t.Errorf("Expected request_id=req-123, got %v", parsed["request_id"])
	}

	// Check that other fields are present
	if parsed["type"] != "request" {
		t.Errorf("Expected type=request, got %v", parsed["type"])
	}
	if parsed["runtime"] == nil {
		t.Error("Expected runtime field to be present")
	}
}

func TestContextFlattenFormatter_Format_WithIndent(t *testing.T) {
	formatter := NewContextFlattenFormatterWithIndent("  ")

	output := &LogOutput{
		Type: "request",
		Context: map[string]interface{}{
			"request_id": "req-123",
			"user_id":    "user-456",
		},
		Runtime: RuntimeInfo{
			Severity:  "INFO",
			StartTime: "2023-10-27T09:59:58.123456+09:00",
			EndTime:   "2023-10-27T10:00:00.223456+09:00",
			Elapsed:   2100,
			Lines: []*LogEntry{
				{
					Timestamp: time.Date(2023, 10, 27, 9, 59, 59, 123456000, time.FixedZone("JST", 9*60*60)),
					Level:     "INFO",
					Message:   "Request processing started",
				},
			},
		},
	}

	result, err := formatter.Format(output)
	if err != nil {
		t.Fatalf("Format() error = %v", err)
	}

	// Check that result contains indentation
	resultStr := string(result)
	if !strings.Contains(resultStr, "\n") || !strings.Contains(resultStr, "  ") {
		t.Error("Expected indented JSON output")
	}
}

func TestContextFlattenFormatter_Format_EmptyContext(t *testing.T) {
	formatter := NewContextFlattenFormatter()

	output := &LogOutput{
		Type:    "request",
		Context: map[string]interface{}{},
		Runtime: RuntimeInfo{
			Severity:  "INFO",
			StartTime: "2023-10-27T09:59:58.123456+09:00",
			EndTime:   "2023-10-27T10:00:00.223456+09:00",
			Elapsed:   0,
			Lines:     []*LogEntry{},
		},
	}

	result, err := formatter.Format(output)
	if err != nil {
		t.Fatalf("Format() error = %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(result, &parsed); err != nil {
		t.Fatalf("Failed to parse result JSON: %v", err)
	}

	// Should still have the basic fields
	if parsed["type"] != "request" {
		t.Errorf("Expected type=request, got %v", parsed["type"])
	}
	if parsed["runtime"] == nil {
		t.Error("Expected runtime field to be present")
	}
}

func TestContextFlattenFormatter_Format_KeyCollision(t *testing.T) {
	formatter := NewContextFlattenFormatter()

	output := &LogOutput{
		Type: "request",
		Context: map[string]interface{}{
			"request_id": "req-123",
			"user_id":    "user-456",
			"type":       "custom_type", // This should override the default type
		},
		Runtime: RuntimeInfo{
			Severity:  "INFO",
			StartTime: "2023-10-27T09:59:58.123456+09:00",
			EndTime:   "2023-10-27T10:00:00.223456+09:00",
			Elapsed:   2100,
			Lines: []*LogEntry{
				{
					Timestamp: time.Date(2023, 10, 27, 9, 59, 59, 123456000, time.FixedZone("JST", 9*60*60)),
					Level:     "INFO",
					Message:   "Request processing started",
				},
			},
		},
	}

	result, err := formatter.Format(output)
	if err != nil {
		t.Fatalf("Format() error = %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(result, &parsed); err != nil {
		t.Fatalf("Failed to parse result JSON: %v", err)
	}

	// Context values should take precedence
	if parsed["type"] != "custom_type" {
		t.Errorf("Expected type=custom_type (from context), got %v", parsed["type"])
	}
	if parsed["runtime"] == nil {
		t.Error("Expected runtime field to be present (no collision)")
	}
	if parsed["user_id"] != "user-456" {
		t.Errorf("Expected user_id=user-456, got %v", parsed["user_id"])
	}
}

func TestContextFlattenFormatter_Format_NilContext(t *testing.T) {
	formatter := NewContextFlattenFormatter()

	output := &LogOutput{
		Type:    "request",
		Context: nil, // nil context
		Runtime: RuntimeInfo{
			Severity:  "INFO",
			StartTime: "2023-10-27T09:59:58.123456+09:00",
			EndTime:   "2023-10-27T10:00:00.223456+09:00",
			Elapsed:   2100,
			Lines:     []*LogEntry{},
		},
	}

	result, err := formatter.Format(output)
	if err != nil {
		t.Fatalf("Format() error = %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(result, &parsed); err != nil {
		t.Fatalf("Failed to parse result JSON: %v", err)
	}

	// Should still have the basic fields
	if parsed["type"] != "request" {
		t.Errorf("Expected type=request, got %v", parsed["type"])
	}
	if parsed["runtime"] == nil {
		t.Error("Expected runtime field to be present")
	}
}
