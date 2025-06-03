package logger

import (
	"path/filepath"
	"runtime"
	"time"
)

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Funcname  string    `json:"funcname,omitempty"`
	Filename  string    `json:"filename,omitempty"`
	Fileline  int       `json:"fileline,omitempty"`
}

// SourceInfo holds source code location information
type SourceInfo struct {
	Funcname string
	Filename string
	Fileline int
}

// getSourceInfo retrieves source information using runtime.Caller
// skip indicates how many stack frames to skip (0 = current function, 1 = caller, etc.)
func getSourceInfo(skip int) *SourceInfo {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return &SourceInfo{}
	}

	// Get function name
	funcName := "unknown"
	if fn := runtime.FuncForPC(pc); fn != nil {
		funcName = fn.Name()
	}

	// Get base filename (without full path)
	filename := filepath.Base(file)

	return &SourceInfo{
		Funcname: funcName,
		Filename: filename,
		Fileline: line,
	}
}
