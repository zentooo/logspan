package logger

import (
	"regexp"
	"testing"
	"time"
)

func TestPasswordMaskingMiddleware_NewPasswordMaskingMiddleware(t *testing.T) {
	pmm := NewPasswordMaskingMiddleware()

	if pmm.MaskString != "***" {
		t.Errorf("Expected default mask string to be '***', got '%s'", pmm.MaskString)
	}

	if len(pmm.PasswordKeys) == 0 {
		t.Error("Expected default password keys to be set")
	}

	if len(pmm.PasswordPatterns) == 0 {
		t.Error("Expected default password patterns to be set")
	}

	// Check if common password keys are included
	expectedKeys := []string{"password", "secret", "token", "api_key"}
	for _, expectedKey := range expectedKeys {
		found := false
		for _, key := range pmm.PasswordKeys {
			if key == expectedKey {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected password key '%s' to be in default keys", expectedKey)
		}
	}
}

func TestPasswordMaskingMiddleware_WithMaskString(t *testing.T) {
	pmm := NewPasswordMaskingMiddleware()
	customMask := "[REDACTED]"

	pmm = pmm.WithMaskString(customMask)

	if pmm.MaskString != customMask {
		t.Errorf("Expected mask string to be '%s', got '%s'", customMask, pmm.MaskString)
	}
}

func TestPasswordMaskingMiddleware_WithPasswordKeys(t *testing.T) {
	pmm := NewPasswordMaskingMiddleware()
	customKeys := []string{"custom_password", "secret_key"}

	pmm = pmm.WithPasswordKeys(customKeys)

	if len(pmm.PasswordKeys) != len(customKeys) {
		t.Errorf("Expected %d password keys, got %d", len(customKeys), len(pmm.PasswordKeys))
	}

	for i, key := range customKeys {
		if pmm.PasswordKeys[i] != key {
			t.Errorf("Expected password key '%s' at index %d, got '%s'", key, i, pmm.PasswordKeys[i])
		}
	}
}

func TestPasswordMaskingMiddleware_AddPasswordKey(t *testing.T) {
	pmm := NewPasswordMaskingMiddleware()
	originalCount := len(pmm.PasswordKeys)
	newKey := "custom_secret"

	pmm = pmm.AddPasswordKey(newKey)

	if len(pmm.PasswordKeys) != originalCount+1 {
		t.Errorf("Expected %d password keys after adding one, got %d", originalCount+1, len(pmm.PasswordKeys))
	}

	if pmm.PasswordKeys[len(pmm.PasswordKeys)-1] != newKey {
		t.Errorf("Expected last password key to be '%s', got '%s'", newKey, pmm.PasswordKeys[len(pmm.PasswordKeys)-1])
	}
}

func TestPasswordMaskingMiddleware_AddPasswordPattern(t *testing.T) {
	pmm := NewPasswordMaskingMiddleware()
	originalCount := len(pmm.PasswordPatterns)
	newPattern := regexp.MustCompile(`custom_pattern=\w+`)

	pmm = pmm.AddPasswordPattern(newPattern)

	if len(pmm.PasswordPatterns) != originalCount+1 {
		t.Errorf("Expected %d password patterns after adding one, got %d", originalCount+1, len(pmm.PasswordPatterns))
	}
}

func TestPasswordMaskingMiddleware_MaskPasswordsInMessage(t *testing.T) {
	pmm := NewPasswordMaskingMiddleware()

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "password with equals",
			input:    "User login with password=secret123",
			expected: "User login with password=***",
		},
		{
			name:     "password with colon",
			input:    "Config: password: mypassword",
			expected: "Config: password: ***",
		},
		{
			name:     "JSON-like password",
			input:    `{"password":"secret123","user":"john"}`,
			expected: `{"password":"***","user":"john"}`,
		},
		{
			name:     "multiple passwords",
			input:    "password=secret123 and token=abc456",
			expected: "password=*** and token=***",
		},
		{
			name:     "case insensitive",
			input:    "PASSWORD=secret123 and Secret=abc456",
			expected: "PASSWORD=*** and Secret=***",
		},
		{
			name:     "no password",
			input:    "Normal log message without sensitive data",
			expected: "Normal log message without sensitive data",
		},
		{
			name:     "api_key pattern",
			input:    "Using api_key=12345abcdef",
			expected: "Using api_key=***",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := pmm.maskPasswordsInMessage(tc.input)
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestPasswordMaskingMiddleware_MaskPasswordsInFields(t *testing.T) {
	pmm := NewPasswordMaskingMiddleware()

	testCases := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name: "simple password field",
			input: map[string]interface{}{
				"username": "john",
				"password": "secret123",
			},
			expected: map[string]interface{}{
				"username": "john",
				"password": "***",
			},
		},
		{
			name: "multiple sensitive fields",
			input: map[string]interface{}{
				"username":    "john",
				"password":    "secret123",
				"api_key":     "abc123",
				"secret":      "mysecret",
				"normal_data": "value",
			},
			expected: map[string]interface{}{
				"username":    "john",
				"password":    "***",
				"api_key":     "***",
				"secret":      "***",
				"normal_data": "value",
			},
		},
		{
			name: "case insensitive keys",
			input: map[string]interface{}{
				"PASSWORD": "secret123",
				"Secret":   "mysecret",
				"API_KEY":  "abc123",
			},
			expected: map[string]interface{}{
				"PASSWORD": "***",
				"Secret":   "***",
				"API_KEY":  "***",
			},
		},
		{
			name: "nested map",
			input: map[string]interface{}{
				"user": map[string]interface{}{
					"name":     "john",
					"password": "secret123",
				},
				"config": map[string]interface{}{
					"api_key": "abc123",
					"timeout": 30,
				},
			},
			expected: map[string]interface{}{
				"user": map[string]interface{}{
					"name":     "john",
					"password": "***",
				},
				"config": map[string]interface{}{
					"api_key": "***",
					"timeout": 30,
				},
			},
		},
		{
			name:     "empty fields",
			input:    map[string]interface{}{},
			expected: map[string]interface{}{},
		},
		{
			name:     "nil fields",
			input:    nil,
			expected: map[string]interface{}{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result map[string]interface{}
			if tc.input == nil {
				result = pmm.maskPasswordsInFields(tc.input)
			} else {
				result = pmm.maskPasswordsInFields(tc.input)
			}

			if tc.input == nil && len(result) != 0 {
				t.Error("Expected empty map for nil input")
				return
			}

			if tc.input != nil {
				compareFields(t, tc.expected, result)
			}
		})
	}
}

