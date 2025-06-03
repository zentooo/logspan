# Changelog
## v0.0.4 (2025-06-03)

### Changes
- Merge pull request #5 from zentooo/renovate/pin-dependencies (3963699)
- chore(deps): pin dependencies (d55af2d)
## v0.0.3 (2025-06-03)

### Changes
- Merge branch 'main' of github.com:zentooo/logspan (5fee439)
- chore: update renovate.json to include package rules for action dependencies with pinned digests (5a4cfb7)
## v0.0.2 (2025-06-03)

### Changes
- Update CI workflows to use gotestsum for running tests and add mutex locking in base_logger for thread safety (4d03aa1)
- Merge branch 'main' of github.com:zentooo/logspan (61ae366)
- Update go.mod to include additional indirect dependencies and modify CI workflow to use gotestsum for running tests (6525726)
- Merge pull request #2 from zentooo/renovate/codecov-codecov-action-5.x (46ac80b)
- Merge pull request #3 from zentooo/renovate/golangci-golangci-lint-action-8.x (ea5abbc)
- fix lint (1b47ea2)
- Merge branch 'main' into renovate/codecov-codecov-action-5.x (435d71c)
- Merge branch 'main' into renovate/golangci-golangci-lint-action-8.x (79c6f12)
- Refactor log output tests to use table-driven approach for log type verification (01bf2dd)
- Enhance README files to emphasize context-based log aggregation and zero-dependency features (d06809b)
- Implement object pooling for LogEntry in logger package to optimize memory allocation (d4b9453)
- Update progress and documentation for logger package enhancements (7bcb9ec)
- Add configurable log type feature to logger package (008095f)
- Fix formatting in error handling example output and update error logging to ignore write errors. (ed7af57)
- Merge branch 'main' of github.com:zentooo/logspan (a3ebd2e)
- Update .gitignore to include test coverage files and binaries, and remove coverage.html and coverage.out files to clean up the repository. (4ea3a56)
- chore(deps): update golangci/golangci-lint-action action to v8 (b2d7d66)
- chore(deps): update codecov/codecov-action action to v5 (4345063)
- Merge pull request #1 from zentooo/renovate/configure (a9e13a8)
- Merge branch 'main' of github.com:zentooo/logspan (bb39d08)
- Merge branch 'package_refactor' (5dd41c5)
- Enhance logger package by implementing comprehensive error handling, updating configuration defaults, and improving documentation. Complete tasks for package quality improvements, including adding benchmark and performance tests. Document progress and outline next steps for security measures. (a9809f6)
- Complete documentation updates by creating CONTRIBUTING.md, LICENSE.md, and CODE_OF_CONDUCT.md. Enhance examples with new middleware and advanced configuration examples, updating the examples README to reflect these additions. Document progress and next steps for package quality improvements. (15485dc)
- Refactor logger package to use pointer receivers for BaseLogger in DirectLogger and ContextLogger, enhancing memory efficiency and performance. Update related tests to ensure functionality remains intact. Document changes and next steps for ongoing improvements. (7f80929)
- Refactor logger package to enhance maintainability and performance by introducing BaseLogger for shared functionality, separating middleware management and formatting utilities into dedicated files, and improving test coverage to nearly 100%. Update documentation to reflect architectural changes and improvements in code organization. (2c0045c)
- Refactor logger package to implement BaseLogger structure, reducing code duplication in DirectLogger and ContextLogger. Update mutex handling for thread safety and enhance overall code maintainability. Document progress and next steps for ongoing improvements. (c2bc6d1)
- Enhance logger package documentation to outline refactoring tasks, including analysis of code structure, identification of code duplication, and planning for improvements. Document current progress and next steps for ongoing refactoring efforts. (8385c39)
- Refactor logger package to unify log level comparison logic, reduce code duplication, and improve code readability. Update tests to align with new API and enhance performance by switching from string to int comparisons for log levels. Document changes and next steps for ongoing improvements. (c30fa4a)
- Add renovate.json (cc657ba)
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