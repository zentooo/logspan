package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/zentooo/logspan/pkg/formatter"
)

func TestDirectLogger_BasicLogging(t *testing.T) {
	// Create a test buffer
	var buf bytes.Buffer
	logger := NewDirectLogger()
	logger.SetOutput(&buf)
	logger.SetLevel(DebugLevel) // Set to output all levels

	// Test each log level
	testCases := []struct {
		name    string
		logFunc func(string, ...interface{})
		level   string
		message string
		args    []interface{}
	}{
		{"Debug", logger.Debugf, "DEBUG", "debug message", nil},
		{"Info", logger.Infof, "INFO", "info message", nil},
		{"Warn", logger.Warnf, "WARN", "warn message", nil},
		{"Error", logger.Errorf, "ERROR", "error message", nil},
		{"Critical", logger.Criticalf, "CRITICAL", "critical message", nil},
		{"WithArgs", logger.Infof, "INFO", "message with %s and %d", []interface{}{"string", 42}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf.Reset()

			if tc.args != nil {
				tc.logFunc(tc.message, tc.args...)
			} else {
				tc.logFunc(tc.message)
			}

			output := buf.String()

			// Verify that log output is generated
			if output == "" {
				t.Error("Expected log output, got empty string")
				return
			}

			// Verify that output can be parsed as JSON
			var logData map[string]interface{}
			if err := json.Unmarshal([]byte(output), &logData); err != nil {
				t.Errorf("Expected valid JSON output, got error: %v, output: %s", err, output)
				return
			}

			// Verify basic structure of structured log
			if logData["type"] != "request" {
				t.Errorf("Expected type 'request', got: %v", logData["type"])
			}

			runtime, ok := logData["runtime"].(map[string]interface{})
			if !ok {
				t.Error("Expected runtime section in log output")
				return
			}

			// Verify that severity is correct
			if runtime["severity"] != tc.level {
				t.Errorf("Expected severity %s, got: %v", tc.level, runtime["severity"])
			}

			// Verify that lines is an array with one entry
			lines, ok := runtime["lines"].([]interface{})
			if !ok {
				t.Error("Expected lines to be an array")
				return
			}

			if len(lines) != 1 {
				t.Errorf("Expected exactly 1 log entry, got %d", len(lines))
				return
			}

			// Verify log entry content
			entry, ok := lines[0].(map[string]interface{})
			if !ok {
				t.Error("Expected log entry to be an object")
				return
			}

			// Verify that message is included
			expectedMessage := tc.message
			if tc.args != nil {
				expectedMessage = "message with string and 42"
			}
			if entry["message"] != expectedMessage {
				t.Errorf("Expected message %s, got: %v", expectedMessage, entry["message"])
			}

			// Verify that level is included
			if entry["level"] != tc.level {
				t.Errorf("Expected level %s, got: %v", tc.level, entry["level"])
			}
		})
	}
}

func TestDirectLogger_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := NewDirectLogger()
	logger.SetOutput(&buf)

	testCases := []struct {
		name         string
		setLevel     LogLevel
		logLevel     LogLevel
		logFunc      func(string, ...interface{})
		message      string
		shouldOutput bool
	}{
		{"DebugLevel_DebugLog", DebugLevel, DebugLevel, logger.Debugf, "debug message", true},
		{"InfoLevel_DebugLog", InfoLevel, DebugLevel, logger.Debugf, "debug message", false},
		{"InfoLevel_InfoLog", InfoLevel, InfoLevel, logger.Infof, "info message", true},
		{"WarnLevel_InfoLog", WarnLevel, InfoLevel, logger.Infof, "info message", false},
		{"WarnLevel_WarnLog", WarnLevel, WarnLevel, logger.Warnf, "warn message", true},
		{"ErrorLevel_WarnLog", ErrorLevel, WarnLevel, logger.Warnf, "warn message", false},
		{"ErrorLevel_ErrorLog", ErrorLevel, ErrorLevel, logger.Errorf, "error message", true},
		{"CriticalLevel_ErrorLog", CriticalLevel, ErrorLevel, logger.Errorf, "error message", false},
		{"CriticalLevel_CriticalLog", CriticalLevel, CriticalLevel, logger.Criticalf, "critical message", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf.Reset()
			logger.SetLevel(tc.setLevel)

			tc.logFunc(tc.message)
			output := buf.String()

			if tc.shouldOutput {
				if output == "" {
					t.Errorf("Expected log output for level %s with min level %s, got empty string",
						tc.logLevel.String(), tc.setLevel.String())
					return
				}

				// Verify that output can be parsed as JSON
				var logData map[string]interface{}
				if err := json.Unmarshal([]byte(output), &logData); err != nil {
					t.Errorf("Expected valid JSON output, got error: %v", err)
					return
				}

				// Verify that message is included
				runtime := logData["runtime"].(map[string]interface{})
				lines := runtime["lines"].([]interface{})
				entry := lines[0].(map[string]interface{})
				if entry["message"] != tc.message {
					t.Errorf("Expected message %s, got: %v", tc.message, entry["message"])
				}
			} else {
				if output != "" {
					t.Errorf("Expected no log output for level %s with min level %s, got: %s",
						tc.logLevel.String(), tc.setLevel.String(), output)
				}
			}
		})
	}
}

