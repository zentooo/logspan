# LogSpan

LogSpan is a structured logging library for Go that provides two logging modes: **Context Logger** for aggregating logs per HTTP request unit, and **Direct Logger** for immediate log output, delivering flexible and user-friendly logging functionality.

## 🎯 Key Features

- **Dual-mode logging**: Context-based and direct logging modes
- **Structured log output**: Consistent JSON-formatted log output
- **Middleware mechanism**: Customizable log processing pipeline
- **Context flattening**: Formatter that expands context fields to top level
- **HTTP middleware**: Automatic log setup for web applications
- **Concurrency-safe**: Goroutine-safe implementation
- **Simple API**: Intuitive and easy-to-use interface

## 📦 Installation

```bash
go get github.com/zentooo/logspan
```

## 🚀 Quick Start

### Direct Logger (Immediate Output)

```go
package main

import "github.com/zentooo/logspan/pkg/logger"

func main() {
    // Use global direct logger
    logger.D.Infof("Application started")
    logger.D.Warnf("Warning: %s", "Configuration file not found")
    logger.D.Errorf("Error: %v", err)
}
```

### Context Logger (Log Aggregation)

```go
package main

import (
    "context"
    "github.com/zentooo/logspan/pkg/logger"
)

func main() {
    // Create context logger
    ctx := context.Background()
    contextLogger := logger.NewContextLogger()
    ctx = logger.WithLogger(ctx, contextLogger)

    // Add context information
    logger.AddContextValue(ctx, "request_id", "req-12345")
    logger.AddContextValue(ctx, "user_id", "user-67890")

    // Record logs
    logger.Infof(ctx, "Request processing started")
    processRequest(ctx)
    logger.Infof(ctx, "Request processing completed")

    // Output aggregated logs
    logger.FlushContext(ctx)
}

func processRequest(ctx context.Context) {
    logger.AddContextValue(ctx, "step", "validation")
    logger.Debugf(ctx, "Validating input parameters")
    logger.Infof(ctx, "Input validation completed")
}
```

## 📖 Detailed Usage

### 1. Initialization and Configuration

#### Global Configuration

```go
import "github.com/zentooo/logspan/pkg/logger"

func init() {
    config := logger.Config{
        MinLevel:         logger.DebugLevel,
        Output:           os.Stdout,
        EnableSourceInfo: true,
    }
    logger.Init(config)
}
```

#### Individual Logger Configuration

```go
// Direct logger configuration
directLogger := logger.NewDirectLogger()
directLogger.SetLevelFromString("WARN")
directLogger.SetOutput(logFile)

// Context logger configuration
contextLogger := logger.NewContextLogger()
contextLogger.SetLevel(logger.InfoLevel)
contextLogger.SetOutput(logFile)
```

### 2. Log Levels

LogSpan supports five log levels:

- `DEBUG`: Detailed debugging information
- `INFO`: General informational messages
- `WARN`: Warning messages
- `ERROR`: Error messages
- `CRITICAL`: Critical error messages

```go
logger.D.Debugf("Debug info: %s", debugInfo)
logger.D.Infof("Info: %s", info)
logger.D.Warnf("Warning: %s", warning)
logger.D.Errorf("Error: %v", err)
logger.D.Criticalf("Critical error: %v", criticalErr)
```

### 3. Context Operations

#### Context Logger Setup

```go
// Create and configure context logger
ctx := context.Background()
contextLogger := logger.NewContextLogger()
ctx = logger.WithLogger(ctx, contextLogger)

// Or automatically retrieve from context (creates new if not exists)
contextLogger := logger.FromContext(ctx)
```

#### Adding Context Fields

```go
// Add single field
logger.AddContextValue(ctx, "user_id", "12345")
logger.AddContextValue(ctx, "session_id", "session-abc")

// Add multiple fields
logger.AddContextValues(ctx, map[string]interface{}{
    "request_id": "req-67890",
    "ip_address": "192.168.1.1",
    "user_agent": "Mozilla/5.0...",
})

// Use logger instance directly
contextLogger := logger.FromContext(ctx)
contextLogger.AddContextValue("operation", "user_login")
contextLogger.AddContextValues(map[string]interface{}{
    "step": "validation",
    "attempt": 1,
})
```

#### Context Logger API

```go
// Log with context
logger.Infof(ctx, "User %s logged in", userID)
logger.Debugf(ctx, "Processing step: %s", step)
logger.Errorf(ctx, "Error occurred during processing: %v", err)

// Output logs (flush aggregated logs at once)
logger.FlushContext(ctx)
```

### 4. HTTP Middleware

Automatic log setup for web applications:

```go
package main

import (
    "net/http"
    "github.com/zentooo/logspan/pkg/http_middleware"
    "github.com/zentooo/logspan/pkg/logger"
)

func main() {
    mux := http.NewServeMux()

    // Apply logging middleware
    handler := http_middleware.LoggingMiddleware(mux)

    mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        // Request information is automatically added
        logger.Infof(ctx, "Fetching user list")

        // Additional context information
        logger.AddContextValue(ctx, "query_params", r.URL.Query())

        // Processing...

        logger.Infof(ctx, "User list fetch completed")
        // FlushContext is called automatically
    })

    http.ListenAndServe(":8080", handler)
}
```

