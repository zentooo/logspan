package logger

import (
	"sync"
	"testing"
	"time"
)

func TestLogEntryPool(t *testing.T) {
	// Reset pools for clean test
	resetPools()

	// Test basic get/put operations
	entry1 := getLogEntry()
	if entry1 == nil {
		t.Fatal("getLogEntry() returned nil")
	}

	// Set some values
	entry1.Timestamp = time.Now()
	entry1.Level = "INFO"
	entry1.Message = "test message"
	entry1.Funcname = "test.func"
	entry1.Filename = "test.go"
	entry1.Fileline = 42

	// Return to pool
	putLogEntry(entry1)

	// Get another entry (should be the same object, but cleared)
	entry2 := getLogEntry()
	if entry2 == nil {
		t.Fatal("getLogEntry() returned nil after put")
	}

	// Verify fields are cleared
	if !entry2.Timestamp.IsZero() {
		t.Error("Timestamp should be cleared")
	}
	if entry2.Level != "" {
		t.Error("Level should be cleared")
	}
	if entry2.Message != "" {
		t.Error("Message should be cleared")
	}
	if entry2.Funcname != "" {
		t.Error("Funcname should be cleared")
	}
	if entry2.Filename != "" {
		t.Error("Filename should be cleared")
	}
	if entry2.Fileline != 0 {
		t.Error("Fileline should be cleared")
	}

	// Clean up
	putLogEntry(entry2)
}

func TestLogEntryPoolConcurrency(t *testing.T) {
	resetPools()

	const numGoroutines = 100
	const operationsPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Run concurrent get/put operations
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < operationsPerGoroutine; j++ {
				entry := getLogEntry()
				if entry == nil {
					t.Errorf("getLogEntry() returned nil in goroutine %d, operation %d", goroutineID, j)
					return
				}

				// Set some values
				entry.Level = "INFO"
				entry.Message = "concurrent test"

				// Return to pool
				putLogEntry(entry)
			}
		}(i)
	}

	wg.Wait()
}

func TestSlicePool(t *testing.T) {
	resetPools()

	// Test basic get/put operations
	slice1 := getLogEntrySlice()
	if slice1 == nil {
		t.Fatal("getLogEntrySlice() returned nil")
	}

	// Verify initial state
	if len(slice1) != 0 {
		t.Error("Initial slice should be empty")
	}
	if cap(slice1) == 0 {
		t.Error("Initial slice should have some capacity")
	}

	// Add some entries
	entry1 := getLogEntry()
	entry2 := getLogEntry()
	slice1 = append(slice1, entry1, entry2)

	if len(slice1) != 2 {
		t.Error("Slice should contain 2 entries")
	}

	// Return to pool
	putLogEntrySlice(slice1)

	// Get another slice (should be empty but may have capacity)
	slice2 := getLogEntrySlice()
	if slice2 == nil {
		t.Fatal("getLogEntrySlice() returned nil after put")
	}

	if len(slice2) != 0 {
		t.Error("Reused slice should be empty")
	}
}

func TestSlicePoolLargeCapacity(t *testing.T) {
	resetPools()

	// Create a slice with large capacity
	largeSlice := make([]*LogEntry, 0, 2000) // Larger than threshold (1024)

	// Add some entries
	for i := 0; i < 10; i++ {
		entry := getLogEntry()
		entry.Level = "INFO"
		largeSlice = append(largeSlice, entry)
	}

	// This should not be returned to pool due to large capacity
	putLogEntrySlice(largeSlice)

	// Get a new slice - should be a fresh one, not the large one
	newSlice := getLogEntrySlice()
	if cap(newSlice) >= 2000 {
		t.Error("Large capacity slice should not be reused")
	}
}

func TestPutLogEntryNil(t *testing.T) {
	// Should not panic
	putLogEntry(nil)
}

func TestPutLogEntrySliceNil(t *testing.T) {
	// Should not panic
	putLogEntrySlice(nil)
}

func TestPoolStats(t *testing.T) {
	stats := GetPoolStats()

	// sync.Pool doesn't expose size, so we expect -1
	if stats.LogEntryPoolSize != -1 {
		t.Error("LogEntryPoolSize should be -1 (unknown)")
	}
	if stats.SlicePoolSize != -1 {
		t.Error("SlicePoolSize should be -1 (unknown)")
	}
}

// Benchmark tests to verify pool effectiveness
func BenchmarkLogEntryWithoutPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entry := &LogEntry{
			Timestamp: time.Now(),
			Level:     "INFO",
			Message:   "benchmark test message",
			Funcname:  "benchmark.func",
			Filename:  "benchmark.go",
			Fileline:  42,
		}
		// Use the entry fields to prevent optimization and unused write warnings
		if entry.Level == "" || entry.Message == "" || entry.Funcname == "" || entry.Filename == "" || entry.Fileline == 0 || entry.Timestamp.IsZero() {
			b.Fatal("Entry fields should be set")
		}
	}
}

func BenchmarkLogEntryWithPool(b *testing.B) {
	resetPools()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entry := getLogEntry()
		entry.Timestamp = time.Now()
		entry.Level = "INFO"
		entry.Message = "benchmark test message"
		entry.Funcname = "benchmark.func"
		entry.Filename = "benchmark.go"
		entry.Fileline = 42
		putLogEntry(entry)
	}
}

func BenchmarkSliceWithoutPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slice := make([]*LogEntry, 0, 16)
		for j := 0; j < 10; j++ {
			entry := &LogEntry{Level: "INFO"}
			slice = append(slice, entry)
		}
		_ = slice // Use the slice to prevent optimization
	}
}

func BenchmarkSliceWithPool(b *testing.B) {
	resetPools()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slice := getLogEntrySlice()
		for j := 0; j < 10; j++ {
			entry := getLogEntry()
			entry.Level = "INFO"
			slice = append(slice, entry)
		}
		putLogEntrySlice(slice)
	}
}
