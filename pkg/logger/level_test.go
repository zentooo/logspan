package logger

import "testing"

func TestLogLevel_String(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{DebugLevel, "DEBUG"},
		{InfoLevel, "INFO"},
		{WarnLevel, "WARN"},
		{ErrorLevel, "ERROR"},
		{CriticalLevel, "CRITICAL"},
		{LogLevel(999), "UNKNOWN"},
	}

	for _, test := range tests {
		result := test.level.String()
		if result != test.expected {
			t.Errorf("LogLevel(%d).String() = %s, expected %s", test.level, result, test.expected)
		}
	}
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected LogLevel
	}{
		{"DEBUG", DebugLevel},
		{"INFO", InfoLevel},
		{"WARN", WarnLevel},
		{"ERROR", ErrorLevel},
		{"CRITICAL", CriticalLevel},
		{"UNKNOWN", InfoLevel}, // Default to INFO
		{"", InfoLevel},        // Default to INFO
	}

	for _, test := range tests {
		result := ParseLogLevel(test.input)
		if result != test.expected {
			t.Errorf("ParseLogLevel(%s) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestLogLevel_GreaterThan(t *testing.T) {
	tests := []struct {
		level1   LogLevel
		level2   LogLevel
		expected bool
	}{
		{ErrorLevel, InfoLevel, true},
		{InfoLevel, ErrorLevel, false},
		{InfoLevel, InfoLevel, false},
		{CriticalLevel, DebugLevel, true},
		{DebugLevel, CriticalLevel, false},
	}

	for _, test := range tests {
		result := test.level1.GreaterThan(test.level2)
		if result != test.expected {
			t.Errorf("%s.GreaterThan(%s) = %v, expected %v",
				test.level1.String(), test.level2.String(), result, test.expected)
		}
	}
}

func TestLogLevel_GreaterThanOrEqual(t *testing.T) {
	tests := []struct {
		level1   LogLevel
		level2   LogLevel
		expected bool
	}{
		{ErrorLevel, InfoLevel, true},
		{InfoLevel, ErrorLevel, false},
		{InfoLevel, InfoLevel, true},
		{CriticalLevel, DebugLevel, true},
		{DebugLevel, CriticalLevel, false},
	}

	for _, test := range tests {
		result := test.level1.GreaterThanOrEqual(test.level2)
		if result != test.expected {
			t.Errorf("%s.GreaterThanOrEqual(%s) = %v, expected %v",
				test.level1.String(), test.level2.String(), result, test.expected)
		}
	}
}

func TestLogLevel_LessThan(t *testing.T) {
	tests := []struct {
		level1   LogLevel
		level2   LogLevel
		expected bool
	}{
		{InfoLevel, ErrorLevel, true},
		{ErrorLevel, InfoLevel, false},
		{InfoLevel, InfoLevel, false},
		{DebugLevel, CriticalLevel, true},
		{CriticalLevel, DebugLevel, false},
	}

	for _, test := range tests {
		result := test.level1.LessThan(test.level2)
		if result != test.expected {
			t.Errorf("%s.LessThan(%s) = %v, expected %v",
				test.level1.String(), test.level2.String(), result, test.expected)
		}
	}
}

func TestLogLevel_LessThanOrEqual(t *testing.T) {
	tests := []struct {
		level1   LogLevel
		level2   LogLevel
		expected bool
	}{
		{InfoLevel, ErrorLevel, true},
		{ErrorLevel, InfoLevel, false},
		{InfoLevel, InfoLevel, true},
		{DebugLevel, CriticalLevel, true},
		{CriticalLevel, DebugLevel, false},
	}

	for _, test := range tests {
		result := test.level1.LessThanOrEqual(test.level2)
		if result != test.expected {
			t.Errorf("%s.LessThanOrEqual(%s) = %v, expected %v",
				test.level1.String(), test.level2.String(), result, test.expected)
		}
	}
}

func TestIsLevelEnabled(t *testing.T) {
	tests := []struct {
		level    LogLevel
		minLevel LogLevel
		expected bool
	}{
		{ErrorLevel, InfoLevel, true},
		{InfoLevel, ErrorLevel, false},
		{InfoLevel, InfoLevel, true},
		{DebugLevel, InfoLevel, false},
		{CriticalLevel, DebugLevel, true},
	}

	for _, test := range tests {
		result := IsLevelEnabled(test.level, test.minLevel)
		if result != test.expected {
			t.Errorf("IsLevelEnabled(%s, %s) = %v, expected %v",
				test.level.String(), test.minLevel.String(), result, test.expected)
		}
	}
}

func TestGetHigherLevel(t *testing.T) {
	tests := []struct {
		level1   LogLevel
		level2   LogLevel
		expected LogLevel
	}{
		{ErrorLevel, InfoLevel, ErrorLevel},
		{InfoLevel, ErrorLevel, ErrorLevel},
		{InfoLevel, InfoLevel, InfoLevel},
		{DebugLevel, CriticalLevel, CriticalLevel},
		{CriticalLevel, DebugLevel, CriticalLevel},
	}

	for _, test := range tests {
		result := GetHigherLevel(test.level1, test.level2)
		if result != test.expected {
			t.Errorf("GetHigherLevel(%s, %s) = %s, expected %s",
				test.level1.String(), test.level2.String(), result.String(), test.expected.String())
		}
	}
}
