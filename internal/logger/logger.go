package logger

// LogLevel represents the severity level of a log entry
type LogLevel int

const (
	// DebugLevel represents debug level logs
	DebugLevel LogLevel = iota
	// InfoLevel represents info level logs
	InfoLevel
	// WarnLevel represents warning level logs
	WarnLevel
	// ErrorLevel represents error level logs
	ErrorLevel
	// CriticalLevel represents critical level logs
	CriticalLevel
)

// Logger is the main interface for logging
type Logger interface {
	// Debug logs a message at debug level
	Debug(msg string, fields ...Field)
	// Info logs a message at info level
	Info(msg string, fields ...Field)
	// Warn logs a message at warning level
	Warn(msg string, fields ...Field)
	// Error logs a message at error level
	Error(msg string, fields ...Field)
	// Critical logs a message at critical level
	Critical(msg string, fields ...Field)
	// With returns a new logger with the given fields added to its context
	With(fields ...Field) Logger
}

// Field represents a key-value pair in a log entry
type Field struct {
	Key   string
	Value interface{}
}

// NewField creates a new Field
func NewField(key string, value interface{}) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}
