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

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case CriticalLevel:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

// ParseLogLevel parses a string into a LogLevel
func ParseLogLevel(level string) LogLevel {
	switch level {
	case "DEBUG":
		return DebugLevel
	case "INFO":
		return InfoLevel
	case "WARN":
		return WarnLevel
	case "ERROR":
		return ErrorLevel
	case "CRITICAL":
		return CriticalLevel
	default:
		return InfoLevel // Default to INFO level
	}
}
