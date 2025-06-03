// Package logspan provides a flexible and extensible logging library for Go applications.
//
// LogSpan is a **zero-dependency** structured logging library for Go that provides
// **context-based log aggregation** by **consolidating multiple log entries into a
// single JSON structure**. This enables unified log management per HTTP request or
// batch processing unit. Unlike traditional scattered logging, LogSpan outputs all
// related logs in a single JSON structure, improving log analysis and troubleshooting
// efficiency.
//
// # Key Features
//
//   - 🔗 Context-based Log Aggregation: Consolidates multiple log entries into a single JSON for unified log management
//   - 🚀 Zero Dependencies: Operates solely with Go standard library, no external dependencies required
//   - 💾 Memory Efficient: Automatic memory pooling and configurable auto-flush to minimize memory footprint
//   - Dual-mode logging: Context-based and direct logging modes
//   - Structured log output: Consistent JSON-formatted log output
//   - Middleware mechanism: Customizable log processing pipeline
//   - Context flattening: Formatter that expands context fields to top level
//   - HTTP middleware: Automatic log setup for web applications
//   - Concurrency-safe: Goroutine-safe implementation
//   - Simple API: Intuitive and easy-to-use interface
//
// # Overview
//
// The logspan library consists of three main packages:
//
//   - github.com/zentooo/logspan/logger: Core logging functionality
//   - github.com/zentooo/logspan/formatter: Output formatting capabilities
//   - github.com/zentooo/logspan/http_middleware: HTTP server integration
//
// # Quick Start
//
// Basic setup and usage:
//
//	import (
//	    "context"
//	    "os"
//	    "github.com/zentooo/logspan/logger"
//	)
//
//	func main() {
//	    // Initialize logger
//	    logger.Init(logger.Config{
//	        MinLevel:     logger.InfoLevel,
//	        Output:       os.Stdout,
//	        PrettifyJSON: true,
//	    })
//
//	    // Direct logging
//	    logger.D.Infof("Application started")
//
//	    // Context-based logging
//	    ctx := context.Background()
//	    contextLogger := logger.NewContextLogger()
//	    ctx = logger.WithLogger(ctx, contextLogger)
//
//	    // Add context information
//	    logger.AddContextValue(ctx, "request_id", "req-12345")
//	    logger.AddContextValue(ctx, "user_id", "user-67890")
//
//	    // Record logs
//	    logger.Infof(ctx, "Request processing started")
//	    logger.Infof(ctx, "Request processing completed")
//
//	    // Output aggregated logs
//	    logger.FlushContext(ctx)
//	}
//
// # Concept
//
// ## Traditional Logging Challenges
// Traditional logging libraries output multiple log entries individually for a single
// request or process, causing related logs to scatter throughout log files. This makes
// it difficult to trace related logs and leads to inefficient debugging and troubleshooting.
//
// ## LogSpan's Approach
// LogSpan solves this problem with **context-based log aggregation**:
//
// ### Unified JSON Structure
//	{
//	  "type": "request",
//	  "context": {
//	    "request_id": "req-12345",
//	    "user_id": "user-67890"
//	  },
//	  "runtime": {
//	    "severity": "INFO",
//	    "startTime": "2023-10-27T09:59:58.123456+09:00",
//	    "endTime": "2023-10-27T10:00:00.223456+09:00",
//	    "elapsed": 150,
//	    "lines": [
//	      {"timestamp": "...", "level": "INFO", "message": "Request processing started"},
//	      {"timestamp": "...", "level": "DEBUG", "message": "Validating parameters"},
//	      {"timestamp": "...", "level": "INFO", "message": "Processing completed"}
//	    ]
//	  }
//	}
//
// ### Zero-Dependency Lightweight Design
//   - No External Dependencies: Uses only Go standard library
//   - Lightweight: Minimal memory footprint with automatic memory pooling
//   - Memory Efficient: Object pooling for LogEntry and slice reuse to reduce GC pressure
//   - Auto-Flush: Configurable automatic flushing to control memory usage
//   - Fast: Efficient log processing with optimized memory management
//   - Secure: No vulnerability risks from external dependencies
//
// ### Benefits
//   - Efficient Log Analysis: All related logs consolidated into a single JSON
//   - Improved Troubleshooting: Context information and processing time visible at a glance
//   - Simplified Operations: No dependency management required
//   - Better Performance: Lightweight and fast processing
//
// # Architecture
//
// The library follows a modular architecture with improved separation of concerns:
//
//	┌─────────────────────────────────────────────────────────┐
//	│                    Application Layer                    │
//	├─────────────────────────────────────────────────────────┤
//	│              HTTP Middleware (Optional)                 │
//	├─────────────────────────────────────────────────────────┤
//	│                    Logger Package                       │
//	│  ┌─────────────────┐    ┌─────────────────────────────┐ │
//	│  │ Context Logger  │    │      Direct Logger         │ │
//	│  │                 │    │                             │ │
//	│  └─────────┬───────┘    └─────────────┬───────────────┘ │
//	│            │                          │                 │
//	│            └──────────┬───────────────┘                 │
//	│                       │                                 │
//	│              ┌────────▼────────┐                        │
//	│              │   Base Logger   │                        │
//	│              │ (Shared Logic)  │                        │
//	│              └─────────────────┘                        │
//	├─────────────────────────────────────────────────────────┤
//	│                 Middleware System                       │
//	│  ┌─────────────────────────────────────────────────────┐ │
//	│  │          Middleware Manager                         │ │
//	│  │     (Global Middleware Management)                  │ │
//	│  └─────────────────────────────────────────────────────┘ │
//	├─────────────────────────────────────────────────────────┤
//	│                 Formatter Package                       │
//	│  ┌─────────────────┐    ┌─────────────────────────────┐ │
//	│  │  JSON Formatter │    │ Context Flatten Formatter  │ │
//	│  └─────────────────┘    └─────────────────────────────┘ │
//	│  ┌─────────────────────────────────────────────────────┐ │
//	│  │            Formatter Utils                          │ │
//	│  │      (Shared Formatting Logic)                     │ │
//	│  └─────────────────────────────────────────────────────┘ │
//	└─────────────────────────────────────────────────────────┘
//
// ## Architecture Improvements
//
// The library has been refactored for better maintainability and performance:
//
// ### Code Deduplication
// - **BaseLogger**: Common functionality shared between DirectLogger and ContextLogger
// - **Unified Methods**: SetOutput, SetLevel, SetFormatter methods consolidated
// - **Consistent Naming**: Standardized mutex naming and thread-safety patterns
//
// ### Responsibility Separation
// - **middleware_manager.go**: Global middleware management isolated
// - **formatter_utils.go**: Formatting utilities extracted to dedicated file
// - **logger.go**: Focused on core interfaces and global instances
//
// ### Enhanced Testing
// - **Comprehensive Coverage**: New test files for all major components
// - **Concurrency Testing**: Robust goroutine safety verification
// - **Edge Case Handling**: Thorough testing of error conditions and edge cases
//
// # Configuration
//
// The library provides flexible configuration options:
//
//	config := logger.Config{
//	    MinLevel:         logger.InfoLevel,    // Filter log levels
//	    Output:           os.Stdout,           // Output destination
//	    EnableSourceInfo: true,                // Include source file info
//	    PrettifyJSON:     true,                // Pretty-print JSON
//	    MaxLogEntries:    1000,                // Auto-flush threshold
//	    LogType:          "request",           // Custom log type
//	    ErrorHandler:     logger.NewDefaultErrorHandler(), // Error handler
//	}
//	logger.Init(config)
//
// # Memory Optimization
//
// LogSpan provides comprehensive memory optimization features:
//
// ## Auto-Flush Feature
//
// Configure auto-flush to control memory usage:
//
//	// Configure auto-flush
//	logger.Init(logger.Config{
//	    MaxLogEntries: 100, // Auto-flush every 100 entries
//	})
//
//	ctx := context.Background()
//	contextLogger := logger.NewContextLogger()
//	ctx = logger.WithLogger(ctx, contextLogger)
//
//	logger.AddContextValue(ctx, "request_id", "req-123")
//
//	// Auto-flush occurs when 100 entries are reached
//	for i := 0; i < 250; i++ {
//	    logger.Infof(ctx, "Processing item %d", i)
//	}
//	// Result: 2 auto-flushes (at 100 and 200 entries)
//	// Remaining 50 entries need manual flush
//
//	logger.FlushContext(ctx) // Output remaining entries
//
// ## Memory Pool Management
//
// LogSpan automatically manages memory pools for optimal performance:
//
//	// Pool statistics (for monitoring)
//	stats := logger.GetPoolStats()
//	fmt.Printf("LogEntry Pool Size: %d\n", stats.LogEntryPoolSize)
//	fmt.Printf("Slice Pool Size: %d\n", stats.SlicePoolSize)
//
//	// Note: Pool management is automatic and internal
//	// LogEntry and []*LogEntry slices are automatically pooled for memory efficiency
//
// # Use Cases
//
// ## Web Applications
//
// Perfect for web applications that need request-scoped logging:
//
//	import "github.com/zentooo/logspan/http_middleware"
//
//	handler := http_middleware.LoggingMiddleware(yourHandler)
//	http.Handle("/", handler)
//
// ## Microservices
//
// Ideal for microservices with structured logging requirements:
//
//	logger.Init(logger.Config{
//	    MinLevel:     logger.InfoLevel,
//	    PrettifyJSON: false, // Compact for production
//	    MaxLogEntries: 1000, // Auto-flush for memory efficiency
//	})
//
// ## CLI Applications
//
// Simple direct logging for command-line tools:
//
//	logger.D.Infof("Processing file: %s", filename)
//	logger.D.Errorf("Failed to process: %v", err)
//
// ## Background Services
//
// Context-based logging for background processing:
//
//	ctx := context.Background()
//	contextLogger := logger.NewContextLogger()
//	ctx = logger.WithLogger(ctx, contextLogger)
//
//	logger.AddContextValue(ctx, "job_id", jobID)
//	logger.Infof(ctx, "Job started")
//	// ... processing ...
//	logger.FlushContext(ctx)
//
// # Security
//
// Built-in security features for sensitive data protection:
//
//	// Password masking middleware
//	passwordMasker := logger.NewPasswordMaskingMiddleware()
//	logger.AddMiddleware(passwordMasker.Middleware())
//
//	// Custom password masking configuration
//	passwordMasker := logger.NewPasswordMaskingMiddleware().
//	    WithMaskString("[REDACTED]").
//	    WithPasswordKeys([]string{"password", "secret", "token"}).
//	    AddPasswordKey("api_key")
//
//	logger.AddMiddleware(passwordMasker.Middleware())
//
// # Error Handling
//
// Comprehensive error handling capabilities:
//
//	// Default error handler (writes to stderr)
//	defaultHandler := logger.NewDefaultErrorHandler()
//	logger.SetGlobalErrorHandler(defaultHandler)
//
//	// Custom output error handler
//	fileHandler := logger.NewDefaultErrorHandlerWithOutput(errorLogFile)
//	logger.SetGlobalErrorHandler(fileHandler)
//
//	// Silent error handler (ignores all errors)
//	silentHandler := &logger.SilentErrorHandler{}
//	logger.SetGlobalErrorHandler(silentHandler)
//
//	// Function-based error handler
//	funcHandler := logger.ErrorHandlerFunc(func(operation string, err error) {
//	    fmt.Printf("Logger error in %s: %v\n", operation, err)
//	})
//	logger.SetGlobalErrorHandler(funcHandler)
//
// # Performance
//
// Optimized for performance with:
//   - Zero external dependencies
//   - Minimal allocation overhead with memory pooling
//   - Efficient JSON marshaling
//   - Configurable auto-flush for memory management
//   - Thread-safe concurrent operations
//   - Object pooling for LogEntry and slice reuse
//
// # Compatibility
//
// The library is compatible with:
//   - Go 1.22+
//   - Standard library HTTP servers
//   - Popular HTTP frameworks (Gin, Echo, Gorilla, etc.)
//   - Cloud-native environments
//   - Container deployments
//
// For detailed documentation of individual packages, see:
//   - logger: Core logging functionality
//   - formatter: Output formatting
//   - http_middleware: HTTP integration
//
// # HTTP Integration
//
// Seamless HTTP server integration with automatic context setup:
//
//	import "github.com/zentooo/logspan/http_middleware"
//
//	mux := http.NewServeMux()
//	handler := http_middleware.LoggingMiddleware(mux)
//
//	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
//	    ctx := r.Context()
//
//	    // Request information is automatically added
//	    logger.Infof(ctx, "Fetching user list")
//
//	    // Additional context information
//	    logger.AddContextValue(ctx, "query_params", r.URL.Query())
//
//	    logger.Infof(ctx, "User list fetch completed")
//	    // FlushContext is called automatically
//	})
//
//	http.ListenAndServe(":8080", handler)

package logspan
