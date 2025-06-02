package main

import (
	"bytes"
	"fmt"

	"github.com/zentooo/logspan/pkg/logger"
)

func main() {
	// Create a new direct logger
	dl := logger.NewDirectLogger()

	// Use a buffer to capture output for demonstration
	var buf bytes.Buffer
	dl.SetOutput(&buf)

	fmt.Println("=== Testing with INFO level (default) ===")
	dl.Debugf("This debug message should NOT appear")
	dl.Infof("This info message should appear")
	dl.Warnf("This warning message should appear")
	dl.Errorf("This error message should appear")
	dl.Criticalf("This critical message should appear")

	fmt.Print(buf.String())
	buf.Reset()

	fmt.Println("\n=== Testing with DEBUG level ===")
	dl.SetLevelFromString("DEBUG")
	dl.Debugf("This debug message should appear now")
	dl.Infof("This info message should appear")
	dl.Warnf("This warning message should appear")

	fmt.Print(buf.String())
	buf.Reset()

	fmt.Println("\n=== Testing with ERROR level ===")
	dl.SetLevelFromString("ERROR")
	dl.Debugf("This debug message should NOT appear")
	dl.Infof("This info message should NOT appear")
	dl.Warnf("This warning message should NOT appear")
	dl.Errorf("This error message should appear")
	dl.Criticalf("This critical message should appear")

	fmt.Print(buf.String())
}
