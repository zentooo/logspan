package main

import (
	"context"
	"fmt"

	"github.com/zentooo/logspan/pkg/formatter"
	"github.com/zentooo/logspan/pkg/logger"
)

func main() {
	fmt.Println("=== DataDog Formatter Example ===")

	// Initialize logger with prettify enabled for better readability
	logger.Init(logger.Config{
		MinLevel:     logger.DebugLevel,
		Output:       nil, // Use default (stdout)
		PrettifyJSON: true,
	})

	// Create a context with a logger
	ctx := context.Background()
	contextLogger := logger.NewContextLogger()

	// Set DataDog formatter
	contextLogger.SetFormatter(formatter.NewDataDogFormatterWithIndent("  "))

	ctx = logger.WithLogger(ctx, contextLogger)

	// Add context fields that will appear as custom attributes in DataDog
	logger.AddContextValue(ctx, "request_id", "req-12345")
	logger.AddContextValue(ctx, "user_id", "user-67890")
	logger.AddContextValue(ctx, "service", "user-service")
	logger.AddContextValue(ctx, "version", "1.2.3")

	// Simulate a request processing
	processRequest(ctx)

	// Flush the logger to output all accumulated logs in DataDog format
	logger.FlushContext(ctx)

	fmt.Println("\n=== Comparison: Standard JSON Format ===")

	// Compare with standard JSON formatter
	ctx2 := context.Background()
	contextLogger2 := logger.NewContextLogger()
	// Use default JSON formatter (no explicit SetFormatter call)
	ctx2 = logger.WithLogger(ctx2, contextLogger2)

	logger.AddContextValue(ctx2, "request_id", "req-67890")
	logger.AddContextValue(ctx2, "comparison", "standard_json")

	logger.Infof(ctx2, "Standard JSON format example")
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

	logger.Infof(ctx, "Request processing completed successfully")
}

func validateInput(ctx context.Context) {
	logger.AddContextValue(ctx, "validation_step", "input_check")
	logger.Debugf(ctx, "Checking input format and constraints")
	logger.Infof(ctx, "Input validation passed")
}

func processData(ctx context.Context) {
	logger.AddContextValue(ctx, "processing_step", "data_transformation")
	logger.Debugf(ctx, "Transforming data according to business rules")
	logger.Warnf(ctx, "Using deprecated API for backward compatibility")
	logger.Infof(ctx, "Data processing completed")
}

func generateResponse(ctx context.Context) {
	logger.AddContextValue(ctx, "response_step", "serialization")
	logger.Debugf(ctx, "Serializing response to JSON format")
	logger.Infof(ctx, "Response generated and ready to send")
}
