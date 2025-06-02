package formatter

import (
	"encoding/json"
)

// DataDogFormatter implements the Formatter interface for DataDog Standard Attributes
type DataDogFormatter struct {
	// Indent specifies the indentation string for pretty printing
	// Empty string means no indentation (compact JSON)
	Indent string
}

// NewDataDogFormatter creates a new DataDogFormatter instance
func NewDataDogFormatter() *DataDogFormatter {
	return &DataDogFormatter{
		Indent: "", // Default to compact JSON
	}
}

// NewDataDogFormatterWithIndent creates a new DataDogFormatter with indentation
func NewDataDogFormatterWithIndent(indent string) *DataDogFormatter {
	return &DataDogFormatter{
		Indent: indent,
	}
}

// Format formats the log output according to DataDog Standard Attributes
func (f *DataDogFormatter) Format(output *LogOutput) ([]byte, error) {
	// Convert to DataDog format
	ddOutput := f.convertToDataDogFormat(output)

	if f.Indent == "" {
		return json.Marshal(ddOutput)
	}
	return json.MarshalIndent(ddOutput, "", f.Indent)
}

// convertToDataDogFormat converts LogOutput to DataDog Standard Attributes format
func (f *DataDogFormatter) convertToDataDogFormat(output *LogOutput) interface{} {
	// Create base DataDog output
	ddOutput := map[string]interface{}{
		"timestamp": output.Runtime.StartTime,
		"status":    f.convertSeverityToStatus(output.Runtime.Severity),
		"message":   f.createSummaryMessage(output),
		"logger":    "logspan",
		"duration":  output.Runtime.Elapsed,
	}

	// Add context fields as custom attributes
	for k, v := range output.Context {
		ddOutput[k] = v
	}

	// Convert log entries
	if len(output.Runtime.Lines) > 0 {
		lines := make([]map[string]interface{}, len(output.Runtime.Lines))
		for i, entry := range output.Runtime.Lines {
			line := map[string]interface{}{
				"timestamp": entry.Timestamp.Format("2006-01-02T15:04:05.000Z07:00"),
				"status":    f.convertSeverityToStatus(entry.Level),
				"message":   entry.Message,
				"logger":    "logspan",
			}

			// Add tags if present
			if len(entry.Tags) > 0 {
				line["tags"] = entry.Tags
			}

			// Add custom fields
			for k, v := range entry.Fields {
				line[k] = v
			}

			lines[i] = line
		}
		ddOutput["lines"] = lines
	}

	return ddOutput
}

// convertSeverityToStatus converts log level to DataDog status
func (f *DataDogFormatter) convertSeverityToStatus(severity string) string {
	switch severity {
	case "DEBUG":
		return "debug"
	case "INFO":
		return "info"
	case "WARN":
		return "warn"
	case "ERROR":
		return "error"
	case "CRITICAL":
		return "critical"
	default:
		return "info"
	}
}

// createSummaryMessage creates a summary message for the log output
func (f *DataDogFormatter) createSummaryMessage(output *LogOutput) string {
	if len(output.Runtime.Lines) == 0 {
		return "Empty log context"
	}
	if len(output.Runtime.Lines) == 1 {
		return output.Runtime.Lines[0].Message
	}
	return "Log context with multiple entries"
}
