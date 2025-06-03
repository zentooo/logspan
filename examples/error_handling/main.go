package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/zentooo/logspan/logger"
)

func main() {
	fmt.Println("=== LogSpan Error Handling Examples ===")

	// Example 1: Default Error Handler
	fmt.Println("Example 1: Default Error Handler")
	demonstrateDefaultErrorHandler()

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// Example 2: Custom Error Handler
	fmt.Println("Example 2: Custom Error Handler")
	demonstrateCustomErrorHandler()

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// Example 3: Silent Error Handler
	fmt.Println("Example 3: Silent Error Handler")
	demonstrateSilentErrorHandler()

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// Example 4: Error Handler with Context Logger
	fmt.Println("Example 4: Error Handler with Context Logger")
	demonstrateContextLoggerErrorHandling()

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// Example 5: Custom Error Handler Function
	fmt.Println("Example 5: Custom Error Handler Function")
	demonstrateCustomErrorHandlerFunction()
}

func demonstrateDefaultErrorHandler() {
	fmt.Println("Using default error handler (outputs to stderr)...")

	// Initialize logger with default configuration
	logger.Init(logger.DefaultConfig())

	// Create a logger with invalid output to trigger errors
	directLogger := logger.NewDirectLogger()

	// Create a broken writer that always returns an error
	brokenWriter := &brokenWriter{}
	directLogger.SetOutput(brokenWriter)

	fmt.Println("Attempting to log with broken output (check stderr for error messages):")
	directLogger.Infof("This will fail to write")

	fmt.Println("Default error handler will output error messages to stderr.")
}

func demonstrateCustomErrorHandler() {
	fmt.Println("Using custom error handler...")

	// Create a buffer to capture error messages
	var errorBuffer bytes.Buffer

	// Create custom error handler that writes to our buffer
	customHandler := logger.NewDefaultErrorHandlerWithOutput(&errorBuffer)

	// Set custom error handler
	logger.SetGlobalErrorHandler(customHandler)

	// Create a logger with broken output
	directLogger := logger.NewDirectLogger()
	brokenWriter := &brokenWriter{}
	directLogger.SetOutput(brokenWriter)

	fmt.Println("Attempting to log with broken output:")
	directLogger.Infof("This will fail to write")

	// Show captured error messages
	fmt.Printf("Captured error messages:\n%s", errorBuffer.String())
}

func demonstrateSilentErrorHandler() {
	fmt.Println("Using silent error handler...")

	// Set silent error handler
	logger.SetGlobalErrorHandler(&logger.SilentErrorHandler{})

	// Create a logger with broken output
	directLogger := logger.NewDirectLogger()
	brokenWriter := &brokenWriter{}
	directLogger.SetOutput(brokenWriter)

	fmt.Println("Attempting to log with broken output (errors will be silently ignored):")
	directLogger.Infof("This will fail to write")

	fmt.Println("No error messages should appear - errors are silently ignored.")
}

func demonstrateContextLoggerErrorHandling() {
	fmt.Println("Using context logger error handling...")

	// Create a buffer to capture error messages
	var errorBuffer bytes.Buffer
	customHandler := logger.NewDefaultErrorHandlerWithOutput(&errorBuffer)
	logger.SetGlobalErrorHandler(customHandler)

	// Create context and logger
	ctx := context.Background()
	contextLogger := logger.NewContextLogger()
	ctx = logger.WithLogger(ctx, contextLogger)

	// Set broken output
	brokenWriter := &brokenWriter{}
	contextLogger.SetOutput(brokenWriter)

	// Add some log entries
	logger.Infof(ctx, "First log entry")
	logger.Infof(ctx, "Second log entry")

	fmt.Println("Attempting to flush context logger with broken output:")
	logger.FlushContext(ctx)

	// Show captured error messages
	fmt.Printf("Captured error messages:\n%s", errorBuffer.String())

	// Test with working output
	fmt.Println("\nTesting flush with working output:")
	errorBuffer.Reset()
	contextLogger.SetOutput(os.Stdout)
	logger.Infof(ctx, "This should work")
	logger.FlushContext(ctx)

	if errorBuffer.Len() == 0 {
		fmt.Println("No errors occurred - logging successful")
	} else {
		fmt.Printf("Unexpected errors: %s", errorBuffer.String())
	}
}

func demonstrateCustomErrorHandlerFunction() {
	fmt.Println("Using custom error handler function...")

	// Create a custom error handler using ErrorHandlerFunc
	var errorCount int
	var lastError error
	var lastOperation string

	customFunc := logger.ErrorHandlerFunc(func(operation string, err error) {
		errorCount++
		lastError = err
		lastOperation = operation
		fmt.Printf("Custom handler: Operation '%s' failed with error: %v\n", operation, err)
	})

	// Set the custom function as the global error handler
	logger.SetGlobalErrorHandler(customFunc)

	// Create a logger with broken output
	directLogger := logger.NewDirectLogger()
	brokenWriter := &brokenWriter{}
	directLogger.SetOutput(brokenWriter)

	fmt.Println("Attempting multiple log operations:")
	directLogger.Infof("First message")
	directLogger.Warnf("Second message")
	directLogger.Errorf("Third message")

	fmt.Printf("\nSummary: %d errors occurred\n", errorCount)
	fmt.Printf("Last operation: %s\n", lastOperation)
	fmt.Printf("Last error: %v\n", lastError)
}

// brokenWriter is a writer that always returns an error
type brokenWriter struct{}

func (w *brokenWriter) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("simulated write error")
}
