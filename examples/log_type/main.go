package main

import (
	"context"
	"os"

	"github.com/zentooo/logspan/logger"
)

func main() {
	// Example 1: Default log type (request)
	logger.Init() // Use default options

	ctx := context.Background()
	contextLogger := logger.NewContextLogger()
	ctx = logger.WithLogger(ctx, contextLogger)

	logger.AddContextValue(ctx, "user_id", "user-123")
	logger.Infof(ctx, "Processing with default log type")
	logger.FlushContext(ctx)

	// Example 2: Custom log type for batch processing
	logger.Init(
		logger.WithMinLevel(logger.InfoLevel),
		logger.WithOutput(os.Stdout),
		logger.WithSourceInfo(false),
		logger.WithPrettifyJSON(true),
		logger.WithMaxLogEntries(0),
		logger.WithLogType("batch_job"),
	)

	ctx2 := context.Background()
	contextLogger2 := logger.NewContextLogger()
	ctx2 = logger.WithLogger(ctx2, contextLogger2)

	logger.AddContextValue(ctx2, "job_id", "job-456")
	logger.AddContextValue(ctx2, "batch_size", 1000)
	logger.Infof(ctx2, "Starting batch processing")
	logger.Infof(ctx2, "Processing 1000 records")
	logger.Infof(ctx2, "Batch processing completed")
	logger.FlushContext(ctx2)

	// Example 3: Custom log type for API operations
	logger.Init(
		logger.WithMinLevel(logger.InfoLevel),
		logger.WithOutput(os.Stdout),
		logger.WithSourceInfo(false),
		logger.WithPrettifyJSON(true),
		logger.WithMaxLogEntries(0),
		logger.WithLogType("api_operation"),
	)

	ctx3 := context.Background()
	contextLogger3 := logger.NewContextLogger()
	ctx3 = logger.WithLogger(ctx3, contextLogger3)

	logger.AddContextValue(ctx3, "operation", "create_user")
	logger.AddContextValue(ctx3, "api_version", "v1")
	logger.Infof(ctx3, "API operation started")
	logger.Infof(ctx3, "Validating input parameters")
	logger.Infof(ctx3, "Creating user record")
	logger.Infof(ctx3, "API operation completed successfully")
	logger.FlushContext(ctx3)

	// Example 4: Direct logger with custom log type
	directLogger := logger.NewDirectLogger()
	directLogger.Infof("Direct log with custom log type")
}
