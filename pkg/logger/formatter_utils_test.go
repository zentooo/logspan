package logger

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/zentooo/logspan/pkg/formatter"
)

func TestFormatterUtils_CreateDefaultFormatter(t *testing.T) {
	// Test with PrettifyJSON = false
	Init(Config{
		PrettifyJSON: false,
	})

	formatter := createDefaultFormatter()
	if formatter == nil {
		t.Error("Expected formatter to be created")
	}

	// Test with PrettifyJSON = true
	Init(Config{
		PrettifyJSON: true,
	})

	formatterPretty := createDefaultFormatter()
	if formatterPretty == nil {
		t.Error("Expected pretty formatter to be created")
	}

	// The formatters should be different types (though we can't easily test this without reflection)
	// We'll test the output format instead
	entries := []*LogEntry{
		{
			Timestamp: time.Now(),
			Level:     "INFO",
			Message:   "test message",
		},
	}

	output1, err1 := formatLogOutput(entries, map[string]interface{}{}, time.Now(), time.Now(), formatter)
	if err1 != nil {
		t.Errorf("Expected no error, got %v", err1)
	}

	output2, err2 := formatLogOutput(entries, map[string]interface{}{}, time.Now(), time.Now(), formatterPretty)
	if err2 != nil {
		t.Errorf("Expected no error, got %v", err2)
	}

	// Pretty formatted output should be longer (has indentation)
	if len(output2) <= len(output1) {
		t.Error("Expected pretty formatted output to be longer than compact output")
	}
}

func TestFormatterUtils_FormatLogOutput(t *testing.T) {
	// Test basic formatting
	startTime := time.Date(2023, 10, 27, 10, 0, 0, 0, time.UTC)
	endTime := startTime.Add(100 * time.Millisecond)

	entries := []*LogEntry{
		{
			Timestamp: startTime.Add(10 * time.Millisecond),
			Level:     "INFO",
			Message:   "first message",
			Funcname:  "TestFunc",
			Filename:  "test.go",
			Fileline:  10,
		},
		{
			Timestamp: startTime.Add(50 * time.Millisecond),
			Level:     "ERROR",
			Message:   "error message",
			Funcname:  "ErrorFunc",
			Filename:  "error.go",
			Fileline:  20,
		},
	}

	contextFields := map[string]interface{}{
		"request_id": "req-123",
		"user_id":    "user-456",
	}

	output, err := formatLogOutput(entries, contextFields, startTime, endTime, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Parse the JSON output to verify structure
	var logOutput map[string]interface{}
	if err := json.Unmarshal(output, &logOutput); err != nil {
		t.Errorf("Failed to parse JSON output: %v", err)
	}

	// Verify top-level structure
	if logOutput["type"] != "request" {
		t.Errorf("Expected type to be 'request', got %v", logOutput["type"])
	}

	// Verify context
	context, ok := logOutput["context"].(map[string]interface{})
	if !ok {
		t.Error("Expected context to be a map")
	}
	if context["request_id"] != "req-123" {
		t.Errorf("Expected request_id to be 'req-123', got %v", context["request_id"])
	}
	if context["user_id"] != "user-456" {
		t.Errorf("Expected user_id to be 'user-456', got %v", context["user_id"])
	}

	// Verify runtime
	runtime, ok := logOutput["runtime"].(map[string]interface{})
	if !ok {
		t.Error("Expected runtime to be a map")
	}
	if runtime["severity"] != "ERROR" {
		t.Errorf("Expected severity to be 'ERROR', got %v", runtime["severity"])
	}
	if runtime["elapsed"] != float64(100) {
		t.Errorf("Expected elapsed to be 100, got %v", runtime["elapsed"])
	}

	// Verify lines
	lines, ok := runtime["lines"].([]interface{})
	if !ok {
		t.Error("Expected lines to be an array")
	}
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}

	// Verify first line
	firstLine, ok := lines[0].(map[string]interface{})
	if !ok {
		t.Error("Expected first line to be a map")
	}
	if firstLine["level"] != "INFO" {
		t.Errorf("Expected first line level to be 'INFO', got %v", firstLine["level"])
	}
	if firstLine["message"] != "first message" {
		t.Errorf("Expected first line message to be 'first message', got %v", firstLine["message"])
	}
}

