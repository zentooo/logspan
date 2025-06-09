# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

LogSpan is a zero-dependency structured logging library for Go that provides context-based log aggregation. It consolidates multiple log entries into a single JSON structure, improving log analysis and troubleshooting efficiency. The library supports both context-based (aggregated) and direct (immediate) logging modes with automatic memory management and middleware extensibility.

## Common Development Commands

### Testing
- **Run all tests**: `go test ./...`
- **Test with coverage**: `go test -cover ./...`
- **Verbose testing**: `go test -v ./...`
- **Single package test**: `go test ./logger/`
- **Benchmark tests**: `go test -bench=. ./...`

### Building and Linting
- **Build all packages**: `go build ./...`
- **Download dependencies**: `go mod download`
- **Run linter**: `golangci-lint run`

### Version Management
- **Show current version**: `./scripts/version.sh current`
- **List all versions**: `./scripts/version.sh list`
- **Create new version**: `./scripts/version.sh create patch|minor|major`

### Running Examples
- **Direct logging**: `go run examples/direct_logger/main.go`
- **Context logging**: `go run examples/context_logger/main.go`
- **Auto flush**: `go run examples/auto_flush/main.go`
- **Middleware**: `go run examples/middleware/main.go`
- **HTTP middleware**: `go run examples/http_middleware_example.go`

## Architecture & Code Organization

### Core Package Structure
- **logger/**: Main logging functionality with BaseLogger, DirectLogger, and ContextLogger
- **formatter/**: Output formatters (JSON, Context Flatten) implementing pluggable interface
- **http_middleware/**: HTTP server integration for automatic request logging

### Key Design Patterns
- **Functional Options Pattern**: Configuration via `logger.WithXXX()` options
- **Interface-based Design**: Pluggable formatters, error handlers, and middleware
- **Memory Pool Management**: Automatic LogEntry and slice pooling for performance
- **Middleware Pipeline**: Extensible log processing chain with built-in password masking

### Logger Types
- **DirectLogger**: Immediate log output for simple use cases
- **ContextLogger**: Aggregates logs within context scope, flushes on completion
- **BaseLogger**: Shared functionality between logger types

### Configuration Management
The library uses functional options for initialization:
```go
logger.Init(
    logger.WithMinLevel(logger.InfoLevel),
    logger.WithOutput(os.Stdout),
    logger.WithMaxLogEntries(1000),
    logger.WithFlushEmpty(true),
)
```

### Testing Strategy
- High test coverage with comprehensive test files for all components
- Table-driven tests and benchmark testing
- Concurrency safety validation
- Test files follow `*_test.go` naming pattern

### Dependencies
- **Go Version**: 1.24+ (specified in go.mod)
- **Zero External Dependencies**: Uses only Go standard library
- **Development Tools**: golangci-lint for linting, version.sh for release management