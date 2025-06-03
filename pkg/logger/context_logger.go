package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

// ContextLogger implements context-based logging with log aggregation
type ContextLogger struct {
	*BaseLogger
	entries    []*LogEntry
	fields     map[string]interface{}
	startTime  time.Time
	maxEntries int // Maximum number of entries before auto-flush
}

// NewContextLogger creates a new ContextLogger instance
func NewContextLogger() *ContextLogger {
	// Get global config to determine formatter settings
	config := GetConfig()

	base := newBaseLogger()
	base.output = os.Stdout // Set default output for ContextLogger
	base.minLevel = config.MinLevel

	return &ContextLogger{
		BaseLogger: &base,
		entries:    make([]*LogEntry, 0),
		fields:     make(map[string]interface{}),
		startTime:  time.Now(),
		maxEntries: config.MaxLogEntries,
	}
}

// addEntry adds a log entry to the context logger
func (l *ContextLogger) addEntry(level LogLevel, message string) {
	if !l.isLevelEnabled(level) {
		return
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     level.String(),
		Message:   message,
	}

	// Add source information if enabled
	config := GetConfig()
	if config.EnableSourceInfo {
		// Determine the appropriate skip level by examining the call stack
		skipLevel := 3 // Default for direct ContextLogger method calls

		// Check the call stack to see if we're being called through context.go functions
		for i := 2; i <= 5; i++ {
			if pc, _, _, ok := runtime.Caller(i); ok {
				if fn := runtime.FuncForPC(pc); fn != nil {
					funcName := fn.Name()
					// If we find a context.go function in the stack, use skip level 4
					if strings.Contains(funcName, "/pkg/logger.Infof") ||
						strings.Contains(funcName, "/pkg/logger.Debugf") ||
						strings.Contains(funcName, "/pkg/logger.Warnf") ||
						strings.Contains(funcName, "/pkg/logger.Errorf") ||
						strings.Contains(funcName, "/pkg/logger.Criticalf") {
						skipLevel = 4
						break
					}
				}
			}
		}

		sourceInfo := getSourceInfo(skipLevel)
		entry.Funcname = sourceInfo.Funcname
		entry.Filename = sourceInfo.Filename
		entry.Fileline = sourceInfo.Fileline
	}

	// Process through global middleware chain
	processWithGlobalMiddleware(entry, func(processedEntry *LogEntry) {
		l.entries = append(l.entries, processedEntry)

		// Check if we need to auto-flush due to entry limit
		if l.maxEntries > 0 && len(l.entries) >= l.maxEntries {
			l.flushInternal()
		}
	})
}

// flushInternal performs the flush operation without acquiring the mutex
// This method assumes the mutex is already held by the caller
func (l *ContextLogger) flushInternal() {
	if l.output == nil || len(l.entries) == 0 {
		return
	}

	endTime := time.Now()

	// Use the formatter (default or explicitly set)
	jsonData, err := formatLogOutput(l.entries, l.fields, l.startTime, endTime, l.formatter)
	if err != nil {
		// Handle formatting error using error handler
		handleError("format", err)
		// Fallback to simple output if formatting fails
		_, writeErr := fmt.Fprintf(l.output, "Error formatting log: %v\n", err)
		if writeErr != nil {
			handleError("write_fallback", writeErr)
		}
		return
	}

	if _, err := fmt.Fprintf(l.output, "%s\n", jsonData); err != nil {
		// Handle write error using error handler
		handleError("write", err)
		// Try to write an error message as fallback
		_, fallbackErr := fmt.Fprintf(l.output, "Error writing log output: %v\n", err)
		if fallbackErr != nil {
			handleError("write_error_fallback", fallbackErr)
		}
	}

	// Clear entries after flushing and reset start time
	l.entries = l.entries[:0] // Clear slice but keep capacity
	l.startTime = time.Now()  // Reset start time for next batch
}

// Flush outputs all accumulated log entries as a single JSON
func (l *ContextLogger) Flush() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.flushInternal()
}

// AddContextValue adds a field to the context
func (l *ContextLogger) AddContextValue(key string, value interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.fields[key] = value
}

// AddContextValues adds multiple fields to the context
func (l *ContextLogger) AddContextValues(fields map[string]interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	for k, v := range fields {
		l.fields[k] = v
	}
}

// Debugf logs a debug message
func (l *ContextLogger) Debugf(format string, args ...interface{}) {
	l.addEntry(DebugLevel, fmt.Sprintf(format, args...))
}

// Infof logs an info message
func (l *ContextLogger) Infof(format string, args ...interface{}) {
	l.addEntry(InfoLevel, fmt.Sprintf(format, args...))
}

// Warnf logs a warning message
func (l *ContextLogger) Warnf(format string, args ...interface{}) {
	l.addEntry(WarnLevel, fmt.Sprintf(format, args...))
}

// Errorf logs an error message
func (l *ContextLogger) Errorf(format string, args ...interface{}) {
	l.addEntry(ErrorLevel, fmt.Sprintf(format, args...))
}

// Criticalf logs a critical message
func (l *ContextLogger) Criticalf(format string, args ...interface{}) {
	l.addEntry(CriticalLevel, fmt.Sprintf(format, args...))
}
