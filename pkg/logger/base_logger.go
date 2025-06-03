package logger

import (
	"io"
	"sync"

	"github.com/zentooo/logspan/pkg/formatter"
)

// BaseLogger contains common fields and methods shared by DirectLogger and ContextLogger
type BaseLogger struct {
	output    io.Writer
	minLevel  LogLevel
	formatter formatter.Formatter
	mutex     sync.Mutex
}

// newBaseLogger creates a new BaseLogger with default settings
func newBaseLogger() BaseLogger {
	return BaseLogger{
		output:    nil, // Will be set by the specific logger
		minLevel:  InfoLevel,
		formatter: createDefaultFormatter(),
	}
}

// SetOutput sets the output writer for the logger
func (b *BaseLogger) SetOutput(w io.Writer) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.output = w
}

// SetLevel sets the minimum log level for filtering
func (b *BaseLogger) SetLevel(level LogLevel) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.minLevel = level
}

// SetFormatter sets the formatter for the logger
func (b *BaseLogger) SetFormatter(f formatter.Formatter) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.formatter = f
}

// SetLevelFromString sets the minimum log level from a string
func (b *BaseLogger) SetLevelFromString(level string) {
	b.SetLevel(ParseLogLevel(level))
}

// isLevelEnabled checks if the given level should be logged
func (b *BaseLogger) isLevelEnabled(level LogLevel) bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return IsLevelEnabled(level, b.minLevel)
}

// getOutput returns the current output writer (thread-safe)
func (b *BaseLogger) getOutput() io.Writer {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.output
}

// getFormatter returns the current formatter (thread-safe)
func (b *BaseLogger) getFormatter() formatter.Formatter {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.formatter
}
