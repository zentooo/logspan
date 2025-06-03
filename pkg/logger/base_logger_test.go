package logger

import (
	"bytes"
	"testing"

	"github.com/zentooo/logspan/pkg/formatter"
)

func TestBaseLogger_SetOutput(t *testing.T) {
	base := newBaseLogger()
	var buf bytes.Buffer

	base.SetOutput(&buf)

	output := base.getOutput()
	if output != &buf {
		t.Error("Expected output to be set to buffer")
	}
}

func TestBaseLogger_SetLevel(t *testing.T) {
	base := newBaseLogger()

	// Test setting different levels
	testCases := []LogLevel{
		DebugLevel,
		InfoLevel,
		WarnLevel,
		ErrorLevel,
		CriticalLevel,
	}

	for _, level := range testCases {
		base.SetLevel(level)
		if base.minLevel != level {
			t.Errorf("Expected level %v, got %v", level, base.minLevel)
		}
	}
}

func TestBaseLogger_SetLevelFromString(t *testing.T) {
	base := newBaseLogger()

	testCases := []struct {
		input    string
		expected LogLevel
	}{
		{"DEBUG", DebugLevel},
		{"INFO", InfoLevel},
		{"WARN", WarnLevel},
		{"ERROR", ErrorLevel},
		{"CRITICAL", CriticalLevel},
		{"INVALID", InfoLevel}, // Default to INFO for invalid input
	}

	for _, tc := range testCases {
		base.SetLevelFromString(tc.input)
		if base.minLevel != tc.expected {
			t.Errorf("For input %s, expected level %v, got %v", tc.input, tc.expected, base.minLevel)
		}
	}
}

func TestBaseLogger_SetFormatter(t *testing.T) {
	base := newBaseLogger()
	jsonFormatter := formatter.NewJSONFormatter()

	base.SetFormatter(jsonFormatter)

	currentFormatter := base.getFormatter()
	if currentFormatter != jsonFormatter {
		t.Error("Expected formatter to be set to jsonFormatter")
	}
}

func TestBaseLogger_IsLevelEnabled(t *testing.T) {
	base := newBaseLogger()
	base.SetLevel(WarnLevel)

	testCases := []struct {
		level    LogLevel
		expected bool
	}{
		{DebugLevel, false},
		{InfoLevel, false},
		{WarnLevel, true},
		{ErrorLevel, true},
		{CriticalLevel, true},
	}

	for _, tc := range testCases {
		result := base.isLevelEnabled(tc.level)
		if result != tc.expected {
			t.Errorf("For level %v with minLevel %v, expected %v, got %v",
				tc.level, base.minLevel, tc.expected, result)
		}
	}
}

func TestBaseLogger_DefaultValues(t *testing.T) {
	base := newBaseLogger()

	// Check default values
	if base.minLevel != InfoLevel {
		t.Errorf("Expected default minLevel to be InfoLevel, got %v", base.minLevel)
	}

	if base.output != nil {
		t.Error("Expected default output to be nil")
	}

	if base.formatter == nil {
		t.Error("Expected default formatter to be set")
	}
}

func TestBaseLogger_ConcurrentAccess(t *testing.T) {
	base := newBaseLogger()
	var buf bytes.Buffer

	// Test concurrent access to ensure thread safety
	done := make(chan bool, 3)

	// Goroutine 1: Set output
	go func() {
		for i := 0; i < 100; i++ {
			base.SetOutput(&buf)
		}
		done <- true
	}()

	// Goroutine 2: Set level
	go func() {
		for i := 0; i < 100; i++ {
			base.SetLevel(DebugLevel)
		}
		done <- true
	}()

	// Goroutine 3: Check level enabled
	go func() {
		for i := 0; i < 100; i++ {
			base.isLevelEnabled(InfoLevel)
		}
		done <- true
	}()

	// Wait for all goroutines to complete
	for i := 0; i < 3; i++ {
		<-done
	}

	// If we reach here without deadlock, the test passes
}

func TestBaseLogger_GettersThreadSafety(t *testing.T) {
	base := newBaseLogger()
	var buf bytes.Buffer
	jsonFormatter := formatter.NewJSONFormatter()

	base.SetOutput(&buf)
	base.SetFormatter(jsonFormatter)

	// Test that getters are thread-safe
	done := make(chan bool, 2)

	// Goroutine 1: Get output
	go func() {
		for i := 0; i < 100; i++ {
			output := base.getOutput()
			if output != &buf {
				t.Error("Unexpected output value")
			}
		}
		done <- true
	}()

	// Goroutine 2: Get formatter
	go func() {
		for i := 0; i < 100; i++ {
			formatter := base.getFormatter()
			if formatter != jsonFormatter {
				t.Error("Unexpected formatter value")
			}
		}
		done <- true
	}()

	// Wait for all goroutines to complete
	for i := 0; i < 2; i++ {
		<-done
	}
}
