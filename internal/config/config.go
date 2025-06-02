package config

import (
	"github.com/zentooo/logspan/internal/logger"
)

// Config represents the configuration for the logger
type Config struct {
	// Level is the minimum log level that will be logged
	Level logger.LogLevel
	// EnableCaller enables logging of the caller's file and line number
	EnableCaller bool
	// EnableStacktrace enables logging of stack traces for error and critical levels
	EnableStacktrace bool
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Level:            logger.InfoLevel,
		EnableCaller:     true,
		EnableStacktrace: true,
	}
}