func TestDirectLogger_SetLevelFromString(t *testing.T) {
	var buf bytes.Buffer
	logger := NewDirectLogger()
	logger.SetOutput(&buf)

	testCases := []struct {
		levelString   string
		expectedLevel LogLevel
		testLogFunc   func(string, ...interface{})
		shouldOutput  bool
	}{
		{"DEBUG", DebugLevel, logger.Debugf, true},
		{"INFO", InfoLevel, logger.Debugf, false},
		{"WARN", WarnLevel, logger.Infof, false},
		{"ERROR", ErrorLevel, logger.Warnf, false},
		{"CRITICAL", CriticalLevel, logger.Errorf, false},
		{"INVALID", InfoLevel, logger.Debugf, false}, // Default to InfoLevel
	}

	for _, tc := range testCases {
		t.Run(tc.levelString, func(t *testing.T) {
			buf.Reset()
			logger.SetLevelFromString(tc.levelString)

			tc.testLogFunc("test message")
			output := buf.String()

			if tc.shouldOutput {
				if output == "" {
					t.Errorf("Expected log output when setting level to %s, got empty string", tc.levelString)
				}
			} else {
				if output != "" {
					t.Errorf("Expected no log output when setting level to %s, got: %s", tc.levelString, output)
				}
			}
		})
	}
}

func TestDirectLogger_SetOutput(t *testing.T) {
	logger := NewDirectLogger()

	// First buffer
	var buf1 bytes.Buffer
	logger.SetOutput(&buf1)
	logger.Infof("message to buffer 1")

	if buf1.String() == "" {
		t.Error("Expected output in buffer 1, got empty string")
	}

	// Verify that output can be parsed as JSON
	var logData1 map[string]interface{}
	if err := json.Unmarshal([]byte(buf1.String()), &logData1); err != nil {
		t.Errorf("Expected valid JSON output in buffer 1, got error: %v", err)
	} else {
		runtime := logData1["runtime"].(map[string]interface{})
		lines := runtime["lines"].([]interface{})
		entry := lines[0].(map[string]interface{})
		if entry["message"] != "message to buffer 1" {
			t.Errorf("Expected buffer 1 to contain message, got: %v", entry["message"])
		}
	}

	// Switch to second buffer
	var buf2 bytes.Buffer
	logger.SetOutput(&buf2)
	logger.Infof("message to buffer 2")

	// Verify that no new message is added to buf1
	if strings.Contains(buf1.String(), "message to buffer 2") {
		t.Error("Buffer 1 should not contain message sent to buffer 2")
	}

	// Verify that new message is output to buf2
	if buf2.String() == "" {
		t.Error("Expected output in buffer 2, got empty string")
	}

	// Verify that output can be parsed as JSON
	var logData2 map[string]interface{}
	if err := json.Unmarshal([]byte(buf2.String()), &logData2); err != nil {
		t.Errorf("Expected valid JSON output in buffer 2, got error: %v", err)
	} else {
		runtime := logData2["runtime"].(map[string]interface{})
		lines := runtime["lines"].([]interface{})
		entry := lines[0].(map[string]interface{})
		if entry["message"] != "message to buffer 2" {
			t.Errorf("Expected buffer 2 to contain message, got: %v", entry["message"])
		}
	}
}