### 5. Middleware Mechanism

Customize the log processing pipeline:

#### Basic Middleware

```go
// Create custom middleware
func customMiddleware(entry *logger.LogEntry, next func(*logger.LogEntry)) {
    // Pre-process log entry
    entry.Message = "[CUSTOM] " + entry.Message

    // Call next middleware or final processing
    next(entry)
}

// Register middleware
logger.AddMiddleware(customMiddleware)

// Middleware management
logger.ClearMiddleware()                    // Clear all middleware
count := logger.GetMiddlewareCount()        // Get middleware count
```

#### Password Masking Middleware

LogSpan includes built-in middleware that automatically masks sensitive information:

```go
// Enable password masking with default settings
passwordMasker := logger.NewPasswordMaskingMiddleware()
logger.AddMiddleware(passwordMasker.Middleware())

// Custom password masking configuration
passwordMasker := logger.NewPasswordMaskingMiddleware().
    WithMaskString("[REDACTED]").                           // Customize mask string
    WithPasswordKeys([]string{"password", "secret"}).       // Set target keys
    AddPasswordKey("api_key").                              // Add additional key
    AddPasswordPattern(regexp.MustCompile(`token=\w+`))     // Custom regex pattern

logger.AddMiddleware(passwordMasker.Middleware())

// Usage example
logger.D.Infof("User login: username=john password=secret123 token=abc123")
// Output: "User login: username=john password=*** token=***"
```

##### Default Masked Keywords
- `password`, `passwd`, `pwd`, `pass`
- `secret`, `token`, `key`, `auth`
- `credential`, `credentials`, `api_key`
- `access_token`, `refresh_token`

##### Supported Patterns
- `key=value` format: `password=secret` → `password=***`
- JSON format: `"password":"secret"` → `"password":"***"`
- Custom regex patterns

### 6. Formatters

#### JSON Formatter (Default)

```go
contextLogger := logger.NewContextLogger()
contextLogger.SetFormatter(formatter.NewJSONFormatter())
```

#### Context Flatten Formatter

```go
import "github.com/zentooo/logspan/pkg/formatter"

contextLogger := logger.NewContextLogger()
contextLogger.SetFormatter(formatter.NewContextFlattenFormatter())
```

## 📋 Log Output Formats

### Default JSON Format

```json
{
  "type": "request",
  "context": {
    "request_id": "req-12345",
    "user_id": "user-67890"
  },
  "runtime": {
    "severity": "INFO",
    "startTime": "2023-10-27T09:59:58.123456+09:00",
    "endTime": "2023-10-27T10:00:00.223456+09:00",
    "elapsed": 150,
    "lines": [
      {
        "timestamp": "2023-10-27T09:59:59.123456+09:00",
        "level": "INFO",
        "message": "Request processing started"
      }
    ]
  },
  "config": {
    "elapsedUnit": "ms"
  }
}
```

### Context Flatten Format

```json
{
  "request_id": "req-12345",
  "user_id": "user-67890",
  "type": "request",
  "runtime": {
    "severity": "INFO",
    "startTime": "2023-10-27T09:59:58.123456+09:00",
    "endTime": "2023-10-27T10:00:00.223456+09:00",
    "elapsed": 150,
    "lines": [
      {
        "timestamp": "2023-10-27T09:59:59.123456+09:00",
        "level": "INFO",
        "message": "Request processing started"
      }
    ]
  },
  "config": {
    "elapsedUnit": "ms"
  }
}
```

## 🔧 Configuration Options

### Config Structure

```go
type Config struct {
    // Minimum log level
    MinLevel LogLevel

    // Output destination
    Output io.Writer

    // Enable source file information
    EnableSourceInfo bool

    // Enable JSON output formatting (indentation)
    PrettifyJSON bool

    // Maximum entries for context logger (0 = no limit)
    MaxLogEntries int
}
```

### Default Configuration

```go
config := logger.DefaultConfig()
// MinLevel: InfoLevel
// Output: os.Stdout
// EnableSourceInfo: false
// PrettifyJSON: false
// MaxLogEntries: 1000
```

### Custom Configuration Examples

```go
// Development environment configuration (formatted JSON output)
logger.Init(logger.Config{
    MinLevel:         logger.DebugLevel,
    Output:           os.Stdout,
    EnableSourceInfo: true,
    PrettifyJSON:     true,  // Pretty-formatted JSON for readability
    MaxLogEntries:    500,   // Auto-flush every 500 entries
})

// Production environment configuration (compact JSON output)
logger.Init(logger.Config{
    MinLevel:         logger.InfoLevel,
    Output:           logFile,
    EnableSourceInfo: false,
    PrettifyJSON:     false,  // Compact JSON
    MaxLogEntries:    1000,   // Auto-flush every 1000 entries
})

// Memory-efficient configuration
logger.Init(logger.Config{
    MinLevel:      logger.InfoLevel,
    Output:        logFile,
    PrettifyJSON:  false,
    MaxLogEntries: 100,  // Frequent auto-flush to reduce memory usage
})

// No limit configuration (manual flush only)
logger.Init(logger.Config{
    MinLevel:      logger.InfoLevel,
    Output:        logFile,
    MaxLogEntries: 0,  // Disable auto-flush
})
```

