// Package http_middleware provides HTTP middleware for automatic request logging.
//
// The http_middleware package offers seamless integration between HTTP servers
// and the logger package, automatically setting up logging context for each
// HTTP request and collecting relevant request/response information.
//
// # Key Features
//
//   - Automatic logging context setup for each HTTP request
//   - Collection of HTTP request information (method, URL, headers, etc.)
//   - Response information capture (status code, duration)
//   - Transparent integration with existing HTTP handlers
//   - Context propagation throughout the request lifecycle
//
// # Basic Usage
//
// Wrap your HTTP handler with the logging middleware:
//
//	import (
//	    "net/http"
//	    "github.com/zentooo/logspan/http_middleware"
//	    "github.com/zentooo/logspan/logger"
//	)
//
//	func main() {
//	    // Initialize logger
//	    logger.Init(logger.DefaultConfig())
//
//	    // Your HTTP handler
//	    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	        // Use logger with automatic context
//	        logger.Infof(r.Context(), "Processing request")
//	        w.WriteHeader(http.StatusOK)
//	        w.Write([]byte("Hello, World!"))
//	    })
//
//	    // Wrap with logging middleware
//	    loggingHandler := http_middleware.LoggingMiddleware(handler)
//
//	    // Start server
//	    http.Handle("/", loggingHandler)
//	    http.ListenAndServe(":8080", nil)
//	}
//
// # Automatic Context Information
//
// The middleware automatically adds the following information to the logging context:
//
//   - method: HTTP method (GET, POST, etc.)
//   - url: Complete request URL
//   - path: URL path
//   - query: Query string parameters
//   - user_agent: User-Agent header
//   - remote_addr: Client's remote address
//   - host: Host header
//   - status_code: HTTP response status code (added after response)
//   - duration_ms: Request processing duration in milliseconds (added after response)
//
// # Request Lifecycle
//
// The middleware follows this lifecycle for each request:
//
//  1. Create a new context logger for the request
//  2. Add HTTP request information to the context
//  3. Log "Request started" message
//  4. Call the next handler in the chain
//  5. Capture response information (status code, duration)
//  6. Add response information to the context
//  7. Log "Request completed" message
//  8. Flush all accumulated logs for the request
//
// # Integration with Custom Handlers
//
// The middleware works seamlessly with any HTTP handler:
//
//	// With http.ServeMux
//	mux := http.NewServeMux()
//	mux.HandleFunc("/api/users", usersHandler)
//	mux.HandleFunc("/api/posts", postsHandler)
//	wrappedMux := http_middleware.LoggingMiddleware(mux)
//
//	// With third-party routers (example with gorilla/mux)
//	router := mux.NewRouter()
//	router.HandleFunc("/api/users", usersHandler)
//	wrappedRouter := http_middleware.LoggingMiddleware(router)
//
// # Custom Logging in Handlers
//
// Within your handlers, you can add custom logging using the request context:
//
//	func myHandler(w http.ResponseWriter, r *http.Request) {
//	    ctx := r.Context()
//
//	    // Add custom context information
//	    logger.AddContextValue(ctx, "user_id", "user-123")
//	    logger.AddContextValue(ctx, "operation", "create_user")
//
//	    // Log custom messages
//	    logger.Infof(ctx, "Starting user creation")
//	    logger.Debugf(ctx, "Validating user data: %+v", userData)
//
//	    // Your business logic here
//	    if err := createUser(userData); err != nil {
//	        logger.Errorf(ctx, "Failed to create user: %v", err)
//	        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
//	        return
//	    }
//
//	    logger.Infof(ctx, "User created successfully")
//	    w.WriteHeader(http.StatusCreated)
//	}
//
// # Output Example
//
// The middleware produces structured log output like:
//
//	{
//	  "type": "request",
//	  "context": {
//	    "method": "POST",
//	    "url": "http://localhost:8080/api/users?format=json",
//	    "path": "/api/users",
//	    "query": "format=json",
//	    "user_agent": "Mozilla/5.0...",
//	    "remote_addr": "127.0.0.1:54321",
//	    "host": "localhost:8080",
//	    "status_code": 201,
//	    "duration_ms": 45,
//	    "user_id": "user-123",
//	    "operation": "create_user"
//	  },
//	  "runtime": {
//	    "severity": "INFO",
//	    "startTime": "2023-10-27T09:59:58.123456+09:00",
//	    "endTime": "2023-10-27T10:00:00.168456+09:00",
//	    "elapsed": 45,
//	    "lines": [
//	      {
//	        "timestamp": "2023-10-27T09:59:58.123456+09:00",
//	        "level": "INFO",
//	        "message": "Request started"
//	      },
//	      {
//	        "timestamp": "2023-10-27T09:59:58.135456+09:00",
//	        "level": "INFO",
//	        "message": "Starting user creation"
//	      },
//	      {
//	        "timestamp": "2023-10-27T10:00:00.165456+09:00",
//	        "level": "INFO",
//	        "message": "User created successfully"
//	      },
//	      {
//	        "timestamp": "2023-10-27T10:00:00.168456+09:00",
//	        "level": "INFO",
//	        "message": "Request completed"
//	      }
//	    ]
//	  }
//	}
//
// # Thread Safety
//
// The middleware is thread-safe and can handle concurrent requests.
// Each request gets its own isolated logging context.
package http_middleware