func TestFormatterUtils_FormatLogOutput_WithCustomFormatter(t *testing.T) {
	// Test with custom formatter
	customFormatter := formatter.NewJSONFormatterWithIndent("    ") // 4 spaces

	entries := []*LogEntry{
		{
			Timestamp: time.Now(),
			Level:     "INFO",
			Message:   "test message",
		},
	}

	output, err := formatLogOutput(entries, map[string]interface{}{}, time.Now(), time.Now(), customFormatter)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify it's valid JSON
	var logOutput map[string]interface{}
	if err := json.Unmarshal(output, &logOutput); err != nil {
		t.Errorf("Failed to parse JSON output: %v", err)
	}

	// Verify it contains indentation (4 spaces)
	if !bytes.Contains(output, []byte("    ")) {
		t.Error("Expected output to contain 4-space indentation")
	}
}

func TestFormatterUtils_FormatLogOutput_SeverityCalculation(t *testing.T) {
	// Test severity calculation with different log levels
	testCases := []struct {
		name             string
		levels           []string
		expectedSeverity string
	}{
		{
			name:             "Single DEBUG",
			levels:           []string{"DEBUG"},
			expectedSeverity: "DEBUG",
		},
		{
			name:             "DEBUG and INFO",
			levels:           []string{"DEBUG", "INFO"},
			expectedSeverity: "INFO",
		},
		{
			name:             "INFO, WARN, ERROR",
			levels:           []string{"INFO", "WARN", "ERROR"},
			expectedSeverity: "ERROR",
		},
		{
			name:             "All levels",
			levels:           []string{"DEBUG", "INFO", "WARN", "ERROR", "CRITICAL"},
			expectedSeverity: "CRITICAL",
		},
		{
			name:             "Mixed order",
			levels:           []string{"ERROR", "DEBUG", "CRITICAL", "INFO"},
			expectedSeverity: "CRITICAL",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entries := make([]*LogEntry, len(tc.levels))
			for i, level := range tc.levels {
				entries[i] = &LogEntry{
					Timestamp: time.Now(),
					Level:     level,
					Message:   "test message",
				}
			}

			output, err := formatLogOutput(entries, map[string]interface{}{}, time.Now(), time.Now(), nil)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			var logOutput map[string]interface{}
			if err := json.Unmarshal(output, &logOutput); err != nil {
				t.Errorf("Failed to parse JSON output: %v", err)
			}

			runtime := logOutput["runtime"].(map[string]interface{})
			if runtime["severity"] != tc.expectedSeverity {
				t.Errorf("Expected severity to be '%s', got %v", tc.expectedSeverity, runtime["severity"])
			}
		})
	}
}

func TestFormatterUtils_FormatLogOutput_EmptyEntries(t *testing.T) {
	// Test with empty entries
	output, err := formatLogOutput([]*LogEntry{}, map[string]interface{}{}, time.Now(), time.Now(), nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	var logOutput map[string]interface{}
	if err := json.Unmarshal(output, &logOutput); err != nil {
		t.Errorf("Failed to parse JSON output: %v", err)
	}

	runtime := logOutput["runtime"].(map[string]interface{})
	lines := runtime["lines"].([]interface{})
	if len(lines) != 0 {
		t.Errorf("Expected 0 lines, got %d", len(lines))
	}

	// With empty entries, severity should be DEBUG (default)
	if runtime["severity"] != "DEBUG" {
		t.Errorf("Expected severity to be 'DEBUG', got %v", runtime["severity"])
	}
}

func TestFormatterUtils_FormatLogOutput_LogType(t *testing.T) {
	// Save original state
	originalConfig := GetConfig()
	originalInitialized := IsInitialized()

	testCases := []struct {
		name         string
		logType      string
		expectedType string
	}{
		{
			name:         "Custom log type",
			logType:      "custom_log_type",
			expectedType: "custom_log_type",
		},
		{
			name:         "Empty log type defaults to request",
			logType:      "",
			expectedType: "request",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			customConfig := Config{
				MinLevel:         InfoLevel,
				Output:           os.Stdout,
				EnableSourceInfo: false,
				PrettifyJSON:     false,
				MaxLogEntries:    0,
				LogType:          tc.logType,
				ErrorHandler:     nil,
			}

			Init(customConfig)

			entries := []*LogEntry{
				{
					Timestamp: time.Now(),
					Level:     "INFO",
					Message:   "test message",
				},
			}

			output, err := formatLogOutput(entries, map[string]interface{}{}, time.Now(), time.Now(), nil)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			var logOutput map[string]interface{}
			if err := json.Unmarshal(output, &logOutput); err != nil {
				t.Errorf("Failed to parse JSON output: %v", err)
			}

			// Verify that expected log type is used
			if logOutput["type"] != tc.expectedType {
				t.Errorf("Expected type to be '%s', got %v", tc.expectedType, logOutput["type"])
			}
		})
	}

	// Restore original state
	if originalInitialized {
		Init(originalConfig)
	}
}
