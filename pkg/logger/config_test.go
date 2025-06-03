package logger

import (
	"bytes"
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.MinLevel != InfoLevel {
		t.Errorf("Expected MinLevel to be InfoLevel, got %v", config.MinLevel)
	}

	if config.Output != os.Stdout {
		t.Errorf("Expected Output to be os.Stdout, got %v", config.Output)
	}

	if config.EnableSourceInfo != false {
		t.Errorf("Expected EnableSourceInfo to be false, got %v", config.EnableSourceInfo)
	}

	if config.PrettifyJSON != false {
		t.Errorf("Expected PrettifyJSON to be false, got %v", config.PrettifyJSON)
	}

	if config.MaxLogEntries != 1000 {
		t.Errorf("Expected MaxLogEntries to be 1000, got %v", config.MaxLogEntries)
	}
}

func TestInit(t *testing.T) {
	// Save original state
	originalConfig := GetConfig()
	originalInitialized := IsInitialized()

	// Test initialization
	testConfig := Config{
		MinLevel:         DebugLevel,
		Output:           os.Stderr,
		EnableSourceInfo: true,
		PrettifyJSON:     true,
		MaxLogEntries:    500,
	}

	Init(testConfig)

	// Check if initialized
	if !IsInitialized() {
		t.Error("Expected IsInitialized to return true after Init")
	}

	// Check if config is set correctly
	config := GetConfig()
	if config.MinLevel != DebugLevel {
		t.Errorf("Expected MinLevel to be DebugLevel, got %v", config.MinLevel)
	}

	if config.Output != os.Stderr {
		t.Errorf("Expected Output to be os.Stderr, got %v", config.Output)
	}

	if config.EnableSourceInfo != true {
		t.Errorf("Expected EnableSourceInfo to be true, got %v", config.EnableSourceInfo)
	}

	if config.PrettifyJSON != true {
		t.Errorf("Expected PrettifyJSON to be true, got %v", config.PrettifyJSON)
	}

	if config.MaxLogEntries != 500 {
		t.Errorf("Expected MaxLogEntries to be 500, got %v", config.MaxLogEntries)
	}

	// Restore original state
	if originalInitialized {
		Init(originalConfig)
	}
}

func TestIsInitialized(t *testing.T) {
	// Save original state
	originalConfig := GetConfig()
	originalInitialized := IsInitialized()

	// Reset initialization state for testing
	configMutex.Lock()
	initialized = false
	configMutex.Unlock()

	// Test before initialization
	if IsInitialized() {
		t.Error("Expected IsInitialized to return false before Init")
	}

	// Test after initialization
	Init(DefaultConfig())
	if !IsInitialized() {
		t.Error("Expected IsInitialized to return true after Init")
	}

	// Restore original state
	if originalInitialized {
		Init(originalConfig)
	} else {
		configMutex.Lock()
		initialized = false
		configMutex.Unlock()
	}
}

func TestGetConfig(t *testing.T) {
	// Save original state
	originalConfig := GetConfig()
	originalInitialized := IsInitialized()

	testConfig := Config{
		MinLevel:         WarnLevel,
		Output:           os.Stderr,
		EnableSourceInfo: true,
		PrettifyJSON:     false,
		MaxLogEntries:    200,
	}

	Init(testConfig)

	config := GetConfig()

	// Verify that GetConfig returns a copy with correct values
	if config.MinLevel != WarnLevel {
		t.Errorf("Expected MinLevel to be WarnLevel, got %v", config.MinLevel)
	}

	if config.Output != os.Stderr {
		t.Errorf("Expected Output to be os.Stderr, got %v", config.Output)
	}

	if config.EnableSourceInfo != true {
		t.Errorf("Expected EnableSourceInfo to be true, got %v", config.EnableSourceInfo)
	}

	if config.PrettifyJSON != false {
		t.Errorf("Expected PrettifyJSON to be false, got %v", config.PrettifyJSON)
	}

	if config.MaxLogEntries != 200 {
		t.Errorf("Expected MaxLogEntries to be 200, got %v", config.MaxLogEntries)
	}

	// Restore original state
	if originalInitialized {
		Init(originalConfig)
	}
}

func TestInit_DirectLoggerUpdate(t *testing.T) {
	// Save original state
	originalConfig := GetConfig()
	originalInitialized := IsInitialized()

	// Create a buffer to capture output
	var buf bytes.Buffer

	testConfig := Config{
		MinLevel:         ErrorLevel,
		Output:           &buf,
		EnableSourceInfo: false,
		PrettifyJSON:     true,
		MaxLogEntries:    100,
	}

	Init(testConfig)

	// Test that the global DirectLogger D is updated
	D.Infof("test message")   // This should not appear because MinLevel is ErrorLevel
	D.Errorf("error message") // This should appear

	output := buf.String()

	// Should not contain info message
	if len(output) == 0 {
		t.Error("Expected error message to be logged")
	}

	// Should contain error message and be pretty-printed (with indentation)
	if !bytes.Contains(buf.Bytes(), []byte("error message")) {
		t.Error("Expected output to contain 'error message'")
	}

	// Check for pretty-printing (indentation)
	if !bytes.Contains(buf.Bytes(), []byte("  ")) {
		t.Error("Expected output to be pretty-printed with indentation")
	}

	// Restore original state
	if originalInitialized {
		Init(originalConfig)
	}
}
