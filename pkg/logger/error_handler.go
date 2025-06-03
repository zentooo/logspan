package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// ErrorHandler defines the interface for handling logger errors
type ErrorHandler interface {
	// HandleError is called when an error occurs in the logger
	HandleError(operation string, err error)
}

// ErrorHandlerFunc is a function type that implements ErrorHandler
type ErrorHandlerFunc func(operation string, err error)

// HandleError implements ErrorHandler interface
func (f ErrorHandlerFunc) HandleError(operation string, err error) {
	f(operation, err)
}

// DefaultErrorHandler provides default error handling behavior
type DefaultErrorHandler struct {
	output io.Writer
	mutex  sync.Mutex
}

// NewDefaultErrorHandler creates a new DefaultErrorHandler
func NewDefaultErrorHandler() *DefaultErrorHandler {
	return &DefaultErrorHandler{
		output: os.Stderr,
	}
}

// NewDefaultErrorHandlerWithOutput creates a new DefaultErrorHandler with custom output
func NewDefaultErrorHandlerWithOutput(output io.Writer) *DefaultErrorHandler {
	return &DefaultErrorHandler{
		output: output,
	}
}

// HandleError handles errors by writing to the configured output
func (h *DefaultErrorHandler) HandleError(operation string, err error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.output != nil {
		_, _ = fmt.Fprintf(h.output, "[LOGGER ERROR] %s: %v\n", operation, err)
	}
}

// SetOutput sets the output writer for error messages
func (h *DefaultErrorHandler) SetOutput(output io.Writer) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.output = output
}

// SilentErrorHandler ignores all errors
type SilentErrorHandler struct{}

// HandleError implements ErrorHandler interface but does nothing
func (h *SilentErrorHandler) HandleError(operation string, err error) {
	// Do nothing - silent error handling
}

// Global error handler
var (
	globalErrorHandler ErrorHandler = NewDefaultErrorHandler()
	errorHandlerMutex  sync.RWMutex
)

// SetGlobalErrorHandler sets the global error handler
func SetGlobalErrorHandler(handler ErrorHandler) {
	errorHandlerMutex.Lock()
	defer errorHandlerMutex.Unlock()
	globalErrorHandler = handler
}

// GetGlobalErrorHandler returns the current global error handler
func GetGlobalErrorHandler() ErrorHandler {
	errorHandlerMutex.RLock()
	defer errorHandlerMutex.RUnlock()
	return globalErrorHandler
}

// handleError handles an error using the global error handler
func handleError(operation string, err error) {
	if err == nil {
		return
	}

	handler := GetGlobalErrorHandler()
	if handler != nil {
		handler.HandleError(operation, err)
	}
}

// LoggerError represents an error that occurred in the logger
type LoggerError struct {
	Operation string
	Err       error
}

// Error implements the error interface
func (e *LoggerError) Error() string {
	return fmt.Sprintf("logger %s error: %v", e.Operation, e.Err)
}

// Unwrap returns the underlying error
func (e *LoggerError) Unwrap() error {
	return e.Err
}

// NewLoggerError creates a new LoggerError
func NewLoggerError(operation string, err error) *LoggerError {
	return &LoggerError{
		Operation: operation,
		Err:       err,
	}
}