func TestPasswordMaskingMiddleware_MaskPasswordsInGenericMap(t *testing.T) {
	pmm := NewPasswordMaskingMiddleware()

	input := map[interface{}]interface{}{
		"username": "john",
		"password": "secret123",
		"api_key":  "abc123",
		123:        "numeric_key",
	}

	result := pmm.maskPasswordsInGenericMap(input)

	if result["username"] != "john" {
		t.Errorf("Expected username to remain 'john', got '%v'", result["username"])
	}

	if result["password"] != "***" {
		t.Errorf("Expected password to be masked, got '%v'", result["password"])
	}

	if result["api_key"] != "***" {
		t.Errorf("Expected api_key to be masked, got '%v'", result["api_key"])
	}

	if result[123] != "numeric_key" {
		t.Errorf("Expected numeric key to remain unchanged, got '%v'", result[123])
	}
}

func TestPasswordMaskingMiddleware_IsPasswordKey(t *testing.T) {
	pmm := NewPasswordMaskingMiddleware()

	testCases := []struct {
		key      string
		expected bool
	}{
		{"password", true},
		{"PASSWORD", true},
		{"Password", true},
		{"secret", true},
		{"SECRET", true},
		{"api_key", true},
		{"API_KEY", true},
		{"token", true},
		{"username", false},
		{"email", false},
		{"data", false},
		{"", false},
	}

	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			result := pmm.isPasswordKey(tc.key)
			if result != tc.expected {
				t.Errorf("Expected isPasswordKey('%s') to be %v, got %v", tc.key, tc.expected, result)
			}
		})
	}
}

