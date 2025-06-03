// Package logspan provides a flexible and extensible logging library for Go applications.
//
// Logspan offers a comprehensive logging solution with dual-mode logging capabilities,
// middleware support, and seamless HTTP integration. It's designed to provide both
// immediate logging for simple use cases and context-based logging for complex
// request processing scenarios.
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
//	    "github.com/zentooo/logspan/logger"
//	)
//
//	func main() {
//	    // Initialize logger
//	    logger.Init(logger.DefaultConfig())
//
//	    // Direct logging
//	    logger.D.Infof("Application started")
//
//	    // Context-based logging
//	    ctx := context.Background()
//	    contextLogger := logger.NewContextLogger()
//	    ctx = logger.WithLogger(ctx, contextLogger)
//
//	    logger.AddContextValue(ctx, "request_id", "req-123")
//	    logger.Infof(ctx, "Processing request")
//	    logger.FlushContext(ctx)
//	}
//
// # Key Features
//
//   - Dual-mode logging: Context-based aggregation and direct immediate output
//   - Transparent context management through HTTP middleware
//   - Extensible middleware system for log processing
//   - Built-in password masking for sensitive information
//   - Multiple output formatters (JSON, Context Flatten)
//   - Configurable log levels and output destinations
//   - Thread-safe operations
//   - Memory-efficient with auto-flush capabilities
//   - Comprehensive HTTP integration
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
//	}
//	logger.Init(config)
//
// # Security
//
// Built-in security features for sensitive data protection:
//
//	// Password masking middleware
//	passwordMasker := logger.NewPasswordMaskingMiddleware()
//	logger.AddMiddleware(passwordMasker.Middleware())
//
// # Performance
//
// Optimized for performance with:
//   - Minimal allocation overhead
//   - Efficient JSON marshaling
//   - Configurable auto-flush for memory management
//   - Thread-safe concurrent operations
//
// # Compatibility
//
// The library is compatible with:
//   - Go 1.21+
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
//	handler := http_middleware.LoggingMiddleware(yourHandler)
//	http.Handle("/", handler)

package logspan
