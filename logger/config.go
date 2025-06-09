package logger

import (
	"io"
	"os"
	"sync"

	"github.com/zentooo/logspan/formatter"
)

// Config holds the configuration for the logger (internal use)
type Config struct {
	// MinLevel is the minimum log level for filtering
	MinLevel LogLevel

	// Output is the output destination for logs
	Output io.Writer

	// EnableSourceInfo enables source file information in log entries
	EnableSourceInfo bool

	// PrettifyJSON enables pretty-printed JSON output
	PrettifyJSON bool

	// MaxLogEntries is the maximum number of log entries before auto-flush
	// 0 means no limit (manual flush only)
	MaxLogEntries int

	// LogType is the type field value in log output
	LogType string

	// ErrorHandler is the error handler for logger errors
	// If nil, the global error handler will be used
	ErrorHandler ErrorHandler

	// FlushEmpty enables flushing even when there are no log entries
	// Useful for HTTP request logging to record request context even without logs
	FlushEmpty bool
}

// Option is a function that configures the logger
type Option func(*Config)

// WithMinLevel sets the minimum log level for filtering
func WithMinLevel(level LogLevel) Option {
	return func(c *Config) {
		c.MinLevel = level
	}
}

// WithOutput sets the output destination for logs
func WithOutput(output io.Writer) Option {
	return func(c *Config) {
		c.Output = output
	}
}

// WithSourceInfo enables or disables source file information in log entries
func WithSourceInfo(enabled bool) Option {
	return func(c *Config) {
		c.EnableSourceInfo = enabled
	}
}

// WithPrettifyJSON enables or disables pretty-printed JSON output
func WithPrettifyJSON(enabled bool) Option {
	return func(c *Config) {
		c.PrettifyJSON = enabled
	}
}

// WithMaxLogEntries sets the maximum number of log entries before auto-flush
// 0 means no limit (manual flush only)
func WithMaxLogEntries(count int) Option {
	return func(c *Config) {
		c.MaxLogEntries = count
	}
}

// WithLogType sets the type field value in log output
func WithLogType(logType string) Option {
	return func(c *Config) {
		c.LogType = logType
	}
}

// WithErrorHandler sets the error handler for logger errors
func WithErrorHandler(handler ErrorHandler) Option {
	return func(c *Config) {
		c.ErrorHandler = handler
	}
}

// WithFlushEmpty enables or disables flushing when there are no log entries
// This is useful for HTTP request logging to record request context even without logs
func WithFlushEmpty(enabled bool) Option {
	return func(c *Config) {
		c.FlushEmpty = enabled
	}
}

// defaultConfig returns a default configuration
func defaultConfig() Config {
	return Config{
		MinLevel:         InfoLevel,
		Output:           os.Stdout,
		EnableSourceInfo: false,
		PrettifyJSON:     false,
		MaxLogEntries:    0,         // No auto-flush by default
		LogType:          "request", // Default log type
		ErrorHandler:     nil,       // Use global error handler
		FlushEmpty:       true,      // Default to true
	}
}

// Global configuration
var (
	globalConfig Config
	configMutex  sync.RWMutex
	initialized  bool
)

// Init initializes the global logger configuration with functional options
func Init(options ...Option) {
	configMutex.Lock()
	defer configMutex.Unlock()

	// Start with default config
	config := defaultConfig()

	// Apply all options
	for _, option := range options {
		option(&config)
	}

	globalConfig = config
	initialized = true

	// Set global error handler if provided
	if config.ErrorHandler != nil {
		SetGlobalErrorHandler(config.ErrorHandler)
	}

	// Update global direct logger with new configuration
	if directLogger, ok := D.(*DirectLogger); ok {
		directLogger.SetOutput(globalConfig.Output)
		directLogger.SetLevel(globalConfig.MinLevel)

		// Update formatter based on PrettifyJSON setting (avoid calling createDefaultFormatter to prevent deadlock)
		var jsonFormatter formatter.Formatter
		if globalConfig.PrettifyJSON {
			jsonFormatter = formatter.NewJSONFormatterWithIndent("  ")
		} else {
			jsonFormatter = formatter.NewJSONFormatter()
		}
		directLogger.SetFormatter(jsonFormatter)
	}
}

// GetConfig returns the current global configuration
// If not initialized, returns default configuration
func GetConfig() Config {
	configMutex.RLock()
	defer configMutex.RUnlock()

	if !initialized {
		return defaultConfig()
	}

	return globalConfig
}

// IsInitialized returns whether the logger has been initialized
func IsInitialized() bool {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return initialized
}
