package formatter

import (
	"encoding/json"
)

// JSONFormatter implements the Formatter interface for JSON output
type JSONFormatter struct {
	// Indent specifies the indentation string for pretty printing
	// Empty string means no indentation (compact JSON)
	Indent string
}

// NewJSONFormatter creates a new JSONFormatter instance
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{
		Indent: "", // Default to compact JSON
	}
}

// NewJSONFormatterWithIndent creates a new JSONFormatter with indentation
func NewJSONFormatterWithIndent(indent string) *JSONFormatter {
	return &JSONFormatter{
		Indent: indent,
	}
}

// Format formats the log output as JSON
func (f *JSONFormatter) Format(output *LogOutput) ([]byte, error) {
	if f.Indent == "" {
		return json.Marshal(output)
	}
	return json.MarshalIndent(output, "", f.Indent)
}
