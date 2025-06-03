package logger

import (
	"io"
	"os"
	"sync"

	"github.com/zentooo/logspan/pkg/formatter"
)

// Config holds the configuration for the logger
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
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		MinLevel:         InfoLevel,
		Output:           os.Stdout,
		EnableSourceInfo: false,
		PrettifyJSON:     false,
		MaxLogEntries:    0,         // No auto-flush by default
		LogType:          "request", // Default log type
		ErrorHandler:     nil,       // Use global error handler
	}
}

// Global configuration
var (
	globalConfig Config
	configMutex  sync.RWMutex
	initialized  bool
)

// Init initializes the global logger configuration
func Init(config Config) {
	configMutex.Lock()
	defer configMutex.Unlock()

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
		return DefaultConfig()
	}

	return globalConfig
}

// IsInitialized returns whether the logger has been initialized
func IsInitialized() bool {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return initialized
}