func TestPasswordMaskingMiddleware_Middleware(t *testing.T) {
	pmm := NewPasswordMaskingMiddleware()
	middleware := pmm.Middleware()

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "User login with password=secret123",
		Fields: map[string]interface{}{
			"username": "john",
			"password": "secret123",
			"api_key":  "abc123",
		},
	}

	var processedEntry *LogEntry
	next := func(e *LogEntry) {
		processedEntry = e
	}

	middleware(entry, next)

	if processedEntry == nil {
		t.Fatal("Expected entry to be processed")
	}

	expectedMessage := "User login with password=***"
	if processedEntry.Message != expectedMessage {
		t.Errorf("Expected message to be '%s', got '%s'", expectedMessage, processedEntry.Message)
	}

	if processedEntry.Fields["username"] != "john" {
		t.Errorf("Expected username to remain 'john', got '%v'", processedEntry.Fields["username"])
	}

	if processedEntry.Fields["password"] != "***" {
		t.Errorf("Expected password to be masked, got '%v'", processedEntry.Fields["password"])
	}

	if processedEntry.Fields["api_key"] != "***" {
		t.Errorf("Expected api_key to be masked, got '%v'", processedEntry.Fields["api_key"])
	}
}

func TestPasswordMaskingMiddleware_CustomMaskString(t *testing.T) {
	pmm := NewPasswordMaskingMiddleware().WithMaskString("[REDACTED]")
	middleware := pmm.Middleware()

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "password=secret123",
		Fields: map[string]interface{}{
			"password": "secret123",
		},
	}

	var processedEntry *LogEntry
	next := func(e *LogEntry) {
		processedEntry = e
	}

	middleware(entry, next)

	expectedMessage := "password=[REDACTED]"
	if processedEntry.Message != expectedMessage {
		t.Errorf("Expected message to be '%s', got '%s'", expectedMessage, processedEntry.Message)
	}

	if processedEntry.Fields["password"] != "[REDACTED]" {
		t.Errorf("Expected password to be '[REDACTED]', got '%v'", processedEntry.Fields["password"])
	}
}

func TestPasswordMaskingMiddleware_NilFields(t *testing.T) {
	pmm := NewPasswordMaskingMiddleware()
	middleware := pmm.Middleware()

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "password=secret123",
		Fields:    nil,
	}

	var processedEntry *LogEntry
	next := func(e *LogEntry) {
		processedEntry = e
	}

	middleware(entry, next)

	if processedEntry == nil {
		t.Fatal("Expected entry to be processed")
	}

	expectedMessage := "password=***"
	if processedEntry.Message != expectedMessage {
		t.Errorf("Expected message to be '%s', got '%s'", expectedMessage, processedEntry.Message)
	}

	// Fields should remain nil
	if processedEntry.Fields != nil {
		t.Error("Expected fields to remain nil")
	}
}

// Helper function to compare field maps
func compareFields(t *testing.T, expected, actual map[string]interface{}) {
	if len(expected) != len(actual) {
		t.Errorf("Expected %d fields, got %d", len(expected), len(actual))
		return
	}

	for key, expectedValue := range expected {
		actualValue, exists := actual[key]
		if !exists {
			t.Errorf("Expected field '%s' to exist", key)
			continue
		}

		// Handle nested maps
		if expectedMap, ok := expectedValue.(map[string]interface{}); ok {
			if actualMap, ok := actualValue.(map[string]interface{}); ok {
				compareFields(t, expectedMap, actualMap)
			} else {
				t.Errorf("Expected field '%s' to be a map, got %T", key, actualValue)
			}
		} else {
			if actualValue != expectedValue {
				t.Errorf("Expected field '%s' to be '%v', got '%v'", key, expectedValue, actualValue)
			}
		}
	}
}
