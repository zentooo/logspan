package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/zentooo/logspan/logger"
)

func main() {
	fmt.Println("=== Middleware Example ===")
	fmt.Println("This example demonstrates how to use middleware, especially password masking middleware.")
	fmt.Println()

	// Initialize logger with basic configuration
	logger.Init(
		logger.WithMinLevel(logger.DebugLevel),
		logger.WithOutput(os.Stdout),
		logger.WithPrettifyJSON(true),
		logger.WithMaxLogEntries(0), // No auto-flush for this example
	)

	// Create password masking middleware with default settings
	passwordMasker := logger.NewPasswordMaskingMiddleware()
	logger.AddMiddleware(passwordMasker.Middleware())

	fmt.Println("1. Basic Password Masking:")
	demonstrateBasicPasswordMasking()

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("2. Custom Password Masking Configuration:")
	demonstrateCustomPasswordMasking()

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("3. Multiple Middleware Chain:")
	demonstrateMultipleMiddleware()
}

func demonstrateBasicPasswordMasking() {
	ctx := context.Background()
	contextLogger := logger.NewContextLogger()
	ctx = logger.WithLogger(ctx, contextLogger)

	// Add context information
	logger.AddContextValue(ctx, "user_id", "user123")
	logger.AddContextValue(ctx, "request_id", "req-456")

	// Log messages with sensitive information
	logger.Infof(ctx, "User login attempt with password=secret123")
	logger.Debugf(ctx, "API call with token=abc123xyz and api_key=super_secret")
	logger.Warnf(ctx, "Authentication failed for user with passwd=wrongpass")
	logger.Errorf(ctx, "Database connection failed: {\"password\":\"dbpass123\", \"host\":\"localhost\"}")

	// Flush the logs
	logger.FlushContext(ctx)
}

func demonstrateCustomPasswordMasking() {
	// Clear existing middleware
	logger.ClearMiddleware()

	// Create custom password masking middleware
	customMasker := logger.NewPasswordMaskingMiddleware().
		WithMaskString("[REDACTED]").
		WithPasswordKeys([]string{"password", "secret", "token", "key"}).
		AddPasswordKey("custom_secret").
		AddPasswordKey("private_key")

	logger.AddMiddleware(customMasker.Middleware())

	ctx := context.Background()
	contextLogger := logger.NewContextLogger()
	ctx = logger.WithLogger(ctx, contextLogger)

	logger.AddContextValue(ctx, "service", "auth-service")

	// Log with custom sensitive keys
	logger.Infof(ctx, "Service started with custom_secret=my_custom_secret")
	logger.Debugf(ctx, "Encryption key loaded: private_key=rsa_private_key_content")
	logger.Warnf(ctx, "Configuration: {\"secret\":\"config_secret\", \"key\":\"encryption_key\"}")

	logger.FlushContext(ctx)
}

func demonstrateMultipleMiddleware() {
	// Clear existing middleware
	logger.ClearMiddleware()

	// Add password masking middleware
	passwordMasker := logger.NewPasswordMaskingMiddleware()
	logger.AddMiddleware(passwordMasker.Middleware())

	// Add custom logging middleware
	logger.AddMiddleware(func(entry *logger.LogEntry, next func(*logger.LogEntry)) {
		// Add timestamp prefix to all messages
		entry.Message = fmt.Sprintf("[%s] %s", entry.Timestamp.Format("15:04:05"), entry.Message)
		next(entry)
	})

	// Add another custom middleware for message transformation
	logger.AddMiddleware(func(entry *logger.LogEntry, next func(*logger.LogEntry)) {
		// Add log level prefix
		entry.Message = fmt.Sprintf("[%s] %s", entry.Level, entry.Message)
		next(entry)
	})

	fmt.Printf("Total middleware count: %d\n", logger.GetMiddlewareCount())

	ctx := context.Background()
	contextLogger := logger.NewContextLogger()
	ctx = logger.WithLogger(ctx, contextLogger)

	logger.AddContextValue(ctx, "component", "middleware-demo")

	// Log messages that will be processed by all middleware
	logger.Infof(ctx, "Processing user login with password=userpass123")
	logger.Debugf(ctx, "System status check completed")
	logger.Warnf(ctx, "Rate limit approaching for token=rate_limit_token")

	logger.FlushContext(ctx)

	// Clear middleware for cleanup
	logger.ClearMiddleware()
	fmt.Printf("Middleware cleared. Current count: %d\n", logger.GetMiddlewareCount())
}
