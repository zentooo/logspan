// Package formatter provides various output formatters for log entries.
//
// The formatter package defines the interface and implementations for formatting
// log output in different formats. It supports JSON formatting with optional
// pretty-printing and context flattening capabilities.
//
// # Available Formatters
//
//   - JSONFormatter: Standard JSON output with optional indentation
//   - ContextFlattenFormatter: Flattens context fields to the top level
//
// # Basic Usage
//
// Using JSON formatter:
//
//	formatter := formatter.NewJSONFormatter()
//	output, err := formatter.Format(logOutput)
//
// Using JSON formatter with pretty-printing:
//
//	formatter := formatter.NewJSONFormatterWithIndent("  ")
//	output, err := formatter.Format(logOutput)
//
// Using context flatten formatter:
//
//	formatter := formatter.NewContextFlattenFormatter()
//	output, err := formatter.Format(logOutput)
//
// # Data Structures
//
// The package defines the following key structures:
//
//   - LogEntry: Represents a single log entry with timestamp, level, and message
//   - LogOutput: Complete log output structure with context and runtime information
//   - RuntimeInfo: Runtime information including severity, timing, and log entries
//   - Formatter: Interface for implementing custom formatters
//
// # JSON Format
//
// The standard JSON format produces output like:
//
//	{
//	  "type": "request",
//	  "context": {
//	    "request_id": "req-123"
//	  },
//	  "runtime": {
//	    "severity": "INFO",
//	    "startTime": "2023-10-27T09:59:58.123456+09:00",
//	    "endTime": "2023-10-27T10:00:00.223456+09:00",
//	    "elapsed": 150,
//	    "lines": [
//	      {
//	        "timestamp": "2023-10-27T09:59:59.123456+09:00",
//	        "level": "INFO",
//	        "message": "Processing started"
//	      }
//	    ]
//	  }
//	}
//
// # Context Flatten Format
//
// The context flatten format moves context fields to the top level:
//
//	{
//	  "request_id": "req-123",
//	  "type": "request",
//	  "runtime": {
//	    "severity": "INFO",
//	    "startTime": "2023-10-27T09:59:58.123456+09:00",
//	    "endTime": "2023-10-27T10:00:00.223456+09:00",
//	    "elapsed": 150,
//	    "lines": [
//	      {
//	        "timestamp": "2023-10-27T09:59:59.123456+09:00",
//	        "level": "INFO",
//	        "message": "Processing started"
//	      }
//	    ]
//	  }
//	}
//
// # Custom Formatters
//
// You can implement custom formatters by implementing the Formatter interface:
//
//	type CustomFormatter struct{}
//
//	func (f *CustomFormatter) Format(output *LogOutput) ([]byte, error) {
//	    // Custom formatting logic
//	    return []byte("custom format"), nil
//	}
//
// # Integration with Logger
//
// Formatters are typically used with the logger package:
//
//	directLogger := logger.NewDirectLogger()
//	directLogger.SetFormatter(formatter.NewContextFlattenFormatter())
//
//	contextLogger := logger.NewContextLogger()
//	contextLogger.SetFormatter(formatter.NewJSONFormatterWithIndent("  "))
package formatter
