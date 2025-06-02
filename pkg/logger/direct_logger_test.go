package logger

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestDirectLogger_BasicLogging(t *testing.T) {
	// テスト用のバッファを作成
	var buf bytes.Buffer
	logger := NewDirectLogger()
	logger.SetOutput(&buf)
	logger.SetLevel(DebugLevel) // すべてのレベルを出力するように設定

	// 各ログレベルをテスト
	testCases := []struct {
		name    string
		logFunc func(string, ...interface{})
		level   string
		message string
		args    []interface{}
	}{
		{"Debug", logger.Debugf, "DEBUG", "debug message", nil},
		{"Info", logger.Infof, "INFO", "info message", nil},
		{"Warn", logger.Warnf, "WARN", "warn message", nil},
		{"Error", logger.Errorf, "ERROR", "error message", nil},
		{"Critical", logger.Criticalf, "CRITICAL", "critical message", nil},
		{"WithArgs", logger.Infof, "INFO", "message with %s and %d", []interface{}{"string", 42}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf.Reset()

			if tc.args != nil {
				tc.logFunc(tc.message, tc.args...)
			} else {
				tc.logFunc(tc.message)
			}

			output := buf.String()

			// ログが出力されていることを確認
			if output == "" {
				t.Error("Expected log output, got empty string")
				return
			}

			// ログレベルが含まれていることを確認
			if !strings.Contains(output, tc.level) {
				t.Errorf("Expected output to contain level %s, got: %s", tc.level, output)
			}

			// メッセージが含まれていることを確認
			expectedMessage := tc.message
			if tc.args != nil {
				expectedMessage = "message with string and 42"
			}
			if !strings.Contains(output, expectedMessage) {
				t.Errorf("Expected output to contain message %s, got: %s", expectedMessage, output)
			}

			// タイムスタンプが含まれていることを確認（RFC3339形式）
			if !strings.Contains(output, time.Now().Format("2006-01-02")) {
				t.Errorf("Expected output to contain timestamp, got: %s", output)
			}
		})
	}
}

func TestDirectLogger_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := NewDirectLogger()
	logger.SetOutput(&buf)

	testCases := []struct {
		name         string
		setLevel     LogLevel
		logLevel     LogLevel
		logFunc      func(string, ...interface{})
		message      string
		shouldOutput bool
	}{
		{"DebugLevel_DebugLog", DebugLevel, DebugLevel, logger.Debugf, "debug message", true},
		{"InfoLevel_DebugLog", InfoLevel, DebugLevel, logger.Debugf, "debug message", false},
		{"InfoLevel_InfoLog", InfoLevel, InfoLevel, logger.Infof, "info message", true},
		{"WarnLevel_InfoLog", WarnLevel, InfoLevel, logger.Infof, "info message", false},
		{"WarnLevel_WarnLog", WarnLevel, WarnLevel, logger.Warnf, "warn message", true},
		{"ErrorLevel_WarnLog", ErrorLevel, WarnLevel, logger.Warnf, "warn message", false},
		{"ErrorLevel_ErrorLog", ErrorLevel, ErrorLevel, logger.Errorf, "error message", true},
		{"CriticalLevel_ErrorLog", CriticalLevel, ErrorLevel, logger.Errorf, "error message", false},
		{"CriticalLevel_CriticalLog", CriticalLevel, CriticalLevel, logger.Criticalf, "critical message", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf.Reset()
			logger.SetLevel(tc.setLevel)

			tc.logFunc(tc.message)
			output := buf.String()

			if tc.shouldOutput {
				if output == "" {
					t.Errorf("Expected log output for level %s with min level %s, got empty string",
						tc.logLevel.String(), tc.setLevel.String())
				}
				if !strings.Contains(output, tc.message) {
					t.Errorf("Expected output to contain message %s, got: %s", tc.message, output)
				}
			} else {
				if output != "" {
					t.Errorf("Expected no log output for level %s with min level %s, got: %s",
						tc.logLevel.String(), tc.setLevel.String(), output)
				}
			}
		})
	}
}

