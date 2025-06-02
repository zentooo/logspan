package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
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

			// JSONとしてパースできることを確認
			var logData map[string]interface{}
			if err := json.Unmarshal([]byte(output), &logData); err != nil {
				t.Errorf("Expected valid JSON output, got error: %v, output: %s", err, output)
				return
			}

			// 構造化ログの基本構造を確認
			if logData["type"] != "request" {
				t.Errorf("Expected type 'request', got: %v", logData["type"])
			}

			runtime, ok := logData["runtime"].(map[string]interface{})
			if !ok {
				t.Error("Expected runtime section in log output")
				return
			}

			// severityが正しいことを確認
			if runtime["severity"] != tc.level {
				t.Errorf("Expected severity %s, got: %v", tc.level, runtime["severity"])
			}

			// linesが配列で1つのエントリを持つことを確認
			lines, ok := runtime["lines"].([]interface{})
			if !ok {
				t.Error("Expected lines to be an array")
				return
			}

			if len(lines) != 1 {
				t.Errorf("Expected exactly 1 log entry, got %d", len(lines))
				return
			}

			// ログエントリの内容を確認
			entry, ok := lines[0].(map[string]interface{})
			if !ok {
				t.Error("Expected log entry to be an object")
				return
			}

			// メッセージが含まれていることを確認
			expectedMessage := tc.message
			if tc.args != nil {
				expectedMessage = "message with string and 42"
			}
			if entry["message"] != expectedMessage {
				t.Errorf("Expected message %s, got: %v", expectedMessage, entry["message"])
			}

			// レベルが含まれていることを確認
			if entry["level"] != tc.level {
				t.Errorf("Expected level %s, got: %v", tc.level, entry["level"])
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
					return
				}

				// JSONとしてパースできることを確認
				var logData map[string]interface{}
				if err := json.Unmarshal([]byte(output), &logData); err != nil {
					t.Errorf("Expected valid JSON output, got error: %v", err)
					return
				}

				// メッセージが含まれていることを確認
				runtime := logData["runtime"].(map[string]interface{})
				lines := runtime["lines"].([]interface{})
				entry := lines[0].(map[string]interface{})
				if entry["message"] != tc.message {
					t.Errorf("Expected message %s, got: %v", tc.message, entry["message"])
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

	// JSONとしてパースできることを確認
	var logData1 map[string]interface{}
	if err := json.Unmarshal([]byte(buf1.String()), &logData1); err != nil {
		t.Errorf("Expected valid JSON output in buffer 1, got error: %v", err)
	} else {
		runtime := logData1["runtime"].(map[string]interface{})
		lines := runtime["lines"].([]interface{})
		entry := lines[0].(map[string]interface{})
		if entry["message"] != "message to buffer 1" {
			t.Errorf("Expected buffer 1 to contain message, got: %v", entry["message"])
		}
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

	// JSONとしてパースできることを確認
	var logData2 map[string]interface{}
	if err := json.Unmarshal([]byte(buf2.String()), &logData2); err != nil {
		t.Errorf("Expected valid JSON output in buffer 2, got error: %v", err)
	} else {
		runtime := logData2["runtime"].(map[string]interface{})
		lines := runtime["lines"].([]interface{})
		entry := lines[0].(map[string]interface{})
		if entry["message"] != "message to buffer 2" {
			t.Errorf("Expected buffer 2 to contain message, got: %v", entry["message"])
		}
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

	// 期待される総メッセージ数を確認（各ログエントリは1行のJSONなので改行で分割）
	lines := strings.Split(strings.TrimSpace(output), "\n")
	expectedMessages := numGoroutines * messagesPerGoroutine

	if len(lines) != expectedMessages {
		t.Errorf("Expected %d messages, got %d", expectedMessages, len(lines))
	}

	// 各行がJSONとしてパースできることを確認
	messageCount := make(map[int]int)
	for _, line := range lines {
		if line == "" {
			continue
		}

		var logData map[string]interface{}
		if err := json.Unmarshal([]byte(line), &logData); err != nil {
			t.Errorf("Expected valid JSON output, got error: %v, line: %s", err, line)
			continue
		}

		runtime := logData["runtime"].(map[string]interface{})
		linesArray := runtime["lines"].([]interface{})
		entry := linesArray[0].(map[string]interface{})
		message := entry["message"].(string)

		// goroutine IDを抽出
		var goroutineID int
		if _, err := fmt.Sscanf(message, "goroutine %d", &goroutineID); err == nil {
			messageCount[goroutineID]++
		}
	}

	// 各goroutineからの期待されるメッセージ数を確認
	for i := 0; i < numGoroutines; i++ {
		if messageCount[i] != messagesPerGoroutine {
			t.Errorf("Expected %d messages from goroutine %d, got %d", messagesPerGoroutine, i, messageCount[i])
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

func TestDirectLogger_StructuredOutput(t *testing.T) {
	var buf bytes.Buffer
	logger := NewDirectLogger()
	logger.SetOutput(&buf)

	logger.Infof("test message")
	output := buf.String()

	// JSONとしてパースできることを確認
	var logData map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logData); err != nil {
		t.Errorf("Expected valid JSON output, got error: %v, output: %s", err, output)
		return
	}

	// 構造化ログの必須フィールドを確認
	expectedFields := []string{"type", "context", "runtime", "config"}
	for _, field := range expectedFields {
		if _, exists := logData[field]; !exists {
			t.Errorf("Expected field %s in log output", field)
		}
	}

	// runtime セクションの詳細確認
	runtime, ok := logData["runtime"].(map[string]interface{})
	if !ok {
		t.Error("Expected runtime to be an object")
		return
	}

	expectedRuntimeFields := []string{"severity", "startTime", "endTime", "elapsed", "lines"}
	for _, field := range expectedRuntimeFields {
		if _, exists := runtime[field]; !exists {
			t.Errorf("Expected field %s in runtime section", field)
		}
	}

	// elapsed が 0 であることを確認（direct loggerの場合）
	if runtime["elapsed"] != float64(0) {
		t.Errorf("Expected elapsed to be 0 for direct logger, got: %v", runtime["elapsed"])
	}

	// lines が配列で1つのエントリを持つことを確認
	lines, ok := runtime["lines"].([]interface{})
	if !ok {
		t.Error("Expected lines to be an array")
		return
	}

	if len(lines) != 1 {
		t.Errorf("Expected exactly 1 log entry, got %d", len(lines))
	}
}
