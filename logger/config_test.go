package logger

import (
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

	if config.MaxLogEntries != 0 {
		t.Errorf("Expected MaxLogEntries to be 0, got %v", config.MaxLogEntries)
	}

	if config.ErrorHandler != nil {
		t.Errorf("Expected ErrorHandler to be nil, got %v", config.ErrorHandler)
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

	// Test configuration
	testConfig := Config{
		MinLevel:         WarnLevel,
		Output:           os.Stderr,
		EnableSourceInfo: true,
		PrettifyJSON:     true,
		MaxLogEntries:    100,
		LogType:          "custom",
		ErrorHandler:     nil,
	}

	// Initialize with test config
	Init(testConfig)

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
		Init(originalConfig)
	}
}

func TestConfig_LogType(t *testing.T) {
	// Test default config
	defaultConfig := DefaultConfig()
	if defaultConfig.LogType != "request" {
		t.Errorf("Expected default LogType to be 'request', got %s", defaultConfig.LogType)
	}

	// Test custom config
	customConfig := Config{
		MinLevel:         InfoLevel,
		Output:           os.Stdout,
		EnableSourceInfo: false,
		PrettifyJSON:     false,
		MaxLogEntries:    0,
		LogType:          "custom_type",
		ErrorHandler:     nil,
	}

	// Save original state
	originalConfig := GetConfig()
	originalInitialized := IsInitialized()

	// Initialize with custom config
	Init(customConfig)

	// Verify that the config was set
	config := GetConfig()
	if config.LogType != "custom_type" {
		t.Errorf("Expected LogType to be 'custom_type', got %s", config.LogType)
	}

	// Restore original state
	if originalInitialized {
		Init(originalConfig)
	}
}
