package main

import (
	"bytes"
	"fmt"

	"github.com/zentooo/logspan/pkg/formatter"
	"github.com/zentooo/logspan/pkg/logger"
)

func main() {
	// Initialize logger with prettify enabled
	logger.Init(logger.Config{
		MinLevel:     logger.DebugLevel,
		Output:       nil, // Use default (stdout)
		PrettifyJSON: true,
	})

	// Create a new direct logger
	dl := logger.NewDirectLogger()

	// Use a buffer to capture output for demonstration
	var buf bytes.Buffer
	dl.SetOutput(&buf)

	fmt.Println("=== Testing with INFO level (default) ===")
	dl.Debugf("This debug message should NOT appear")
	dl.Infof("This info message should appear")
	dl.Warnf("This warning message should appear")
	dl.Errorf("This error message should appear")
	dl.Criticalf("This critical message should appear")

	fmt.Print(buf.String())
	buf.Reset()

	fmt.Println("\n=== Testing with DEBUG level ===")
	dl.SetLevelFromString("DEBUG")
	dl.Debugf("This debug message should appear now")
	dl.Infof("This info message should appear")
	dl.Warnf("This warning message should appear")

	fmt.Print(buf.String())
	buf.Reset()

	fmt.Println("\n=== Testing with ERROR level ===")
	dl.SetLevelFromString("ERROR")
	dl.Debugf("This debug message should NOT appear")
	dl.Infof("This info message should NOT appear")
	dl.Warnf("This warning message should NOT appear")
	dl.Errorf("This error message should appear")
	dl.Criticalf("This critical message should appear")

	fmt.Print(buf.String())

	fmt.Println("\n=== DataDog Formatter Example ===")

	// Create a new direct logger with DataDog formatter
	ddLogger := logger.NewDirectLogger()
	ddLogger.SetFormatter(formatter.NewDataDogFormatterWithIndent("  "))

	// Use a buffer to capture DataDog formatted output
	var ddBuf bytes.Buffer
	ddLogger.SetOutput(&ddBuf)

	fmt.Println("Testing DataDog formatter with different log levels:")
	ddLogger.SetLevelFromString("DEBUG")
	ddLogger.Debugf("Debug message in DataDog format")
	ddLogger.Infof("Info message in DataDog format")
	ddLogger.Warnf("Warning message in DataDog format")
	ddLogger.Errorf("Error message in DataDog format")
	ddLogger.Criticalf("Critical message in DataDog format")

	fmt.Print(ddBuf.String())
}
