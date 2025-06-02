package logger

import (
	"time"
)

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
	Fields    []Field
}

// NewLogEntry creates a new LogEntry
func NewLogEntry(level LogLevel, msg string, fields ...Field) *LogEntry {
	return &LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Fields:    fields,
	}
}

// AddField adds a field to the log entry
func (e *LogEntry) AddField(field Field) {
	e.Fields = append(e.Fields, field)
}

// AddFields adds multiple fields to the log entry
func (e *LogEntry) AddFields(fields ...Field) {
	e.Fields = append(e.Fields, fields...)
}
