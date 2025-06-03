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
