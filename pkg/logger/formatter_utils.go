package logger

import (
	"time"

	"github.com/zentooo/logspan/pkg/formatter"
)

// createDefaultFormatter creates a default formatter based on global configuration
func createDefaultFormatter() formatter.Formatter {
	config := GetConfig()
	if config.PrettifyJSON {
		return formatter.NewJSONFormatterWithIndent("  ")
	}
	return formatter.NewJSONFormatter()
}

// formatLogOutput creates a LogOutput structure and formats it using the given formatter
// If formatter is nil, uses default JSONFormatter
func formatLogOutput(entries []*LogEntry, contextFields map[string]interface{}, startTime, endTime time.Time, f formatter.Formatter) ([]byte, error) {
	elapsed := endTime.Sub(startTime).Milliseconds()

	// Find the highest severity level
	maxSeverity := DebugLevel
	for _, entry := range entries {
		entryLevel := ParseLogLevel(entry.Level)
		maxSeverity = GetHigherLevel(entryLevel, maxSeverity)
	}

	// Convert logger.LogEntry to formatter.LogEntry
	formatterEntries := make([]*formatter.LogEntry, len(entries))
	for i, entry := range entries {
		formatterEntries[i] = &formatter.LogEntry{
			Timestamp: entry.Timestamp,
			Level:     entry.Level,
			Message:   entry.Message,
			Funcname:  entry.Funcname,
			Filename:  entry.Filename,
			Fileline:  entry.Fileline,
		}
	}

	// Create LogOutput structure
	logOutput := &formatter.LogOutput{
		Type:    "request",
		Context: contextFields,
		Runtime: formatter.RuntimeInfo{
			Severity:  maxSeverity.String(),
			StartTime: startTime.Format(time.RFC3339Nano),
			EndTime:   endTime.Format(time.RFC3339Nano),
			Elapsed:   elapsed,
			Lines:     formatterEntries,
		},
	}

	// Use provided formatter or default JSONFormatter
	if f == nil {
		// Use default JSONFormatter with prettify setting from global config
		f = createDefaultFormatter()
	}

	return f.Format(logOutput)
}
