package logger

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := defaultConfig()

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

	if config.MaxLogEntries != 0 {
		t.Errorf("Expected MaxLogEntries to be 0, got %v", config.MaxLogEntries)
	}

	if config.ErrorHandler != nil {
		t.Errorf("Expected ErrorHandler to be nil, got %v", config.ErrorHandler)
	}
}

func TestInit_DefaultOptions(t *testing.T) {
	// Save original state
	originalConfig := GetConfig()
	originalInitialized := IsInitialized()

	// Test initialization with no options (should use defaults)
	Init()

	// Check if initialized
	if !IsInitialized() {
		t.Error("Expected IsInitialized to return true after Init")
	}

	// Check if config is set correctly with defaults
	config := GetConfig()
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

	if config.MaxLogEntries != 0 {
		t.Errorf("Expected MaxLogEntries to be 0, got %v", config.MaxLogEntries)
	}

	// Restore original state
	if originalInitialized {
		restoreConfig(originalConfig)
	}
}

func TestInit_WithOptions(t *testing.T) {
	// Save original state
	originalConfig := GetConfig()
	originalInitialized := IsInitialized()

	// Test initialization with options
	Init(
		WithMinLevel(DebugLevel),
		WithOutput(os.Stderr),
		WithSourceInfo(true),
		WithPrettifyJSON(true),
		WithMaxLogEntries(500),
		WithLogType("test"),
	)

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

	if config.LogType != "test" {
		t.Errorf("Expected LogType to be 'test', got %v", config.LogType)
	}

	// Restore original state
	if originalInitialized {
		restoreConfig(originalConfig)
	}
}

func TestWithMinLevel(t *testing.T) {
	option := WithMinLevel(ErrorLevel)
	config := defaultConfig()
	option(&config)

	if config.MinLevel != ErrorLevel {
		t.Errorf("Expected MinLevel to be ErrorLevel, got %v", config.MinLevel)
	}
}

func TestWithOutput(t *testing.T) {
	option := WithOutput(os.Stderr)
	config := defaultConfig()
	option(&config)

	if config.Output != os.Stderr {
		t.Errorf("Expected Output to be os.Stderr, got %v", config.Output)
	}
}

func TestWithSourceInfo(t *testing.T) {
	option := WithSourceInfo(true)
	config := defaultConfig()
	option(&config)

	if config.EnableSourceInfo != true {
		t.Errorf("Expected EnableSourceInfo to be true, got %v", config.EnableSourceInfo)
	}
}

func TestWithPrettifyJSON(t *testing.T) {
	option := WithPrettifyJSON(true)
	config := defaultConfig()
	option(&config)

	if config.PrettifyJSON != true {
		t.Errorf("Expected PrettifyJSON to be true, got %v", config.PrettifyJSON)
	}
}

func TestWithMaxLogEntries(t *testing.T) {
	option := WithMaxLogEntries(1000)
	config := defaultConfig()
	option(&config)

	if config.MaxLogEntries != 1000 {
		t.Errorf("Expected MaxLogEntries to be 1000, got %v", config.MaxLogEntries)
	}
}

func TestWithLogType(t *testing.T) {
	option := WithLogType("custom")
	config := defaultConfig()
	option(&config)

	if config.LogType != "custom" {
		t.Errorf("Expected LogType to be 'custom', got %v", config.LogType)
	}
}

func TestWithErrorHandler(t *testing.T) {
	handler := &SilentErrorHandler{}
	option := WithErrorHandler(handler)
	config := defaultConfig()
	option(&config)

	if config.ErrorHandler != handler {
		t.Errorf("Expected ErrorHandler to be set, got %v", config.ErrorHandler)
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
	Init()
	if !IsInitialized() {
		t.Error("Expected IsInitialized to return true after Init")
	}

	// Restore original state
	if originalInitialized {
		restoreConfig(originalConfig)
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

	Init(
		WithMinLevel(WarnLevel),
		WithOutput(os.Stderr),
		WithSourceInfo(true),
		WithPrettifyJSON(false),
		WithMaxLogEntries(200),
	)

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
		restoreConfig(originalConfig)
	}
}

func TestInit_DirectLoggerUpdate(t *testing.T) {
	// Save original state
	originalConfig := GetConfig()
	originalInitialized := IsInitialized()

	// Test configuration
	Init(
		WithMinLevel(WarnLevel),
		WithOutput(os.Stderr),
		WithSourceInfo(true),
		WithPrettifyJSON(true),
		WithMaxLogEntries(100),
		WithLogType("custom"),
	)

	// Verify that global direct logger was updated
	if directLogger, ok := D.(*DirectLogger); ok {
		// Check that level was updated
		if !directLogger.isLevelEnabled(WarnLevel) {
			t.Error("Expected WarnLevel to be enabled after Init")
		}
		if directLogger.isLevelEnabled(InfoLevel) {
			t.Error("Expected InfoLevel to be disabled after Init")
		}

		// Check that output was updated (we can't directly test this, but we can verify it doesn't panic)
		directLogger.Warnf("test message")
	} else {
		t.Error("Expected D to be a DirectLogger")
	}

	// Restore original state
	if originalInitialized {
		restoreConfig(originalConfig)
	}
}

func TestInit_MultipleOptions(t *testing.T) {
	// Save original state
	originalConfig := GetConfig()
	originalInitialized := IsInitialized()

	// Test multiple options applied in order
	Init(
		WithMinLevel(DebugLevel),
		WithMinLevel(ErrorLevel), // Should override the previous one
		WithPrettifyJSON(false),
		WithPrettifyJSON(true), // Should override the previous one
	)

	config := GetConfig()

	// The last option should take precedence
	if config.MinLevel != ErrorLevel {
		t.Errorf("Expected MinLevel to be ErrorLevel (last option), got %v", config.MinLevel)
	}

	if config.PrettifyJSON != true {
		t.Errorf("Expected PrettifyJSON to be true (last option), got %v", config.PrettifyJSON)
	}

	// Restore original state
	if originalInitialized {
		restoreConfig(originalConfig)
	}
}

// Helper function to restore original configuration
func restoreConfig(config Config) {
	Init(
		WithMinLevel(config.MinLevel),
		WithOutput(config.Output),
		WithSourceInfo(config.EnableSourceInfo),
		WithPrettifyJSON(config.PrettifyJSON),
		WithMaxLogEntries(config.MaxLogEntries),
		WithLogType(config.LogType),
		WithErrorHandler(config.ErrorHandler),
	)
}