func TestDirectLogger_ConcurrentSafety(t *testing.T) {
	var buf bytes.Buffer
	logger := NewDirectLogger()
	logger.SetOutput(&buf)
	logger.SetLevel(DebugLevel)

	// Use compact JSON formatter to avoid multi-line JSON issues in concurrent tests
	logger.SetFormatter(formatter.NewJSONFormatter())

	const numGoroutines = 100
	const messagesPerGoroutine = 10

	// Run logging concurrently
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer func() { done <- true }()

			for j := 0; j < messagesPerGoroutine; j++ {
				logger.Infof("goroutine %d message %d", goroutineID, j)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	output := buf.String()

	// Verify that output is not empty
	if output == "" {
		t.Error("Expected log output from concurrent goroutines, got empty string")
		return
	}

	// Verify expected total number of messages (each log entry is one line of JSON, so split by newline)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	expectedMessages := numGoroutines * messagesPerGoroutine

	if len(lines) != expectedMessages {
		t.Errorf("Expected %d messages, got %d", expectedMessages, len(lines))
	}

	// Verify that each line can be parsed as JSON
	messageCount := make(map[int]int)
	for i, line := range lines {
		if line == "" {
			continue
		}

		var logData map[string]interface{}
		if err := json.Unmarshal([]byte(line), &logData); err != nil {
			t.Errorf("Line %d: Expected valid JSON output, got error: %v, line: %s", i, err, line)
			continue
		}

		// Verify structure
		runtime, ok := logData["runtime"].(map[string]interface{})
		if !ok {
			t.Errorf("Line %d: Expected runtime to be an object", i)
			continue
		}

		linesArray, ok := runtime["lines"].([]interface{})
		if !ok {
			t.Errorf("Line %d: Expected lines to be an array", i)
			continue
		}

		if len(linesArray) != 1 {
			t.Errorf("Line %d: Expected exactly 1 log entry, got %d", i, len(linesArray))
			continue
		}

		entry, ok := linesArray[0].(map[string]interface{})
		if !ok {
			t.Errorf("Line %d: Expected log entry to be an object", i)
			continue
		}

		message, ok := entry["message"].(string)
		if !ok {
			t.Errorf("Line %d: Expected message to be a string", i)
			continue
		}

		// Extract goroutine ID
		var goroutineID int
		if _, err := fmt.Sscanf(message, "goroutine %d", &goroutineID); err == nil {
			messageCount[goroutineID]++
		}
	}

	// Verify expected number of messages from each goroutine
	for i := 0; i < numGoroutines; i++ {
		if messageCount[i] != messagesPerGoroutine {
			t.Errorf("Expected %d messages from goroutine %d, got %d", messagesPerGoroutine, i, messageCount[i])
		}
	}
}

func TestDirectLogger_ErrorCases(t *testing.T) {
	logger := NewDirectLogger()

	t.Run("NilOutput", func(t *testing.T) {
		// Test behavior when nil is set as output destination
		// Current implementation should not panic, though some implementations might
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Setting nil output should not panic, got: %v", r)
			}
		}()

		logger.SetOutput(nil)
		// Verify that outputting to nil does not panic
		logger.Infof("test message")
	})

	t.Run("InvalidLogLevel", func(t *testing.T) {
		// Test invalid log level string
		var buf bytes.Buffer
		logger.SetOutput(&buf)

		// Set invalid level (should default to InfoLevel)
		logger.SetLevelFromString("INVALID_LEVEL")

		// Debug level should not be output
		logger.Debugf("debug message")
		if buf.String() != "" {
			t.Error("Debug message should not be output with default INFO level")
		}

		buf.Reset()
		// Info level should be output
		logger.Infof("info message")
		if buf.String() == "" {
			t.Error("Info message should be output with default INFO level")
		}
	})
}

func TestDirectLogger_StructuredOutput(t *testing.T) {
	var buf bytes.Buffer
	logger := NewDirectLogger()
	logger.SetOutput(&buf)

	logger.Infof("test message")
	output := buf.String()

	// Verify that output can be parsed as JSON
	var logData map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logData); err != nil {
		t.Errorf("Expected valid JSON output, got error: %v, output: %s", err, output)
		return
	}

	// Verify required fields in structured log
	expectedFields := []string{"type", "context", "runtime"}
	for _, field := range expectedFields {
		if _, exists := logData[field]; !exists {
			t.Errorf("Expected field %s in log output", field)
		}
	}

	// Detailed verification of runtime section
	runtime, ok := logData["runtime"].(map[string]interface{})
	if !ok {
		t.Error("Expected runtime to be an object")
		return
	}

	expectedRuntimeFields := []string{"severity", "startTime", "endTime", "elapsed", "lines"}
	for _, field := range expectedRuntimeFields {
		if _, exists := runtime[field]; !exists {
			t.Errorf("Expected field %s in runtime section", field)
		}
	}

	// Verify that elapsed is 0 (for direct logger)
	if runtime["elapsed"] != float64(0) {
		t.Errorf("Expected elapsed to be 0 for direct logger, got: %v", runtime["elapsed"])
	}

	// Verify that lines is an array with one entry
	lines, ok := runtime["lines"].([]interface{})
	if !ok {
		t.Error("Expected lines to be an array")
		return
	}

	if len(lines) != 1 {
		t.Errorf("Expected exactly 1 log entry, got %d", len(lines))
	}
}
