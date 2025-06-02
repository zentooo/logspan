package logger

import (
	"io"
	"os"
	"sync"

	"github.com/zentooo/logspan/pkg/formatter"
)

// Config holds the global configuration for the logger
type Config struct {
	// MinLevel sets the minimum log level for filtering
	MinLevel LogLevel

	// Output sets the output writer for logs
	Output io.Writer

	// EnableSourceInfo enables source file and line information in logs
	EnableSourceInfo bool

	// PrettifyJSON enables pretty-printing of JSON output with indentation
	PrettifyJSON bool
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		MinLevel:         InfoLevel,
		Output:           os.Stdout,
		EnableSourceInfo: false,
		PrettifyJSON:     false,
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

	// Update global direct logger with new configuration
	if directLogger, ok := D.(*DirectLogger); ok {
		directLogger.SetOutput(globalConfig.Output)
		directLogger.SetLevel(globalConfig.MinLevel)

		// Update formatter based on PrettifyJSON setting
		var jsonFormatter formatter.Formatter
		if globalConfig.PrettifyJSON {
			jsonFormatter = formatter.NewJSONFormatterWithIndent("  ")
		} else {
			jsonFormatter = formatter.NewJSONFormatter()
		}
		directLogger.SetFormatter(jsonFormatter)
	}

	initialized = true
}

// GetConfig returns a copy of the current global configuration
func GetConfig() Config {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return globalConfig
}

// IsInitialized returns whether the logger has been initialized
func IsInitialized() bool {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return initialized
}