### Configuration Verification

```go
// Check if logger is initialized
if logger.IsInitialized() {
    config := logger.GetConfig()
    fmt.Printf("Current log level: %s\n", config.MinLevel.String())
    fmt.Printf("Pretty JSON enabled: %t\n", config.PrettifyJSON)
    fmt.Printf("Max log entries: %d\n", config.MaxLogEntries)
}
```

## 🚀 Memory Optimization

### Auto-Flush Feature

LogSpan provides an auto-flush feature to control memory usage:

#### Basic Operation

```go
// Configure auto-flush
logger.Init(logger.Config{
    MaxLogEntries: 100, // Auto-flush every 100 entries
})

ctx := context.Background()
contextLogger := logger.NewContextLogger()
ctx = logger.WithLogger(ctx, contextLogger)

logger.AddContextValue(ctx, "request_id", "req-123")

// Auto-flush occurs when 100 entries are reached
for i := 0; i < 250; i++ {
    logger.Infof(ctx, "Processing item %d", i)
}
// Result: 2 auto-flushes (at 100 and 200 entries)
// Remaining 50 entries need manual flush

logger.FlushContext(ctx) // Output remaining entries
```

#### Auto-Flush Features

- **Entry Counting**: Only entries that pass the log level filter are counted
- **Batch Processing**: Each auto-flush outputs as an independent log batch
- **Context Preservation**: Context fields are preserved after auto-flush
- **Memory Release**: Entries are automatically cleared after flush to free memory

#### Memory-Efficient Usage Example

```go
// Configuration for large-scale log processing
logger.Init(logger.Config{
    MinLevel:      logger.InfoLevel,
    MaxLogEntries: 50,    // Small batch size
    PrettifyJSON:  false, // Compact output
})

ctx := context.Background()
contextLogger := logger.NewContextLogger()
ctx = logger.WithLogger(ctx, contextLogger)

logger.AddContextValue(ctx, "batch_id", "batch-001")

// Process large amount of data
for i := 0; i < 10000; i++ {
    logger.Infof(ctx, "Processing record %d", i)

    if i%1000 == 0 {
        // Add progress to context
        logger.AddContextValue(ctx, "progress", fmt.Sprintf("%d/10000", i))
    }
}
// Memory usage remains constant due to auto-flush

logger.FlushContext(ctx) // Output final remaining entries
```

#### Disable Option

```go
// Disable auto-flush (traditional behavior)
logger.Init(logger.Config{
    MaxLogEntries: 0, // 0 = no limit
})

// In this case, entries accumulate until manual FlushContext() call
```

## 📚 Sample Code

Detailed sample code is available in the `examples/` directory:

```bash
# Direct logger sample
go run examples/direct_logger/main.go

# Context logger sample
go run examples/context_logger/main.go

# Auto-flush feature sample
go run examples/auto_flush/main.go

# HTTP middleware sample
go run examples/http_middleware_example.go
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Test with coverage
go test -cover ./...

# Verbose test output
go test -v ./...
```

## 🏗️ Architecture

### Package Structure

```
pkg/
├── logger/                          # Main logger package
│   ├── logger.go                   # Core interface and API
│   ├── context_logger.go           # Context logger implementation
│   ├── direct_logger.go            # Direct logger implementation
│   ├── config.go                   # Configuration management
│   ├── entry.go                    # Log entry structure
│   ├── middleware.go               # Middleware mechanism
│   ├── context.go                  # Context helpers
│   ├── level.go                    # Log level definitions
│   └── password_masking_middleware.go # Password masking
├── formatter/                       # Formatters
│   ├── interface.go                # Formatter interface
│   ├── json_formatter.go           # JSON formatter
│   └── context_flatten_formatter.go # Context flatten formatter
├── http_middleware/                 # HTTP middleware
│   └── middleware.go               # HTTP request logging
└── examples/                        # Usage examples
    ├── context_logger/             # Context logger examples
    ├── direct_logger/              # Direct logger examples
    ├── context_flatten_formatter/  # Context flatten formatter examples
    └── http_middleware_example.go  # HTTP middleware examples
```

### Design Principles

1. **Simple API**: Intuitive and easy-to-use interface
2. **Flexibility**: Design that accommodates various use cases
3. **Extensibility**: Customization through middleware
4. **Performance**: Efficient log processing
5. **Concurrency Safety**: Goroutine-safe implementation

## 🤝 Contributing

1. Fork this repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Create a Pull Request

## 📄 License

This project is released under the MIT License. See the [LICENSE](LICENSE) file for details.

## 🔗 Related Links

- [Go Documentation](https://pkg.go.dev/github.com/zentooo/logspan)
- [Examples](./examples/)
- [Design Document](./design.md)

## 📞 Support

If you have questions or issues, please create an [Issue](https://github.com/zentooo/logspan/issues).