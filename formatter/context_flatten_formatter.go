package formatter

import (
	"encoding/json"
)

// ContextFlattenFormatter implements the Formatter interface for flattened JSON output
// It flattens the context fields to the top level of the JSON output
type ContextFlattenFormatter struct {
	// Indent specifies the indentation string for pretty printing
	// Empty string means no indentation (compact JSON)
	Indent string
}

// NewContextFlattenFormatter creates a new ContextFlattenFormatter instance
func NewContextFlattenFormatter() *ContextFlattenFormatter {
	return &ContextFlattenFormatter{
		Indent: "", // Default to compact JSON
	}
}

// NewContextFlattenFormatterWithIndent creates a new ContextFlattenFormatter with indentation
func NewContextFlattenFormatterWithIndent(indent string) *ContextFlattenFormatter {
	return &ContextFlattenFormatter{
		Indent: indent,
	}
}

// Format formats the log output with flattened context fields
func (f *ContextFlattenFormatter) Format(output *LogOutput) ([]byte, error) {
	// Create a new map to hold the flattened output
	flattened := make(map[string]interface{})

	// First, add all context fields to the top level
	for key, value := range output.Context {
		flattened[key] = value
	}

	// Then add the other top-level fields
	// Note: If context contains keys that conflict with these, context values take precedence
	if _, exists := flattened["type"]; !exists {
		flattened["type"] = output.Type
	}
	if _, exists := flattened["runtime"]; !exists {
		flattened["runtime"] = output.Runtime
	}

	// Marshal the flattened structure
	if f.Indent == "" {
		return json.Marshal(flattened)
	}
	return json.MarshalIndent(flattened, "", f.Indent)
}
