package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/zentooo/logspan/formatter"
	"github.com/zentooo/logspan/logger"
)

func main() {
	fmt.Println("=== Advanced Configuration Example ===")
	fmt.Println("This example demonstrates various advanced configuration options.")
	fmt.Println()

	// Demonstrate different configuration scenarios
	demonstrateBasicConfig()
	fmt.Println("\n" + strings.Repeat("=", 50))

	demonstrateCustomOutput()
	fmt.Println("\n" + strings.Repeat("=", 50))

	demonstrateAutoFlushConfig()
	fmt.Println("\n" + strings.Repeat("=", 50))

	demonstrateFormatterConfig()
	fmt.Println("\n" + strings.Repeat("=", 50))

	demonstrateDirectLoggerConfig()
}

func demonstrateBasicConfig() {
	fmt.Println("1. Basic Configuration with Different Log Levels:")

	// Configure with INFO level
	logger.Init(logger.Config{
		MinLevel:         logger.InfoLevel,
		Output:           os.Stdout,
		EnableSourceInfo: true,
		PrettifyJSON:     true,
		MaxLogEntries:    0,
	})

	ctx := context.Background()
	contextLogger := logger.NewContextLogger()
	ctx = logger.WithLogger(ctx, contextLogger)

	logger.AddContextValue(ctx, "config_demo", "basic")

	// These will be logged (INFO and above)
	logger.Infof(ctx, "This INFO message will be shown")
	logger.Warnf(ctx, "This WARN message will be shown")
	logger.Errorf(ctx, "This ERROR message will be shown")

	// This will be filtered out
	logger.Debugf(ctx, "This DEBUG message will be filtered out")

	logger.FlushContext(ctx)
}

func demonstrateCustomOutput() {
	fmt.Println("2. Custom Output Configuration:")

	// Create a temporary file for logging
	logFile, err := os.CreateTemp("", "logspan_example_*.log")
	if err != nil {
		fmt.Printf("Failed to create temp file: %v\n", err)
		return
	}
	defer os.Remove(logFile.Name())
	defer logFile.Close()

	// Configure to write to file
	logger.Init(logger.Config{
		MinLevel:      logger.DebugLevel,
		Output:        logFile,
		PrettifyJSON:  false, // Compact JSON for file output
		MaxLogEntries: 0,
	})

	ctx := context.Background()
	contextLogger := logger.NewContextLogger()
	ctx = logger.WithLogger(ctx, contextLogger)

	logger.AddContextValue(ctx, "output", "file")
	logger.AddContextValue(ctx, "filename", logFile.Name())

	logger.Infof(ctx, "This message is written to file: %s", logFile.Name())
	logger.Debugf(ctx, "File logging configuration active")

	logger.FlushContext(ctx)

	// Read and display file contents
	logFile.Seek(0, 0)
	content := make([]byte, 1024)
	n, _ := logFile.Read(content)
	fmt.Printf("File contents:\n%s\n", string(content[:n]))

	// Reset to stdout for other examples
	logger.Init(logger.Config{
		MinLevel:     logger.DebugLevel,
		Output:       os.Stdout,
		PrettifyJSON: true,
	})
}

func demonstrateAutoFlushConfig() {
	fmt.Println("3. Auto-Flush Configuration:")

	// Configure with small MaxLogEntries for demonstration
	logger.Init(logger.Config{
		MinLevel:      logger.DebugLevel,
		Output:        os.Stdout,
		PrettifyJSON:  true,
		MaxLogEntries: 3, // Auto-flush after 3 entries
	})

	ctx := context.Background()
	contextLogger := logger.NewContextLogger()
	ctx = logger.WithLogger(ctx, contextLogger)

	logger.AddContextValue(ctx, "auto_flush", "demo")

	fmt.Println("Logging 5 messages (auto-flush at 3 entries):")

	// These will trigger auto-flush
	logger.Infof(ctx, "Message 1")
	logger.Infof(ctx, "Message 2")
	logger.Infof(ctx, "Message 3") // Auto-flush happens here

	fmt.Println("--- Auto-flush occurred ---")

	logger.Infof(ctx, "Message 4")
	logger.Infof(ctx, "Message 5")

	// Manual flush for remaining messages
	fmt.Println("--- Manual flush for remaining messages ---")
	logger.FlushContext(ctx)
}

func demonstrateFormatterConfig() {
	fmt.Println("4. Custom Formatter Configuration:")

	// Reset to basic config
	logger.Init(logger.Config{
		MinLevel:     logger.InfoLevel,
		Output:       os.Stdout,
		PrettifyJSON: true,
	})

	// Create context logger with custom formatter
	contextLogger := logger.NewContextLogger()
	contextLogger.SetFormatter(formatter.NewContextFlattenFormatterWithIndent("  "))

	ctx := context.Background()
	ctx = logger.WithLogger(ctx, contextLogger)

	logger.AddContextValue(ctx, "formatter", "context_flatten")
	logger.AddContextValue(ctx, "feature", "custom_formatting")

	fmt.Println("Using Context Flatten Formatter:")
	logger.Infof(ctx, "This message uses context flatten formatter")
	logger.Warnf(ctx, "Context fields are flattened to top level")

	logger.FlushContext(ctx)
}

func demonstrateDirectLoggerConfig() {
	fmt.Println("5. Direct Logger Configuration:")

	// Create and configure direct logger
	directLogger := logger.NewDirectLogger()
	directLogger.SetLevel(logger.WarnLevel)
	directLogger.SetFormatter(formatter.NewJSONFormatterWithIndent("  "))

	fmt.Println("Direct logger with WARN level filter:")

	// These will be filtered out
	directLogger.Debugf("This DEBUG message will be filtered")
	directLogger.Infof("This INFO message will be filtered")

	// These will be shown
	directLogger.Warnf("This WARN message will be shown")
	directLogger.Errorf("This ERROR message will be shown")
	directLogger.Criticalf("This CRITICAL message will be shown")

	// Demonstrate level setting from string
	fmt.Println("\nChanging level to DEBUG via string:")
	directLogger.SetLevelFromString("DEBUG")
	directLogger.Debugf("Now DEBUG messages are shown")
	directLogger.Infof("And INFO messages too")
}
