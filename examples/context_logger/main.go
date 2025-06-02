package main

import (
	"context"
	"fmt"

	"github.com/zentooo/logspan/pkg/formatter"
	"github.com/zentooo/logspan/pkg/logger"
)

func main() {
	fmt.Println("=== Context Logger Example ===")

	// Initialize logger with prettify enabled
	logger.Init(logger.Config{
		MinLevel:     logger.DebugLevel,
		Output:       nil, // Use default (stdout)
		PrettifyJSON: true,
	})

	// Create a context with a logger
	ctx := context.Background()
	contextLogger := logger.NewContextLogger()
	ctx = logger.WithLogger(ctx, contextLogger)

	// Add context fields
	logger.AddContextValue(ctx, "request_id", "req-12345")
	logger.AddContextValue(ctx, "user_id", "user-67890")

	// Simulate a request processing
	processRequest(ctx)

	// Flush the logger to output all accumulated logs
	logger.FlushContext(ctx)

	fmt.Println("\n=== Direct Logger Example for Comparison ===")

	// Compare with direct logger
	logger.D.Infof("Direct logger message 1")
	logger.D.Warnf("Direct logger message 2")
	logger.D.Errorf("Direct logger message 3")

	fmt.Println("\n=== DataDog Formatter Example ===")

	// Create a context with a logger using DataDog formatter
	ctx2 := context.Background()
	contextLogger2 := logger.NewContextLogger()

	// Set DataDog formatter for this logger
	contextLogger2.SetFormatter(formatter.NewDataDogFormatterWithIndent("  "))

	ctx2 = logger.WithLogger(ctx2, contextLogger2)

	// Add context fields that will appear as custom attributes in DataDog
	logger.AddContextValue(ctx2, "request_id", "req-datadog-123")
	logger.AddContextValue(ctx2, "service", "example-service")
	logger.AddContextValue(ctx2, "environment", "development")

	// Log some messages
	logger.Infof(ctx2, "DataDog formatter example started")
	logger.AddContextValue(ctx2, "operation", "demo")
	logger.Debugf(ctx2, "This is a debug message in DataDog format")
	logger.Warnf(ctx2, "This is a warning message in DataDog format")
	logger.Infof(ctx2, "DataDog formatter example completed")

	// Flush the logger to output in DataDog Standard Attributes format
	logger.FlushContext(ctx2)
}

func processRequest(ctx context.Context) {
	logger.Infof(ctx, "Starting request processing")

	// Add more context during processing
	logger.AddContextValue(ctx, "step", "validation")
	logger.Debugf(ctx, "Validating input parameters")

	// Simulate some processing steps
	validateInput(ctx)
	processData(ctx)
	generateResponse(ctx)

	logger.Infof(ctx, "Request processing completed")
}

func validateInput(ctx context.Context) {
	logger.AddContextValue(ctx, "validation_step", "input_check")
	logger.Debugf(ctx, "Checking input format")
	logger.Infof(ctx, "Input validation passed")
}

func processData(ctx context.Context) {
	logger.AddContextValue(ctx, "processing_step", "data_transformation")
	logger.Debugf(ctx, "Transforming data")
	logger.Warnf(ctx, "Using deprecated API for compatibility")
	logger.Infof(ctx, "Data processing completed")
}

func generateResponse(ctx context.Context) {
	logger.AddContextValue(ctx, "response_step", "serialization")
	logger.Debugf(ctx, "Serializing response")
	logger.Infof(ctx, "Response generated successfully")
}
