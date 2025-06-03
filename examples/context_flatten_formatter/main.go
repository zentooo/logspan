package main

import (
	"fmt"
	"time"

	"github.com/zentooo/logspan/pkg/formatter"
)

func main() {
	fmt.Println("=== ContextFlattenFormatter Example ===")

	// Create sample log output
	output1 := &formatter.LogOutput{
		Type: "request",
		Context: map[string]interface{}{
			"request_id": "req-12345",
			"user_id":    "user-67890",
			"method":     "GET",
			"path":       "/api/users",
		},
		Runtime: formatter.RuntimeInfo{
			Severity:  "INFO",
			StartTime: "2023-10-27T09:59:58.123456+09:00",
			EndTime:   "2023-10-27T10:00:00.223456+09:00",
			Elapsed:   2100,
			Lines: []*formatter.LogEntry{
				{
					Timestamp: time.Date(2023, 10, 27, 9, 59, 59, 123456000, time.FixedZone("JST", 9*60*60)),
					Level:     "INFO",
					Message:   "Request processing started",
				},
				{
					Timestamp: time.Date(2023, 10, 27, 9, 59, 59, 500000000, time.FixedZone("JST", 9*60*60)),
					Level:     "DEBUG",
					Message:   "Validating input parameters",
				},
				{
					Timestamp: time.Date(2023, 10, 27, 10, 0, 0, 223456000, time.FixedZone("JST", 9*60*60)),
					Level:     "INFO",
					Message:   "Request processing completed",
				},
			},
		},
	}

	// Example 1: Compact JSON output
	fmt.Println("1. Compact JSON output:")
	compactFormatter := formatter.NewContextFlattenFormatter()
	result, err := compactFormatter.Format(output1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%s\n\n", result)

	// Example 2: Pretty-printed JSON output
	fmt.Println("2. Pretty-printed JSON output:")
	prettyFormatter := formatter.NewContextFlattenFormatterWithIndent("  ")
	result, err = prettyFormatter.Format(output1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%s\n\n", result)

	// Example 3: Key collision example
	fmt.Println("3. Key collision example (context overrides top-level fields):")
	output2 := &formatter.LogOutput{
		Type: "request",
		Context: map[string]interface{}{
			"request_id": "req-99999",
			"type":       "custom_request_type", // This will override the top-level "type"
			"runtime":    "custom_runtime_info", // This will override the top-level "runtime"
		},
		Runtime: formatter.RuntimeInfo{
			Severity:  "WARN",
			StartTime: "2023-10-27T10:01:00.000000+09:00",
			EndTime:   "2023-10-27T10:01:01.500000+09:00",
			Elapsed:   1500,
			Lines: []*formatter.LogEntry{
				{
					Timestamp: time.Date(2023, 10, 27, 10, 1, 0, 500000000, time.FixedZone("JST", 9*60*60)),
					Level:     "WARN",
					Message:   "Potential issue detected",
				},
			},
		},
	}
	result, err = prettyFormatter.Format(output2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%s\n\n", result)

	// Example 4: Empty context
	fmt.Println("4. Empty context example:")
	output3 := &formatter.LogOutput{
		Type:    "request",
		Context: map[string]interface{}{},
		Runtime: formatter.RuntimeInfo{
			Severity:  "ERROR",
			StartTime: "2023-10-27T10:02:00.000000+09:00",
			EndTime:   "2023-10-27T10:02:00.100000+09:00",
			Elapsed:   100,
			Lines: []*formatter.LogEntry{
				{
					Timestamp: time.Date(2023, 10, 27, 10, 2, 0, 50000000, time.FixedZone("JST", 9*60*60)),
					Level:     "ERROR",
					Message:   "Critical error occurred",
				},
			},
		},
	}
	result, err = prettyFormatter.Format(output3)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%s\n\n", result)
}
