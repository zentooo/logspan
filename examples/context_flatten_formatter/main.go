package main

import (
	"fmt"
	"time"

	"github.com/zentooo/logspan/pkg/formatter"
)

func main() {
	fmt.Println("=== ContextFlattenFormatter Example ===")

	// Create sample log output
	output := &formatter.LogOutput{
		Type: "request",
		Context: map[string]interface{}{
			"user_id":    "user-12345",
			"request_id": "req-67890",
			"session_id": "sess-abcdef",
			"method":     "POST",
			"path":       "/api/users",
		},
		Runtime: formatter.RuntimeInfo{
			Severity:  "INFO",
			StartTime: "2023-10-27T10:00:00.123456+09:00",
			EndTime:   "2023-10-27T10:00:01.234567+09:00",
			Elapsed:   1111,
			Lines: []*formatter.LogEntry{
				{
					Timestamp: time.Date(2023, 10, 27, 10, 0, 0, 123456000, time.FixedZone("JST", 9*3600)),
					Level:     "INFO",
					Message:   "Request started",
				},
				{
					Timestamp: time.Date(2023, 10, 27, 10, 0, 1, 234567000, time.FixedZone("JST", 9*3600)),
					Level:     "INFO",
					Message:   "Request completed successfully",
				},
			},
		},
		Config: formatter.ConfigInfo{
			ElapsedUnit: "ms",
		},
	}

	// Example 1: Compact JSON output
	fmt.Println("1. Compact JSON output:")
	compactFormatter := formatter.NewContextFlattenFormatter()
	result, err := compactFormatter.Format(output)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%s\n\n", result)

	// Example 2: Pretty-printed JSON output
	fmt.Println("2. Pretty-printed JSON output:")
	prettyFormatter := formatter.NewContextFlattenFormatterWithIndent("  ")
	result, err = prettyFormatter.Format(output)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%s\n\n", result)

	// Example 3: Key collision example
	fmt.Println("3. Key collision example (context overrides top-level fields):")
	collisionOutput := &formatter.LogOutput{
		Type: "request",
		Context: map[string]interface{}{
			"type":    "custom_request_type", // This will override the top-level "type"
			"user_id": "user-99999",
		},
		Runtime: formatter.RuntimeInfo{
			Severity: "WARN",
			Lines:    []*formatter.LogEntry{},
		},
		Config: formatter.ConfigInfo{
			ElapsedUnit: "ms",
		},
	}

	result, err = prettyFormatter.Format(collisionOutput)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%s\n\n", result)

	// Example 4: Empty context
	fmt.Println("4. Empty context example:")
	emptyContextOutput := &formatter.LogOutput{
		Type:    "request",
		Context: map[string]interface{}{},
		Runtime: formatter.RuntimeInfo{
			Severity: "DEBUG",
			Lines:    []*formatter.LogEntry{},
		},
		Config: formatter.ConfigInfo{
			ElapsedUnit: "ms",
		},
	}

	result, err = prettyFormatter.Format(emptyContextOutput)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%s\n", result)
}
