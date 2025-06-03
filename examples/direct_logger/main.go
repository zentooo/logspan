package main

import (
	"bytes"
	"fmt"

	"github.com/zentooo/logspan/logger"
)

func main() {
	// Initialize logger with prettify enabled
	logger.Init(logger.Config{
		MinLevel:     logger.DebugLevel,
		Output:       nil, // Use default (stdout)
		PrettifyJSON: true,
	})

	fmt.Println("=== Basic Usage with logger.D (Recommended) ===")
	demonstrateBasicUsage()

	fmt.Println("\n=== Advanced Usage with Custom DirectLogger ===")
	demonstrateAdvancedUsage()
}

func demonstrateBasicUsage() {
	fmt.Println("Most common way to use direct logging:")

	// Using the global direct logger instance
	logger.D.Debugf("This is a debug message")
	logger.D.Infof("Application started successfully")
	logger.D.Warnf("Warning: Configuration file not found, using defaults")
	logger.D.Errorf("Error: Failed to connect to database: %v", "connection timeout")
	logger.D.Criticalf("Critical: System is running out of memory")

	fmt.Println("\nThe global logger.D instance uses the configuration from logger.Init()")
}

func demonstrateAdvancedUsage() {
	fmt.Println("For advanced scenarios where you need custom configuration:")

	// Create a new direct logger with custom settings
	dl := logger.NewDirectLogger()

	// Use a buffer to capture output for demonstration
	var buf bytes.Buffer
	dl.SetOutput(&buf)

	fmt.Println("\n--- Testing with INFO level (default) ---")
	dl.Debugf("This debug message should NOT appear")
	dl.Infof("This info message should appear")
	dl.Warnf("This warning message should appear")
	dl.Errorf("This error message should appear")
	dl.Criticalf("This critical message should appear")

	fmt.Print(buf.String())
	buf.Reset()

	fmt.Println("\n--- Testing with DEBUG level ---")
	dl.SetLevelFromString("DEBUG")
	dl.Debugf("This debug message should appear now")
	dl.Infof("This info message should appear")
	dl.Warnf("This warning message should appear")

	fmt.Print(buf.String())
	buf.Reset()

	fmt.Println("\n--- Testing with ERROR level ---")
	dl.SetLevelFromString("ERROR")
	dl.Debugf("This debug message should NOT appear")
	dl.Infof("This info message should NOT appear")
	dl.Warnf("This warning message should NOT appear")
	dl.Errorf("This error message should appear")
	dl.Criticalf("This critical message should appear")

	fmt.Print(buf.String())

	fmt.Println("\nUse NewDirectLogger() when you need:")
	fmt.Println("- Different log levels for different loggers")
	fmt.Println("- Different output destinations")
	fmt.Println("- Different formatters")
	fmt.Println("- Multiple logger instances with different configurations")
}
