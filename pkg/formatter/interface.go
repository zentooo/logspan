package formatter

import (
	"time"
)

// LogEntry represents a single log entry for formatting
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Funcname  string    `json:"funcname,omitempty"`
	Filename  string    `json:"filename,omitempty"`
	Fileline  int       `json:"fileline,omitempty"`
}

// LogOutput represents the complete log output structure
type LogOutput struct {
	Type    string                 `json:"type"`
	Context map[string]interface{} `json:"context"`
	Runtime RuntimeInfo            `json:"runtime"`
}

// RuntimeInfo contains runtime information for the log output
type RuntimeInfo struct {
	Severity  string      `json:"severity"`
	StartTime string      `json:"startTime"`
	EndTime   string      `json:"endTime"`
	Elapsed   int64       `json:"elapsed"`
	Lines     []*LogEntry `json:"lines"`
}

// Formatter defines the interface for log formatters
type Formatter interface {
	// Format formats the log output and returns the formatted bytes
	Format(output *LogOutput) ([]byte, error)
}
