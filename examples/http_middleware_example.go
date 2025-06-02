package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zentooo/logspan/pkg/http_middleware"
	"github.com/zentooo/logspan/pkg/logger"
)

func main() {
	// Create a simple HTTP handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The logger is automatically available in the request context
		// thanks to the LoggingMiddleware

		// Log some information during request processing
		logger.Infof(r.Context(), "Processing request for path: %s", r.URL.Path)

		// Add custom fields to the request context
		logger.AddField(r.Context(), "user_id", "12345")
		logger.AddField(r.Context(), "operation", "get_user_profile")

		// Simulate some processing time
		time.Sleep(100 * time.Millisecond)

		// Log before sending response
		logger.Infof(r.Context(), "Sending response")

		// Send response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Hello, World!", "user_id": "12345"}`))
	})

	// Example 1: Basic logging middleware
	fmt.Println("Example 1: Basic HTTP middleware")
	basicMiddleware := http_middleware.LoggingMiddleware(handler)

	// Example 2: Another handler with different processing
	fmt.Println("Example 2: Different handler with logging")
	anotherHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Infof(r.Context(), "Processing different endpoint")

		// Add different fields
		logger.AddField(r.Context(), "endpoint", "api_status")
		logger.AddField(r.Context(), "version", "v1.0")

		// Simulate different processing
		time.Sleep(50 * time.Millisecond)

		logger.Infof(r.Context(), "Status check completed")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "version": "v1.0"}`))
	})
	anotherMiddleware := http_middleware.LoggingMiddleware(anotherHandler)

	// Example 3: Error handling with logging
	fmt.Println("Example 3: Error handling with logging")
	errorHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Infof(r.Context(), "Processing error endpoint")

		// Add error context
		logger.AddField(r.Context(), "endpoint", "error_test")
		logger.AddField(r.Context(), "error_type", "simulated")

		// Log an error
		logger.Errorf(r.Context(), "Simulated error occurred")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal server error", "code": 500}`))
	})
	errorMiddleware := http_middleware.LoggingMiddleware(errorHandler)

	// Create HTTP server with different middleware configurations
	mux := http.NewServeMux()

	// Route with basic middleware
	mux.Handle("/basic", basicMiddleware)

	// Route with different handler
	mux.Handle("/status", anotherMiddleware)

	// Route with error handling
	mux.Handle("/error", errorMiddleware)

	// Health check endpoint without logging middleware
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	fmt.Println("Starting server on :8080")
	fmt.Println("Try these endpoints:")
	fmt.Println("  curl http://localhost:8080/basic")
	fmt.Println("  curl http://localhost:8080/status")
	fmt.Println("  curl http://localhost:8080/error")
	fmt.Println("  curl http://localhost:8080/health")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
