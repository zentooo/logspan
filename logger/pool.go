package logger

import (
	"sync"
	"time"
)

// LogEntry pool for memory optimization
var logEntryPool = sync.Pool{
	New: func() interface{} {
		return &LogEntry{}
	},
}

// Slice pool for []*LogEntry slices with different capacities
var slicePool = sync.Pool{
	New: func() interface{} {
		return make([]*LogEntry, 0, 16) // Default capacity of 16
	},
}

// getLogEntry retrieves a LogEntry from the pool
func getLogEntry() *LogEntry {
	return logEntryPool.Get().(*LogEntry)
}

// putLogEntry returns a LogEntry to the pool after clearing its fields
func putLogEntry(entry *LogEntry) {
	if entry == nil {
		return
	}

	// Clear all fields to prevent memory leaks and data contamination
	entry.Timestamp = time.Time{}
	entry.Level = ""
	entry.Message = ""
	entry.Funcname = ""
	entry.Filename = ""
	entry.Fileline = 0

	logEntryPool.Put(entry)
}

// getLogEntrySlice retrieves a []*LogEntry slice from the pool
func getLogEntrySlice() []*LogEntry {
	slice := slicePool.Get().([]*LogEntry)
	// Ensure the slice is empty but keep its capacity
	return slice[:0]
}

// putLogEntrySlice returns a []*LogEntry slice to the pool after clearing it
func putLogEntrySlice(slice []*LogEntry) {
	if slice == nil {
		return
	}

	// Return all LogEntry objects to their pool
	for _, entry := range slice {
		putLogEntry(entry)
	}

	// Clear the slice but keep capacity for reuse
	slice = slice[:0]

	// Only return slices with reasonable capacity to avoid memory bloat
	if cap(slice) <= 1024 { // Configurable threshold
		slicePool.Put(slice) //nolint:staticcheck // slice reuse is more important than avoiding small allocation
	}
	// If capacity is too large, let it be garbage collected
}

// Pool statistics for monitoring (optional)
type PoolStats struct {
	LogEntryPoolSize int
	SlicePoolSize    int
}

// GetPoolStats returns current pool statistics (for debugging/monitoring)
func GetPoolStats() PoolStats {
	// Note: sync.Pool doesn't provide direct size information
	// This is a placeholder for potential future monitoring
	return PoolStats{
		LogEntryPoolSize: -1, // Unknown - sync.Pool doesn't expose size
		SlicePoolSize:    -1, // Unknown - sync.Pool doesn't expose size
	}
}

// resetPools clears all pools (mainly for testing)
func resetPools() {
	logEntryPool = sync.Pool{
		New: func() interface{} {
			return &LogEntry{}
		},
	}
	slicePool = sync.Pool{
		New: func() interface{} {
			return make([]*LogEntry, 0, 16)
		},
	}
}
