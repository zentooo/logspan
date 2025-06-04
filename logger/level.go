package logger

// LogLevel represents the severity level of a log entry
type LogLevel int

const (
	// DebugLevel is the lowest level, used for detailed debugging information
	DebugLevel LogLevel = iota
	// InfoLevel is used for general informational messages
	InfoLevel
	// WarnLevel is used for warning messages
	WarnLevel
	// ErrorLevel is used for error messages
	ErrorLevel
	// CriticalLevel is the highest level, used for critical error messages
	CriticalLevel
)

// Log level string constants
const (
	debugLevelString    = "DEBUG"
	infoLevelString     = "INFO"
	warnLevelString     = "WARN"
	errorLevelString    = "ERROR"
	criticalLevelString = "CRITICAL"
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return debugLevelString
	case InfoLevel:
		return infoLevelString
	case WarnLevel:
		return warnLevelString
	case ErrorLevel:
		return errorLevelString
	case CriticalLevel:
		return criticalLevelString
	default:
		return "UNKNOWN"
	}
}

// ParseLogLevel parses a string into a LogLevel
// Supported values: "DEBUG", "INFO", "WARN", "ERROR", "CRITICAL"
// Returns InfoLevel for any unrecognized input
func ParseLogLevel(level string) LogLevel {
	switch level {
	case debugLevelString:
		return DebugLevel
	case infoLevelString:
		return InfoLevel
	case warnLevelString:
		return WarnLevel
	case errorLevelString:
		return ErrorLevel
	case criticalLevelString:
		return CriticalLevel
	default:
		return InfoLevel // Default to INFO level
	}
}

// GreaterThan returns true if this log level is greater than the other level
func (l LogLevel) GreaterThan(other LogLevel) bool {
	return l > other
}

// GreaterThanOrEqual returns true if this log level is greater than or equal to the other level
func (l LogLevel) GreaterThanOrEqual(other LogLevel) bool {
	return l >= other
}

// LessThan returns true if this log level is less than the other level
func (l LogLevel) LessThan(other LogLevel) bool {
	return l < other
}

// LessThanOrEqual returns true if this log level is less than or equal to the other level
func (l LogLevel) LessThanOrEqual(other LogLevel) bool {
	return l <= other
}

// IsLevelEnabled checks if the given level should be logged based on the minimum level
func IsLevelEnabled(level, minLevel LogLevel) bool {
	return level.GreaterThanOrEqual(minLevel)
}

// GetHigherLevel returns the higher of two log levels
func GetHigherLevel(level1, level2 LogLevel) LogLevel {
	if level1.GreaterThan(level2) {
		return level1
	}
	return level2
}
