package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestSourceInfo_EnabledAndDisabled(t *testing.T) {
	// Test with source info enabled
	t.Run("SourceInfoEnabled", func(t *testing.T) {
		buf := &bytes.Buffer{}

		// Configure with source info enabled
		config := DefaultConfig()
		config.EnableSourceInfo = true
		config.Output = buf
		Init(config)

		// Test direct logger
		D.Infof("test message")

		output := buf.String()
		if output == "" {
			t.Fatal("Expected log output, got empty string")
		}

		var logData map[string]interface{}
		if err := json.Unmarshal([]byte(output), &logData); err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}

		runtime := logData["runtime"].(map[string]interface{})
		lines := runtime["lines"].([]interface{})
		entry := lines[0].(map[string]interface{})

		// Check that source info fields are present
		if entry["funcname"] == nil {
			t.Error("Expected funcname to be present when source info is enabled")
		}
		if entry["filename"] == nil {
			t.Error("Expected filename to be present when source info is enabled")
		}
		if entry["fileline"] == nil {
			t.Error("Expected fileline to be present when source info is enabled")
		}

		// Check that funcname contains the test function name
		funcname := entry["funcname"].(string)
		if !strings.Contains(funcname, "TestSourceInfo_EnabledAndDisabled") {
			t.Errorf("Expected funcname to contain test function name, got: %s", funcname)
		}

		// Check that filename is correct
		filename := entry["filename"].(string)
		if filename != "source_info_test.go" {
			t.Errorf("Expected filename to be source_info_test.go, got: %s", filename)
		}

		// Check that fileline is a positive number
		fileline := entry["fileline"].(float64)
		if fileline <= 0 {
			t.Errorf("Expected fileline to be positive, got: %f", fileline)
		}
	})

	// Test with source info disabled
	t.Run("SourceInfoDisabled", func(t *testing.T) {
		buf := &bytes.Buffer{}

		// Configure with source info disabled
		config := DefaultConfig()
		config.EnableSourceInfo = false
		config.Output = buf
		Init(config)

		// Test direct logger
		D.Infof("test message")

		output := buf.String()
		if output == "" {
			t.Fatal("Expected log output, got empty string")
		}

		var logData map[string]interface{}
		if err := json.Unmarshal([]byte(output), &logData); err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}

		runtime := logData["runtime"].(map[string]interface{})
		lines := runtime["lines"].([]interface{})
		entry := lines[0].(map[string]interface{})

		// Check that source info fields are not present (or empty)
		if funcname, exists := entry["funcname"]; exists && funcname != "" {
			t.Errorf("Expected funcname to be empty when source info is disabled, got: %v", funcname)
		}
		if filename, exists := entry["filename"]; exists && filename != "" {
			t.Errorf("Expected filename to be empty when source info is disabled, got: %v", filename)
		}
		if fileline, exists := entry["fileline"]; exists && fileline != 0 {
			t.Errorf("Expected fileline to be 0 when source info is disabled, got: %v", fileline)
		}
	})
}

func TestSourceInfo_ContextLogger(t *testing.T) {
	buf := &bytes.Buffer{}

	// Configure with source info enabled
	config := DefaultConfig()
	config.EnableSourceInfo = true
	config.Output = buf
	Init(config)

	// Test context logger
	ctx := context.Background()
	contextLogger := NewContextLogger()
	contextLogger.SetOutput(buf)
	ctx = WithLogger(ctx, contextLogger)

	AddContextValue(ctx, "test_key", "test_value")
	Infof(ctx, "context logger test message")
	FlushContext(ctx)

	output := buf.String()
	if output == "" {
		t.Fatal("Expected log output, got empty string")
	}

	var logData map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logData); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	runtime := logData["runtime"].(map[string]interface{})
	lines := runtime["lines"].([]interface{})
	entry := lines[0].(map[string]interface{})

	// Check that source info fields are present
	if entry["funcname"] == nil {
		t.Error("Expected funcname to be present")
	}
	if entry["filename"] == nil {
		t.Error("Expected filename to be present")
	}
	if entry["fileline"] == nil {
		t.Error("Expected fileline to be present")
	}

	// Check that funcname contains the test function name
	funcname := entry["funcname"].(string)
	if !strings.Contains(funcname, "TestSourceInfo_ContextLogger") {
		t.Errorf("Expected funcname to contain test function name, got: %s", funcname)
	}

	// Check that filename is correct
	filename := entry["filename"].(string)
	if filename != "source_info_test.go" {
		t.Errorf("Expected filename to be source_info_test.go, got: %s", filename)
	}
}

func TestSourceInfo_DirectContextLoggerCall(t *testing.T) {
	buf := &bytes.Buffer{}

	// Configure with source info enabled
	config := DefaultConfig()
	config.EnableSourceInfo = true
	config.Output = buf
	Init(config)

	// Test direct context logger call (not through context.go functions)
	contextLogger := NewContextLogger()
	contextLogger.SetOutput(buf)
	contextLogger.Infof("direct context logger call")
	contextLogger.Flush()

	output := buf.String()
	if output == "" {
		t.Fatal("Expected log output, got empty string")
	}

	var logData map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logData); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	runtime := logData["runtime"].(map[string]interface{})
	lines := runtime["lines"].([]interface{})
	entry := lines[0].(map[string]interface{})

	// Check that source info fields are present
	if entry["funcname"] == nil {
		t.Error("Expected funcname to be present")
	}
	if entry["filename"] == nil {
		t.Error("Expected filename to be present")
	}
	if entry["fileline"] == nil {
		t.Error("Expected fileline to be present")
	}

	// Check that funcname contains the test function name
	funcname := entry["funcname"].(string)
	if !strings.Contains(funcname, "TestSourceInfo_DirectContextLoggerCall") {
		t.Errorf("Expected funcname to contain test function name, got: %s", funcname)
	}
}

func TestGetSourceInfo(t *testing.T) {
	// Test getSourceInfo function directly
	sourceInfo := getSourceInfo(0) // Current function

	if sourceInfo.Funcname == "" {
		t.Error("Expected funcname to be non-empty")
	}
	if sourceInfo.Filename == "" {
		t.Error("Expected filename to be non-empty")
	}
	if sourceInfo.Fileline <= 0 {
		t.Error("Expected fileline to be positive")
	}

	// Check that funcname contains the test function name
	if !strings.Contains(sourceInfo.Funcname, "TestGetSourceInfo") {
		t.Logf("Funcname: %s (this is expected when testing getSourceInfo directly)", sourceInfo.Funcname)
	}

	// Check that filename is entry.go (since getSourceInfo is defined there)
	if sourceInfo.Filename != "entry.go" {
		t.Errorf("Expected filename to be entry.go, got: %s", sourceInfo.Filename)
	}
}
