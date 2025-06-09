package main

import (
	"context"
	"fmt"

	"github.com/zentooo/logspan/logger"
)

func main() {
	fmt.Println("=== Context Logger Example ===")

	// Initialize logger with prettify enabled
	logger.Init(
		logger.WithMinLevel(logger.DebugLevel),
		// Output defaults to stdout when not specified
		logger.WithPrettifyJSON(true),
	)

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
