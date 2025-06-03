// Package logger provides a flexible and extensible logging library for Go applications.
//
// The logger package offers two main logging approaches:
//   - Context-based logging: Accumulates log entries and flushes them together with shared context
//   - Direct logging: Immediately outputs individual log entries
//
// # Key Features
//
//   - Dual-mode logging (context-based and direct)
//   - Transparent context management through HTTP middleware
//   - Extensible middleware system for log processing
//   - Built-in password masking for sensitive information
//   - Multiple output formatters (JSON, Context Flatten)
//   - Configurable log levels and output destinations
//   - Thread-safe operations
//   - Memory-efficient with auto-flush capabilities
//
// # Basic Usage
//
// Initialize the logger with configuration:
//
//	logger.Init(logger.Config{
//	    MinLevel:      logger.InfoLevel,
//	    Output:        os.Stdout,
//	    PrettifyJSON:  true,
//	    MaxLogEntries: 1000,
//	})
//
// # Context-based Logging
//
// Context-based logging accumulates log entries and flushes them together:
//
//	ctx := context.Background()
//	contextLogger := logger.NewContextLogger()
//	ctx = logger.WithLogger(ctx, contextLogger)
//
//	// Add context information
//	logger.AddContextValue(ctx, "request_id", "req-123")
//	logger.AddContextValue(ctx, "user_id", "user-456")
//
//	// Log entries
//	logger.Infof(ctx, "Processing started")
//	logger.Debugf(ctx, "Validating input: %s", input)
//	logger.Infof(ctx, "Processing completed")
//
//	// Flush accumulated logs
//	logger.FlushContext(ctx)
//
// # Direct Logging
//
// Direct logging outputs entries immediately:
//
//	// Using global direct logger
//	logger.D.Infof("Application started")
//	logger.D.Errorf("Error occurred: %v", err)
//
//	// Using custom direct logger
//	directLogger := logger.NewDirectLogger()
//	directLogger.SetLevel(logger.WarnLevel)
//	directLogger.Infof("This won't be logged (below warn level)")
//	directLogger.Warnf("This will be logged")
//
// # Middleware System
//
// The logger supports a middleware system for processing log entries:
//
//	// Add password masking middleware
//	passwordMasker := logger.NewPasswordMaskingMiddleware()
//	logger.AddMiddleware(passwordMasker.Middleware())
//
//	// Add custom middleware
//	logger.AddMiddleware(func(entry *logger.LogEntry, next func(*logger.LogEntry)) {
//	    // Process entry before logging
//	    entry.Message = strings.ToUpper(entry.Message)
//	    next(entry)
//	})
//
// # HTTP Integration
//
// Automatic HTTP request logging with context setup:
//
//	import "github.com/zentooo/logspan/pkg/http_middleware"
//
//	handler := http_middleware.LoggingMiddleware(yourHandler)
//	http.Handle("/", handler)
//
// # Configuration Options
//
// The Config struct provides various configuration options:
//
//   - MinLevel: Minimum log level for filtering (DebugLevel, InfoLevel, WarnLevel, ErrorLevel, CriticalLevel)
//   - Output: Output destination (io.Writer)
//   - EnableSourceInfo: Enable source file information in logs
//   - PrettifyJSON: Enable pretty-printed JSON output
//   - MaxLogEntries: Maximum log entries before auto-flush (0 = no limit)
//
// # Log Levels
//
// Available log levels in order of severity:
//   - DebugLevel: Detailed debugging information
//   - InfoLevel: General informational messages
//   - WarnLevel: Warning messages
//   - ErrorLevel: Error messages
//   - CriticalLevel: Critical error messages
//
// # Memory Management
//
// The context logger supports auto-flush to manage memory usage:
//
//	logger.Init(logger.Config{
//	    MaxLogEntries: 100, // Auto-flush after 100 entries
//	})
//
// When the number of accumulated log entries reaches MaxLogEntries,
// the logger automatically flushes the entries and continues accumulating new ones.
//
// # Thread Safety
//
// All logger operations are thread-safe and can be used concurrently
// from multiple goroutines without additional synchronization.
//
// # Output Format
//
// The default output format is JSON with the following structure:
//
//	{
//	  "type": "request",
//	  "context": {
//	    "request_id": "req-123",
//	    "user_id": "user-456"
//	  },
//	  "runtime": {
//	    "severity": "INFO",
//	    "startTime": "2023-10-27T09:59:58.123456+09:00",
//	    "endTime": "2023-10-27T10:00:00.223456+09:00",
//	    "elapsed": 150,
//	    "lines": [
//	      {
//	        "timestamp": "2023-10-27T09:59:59.123456+09:00",
//	        "level": "INFO",
//	        "message": "Processing started"
//	      }
//	    ]
//	  }
//	}
package logger
