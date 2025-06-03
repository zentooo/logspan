package http_middleware

import (
	"net/http"
	"time"

	"github.com/zentooo/logspan/logger"
)

// LoggingMiddleware creates an HTTP middleware that automatically sets up logging context
// for each request and collects basic HTTP request information
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a new context logger for this request
		contextLogger := logger.NewContextLogger()

		// Add basic HTTP request information to the context
		contextLogger.AddContextValues(map[string]interface{}{
			"method":      r.Method,
			"url":         r.URL.String(),
			"path":        r.URL.Path,
			"query":       r.URL.RawQuery,
			"user_agent":  r.UserAgent(),
			"remote_addr": r.RemoteAddr,
			"host":        r.Host,
		})

		// Create a response writer wrapper to capture response information
		wrappedWriter := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default to 200
		}

		// Add the logger to the request context
		ctx := logger.WithLogger(r.Context(), contextLogger)
		r = r.WithContext(ctx)

		// Log the start of the request
		logger.Infof(ctx, "Request started")

		// Record start time for duration calculation
		startTime := time.Now()

		// Call the next handler
		next.ServeHTTP(wrappedWriter, r)

		// Calculate request duration
		duration := time.Since(startTime)

		// Add response information to the context
		contextLogger.AddContextValues(map[string]interface{}{
			"status_code": wrappedWriter.statusCode,
			"duration_ms": duration.Milliseconds(),
		})

		// Log the completion of the request
		logger.Infof(ctx, "Request completed")

		// Flush the accumulated logs
		logger.FlushContext(ctx)
	})
}

// responseWriter wraps http.ResponseWriter to capture response information
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
