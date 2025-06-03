# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release preparation
- Comprehensive CI/CD pipeline with GitHub Actions
- English documentation (README_EN.md)
- Release automation workflow

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

### Testing
- Extensive test coverage (80.4% for core logger, 100% for formatters and HTTP middleware)
- Benchmark tests for performance validation
- Security scanning with gosec
- Go 1.24 compatibility testing

### Infrastructure
- GitHub Actions CI/CD pipeline
- Automated testing and linting
- Security vulnerability scanning
- Code coverage reporting with Codecov
- Release automation with semantic versioning