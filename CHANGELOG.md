# Changelog
## v0.0.1 (2025-06-03)

### Initial Release
- Initial release of LogSpan logging library


All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release preparation
- Comprehensive CI/CD pipeline with GitHub Actions
- English documentation (README_EN.md)
- Release automation workflow
- BaseLogger for shared functionality between DirectLogger and ContextLogger
- New test files for improved coverage: base_logger_test.go, config_test.go, context_test.go
- Middleware manager for centralized global middleware management
- Formatter utilities for shared formatting logic

### Changed
- **BREAKING**: Refactored logger package structure for better maintainability
- Separated middleware management into dedicated middleware_manager.go file
- Extracted formatting utilities into formatter_utils.go file
- Improved code organization with BaseLogger pattern to eliminate duplication
- Enhanced test coverage from 92.4% to near 100% for critical functions
- Standardized mutex naming conventions across all logger implementations
- Improved documentation with updated architecture diagrams and examples

### Fixed
- Resolved concurrent safety issues in DirectLogger tests
- Fixed inconsistent mutex naming between logger implementations
- Improved thread-safety patterns across all components

### Removed
- Duplicate code between DirectLogger and ContextLogger implementations
- Redundant SetOutput, SetLevel, and SetFormatter method implementations
- Japanese comments replaced with English for international compatibility

### Features
- Context-based logger for request-scoped log aggregation
- Direct logger for immediate log output
- Middleware system for extensible log processing
- Password masking middleware for sensitive data protection
- Multiple formatters (JSON, Context Flatten)
- HTTP middleware for automatic request logging
- Auto-flush functionality for memory optimization
- Configurable log levels and output destinations
- Thread-safe implementation

### Documentation
- Comprehensive API documentation
- Usage examples for all major features
- Architecture and design principles documentation
- Japanese and English README files
- Updated package structure documentation
- Enhanced test coverage documentation

### Testing
- Extensive test coverage (improved from 92.4% to near 100% for core functions)
- New comprehensive test files for all major components
- Enhanced concurrent safety testing
- Benchmark tests for performance validation
- Security scanning with gosec
- Go 1.24 compatibility testing

### Infrastructure
- GitHub Actions CI/CD pipeline
- Automated testing and linting
- Security vulnerability scanning
- Code coverage reporting with Codecov
- Release automation with semantic versioning