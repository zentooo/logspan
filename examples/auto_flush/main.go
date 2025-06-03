package main

import (
	"context"
	"fmt"
	"time"

	"github.com/zentooo/logspan/pkg/logger"
)

func main() {
	fmt.Println("=== Auto-Flush Example ===")

	// Initialize logger with small MaxLogEntries for demonstration
	logger.Init(logger.Config{
		MinLevel:      logger.DebugLevel,
		Output:        nil, // Use default (stdout)
		PrettifyJSON:  true,
		MaxLogEntries: 3, // Auto-flush after 3 entries
	})

	// Create a context with a logger
	ctx := context.Background()
	contextLogger := logger.NewContextLogger()
	ctx = logger.WithLogger(ctx, contextLogger)

	// Add context fields
	logger.AddContextValue(ctx, "demo_id", "auto-flush-demo")
	logger.AddContextValue(ctx, "batch", 1)

	fmt.Println("\n--- Adding 3 log entries (should trigger auto-flush) ---")

	// Add entries - should auto-flush after 3rd entry
	logger.Infof(ctx, "Entry 1: Starting process")
	time.Sleep(100 * time.Millisecond) // Small delay to show different timestamps

	logger.Warnf(ctx, "Entry 2: Warning occurred")
	time.Sleep(100 * time.Millisecond)

	logger.Infof(ctx, "Entry 3: Process continuing") // This should trigger auto-flush
	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n--- Adding more entries for second batch ---")

	// Update context for second batch
	logger.AddContextValue(ctx, "batch", 2)

	// Add more entries - should auto-flush after 3rd entry again
	logger.Debugf(ctx, "Entry 4: Debug information")
	time.Sleep(100 * time.Millisecond)

	logger.Errorf(ctx, "Entry 5: Error occurred")
	time.Sleep(100 * time.Millisecond)

	logger.Criticalf(ctx, "Entry 6: Critical issue") // This should trigger second auto-flush
	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n--- Adding partial batch (no auto-flush) ---")

	// Update context for third batch
	logger.AddContextValue(ctx, "batch", 3)

	// Add only 2 entries - should not auto-flush
	logger.Infof(ctx, "Entry 7: Partial batch entry 1")
	time.Sleep(100 * time.Millisecond)

	logger.Infof(ctx, "Entry 8: Partial batch entry 2")
	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n--- Manual flush for remaining entries ---")

	// Manual flush for remaining entries
	logger.FlushContext(ctx)

	fmt.Println("\n=== Demo completed ===")
	fmt.Println("Notice how:")
	fmt.Println("1. First 3 entries were auto-flushed as one batch")
	fmt.Println("2. Next 3 entries were auto-flushed as second batch")
	fmt.Println("3. Last 2 entries required manual flush")
	fmt.Println("4. Each batch has its own start/end time and elapsed duration")
	fmt.Println("5. Context fields are preserved across auto-flushes")
}