func TestDirectLogger_SetLevelFromString(t *testing.T) {
	var buf bytes.Buffer
	logger := NewDirectLogger()
	logger.SetOutput(&buf)

	testCases := []struct {
		levelString   string
		expectedLevel LogLevel
		testLogFunc   func(string, ...interface{})
		shouldOutput  bool
	}{
		{"DEBUG", DebugLevel, logger.Debugf, true},
		{"INFO", InfoLevel, logger.Debugf, false},
		{"WARN", WarnLevel, logger.Infof, false},
		{"ERROR", ErrorLevel, logger.Warnf, false},
		{"CRITICAL", CriticalLevel, logger.Errorf, false},
		{"INVALID", InfoLevel, logger.Debugf, false}, // デフォルトでInfoLevel
	}

	for _, tc := range testCases {
		t.Run(tc.levelString, func(t *testing.T) {
			buf.Reset()
			logger.SetLevelFromString(tc.levelString)

			tc.testLogFunc("test message")
			output := buf.String()

			if tc.shouldOutput {
				if output == "" {
					t.Errorf("Expected log output when setting level to %s, got empty string", tc.levelString)
				}
			} else {
				if output != "" {
					t.Errorf("Expected no log output when setting level to %s, got: %s", tc.levelString, output)
				}
			}
		})
	}
}

func TestDirectLogger_SetOutput(t *testing.T) {
	logger := NewDirectLogger()

	// 最初のバッファ
	var buf1 bytes.Buffer
	logger.SetOutput(&buf1)
	logger.Infof("message to buffer 1")

	if buf1.String() == "" {
		t.Error("Expected output in buffer 1, got empty string")
	}
	if !strings.Contains(buf1.String(), "message to buffer 1") {
		t.Errorf("Expected buffer 1 to contain message, got: %s", buf1.String())
	}

	// 2番目のバッファに切り替え
	var buf2 bytes.Buffer
	logger.SetOutput(&buf2)
	logger.Infof("message to buffer 2")

	// buf1には新しいメッセージが追加されていないことを確認
	if strings.Contains(buf1.String(), "message to buffer 2") {
		t.Error("Buffer 1 should not contain message sent to buffer 2")
	}

	// buf2に新しいメッセージが出力されていることを確認
	if buf2.String() == "" {
		t.Error("Expected output in buffer 2, got empty string")
	}
	if !strings.Contains(buf2.String(), "message to buffer 2") {
		t.Errorf("Expected buffer 2 to contain message, got: %s", buf2.String())
	}
}

func TestDirectLogger_ConcurrentSafety(t *testing.T) {
	var buf bytes.Buffer
	logger := NewDirectLogger()
	logger.SetOutput(&buf)
	logger.SetLevel(DebugLevel)

	const numGoroutines = 100
	const messagesPerGoroutine = 10

	// 並行してログを出力
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer func() { done <- true }()

			for j := 0; j < messagesPerGoroutine; j++ {
				logger.Infof("goroutine %d message %d", goroutineID, j)
			}
		}(i)
	}

	// すべてのgoroutineの完了を待つ
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	output := buf.String()

	// 出力が空でないことを確認
	if output == "" {
		t.Error("Expected log output from concurrent goroutines, got empty string")
		return
	}

	// 期待される総メッセージ数を確認
	expectedMessages := numGoroutines * messagesPerGoroutine
	actualMessages := strings.Count(output, "goroutine")

	if actualMessages != expectedMessages {
		t.Errorf("Expected %d messages, got %d", expectedMessages, actualMessages)
	}

	// 各goroutineからのメッセージが含まれていることを確認
	for i := 0; i < numGoroutines; i++ {
		expectedPattern := fmt.Sprintf("goroutine %d", i)
		if !strings.Contains(output, expectedPattern) {
			t.Errorf("Expected output to contain messages from goroutine %d", i)
		}
	}
}

func TestDirectLogger_ErrorCases(t *testing.T) {
	logger := NewDirectLogger()

	t.Run("NilOutput", func(t *testing.T) {
		// nilを出力先に設定した場合の動作をテスト
		// 実装によってはpanicするかもしれないが、現在の実装では問題ない
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Setting nil output should not panic, got: %v", r)
			}
		}()

		logger.SetOutput(nil)
		// nilに出力してもpanicしないことを確認
		logger.Infof("test message")
	})

	t.Run("InvalidLogLevel", func(t *testing.T) {
		// 無効なログレベル文字列のテスト
		var buf bytes.Buffer
		logger.SetOutput(&buf)

		// 無効なレベルを設定（デフォルトのInfoLevelになるはず）
		logger.SetLevelFromString("INVALID_LEVEL")

		// Debugレベルは出力されないはず
		logger.Debugf("debug message")
		if buf.String() != "" {
			t.Error("Debug message should not be output with default INFO level")
		}

		buf.Reset()
		// Infoレベルは出力されるはず
		logger.Infof("info message")
		if buf.String() == "" {
			t.Error("Info message should be output with default INFO level")
		}
	})
}
