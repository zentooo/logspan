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
}

func TestPasswordMaskingMiddleware_CustomMaskString(t *testing.T) {
	pmm := NewPasswordMaskingMiddleware().WithMaskString("[REDACTED]")
	middleware := pmm.Middleware()

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "password=secret123",
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
}

func TestPasswordMaskingMiddleware_MessageOnly(t *testing.T) {
	pmm := NewPasswordMaskingMiddleware()
	middleware := pmm.Middleware()

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "password=secret123",
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
}
