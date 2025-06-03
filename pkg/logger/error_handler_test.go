package logger

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestDefaultErrorHandler(t *testing.T) {
	t.Run("HandleError", func(t *testing.T) {
		var buf bytes.Buffer
		handler := NewDefaultErrorHandlerWithOutput(&buf)

		err := errors.New("test error")
		handler.HandleError("test_operation", err)

		output := buf.String()
		if !strings.Contains(output, "[LOGGER ERROR]") {
			t.Error("Expected output to contain '[LOGGER ERROR]'")
		}
		if !strings.Contains(output, "test_operation") {
			t.Error("Expected output to contain operation name")
		}
		if !strings.Contains(output, "test error") {
			t.Error("Expected output to contain error message")
		}
	})

	t.Run("SetOutput", func(t *testing.T) {
		handler := NewDefaultErrorHandler()
		var buf bytes.Buffer
		handler.SetOutput(&buf)

		err := errors.New("test error")
		handler.HandleError("test", err)

		if buf.String() == "" {
			t.Error("Expected output after setting custom output")
		}
	})

	t.Run("NilOutput", func(t *testing.T) {
		handler := NewDefaultErrorHandlerWithOutput(nil)

		// Should not panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("HandleError with nil output should not panic: %v", r)
			}
		}()

		err := errors.New("test error")
		handler.HandleError("test", err)
	})
}

func TestSilentErrorHandler(t *testing.T) {
	handler := &SilentErrorHandler{}

	// Should not panic or produce any output
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("SilentErrorHandler should not panic: %v", r)
		}
	}()

	err := errors.New("test error")
	handler.HandleError("test", err)
}

func TestErrorHandlerFunc(t *testing.T) {
	var capturedOperation string
	var capturedError error

	handlerFunc := ErrorHandlerFunc(func(operation string, err error) {
		capturedOperation = operation
		capturedError = err
	})

	testErr := errors.New("test error")
	handlerFunc.HandleError("test_op", testErr)

	if capturedOperation != "test_op" {
		t.Errorf("Expected operation 'test_op', got '%s'", capturedOperation)
	}
	if capturedError != testErr {
		t.Errorf("Expected error to be captured correctly")
	}
}

func TestGlobalErrorHandler(t *testing.T) {
	// Save original handler
	originalHandler := GetGlobalErrorHandler()
	defer SetGlobalErrorHandler(originalHandler)

	t.Run("SetAndGetGlobalErrorHandler", func(t *testing.T) {
		var buf bytes.Buffer
		newHandler := NewDefaultErrorHandlerWithOutput(&buf)

		SetGlobalErrorHandler(newHandler)
		retrieved := GetGlobalErrorHandler()

		if retrieved != newHandler {
			t.Error("Expected to retrieve the same handler that was set")
		}
	})

	t.Run("HandleErrorFunction", func(t *testing.T) {
		var buf bytes.Buffer
		handler := NewDefaultErrorHandlerWithOutput(&buf)
		SetGlobalErrorHandler(handler)

		err := errors.New("test error")
		handleError("test_operation", err)

		output := buf.String()
		if !strings.Contains(output, "test_operation") {
			t.Error("Expected handleError to use global handler")
		}
	})

	t.Run("HandleErrorWithNilError", func(t *testing.T) {
		var buf bytes.Buffer
		handler := NewDefaultErrorHandlerWithOutput(&buf)
		SetGlobalErrorHandler(handler)

		// Should not call handler for nil error
		handleError("test", nil)

		if buf.String() != "" {
			t.Error("Expected no output for nil error")
		}
	})

	t.Run("HandleErrorWithNilHandler", func(t *testing.T) {
		SetGlobalErrorHandler(nil)

		// Should not panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("handleError with nil handler should not panic: %v", r)
			}
		}()

		err := errors.New("test error")
		handleError("test", err)
	})
}

func TestLoggerError(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		originalErr := errors.New("original error")
		loggerErr := NewLoggerError("test_operation", originalErr)

		errorMsg := loggerErr.Error()
		if !strings.Contains(errorMsg, "test_operation") {
			t.Error("Expected error message to contain operation")
		}
		if !strings.Contains(errorMsg, "original error") {
			t.Error("Expected error message to contain original error")
		}
	})

	t.Run("Unwrap", func(t *testing.T) {
		originalErr := errors.New("original error")
		loggerErr := NewLoggerError("test_operation", originalErr)

		unwrapped := loggerErr.Unwrap()
		if unwrapped != originalErr {
			t.Error("Expected Unwrap to return original error")
		}
	})

	t.Run("ErrorsIs", func(t *testing.T) {
		originalErr := errors.New("original error")
		loggerErr := NewLoggerError("test_operation", originalErr)

		if !errors.Is(loggerErr, originalErr) {
			t.Error("Expected errors.Is to work with LoggerError")
		}
	})
}
